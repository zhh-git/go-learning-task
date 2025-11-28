package service

import (
	"Personal-blog/internal/model"
	"Personal-blog/internal/pkg/logger"
	"Personal-blog/internal/repository"
)

type CommentService struct {
	commentRepo *repository.CommentRepoImpl
}

func NewCommentService() *CommentService {
	return &CommentService{
		commentRepo: repository.NewCommentRepo(),
	}
}

func (s *CommentService) CreateComment(comment *model.Comment) error {
	err := s.commentRepo.Create(comment)
	if err != nil {
		logger.Error("Failed to create comment", err)
		return err
	}
	return nil
}

func (s *CommentService) GetCommentsByPostID(postId uint) ([]*model.Comment, error) {
	comments, err := s.commentRepo.FindAllByPostId(postId)
	if err != nil {
		logger.Error("Failed to get comments by post ID", err)
		return nil, err
	}
	return comments, nil
}

func (s *CommentService) UpdateComment(comment *model.Comment) error {
	err := s.commentRepo.Update(comment)
	if err != nil {
		logger.Error("Failed to update comment", err)
		return err
	}
	return nil
}

func (s *CommentService) DeleteComment(id uint) error {
	err := s.commentRepo.Delete(id)
	if err != nil {
		logger.Error("Failed to delete comment", err)
		return err
	}
	return nil
}
