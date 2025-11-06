package repository

import (
	"project_sdu/model"

	"gorm.io/gorm"
)

type PostRepository interface {
	Store(post *model.Post) error
	FindAll(page, limit int, q string) ([]model.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepo(db *gorm.DB) *postRepository {
	return &postRepository{db}
}

func (r *postRepository) Store(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindAll(page, limit int, q string) ([]model.Post, error) {
	var posts []model.Post
	offset := (page - 1) * limit
	query := r.db.Order("created_at DESC").Limit(limit).Offset(offset)

	if q != "" {
		query = query.Where("title ILIKE ?", "%"+q+"%")
	}

	err := query.Find(&posts).Error
	return posts, err
}
