package api

import (
	"net/http"
	"project_sdu/model"
	"project_sdu/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RequirementAPI interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
}

type requirementAPI struct {
	requirementService service.RequirementService
}

func NewRequirementAPI(requirementService service.RequirementService) RequirementAPI {
	return &requirementAPI{requirementService}
}

func (a *requirementAPI) Create(c *gin.Context) {
	var requirement model.Requirement
	if err := c.ShouldBindJSON(&requirement); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid JSON",
			Errors:  map[string]string{"body": err.Error()},
		})
		return
	}

	if err := a.requirementService.Create(&requirement); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to create requirement",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Requirement created successfully",
		Data:    requirement,
	})
}

func (a *requirementAPI) GetAll(c *gin.Context) {
	requirements, err := a.requirementService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch requirements",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Requirements fetched successfully",
		Data:    requirements,
	})
}

func (a *requirementAPI) Update(c *gin.Context) {
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

	var requirement model.Requirement
	if err := c.ShouldBindJSON(&requirement); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid JSON",
			Errors:  map[string]string{"body": err.Error()},
		})
		return
	}

	if err := a.requirementService.Update(id, &requirement); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to update requirement",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Requirement updated successfully",
		Data:    requirement,
	})
}

func (a *requirementAPI) Delete(c *gin.Context) {
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

	if err := a.requirementService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete requirement",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Requirement deleted successfully",
	})
}

func (a *requirementAPI) GetByID(c *gin.Context) {
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

	requirement, err := a.requirementService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch requirement",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Requirement fetched successfully",
		Data:    requirement,
	})
}
