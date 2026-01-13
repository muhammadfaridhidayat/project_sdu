package api

import (
	"net/http"
	"project_sdu/model"
	"project_sdu/repository"
	"project_sdu/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StudentAPI interface {
	CreateStudent(c *gin.Context)
	GetStudentByID(c *gin.Context)
	GetAllStudents(c *gin.Context)
	UpdateStudent(c *gin.Context)
	DeleteStudent(c *gin.Context)
	CreateManyStudents(c *gin.Context)
}

type studentAPI struct {
	studentService service.StudentService
}

func NewStudentAPI(studentService service.StudentService) *studentAPI {
	return &studentAPI{studentService}
}

// ====================
// CREATE STUDENT
// ====================
func (s *studentAPI) CreateStudent(c *gin.Context) {
	var student model.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors: map[string]string{
				"body": "Invalid JSON format",
			},
		})
		return
	}

	// Minimal validation
	errors := make(map[string]string)
	if student.FullName == "" {
		errors["full_name"] = "Full name is required"
	}
	if student.Gender == "" {
		errors["gender"] = "Gender is required"
	}
	if *student.Agama == "" {
		errors["agama"] = "Agama is required"
	}
	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors:  errors,
		})
		return
	}

	if err := s.studentService.CreateStudent(&student); err != nil {
		switch err {
		case repository.ErrNIKExists:
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Message: "Validation failed",
				Errors:  map[string]string{"nik": repository.ErrNIKExists.Error()},
			})
			return

		case repository.ErrNISNExists:
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Message: "Validation failed",
				Errors:  map[string]string{"nisn": repository.ErrNISNExists.Error()},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to create student",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Student created successfully",
		Data:    student,
	})
}

// ====================
// BULK ADD STUDENT
// ====================
func (s *studentAPI) CreateManyStudents(c *gin.Context) {
	var students []model.Student
	if err := c.ShouldBindJSON(&students); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, student := range students {
		if err := s.studentService.CreateStudent(&student); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "All students created successfully"})
}

// ====================
// GET STUDENT BY ID
// ====================
func (s *studentAPI) GetStudentByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid student ID",
		})
		return
	}

	student, err := s.studentService.GetStudentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Success: false,
			Status:  http.StatusNotFound,
			Message: "Student not found",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Student retrieved successfully",
		Data:    student,
	})
}

// ====================
// GET ALL STUDENTS
// ====================
func (s *studentAPI) GetAllStudents(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "10")
	pageParam := c.DefaultQuery("page", "1")
	batchparam := c.Query("batch")
	q := c.Query("q")
	isAcceptedParam := c.Query("is_accepted") // <-- TAMBAHAN

	limit, _ := strconv.Atoi(limitParam)
	page, _ := strconv.Atoi(pageParam)
	batch, _ := strconv.Atoi(batchparam)

	// Handle filter is_accepted
	var isAccepted *bool = nil
	if isAcceptedParam != "" {
		parsed, err := strconv.ParseBool(isAcceptedParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Message: "Invalid value for is_accepted (use true or false)",
			})
			return
		}
		isAccepted = &parsed
	}

	students, err := s.studentService.GetAllStudents(limit, page, q, &batch, isAccepted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve students",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Students retrieved successfully",
		Data:    students,
		Meta: gin.H{
			"limit":       limit,
			"page":        page,
			"batch":       batch,
			"is_accepted": isAccepted,
		},
	})
}

// ====================
// UPDATE STUDENT
// ====================
func (s *studentAPI) UpdateStudent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid student ID",
		})
		return
	}

	var student model.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid JSON format",
			Errors:  map[string]string{"body": err.Error()},
		})
		return
	}

	if err := s.studentService.UpdateStudent(id, &student); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to update student",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Student updated successfully",
		Data:    student,
	})
}

// ====================
// DELETE STUDENT
// ====================
func (s *studentAPI) DeleteStudent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid student ID",
		})
		return
	}

	if err := s.studentService.DeleteStudent(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete student",
			Errors:  map[string]string{"server": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Student deleted successfully",
	})
}
