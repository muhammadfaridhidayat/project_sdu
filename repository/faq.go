package repository

import (
	"project_sdu/model"

	"gorm.io/gorm"
)

type FaqRepository interface {
	Create(faq *model.Faq) error
	GetAll() ([]model.Faq, error)
	Update(id int, faq *model.Faq) error
	Delete(id int) error
	GetByID(id int) (*model.Faq, error)
}

type faqRepository struct {
	db *gorm.DB
}

func NewFaqRepository(db *gorm.DB) FaqRepository {
	return &faqRepository{db}
}

func (r *faqRepository) Create(faq *model.Faq) error {
	return r.db.Create(faq).Error
}

func (r *faqRepository) GetAll() ([]model.Faq, error) {
	var faqs []model.Faq
	err := r.db.Order("created_at ASC").Find(&faqs).Error
	return faqs, err
}

func (r *faqRepository) Update(id int, faq *model.Faq) error {
	return r.db.Model(&model.Faq{}).Where("id = ?", id).Updates(faq).Error
}

func (r *faqRepository) Delete(id int) error {
	return r.db.Delete(&model.Faq{}, id).Error
}

func (r *faqRepository) GetByID(id int) (*model.Faq, error) {
	var faq model.Faq
	err := r.db.First(&faq, id).Error
	if err != nil {
		return nil, err
	}
	return &faq, nil
}
