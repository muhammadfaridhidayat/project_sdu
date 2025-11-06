package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"project_sdu/api"
	"project_sdu/db"
	"project_sdu/middleware"
	"project_sdu/model"
	repo "project_sdu/repository"
	"project_sdu/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	// "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type APIHandler struct {
	UserAPIHandler    api.UserAPI
	StudentAPIHandler api.StudentAPI
	PostAPIHandler    api.PostAPI
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] \"%s %s %s\"\n",
			param.TimeStamp.Format(time.RFC822),
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	// --- CORS SETUP HERE ---
	allowedOrigins := []string{"http://localhost:3000"}

	envOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if envOrigins != "" {
		for _, origin := range strings.Split(envOrigins, ",") {
			o := strings.TrimSpace(origin)
			if o != "" {
				allowedOrigins = append(allowedOrigins, o)
			}
		}
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// --- END CORS ---

	// // Koneksi ke DB Postgres
	// dbCredential := model.Credential{
	// 	Host:         "localhost",
	// 	Username:     "postgres",
	// 	Password:     "farid123",
	// 	DatabaseName: "tes_sdu",
	// 	Port:         5432,
	// 	Schema:       "public",
	// }

	// database := db.NewDB()
	// conn, err := database.Connect(&dbCredential)
	// if err != nil {
	// 	panic(err)
	// }

	// // Migrasi tabel
	// conn.AutoMigrate(&model.User{}, &model.Student{})

	// // Daftarkan semua route dan handler
	// router = RunServer(router, conn)

	// fmt.Println("✅ Server is running on port 8080")
	// if err := router.Run(":8080"); err != nil {
	// 	panic(err)
	// }

	// Ambil DATABASE_URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		panic("DATABASE_URL tidak ditemukan. Pastikan environment variable sudah diatur di Railway.")
	}

	// Koneksi ke DB
	database := db.NewDB()
	conn, err := database.ConnectURL(databaseURL)
	if err != nil {
		panic(err)
	}

	// Buat ENUM type untuk PostCategory jika belum ada
	if err := database.CreatePostCategoryEnum(conn); err != nil {
		panic(fmt.Errorf("failed to create post_category enum: %w", err))
	}

	// Migrasi tabel (opsional, tergantung kebutuhan)
	conn.AutoMigrate(&model.User{}, &model.Student{}, &model.Post{})

	// Daftarkan semua route dan handler
	router = RunServer(router, conn)

	// Ambil PORT dari Railway
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback jika dijalankan lokal
	}

	fmt.Printf("✅ Server is running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		panic(err)
	}

}

func RunServer(r *gin.Engine, conn interface{}) *gin.Engine {
	// Pastikan conn bertipe *gorm.DB
	dbConn := conn.(*gorm.DB)

	// Inisialisasi Repository (pakai GORM)
	userRepo := repo.NewUserRepo(dbConn)
	studentRepo := repo.NewStudentRepo(dbConn)
	postRepo := repo.NewPostRepo(dbConn)

	// Inisialisasi Service
	userService := service.NewUserService(userRepo)
	studentService := service.NewStudentService(studentRepo)
	postService := service.NewPostService(postRepo)

	// Inisialisasi API Handler
	userAPIHandler := api.NewUserAPI(userService)
	studentAPIHandler := api.NewStudentAPI(studentService)
	postAPIHandler := api.NewPostAPI(postService)

	apiHandler := APIHandler{
		UserAPIHandler:    userAPIHandler,
		StudentAPIHandler: studentAPIHandler,
		PostAPIHandler:    postAPIHandler,
	}

	// ROUTES

	// User routes
	user := r.Group("/user")
	{
		user.POST("/register", apiHandler.UserAPIHandler.Register)
		user.POST("/login", apiHandler.UserAPIHandler.Login)
		user.POST("/logout", apiHandler.UserAPIHandler.Logout)

		user.Use(middleware.Auth())
		user.GET("/profile", apiHandler.UserAPIHandler.GetUserProfile)
	}

	// Student routes
	students := r.Group("/students")
	{
		students.Use(middleware.Auth())
		students.GET("", apiHandler.StudentAPIHandler.FetchAllStudent)
		students.GET("/:id", apiHandler.StudentAPIHandler.FetchStudentByID)
		students.POST("", apiHandler.StudentAPIHandler.StoreStudent)
		students.PUT("/:id", apiHandler.StudentAPIHandler.UpdateStudent)
		students.DELETE("/:id", apiHandler.StudentAPIHandler.DeleteStudent)
		students.GET("/class", apiHandler.StudentAPIHandler.FetchStudentWithClass)
	}

	// Post routes
	posts := r.Group("/posts")
	{
		posts.Use(middleware.Auth())
		posts.POST("", apiHandler.PostAPIHandler.CreatePost)
		posts.GET("", apiHandler.PostAPIHandler.GetPosts)
	}

	return r
}
