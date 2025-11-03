package repository

import (
	"project_sdu/model"

	"gorm.io/gorm"
)

type PostRepository interface {
	Store(post *model.Post) error
	FindAll() ([]model.Post, error)
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

func (r *postRepository) FindAll() ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Find(&posts).Error
	return posts, err
}
