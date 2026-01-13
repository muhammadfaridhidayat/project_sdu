package api

import (
	"net/http"
	"project_sdu/model"
	"project_sdu/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CurriculumAPI interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	GetAll(c *gin.Context)
	GetByCategory(c *gin.Context)
}

type curriculumAPI struct {
	curriculumService service.CurriculumService
}

func NewCurriculumAPI(curriculumService service.CurriculumService) *curriculumAPI {
	return &curriculumAPI{curriculumService}
}

// ====================
// CREATE
// ====================
func (e *curriculumAPI) Create(c *gin.Context) {
	var ex model.Curriculum
	if err := c.ShouldBindJSON(&ex); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors:  map[string]string{"body": "Invalid JSON format"},
		})
		return
	}

	if ex.Name == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors:  map[string]string{"name": "name is required"},
		})
		return
	}

	if err := e.curriculumService.Create(&ex); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to create curriculum",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{
		Success: true,
		Status:  http.StatusCreated,
		Message: "curriculum created successfully",
		Data:    ex,
	})
}

// ====================
// UPDATE
// ====================
func (e *curriculumAPI) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
		return
	}

	var ex model.Curriculum
	if err := c.ShouldBindJSON(&ex); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid JSON format",
			Errors:  map[string]string{"body": err.Error()},
		})
		return
	}

	if err := e.curriculumService.Update(id, &ex); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to update curriculum",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "curriculum updated successfully",
		Data:    ex,
	})
}

// ====================
// DELETE
// ====================
func (e *curriculumAPI) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
		return
	}

	if err := e.curriculumService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete curriculum",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "curriculum deleted successfully",
	})
}

// ====================
// GET BY ID
// ====================
func (e *curriculumAPI) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
		return
	}

	ex, err := e.curriculumService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Success: false,
			Status:  http.StatusNotFound,
			Message: "curriculum not found",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "curriculum retrieved successfully",
		Data:    ex,
	})
}

// ====================
// GET ALL
// ====================
func (e *curriculumAPI) GetAll(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "10")
	pageParam := c.DefaultQuery("page", "1")
	category := c.Query("category")

	limit, _ := strconv.Atoi(limitParam)
	page, _ := strconv.Atoi(pageParam)

	var data []model.Curriculum
	var err error

	if category != "" {
		data, err = e.curriculumService.GetByCategory(category, limit, page)
	} else {
		data, err = e.curriculumService.GetAll(limit, page)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve curriculum",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	responseMeta := gin.H{
		"page":  page,
		"limit": limit,
	}

	if category != "" {
		responseMeta["category"] = category
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "curriculum retrieved successfully",
		Data:    data,
		Meta:    responseMeta,
	})
}

// ====================
// GET BY CATEGORY
// ====================
func (e *curriculumAPI) GetByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Category is required",
			Errors:  map[string]string{"category": "category cannot be empty"},
		})
		return
	}

	limitParam := c.DefaultQuery("limit", "10")
	pageParam := c.DefaultQuery("page", "1")

	limit, _ := strconv.Atoi(limitParam)
	page, _ := strconv.Atoi(pageParam)

	data, err := e.curriculumService.GetByCategory(category, limit, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve curriculum by category",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "curriculum by category retrieved successfully",
		Data:    data,
		Meta: gin.H{
			"page":     page,
			"limit":    limit,
			"category": category,
		},
	})
}
