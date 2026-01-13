package api

import (
	"net/http"
	"project_sdu/model"
	"project_sdu/repository"
	"project_sdu/service"

	"github.com/gin-gonic/gin"
)

type PPDBAPI interface {
	Register(c *gin.Context)
}

type ppdbAPI struct {
	studentService service.StudentService
}

func NewPPDBAPI(studentService service.StudentService) *ppdbAPI {
	return &ppdbAPI{studentService}
}

// ====================
// REGISTER (PUBLIC)
// ====================
func (p *ppdbAPI) Register(c *gin.Context) {
	var student model.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Format data tidak valid",
			Errors: map[string]string{
				"body": "Invalid JSON format",
			},
		})
		return
	}

	// Minimal validation
	errorsMap := make(map[string]string)
	if student.FullName == "" {
		errorsMap["full_name"] = "Nama lengkap wajib diisi"
	}
	if student.Gender == "" {
		errorsMap["gender"] = "Jenis kelamin wajib diisi"
	}
	// Add more as needed...

	if len(errorsMap) > 0 {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validasi gagal",
			Errors:  errorsMap,
		})
		return
	}

	if err := p.studentService.RegisterPPDB(&student); err != nil {
		switch err {
		case repository.ErrNIKExists:
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Message: "NIK sudah terdaftar",
				Errors:  map[string]string{"nik": "NIK already exists"},
			})
			return

		case repository.ErrNISNExists:
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Message: "NISN sudah terdaftar",
				Errors:  map[string]string{"nisn": "NISN already exists"},
			})
			return
		}

		c.JSON(http.StatusBadRequest, model.ErrorResponse{ // StatusBadRequest often better for business logic errors like "date invalid"
			Success: false,
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Pendaftaran berhasil! Data Anda telah kami terima.",
		Data:    student,
	})
}
