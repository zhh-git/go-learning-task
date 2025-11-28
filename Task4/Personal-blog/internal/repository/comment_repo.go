package repository

import (
	"Personal-blog/internal/model"
	"Personal-blog/internal/pkg/db"
)

type CommentRepo interface {
	Create(comment *model.Comment) error
	FindAllByPostId(postId uint) ([]*model.Comment, error)
	Update(comment *model.Comment) error
	Delete(id uint) error
}

type CommentRepoImpl struct{}

func NewCommentRepo() *CommentRepoImpl {
	return &CommentRepoImpl{}
}

func (r *CommentRepoImpl) Create(comment *model.Comment) error {
	return db.DB.Create(comment).Error
}

func (r *CommentRepoImpl) FindAllByPostId(postId uint) ([]*model.Comment, error) {
	var comments []*model.Comment
	if err := db.DB.Where("post_id = ?", postId).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepoImpl) Update(comment *model.Comment) error {
	return db.DB.Save(comment).Error
}

func (r *CommentRepoImpl) Delete(id uint) error {
	return db.DB.Delete(&model.Comment{}, id).Error
}
