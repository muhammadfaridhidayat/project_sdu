package service

import (
	"project_sdu/model"
	"project_sdu/repository"
)

type PostService interface {
	CreatePost(post *model.Post) error
	FindAll(page, limit int) ([]model.Post, error)
}

type postService struct {
	postRepository repository.PostRepository
}

func NewPostService(postRepository repository.PostRepository) *postService {
	return &postService{postRepository}
}

func (s *postService) CreatePost(post *model.Post) error {
	return s.postRepository.Store(post)
}

func (s *postService) FindAll(page, limit int) ([]model.Post, error) {
	return s.postRepository.FindAll(page, limit)
}
