package service

import (
	"project_sdu/model"
	"project_sdu/repository"
)

type RequirementService interface {
	Create(requirement *model.Requirement) error
	GetAll() ([]model.Requirement, error)
	Update(id int, requirement *model.Requirement) error
	Delete(id int) error
	GetByID(id int) (*model.Requirement, error)
}

type requirementService struct {
	requirementRepo repository.RequirementRepository
}

func NewRequirementService(requirementRepo repository.RequirementRepository) RequirementService {
	return &requirementService{requirementRepo}
}

func (s *requirementService) Create(requirement *model.Requirement) error {
	return s.requirementRepo.Create(requirement)
}

func (s *requirementService) GetAll() ([]model.Requirement, error) {
	return s.requirementRepo.GetAll()
}

func (s *requirementService) Update(id int, requirement *model.Requirement) error {
	return s.requirementRepo.Update(id, requirement)
}

func (s *requirementService) Delete(id int) error {
	return s.requirementRepo.Delete(id)
}

func (s *requirementService) GetByID(id int) (*model.Requirement, error) {
	return s.requirementRepo.GetByID(id)
}
