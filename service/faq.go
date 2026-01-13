package service

import (
	"project_sdu/model"
	"project_sdu/repository"
)

type FaqService interface {
	Create(faq *model.Faq) error
	GetAll() ([]model.Faq, error)
	Update(id int, faq *model.Faq) error
	Delete(id int) error
	GetByID(id int) (*model.Faq, error)
}

type faqService struct {
	faqRepo repository.FaqRepository
}

func NewFaqService(faqRepo repository.FaqRepository) FaqService {
	return &faqService{faqRepo}
}

func (s *faqService) Create(faq *model.Faq) error {
	return s.faqRepo.Create(faq)
}

func (s *faqService) GetAll() ([]model.Faq, error) {
	return s.faqRepo.GetAll()
}

func (s *faqService) Update(id int, faq *model.Faq) error {
	return s.faqRepo.Update(id, faq)
}

func (s *faqService) Delete(id int) error {
	return s.faqRepo.Delete(id)
}

func (s *faqService) GetByID(id int) (*model.Faq, error) {
	return s.faqRepo.GetByID(id)
}
