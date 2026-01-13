package main

import (
	"fmt"
	"log"
	"os"
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
	UserAPIHandler       api.UserAPI
	StudentAPIHandler    api.StudentAPI
	ParentAPIHandler     api.ParentAPI
	PostAPIHandler       api.PostAPI
	CurriculumAPIHandler api.CurriculumAPI
	FacilityAPIHandler   api.FacilityAPI
	BatchAPIHandler      api.BatchAPI
	DashboardAPIHanlder  api.DashboardAPI
	PPDBAPIHandler       api.PPDBAPI
	RequirementAPIHandler api.RequirementAPI
	FaqAPIHandler        api.FaqAPI
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
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("CORS_ALLOWED_ORIGINS")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// --- END CORS ---

	// Get DATABASE_URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		panic("DATABASE_URL tidak ditemukan.")
	}

	// Connect to DB
	database := db.NewDB()
	conn, err := database.ConnectURL(databaseURL)
	if err != nil {
		panic(err)
	}

	// Migration
	conn.AutoMigrate(
		&model.User{}, &model.Student{}, &model.Parent{}, &model.Post{}, &model.Curriculum{}, &model.Facility{}, &model.Batch{}, &model.Requirement{}, &model.Faq{},
	)
	
	// Seed
	SeedRequirements(conn)
	SeedFaqs(conn)

	// Route
	router = RunServer(router, conn)

	// Get Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("✅ Server is running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		panic(err)
	}

}

func RunServer(r *gin.Engine, conn *gorm.DB) *gin.Engine {
	dbConn := conn

	userRepo := repo.NewUserRepo(dbConn)
	studentRepo := repo.NewStudentRepo(dbConn)
	parentRepo := repo.NewParentRepo(dbConn)
	postRepo := repo.NewPostRepo(dbConn)
	curriculumRepo := repo.NewCurriculumRepository(dbConn)
	facilityRepo := repo.NewFacilityRepository(dbConn)
	batchRepo := repo.NewBatchRepository(dbConn)
	requirementRepo := repo.NewRequirementRepository(dbConn)
	faqRepo := repo.NewFaqRepository(dbConn)

	userService := service.NewUserService(userRepo)
	studentService := service.NewStudentService(studentRepo, parentRepo, batchRepo)
	parentService := service.NewParentService(parentRepo)
	postService := service.NewPostService(postRepo)
	curriculumService := service.NewCurriculumService(curriculumRepo)
	facilityService := service.NewfacilityService(facilityRepo)
	batchService := service.NewBatchService(batchRepo)
	dashboardService := service.NewDashboardService(studentRepo, postRepo, batchRepo)
	requirementService := service.NewRequirementService(requirementRepo)
	faqService := service.NewFaqService(faqRepo)

	userAPIHandler := api.NewUserAPI(userService)
	studentAPIHandler := api.NewStudentAPI(studentService)
	parentAPIHandler := api.NewParentAPI(parentService)
	postAPIHandler := api.NewPostAPI(postService)
	curriculumAPIHandler := api.NewCurriculumAPI(curriculumService)
	facilityAPIHandler := api.NewFacilityAPI(facilityService)
	batchAPIHandler := api.NewBatchAPI(batchService)
	dashboardAPIHanlder := api.NewDashboardAPI(dashboardService)
	ppdbAPIHandler := api.NewPPDBAPI(studentService)
	requirementAPIHandler := api.NewRequirementAPI(requirementService)
	faqAPIHandler := api.NewFaqAPI(faqService)

	apiHandler := APIHandler{
		UserAPIHandler:       userAPIHandler,
		StudentAPIHandler:    studentAPIHandler,
		ParentAPIHandler:     parentAPIHandler,
		PostAPIHandler:       postAPIHandler,
		CurriculumAPIHandler: curriculumAPIHandler,
		FacilityAPIHandler:   facilityAPIHandler,
		BatchAPIHandler:      batchAPIHandler,
		DashboardAPIHanlder:  dashboardAPIHanlder,
		PPDBAPIHandler:       ppdbAPIHandler,
		RequirementAPIHandler: requirementAPIHandler,
		FaqAPIHandler:        faqAPIHandler,
	}

	// ROUTES //

	// User routes
	user := r.Group("/user")
	{
		user.POST("/register", apiHandler.UserAPIHandler.Register)
		user.POST("/login", apiHandler.UserAPIHandler.Login)
		user.POST("/logout", apiHandler.UserAPIHandler.Logout)

		user.Use(middleware.Auth())
		user.GET("/profile", apiHandler.UserAPIHandler.GetUserProfile)
	}

	// PPDB routes
	ppdb := r.Group("/ppdb")
	{
		ppdb.POST("/add", apiHandler.PPDBAPIHandler.Register)
	}

	// Student routes
	student := r.Group("/student")
	{
		student.Use(middleware.Auth())
		student.POST("/add", apiHandler.StudentAPIHandler.CreateStudent)
		student.POST("/bulk-add", apiHandler.StudentAPIHandler.CreateManyStudents)
		student.GET("/get/:id", apiHandler.StudentAPIHandler.GetStudentByID)
		student.GET("/get-all", apiHandler.StudentAPIHandler.GetAllStudents)
		student.PUT("/update/:id", apiHandler.StudentAPIHandler.UpdateStudent)
		student.DELETE("/delete/:id", apiHandler.StudentAPIHandler.DeleteStudent)



	}

	// Parent routes
	parent := r.Group("/parent")
	{
		parent.Use(middleware.Auth())
		parent.POST("/add", apiHandler.ParentAPIHandler.CreateParent)
		parent.GET("/get-all", apiHandler.ParentAPIHandler.GetAllParents)
		parent.GET("/get/:id", apiHandler.ParentAPIHandler.GetParentByID)
		parent.PUT("/update/:id", apiHandler.ParentAPIHandler.UpdateParent)
		parent.DELETE("/delete/:id", apiHandler.ParentAPIHandler.DeleteParent)
	}

	// Post routes
	post := r.Group("/post")
	{
		post.GET("/get/:slug", apiHandler.PostAPIHandler.GetPostBySlug)
		post.GET("/get-all", apiHandler.PostAPIHandler.GetAllPosts)

		post.Use(middleware.Auth())
		post.POST("/add", apiHandler.PostAPIHandler.CreatePost)
		post.PUT("/update/:slug", apiHandler.PostAPIHandler.UpdatePost)
		post.DELETE("/delete/:slug", apiHandler.PostAPIHandler.DeletePost)
	}

	// Curriculum routes
	extracurricular := r.Group("/curriculum")
	{
		extracurricular.GET("/get-all", apiHandler.CurriculumAPIHandler.GetAll)
		extracurricular.GET("/get/:id", apiHandler.CurriculumAPIHandler.GetByID)
		extracurricular.GET("/category/:category", apiHandler.CurriculumAPIHandler.GetByCategory)

		extracurricular.Use(middleware.Auth())
		extracurricular.POST("/add", apiHandler.CurriculumAPIHandler.Create)
		extracurricular.PUT("/update/:id", apiHandler.CurriculumAPIHandler.Update)
		extracurricular.DELETE("/delete/:id", apiHandler.CurriculumAPIHandler.Delete)
	}

	// Facility routes
	facility := r.Group("/facility")
	{
		facility.GET("/get-all", apiHandler.FacilityAPIHandler.GetAllFacilities)
		facility.GET("/get/:id", apiHandler.FacilityAPIHandler.GetFacilityByID)

		facility.Use(middleware.Auth())
		facility.POST("/add", apiHandler.FacilityAPIHandler.CreateFacility)
		facility.PUT("/update/:id", apiHandler.FacilityAPIHandler.UpdateFacility)
		facility.DELETE("/delete/:id", apiHandler.FacilityAPIHandler.DeleteFacility)
	}

	// Batch routes
	batch := r.Group("/batch")
	{
		batch.GET("/get-active", apiHandler.BatchAPIHandler.GetActiveBatch)
		batch.Use(middleware.Auth())
		batch.GET("/get-all", apiHandler.BatchAPIHandler.GetAll)
		batch.GET("/get/:id", apiHandler.BatchAPIHandler.GetByID)
		batch.POST("/add", apiHandler.BatchAPIHandler.Create)
		batch.PUT("/update/:id", apiHandler.BatchAPIHandler.Update)
		batch.DELETE("/delete/:id", apiHandler.BatchAPIHandler.Delete)
	}

	dashboard := r.Group("/dashboard")
	{
		dashboard.Use(middleware.Auth())
		dashboard.GET("/", apiHandler.DashboardAPIHanlder.GetDashboard)
	}

	// Requirement routes
	requirement := r.Group("/requirement")
	{
		requirement.GET("/get-all", apiHandler.RequirementAPIHandler.GetAll)
		requirement.GET("/get/:id", apiHandler.RequirementAPIHandler.GetByID)

		requirement.Use(middleware.Auth())
		requirement.POST("/add", apiHandler.RequirementAPIHandler.Create)
		requirement.PUT("/update/:id", apiHandler.RequirementAPIHandler.Update)
		requirement.DELETE("/delete/:id", apiHandler.RequirementAPIHandler.Delete)
	}

	// Faq routes
	faq := r.Group("/faq")
	{
		faq.GET("/get-all", apiHandler.FaqAPIHandler.GetAll)
		faq.GET("/get/:id", apiHandler.FaqAPIHandler.GetByID)

		faq.Use(middleware.Auth())
		faq.POST("/add", apiHandler.FaqAPIHandler.Create)
		faq.PUT("/update/:id", apiHandler.FaqAPIHandler.Update)
		faq.DELETE("/delete/:id", apiHandler.FaqAPIHandler.Delete)
	}

	return r
}

func SeedRequirements(db *gorm.DB) {
	var count int64
	db.Model(&model.Requirement{}).Count(&count)
	if count > 0 {
		return
	}

	requirements := []model.Requirement{
		{Description: "Mengisi formulir pendaftaran."},
		{Description: "Foto copy paspor ukuran 3x4 (2 lembar)."},
		{Description: "Fotokopi Akta Kelahiran."},
		{Description: "Fotokopi Kartu Keluarga (KK)."},
		{Description: "Sertifikat Prestasi (Jika ada)."},
	}

	for _, req := range requirements {
		db.Create(&req)
	}
	fmt.Println("✅ Default requirements seeded")
}

func SeedFaqs(db *gorm.DB) {
	var count int64
	db.Model(&model.Faq{}).Count(&count)
	if count > 0 {
		return
	}

	faqs := []model.Faq{
		{Question: "Apakah ada tes masuk?", Answer: "Ya, terdapat tes potensi akademik dan wawancara untuk calon siswa dan orang tua."},
		{Question: "Apakah menyediakan beasiswa?", Answer: "Kami menyediakan beasiswa prestasi dan beasiswa bagi siswa kurang mampu yang memenuhi syarat."},
		{Question: "Berapa biaya pendaftarannya?", Answer: "Biaya pendaftaran dapat dilihat pada menu rincian biaya saat mengisi formulir."},
	}

	for _, faq := range faqs {
		db.Create(&faq)
	}
	fmt.Println("✅ Default faqs seeded")
}
