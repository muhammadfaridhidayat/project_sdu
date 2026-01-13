package repository

import (
	"project_sdu/model"

	"gorm.io/gorm"
)

type RequirementRepository interface {
	Create(requirement *model.Requirement) error
	GetAll() ([]model.Requirement, error)
	Update(id int, requirement *model.Requirement) error
	Delete(id int) error
	GetByID(id int) (*model.Requirement, error)
}

type requirementRepository struct {
	db *gorm.DB
}

func NewRequirementRepository(db *gorm.DB) RequirementRepository {
	return &requirementRepository{db}
}

func (r *requirementRepository) Create(requirement *model.Requirement) error {
	return r.db.Create(requirement).Error
}

func (r *requirementRepository) GetAll() ([]model.Requirement, error) {
	var requirements []model.Requirement
	err := r.db.Order("created_at ASC").Find(&requirements).Error
	return requirements, err
}

func (r *requirementRepository) Update(id int, requirement *model.Requirement) error {
	return r.db.Model(&model.Requirement{}).Where("id = ?", id).Updates(requirement).Error
}

func (r *requirementRepository) Delete(id int) error {
	return r.db.Delete(&model.Requirement{}, id).Error
}

func (r *requirementRepository) GetByID(id int) (*model.Requirement, error) {
	var requirement model.Requirement
	err := r.db.First(&requirement, id).Error
	if err != nil {
		return nil, err
	}
	return &requirement, nil
}
