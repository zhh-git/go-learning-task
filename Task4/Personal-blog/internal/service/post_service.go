package service

import (
	"Personal-blog/internal/model"
	"Personal-blog/internal/pkg/logger"
	"Personal-blog/internal/repository"
)

type PostService struct {
	postRepo *repository.PostRepoImpl
}

func NewPostRepoService() *PostService {
	return &PostService{
		postRepo: repository.NewPostRepo(),
	}
}

func (s *PostService) CreatePost(post *model.Post) error {
	err := s.postRepo.Create(post)
	if err != nil {
		logger.Error("Failed to create post:", err)
		return err
	}
	return nil
}

func (s *PostService) GetPostByID(id uint) (*model.Post, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		logger.Error("Failed to get post by ID:", err)
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetAllPosts(filter *model.Post) ([]*model.Post, error) {
	posts, err := s.postRepo.FindAll(filter)
	if err != nil {
		logger.Error("Failed to get all posts:", err)
		return nil, err
	}
	return posts, nil
}

func (s *PostService) UpdatePost(post *model.Post) error {
	err := s.postRepo.Update(post)
	if err != nil {
		logger.Error("Failed to update post:", err)
		return err
	}
	return nil
}

func (s *PostService) DeletePost(id uint) error {
	err := s.postRepo.Delete(id)
	if err != nil {
		logger.Error("Failed to delete post:", err)
		return err
	}
	return nil
}
