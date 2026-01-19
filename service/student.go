package service

import (
	"errors"
	"project_sdu/model"
	"project_sdu/repository"
	"time"
)

type StudentService interface {
	CreateStudent(student *model.Student) error
	RegisterPPDB(student *model.Student) error
	GetStudentByID(id int) (*model.Student, error)
	GetAllStudents(limit int, page int, q string, batchID *int, isAccepted *bool) ([]model.Student, error)
	UpdateStudent(id int, student *model.Student) error
	DeleteStudent(id int) error
	GetStudentsByBatchID(batchID int) ([]model.Student, error)
}

type studentService struct {
	studentRepo repository.StudentRepository
	parentRepo  repository.ParentRepository
	batchRepo   repository.BatchRepository
}

func NewStudentService(studentRepo repository.StudentRepository, parentRepo repository.ParentRepository, batchRepo repository.BatchRepository) StudentService {
	return &studentService{
		studentRepo: studentRepo,
		parentRepo:  parentRepo,
		batchRepo:   batchRepo,
	}
}

func (s *studentService) CreateStudent(student *model.Student) error {
	var parentCreated bool
	if student.Parent != nil {
		if err := s.parentRepo.Create(student.Parent); err != nil {
			return err
		}
		parentCreated = true
		student.ParentId = &student.Parent.ID
	}

	// If BatchId is not provided, try to find active batch (Admin convenience, or default behavior)
	if student.BatchId == nil || *student.BatchId == 0 {
		activeBatch, err := s.batchRepo.GetActiveBatch()
		if err == nil {
			student.BatchId = &activeBatch.ID
		}
		// If no active batch and no ID provided, we permit for Admin (it will be null)? 
		// Or maybe we just proceed.
	}
	student.Batch = nil

	if err := s.studentRepo.Create(student); err != nil {
		if parentCreated && student.Parent != nil {
			_ = s.parentRepo.Delete(student.Parent.ID)
		}
		return err
	}

	return nil
}

func (s *studentService) RegisterPPDB(student *model.Student) error {
	var parentCreated bool
	if student.Parent != nil {
		if err := s.parentRepo.Create(student.Parent); err != nil {
			return err
		}
		parentCreated = true
		student.ParentId = &student.Parent.ID
	}

	// check if there is an active batch
	activeBatch, err := s.batchRepo.GetActiveBatch()
	if err != nil {
		if parentCreated && student.Parent != nil {
			_ = s.parentRepo.Delete(student.Parent.ID)
		}
		return errors.New("mohon maaf, tidak ada gelombang pendaftaran yang aktif saat ini")
	}

	// validate the date
	now := time.Now()

	if activeBatch.StartDate == nil || activeBatch.EndDate == nil {
		// If dates are null but it is active, maybe we allow it? 
		// ORIGINAL ERROR was "batch has invalid..." so I will keep stricter check OR relax it if user wants to fix the data.
		// User said: "batchnya berdasarkan yang aktif".
		// I'll assume dates MUST be valid for PPDB.
		if parentCreated && student.Parent != nil {
			_ = s.parentRepo.Delete(student.Parent.ID)
		}
		return errors.New("konfigurasi gelombang pendaftaran tidak valid (tanggal mulai/selesai belum diatur)")
	}

	if now.Before(*activeBatch.StartDate) {
		if parentCreated && student.Parent != nil {
			_ = s.parentRepo.Delete(student.Parent.ID)
		}
		return errors.New("pendaftaran belum dibuka")
	}
	
	if now.After(*activeBatch.EndDate) {
		if parentCreated && student.Parent != nil {
			_ = s.parentRepo.Delete(student.Parent.ID)
		}
		return errors.New("pendaftaran sudah ditutup")
	}

	student.BatchId = &activeBatch.ID
	student.Batch = nil

	if err := s.studentRepo.Create(student); err != nil {
		if parentCreated && student.Parent != nil {
			_ = s.parentRepo.Delete(student.Parent.ID)
		}
		return err
	}

	student.Batch = activeBatch

	return nil
}

func (s *studentService) GetStudentByID(id int) (*model.Student, error) {
	student, err := s.studentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return student, nil
}

func (s *studentService) GetAllStudents(limit int, page int, q string, batchID *int, isAccepted *bool) ([]model.Student, error) {
	return s.studentRepo.GetAll(limit, page, q, batchID, isAccepted)
}

func (s *studentService) UpdateStudent(id int, student *model.Student) error {
	if err := s.studentRepo.Update(id, student); err != nil {
		return err
	}

	return nil
}

func (s *studentService) DeleteStudent(id int) error {
	if err := s.studentRepo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (s *studentService) GetStudentsByBatchID(batchID int) ([]model.Student, error) {
	return s.studentRepo.GetAllByBatchID(batchID)
}
