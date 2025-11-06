package api

import (
	"errors"
	"net/http"
	"project_sdu/model"
	"project_sdu/service"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PostAPI interface {
	CreatePost(c *gin.Context)
	GetPosts(c *gin.Context)
}

type postAPI struct {
	postService service.PostService
}

func NewPostAPI(postService service.PostService) *postAPI {
	return &postAPI{postService}
}

func (api *postAPI) GetPosts(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid page number. Please provide a positive integer.",
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid limit value. Please provide a positive integer.",
		})
		return
	}

	q := c.Query("q")

	posts, err := api.postService.FindAll(page, limit, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve posts due to an internal error",
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Posts retrieved successfully",
		Data: posts,
		Meta: gin.H{
			"page":  page,
			"limit": limit,
		},
	})
}

type PostRequest struct {
	Title       string              `json:"title" binding:"required,min=10"`
	Content     string              `json:"content" binding:"required,min=20"`
	Thumbnail   string              `json:"thumbnail"`
	Description *string             `json:"description"`
	Category    *model.PostCategory `json:"category" binding:"oneof=BERITA ARTIKEL INFORMASI"`
	Published   bool                `json:"published"`
	PublishedAt *string             `json:"publishedAt"`
}

func slugify(input string) string {
	input = strings.ToLower(strings.TrimSpace(input))
	re := regexp.MustCompile(`[^a-z0-9\s-]`)
	input = re.ReplaceAllString(input, "")
	re = regexp.MustCompile(`[\s\-_]+`)
	input = re.ReplaceAllString(input, "-")
	input = strings.Trim(input, "-")
	return input
}

func (api *postAPI) CreatePost(c *gin.Context) {
	var req PostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorMessages := make(map[string]string)
			customErrorMessages := map[string]string{
				"required": "This field is required",
				"min":      "This field is too short",
				"oneof":    "Invalid value for this field",
			}
			for _, e := range ve {
				fieldName := e.Field()
				tag := e.ActualTag()
				if customMessage, ok := customErrorMessages[tag]; ok {
					errorMessages[fieldName] = customMessage
				} else {
					errorMessages[fieldName] = tag
				}
			}
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Message: "Input validation failed",
				Errors:  errorMessages,
			})
			return
		}
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "The request body is invalid",
			Errors:  map[string]string{"error": err.Error()},
		})
		return
	}

	var thumbnailPath *string
	if req.Thumbnail != "" {
		thumbnailPath = &req.Thumbnail
	}

	post := model.Post{
		Title:       req.Title,
		Slug:        slugify(req.Title),
		Content:     req.Content,
		Thumbnail:   thumbnailPath,
		Description: req.Description,
		Category:    req.Category,
		Published:   req.Published,
	}

	if req.Published {
		if req.PublishedAt != nil && *req.PublishedAt != "" {
			parsedTime, err := time.Parse("2006-01-02T15:04", *req.PublishedAt)
			if err != nil {
				c.JSON(http.StatusBadRequest, model.ErrorResponse{
					Success: false,
					Status:  http.StatusBadRequest,
					Message: "The published date format is invalid. Please use YYYY-MM-DDTHH:MM.",
					Errors:  map[string]string{"publishedAt": err.Error()},
				})
				return
			}
			post.PublishedAt = &parsedTime
		} else {
			now := time.Now()
			post.PublishedAt = &now
		}
	} else {
		post.PublishedAt = nil
	}

	err := api.postService.CreatePost(&post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to create the post due to an internal error",
		})
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Post created successfully",
		Data:    post,
	})
}
