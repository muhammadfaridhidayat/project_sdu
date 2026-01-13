package api

import (
	"net/http"
	"project_sdu/model"
	"project_sdu/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FaqAPI interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
}

type faqAPI struct {
	faqService service.FaqService
}

func NewFaqAPI(faqService service.FaqService) FaqAPI {
	return &faqAPI{faqService}
}

func (a *faqAPI) Create(c *gin.Context) {
	var faq model.Faq
	if err := c.ShouldBindJSON(&faq); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid JSON",
			Errors:  map[string]string{"body": err.Error()},
		})
		return
	}

	if err := a.faqService.Create(&faq); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to create faq",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Faq created successfully",
		Data:    faq,
	})
}

func (a *faqAPI) GetAll(c *gin.Context) {
	faqs, err := a.faqService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch faqs",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Faqs fetched successfully",
		Data:    faqs,
	})
}

func (a *faqAPI) Update(c *gin.Context) {
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

	var faq model.Faq
	if err := c.ShouldBindJSON(&faq); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid JSON",
			Errors:  map[string]string{"body": err.Error()},
		})
		return
	}

	if err := a.faqService.Update(id, &faq); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to update faq",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Faq updated successfully",
		Data:    faq,
	})
}

func (a *faqAPI) Delete(c *gin.Context) {
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

	if err := a.faqService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete faq",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Faq deleted successfully",
	})
}

func (a *faqAPI) GetByID(c *gin.Context) {
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

	faq, err := a.faqService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch faq",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Faq fetched successfully",
		Data:    faq,
	})
}
