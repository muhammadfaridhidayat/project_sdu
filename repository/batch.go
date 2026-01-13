package repository

import (
	"project_sdu/model"

	"gorm.io/gorm"
)

type BatchRepository interface {
	Create(batch *model.Batch) error
	GetAll(limit, page int, q string) ([]model.Batch, error)
	GetActiveBatch() (*model.Batch, error)
	GetByID(id int) (*model.Batch, error)
	Update(id int, batch *model.Batch) error
	Delete(id int) error
	CountAll() (int, error)
}

type batchRepository struct {
	db *gorm.DB
}

func NewBatchRepository(db *gorm.DB) BatchRepository {
	return &batchRepository{db}
}

func (r *batchRepository) Create(batch *model.Batch) error {
	// Removed year check as requested
	if err := r.db.Create(batch).Error; err != nil {
		return err
	}
	return nil
}

func (r *batchRepository) GetAll(limit, page int, q string) ([]model.Batch, error) {
	var batches []model.Batch

	offset := (page - 1) * limit

	db := r.db

	// Filter search name
	if q != "" {
		db = db.Where("name ILIKE ?", "%"+q+"%")
	}

	err := db.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&batches).Error

	return batches, err
}

func (r *batchRepository) GetByID(id int) (*model.Batch, error) {
	var batch model.Batch
	err := r.db.
		Where("id = ?", id).
		First(&batch).Error

	if err != nil {
		return nil, err
	}

	return &batch, nil
}

func (r *batchRepository) GetActiveBatch() (*model.Batch, error) {
	var batch model.Batch

	err := r.db.
		Where("is_active = ?", true).
		First(&batch).Error

	if err != nil {
		return nil, err
	}

	return &batch, nil
}

func (r *batchRepository) Update(id int, batch *model.Batch) error {
	return r.db.Model(&model.Batch{}).
		Where("id = ?", id).
		Updates(batch).
		Error
}

func (r *batchRepository) Delete(id int) error {
	return r.db.Delete(&model.Batch{}, id).Error
}

func (r *batchRepository) CountAll() (int, error) {
	var count int64

	err := r.db.Model(&model.Batch{}).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

