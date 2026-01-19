package repository

import (
	"errors"
	"project_sdu/model"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrNIKExists  = errors.New("NIK sudah terdaftar")
	ErrNISNExists = errors.New("NISN sudah terdaftar")
)

type StudentRepository interface {
	Create(student *model.Student) error
	GetStudentsByBatchID(batchID int, limit int, page int, q string) ([]model.Student, error)
	GetByID(id int) (*model.Student, error)
	GetAll(limit int, page int, q string, batchID *int, isAccepted *bool) ([]model.Student, error)
	Update(id int, student *model.Student) error
	Delete(id int) error
	CountAll() (int, error)
	CountByBatchID(batchID int) (int, error)
	GetAllByBatchID(batchID int) ([]model.Student, error)
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepo(db *gorm.DB) StudentRepository {
	return &studentRepository{db}
}

func (r *studentRepository) Create(student *model.Student) error {
	err := r.db.Create(student).Error
	if err != nil {
		if strings.Contains(err.Error(), "idx_students_nik") {
			return ErrNIKExists
		}
		if strings.Contains(err.Error(), "idx_students_nisn") {
			return ErrNISNExists
		}
		return err
	}
	return nil
}

func (r *studentRepository) GetStudentsByBatchID(batchID int, limit int, page int, q string) ([]model.Student, error) {
	var students []model.Student

	offset := (page - 1) * limit
	db := r.db

	db = db.Where("batch_id = ?", batchID)

	if q != "" {
		db = db.Where("full_name ILIKE ?", "%"+q+"%")
	}

	err := db.
		Preload("Parent").
		Order("full_name ASC").
		Limit(limit).
		Offset(offset).
		Find(&students).Error

	if err != nil {
		return nil, err
	}

	return students, nil
}

func (r *studentRepository) GetByID(id int) (*model.Student, error) {
	var student model.Student
	err := r.db.
		Preload("Parent").
		Preload("Batch").
		First(&student, id).
		Error

	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (r *studentRepository) GetAll(limit int, page int, q string, batchID *int, isAccepted *bool) ([]model.Student, error) {
	var students []model.Student

	offset := (page - 1) * limit

	db := r.db

	// Filter search name
	if q != "" {
		db = db.Where("full_name ILIKE ?", "%"+q+"%")
	}

	// Filter based on batch id
	if batchID != nil && *batchID != 0 {
		db = db.Where("batch_id = ?", *batchID)
	}

	// Filter by is_accepted (true/false)
	if isAccepted != nil {
		db = db.Where("is_accepted = ?", *isAccepted)
	}

	err := db.
		Preload("Parent").
		Preload("Batch").
		Order("full_name ASC").
		Limit(limit).
		Offset(offset).
		Find(&students).
		Error

	if err != nil {
		return nil, err
	}

	return students, nil
}

func (r *studentRepository) Update(id int, student *model.Student) error {
	var existingStudent model.Student
	if err := r.db.Preload("Parent").First(&existingStudent, id).Error; err != nil {
		return err
	}

	if err := r.db.Model(&existingStudent).Updates(student).Error; err != nil {
		return err
	}

	if student.Parent != nil {
		if err := r.db.Model(&existingStudent.Parent).Updates(student.Parent).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *studentRepository) Delete(id int) error {
	return r.db.Delete(&model.Student{}, id).Error
}

func (r *studentRepository) CountAll() (int, error) {
	var count int64
	err := r.db.Model(&model.Student{}).Count(&count).Error
	return int(count), err
}

func (r *studentRepository) CountByBatchID(batchID int) (int, error) {
	var count int64
	err := r.db.Model(&model.Student{}).
		Where("batch_id = ?", batchID).
		Count(&count).Error
	return int(count), err
}

func (r *studentRepository) GetAllByBatchID(batchID int) ([]model.Student, error) {
	var students []model.Student

	err := r.db.
		Where("batch_id = ?", batchID).
		Preload("Parent").
		Preload("Batch").
		Order("full_name ASC").
		Find(&students).Error

	return students, err
}
