package api

import (
	"net/http"
	"project_sdu/model"
	"project_sdu/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type BatchAPI interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	GetAll(c *gin.Context)
	GetActiveBatch(c *gin.Context)
}

type batchAPI struct {
	batchService service.BatchService
}

func NewBatchAPI(batchService service.BatchService) *batchAPI {
	return &batchAPI{batchService}
}

// ====================
// CREATE
// ====================
func (b *batchAPI) Create(c *gin.Context) {
	var batch model.Batch
	if err := c.ShouldBindJSON(&batch); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors:  map[string]string{"body": "Invalid JSON format"},
		})
		return
	}

	// Validation minimal
	if batch.Name == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors:  map[string]string{"name": "name is required"},
		})
		return
	}

	if batch.Jalur == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors:  map[string]string{"jalur": "jalur is required"},
		})
		return
	}

	if err := b.batchService.Create(&batch); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to create batch",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Batch created successfully",
		Data:    batch,
	})
}

// ====================
// UPDATE
// ====================
func (b *batchAPI) Update(c *gin.Context) {
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

	var batch model.Batch
	if err := c.ShouldBindJSON(&batch); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid JSON format",
			Errors:  map[string]string{"body": err.Error()},
		})
		return
	}

	if err := b.batchService.Update(id, &batch); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to update batch",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Batch updated successfully",
		Data:    batch,
	})
}

// ====================
// DELETE
// ====================
func (b *batchAPI) Delete(c *gin.Context) {
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

	if err := b.batchService.Delete(id); err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Message: "Batch cannot be deleted because it has associated students",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete batch",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Batch deleted successfully",
	})
}

// ====================
// GET BY ID
// ====================
func (b *batchAPI) GetByID(c *gin.Context) {
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

	batch, err := b.batchService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Success: false,
			Status:  http.StatusNotFound,
			Message: "Batch not found",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Batch retrieved successfully",
		Data:    batch,
	})
}

// ====================
// GET ALL
// ====================
func (b *batchAPI) GetAll(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "10")
	pageParam := c.DefaultQuery("page", "1")
	q := c.Query("q")

	limit, _ := strconv.Atoi(limitParam)
	page, _ := strconv.Atoi(pageParam)

	data, err := b.batchService.GetAll(limit, page, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve batches",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Batch list retrieved successfully",
		Data:    data,
		Meta: gin.H{
			"page":  page,
			"limit": limit,
		},
	})
}

// ====================
// GET Active Batch
// ====================
func (b *batchAPI) GetActiveBatch(c *gin.Context) {
	// idParam := c.Param("id")
	// id, err := strconv.Atoi(idParam)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, model.ErrorResponse{
	// 		Success: false,
	// 		Status:  http.StatusBadRequest,
	// 		Message: "Invalid ID",
	// 	})
	// 	return
	// }

	batch, err := b.batchService.GetActiveBatch()
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Success: false,
			Status:  http.StatusNotFound,
			Message: "Batch not found",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Batch retrieved successfully",
		Data:    batch,
	})
}
