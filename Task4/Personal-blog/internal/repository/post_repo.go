package repository

import (
	"Personal-blog/internal/model"
	"Personal-blog/internal/pkg/db"
)

type PostRepo interface {
	Create(post *model.Post) error
	FindByID(id uint) (*model.Post, error)
	FindAll() ([]*model.Post, error)
	Update(post *model.Post) error
	Delete(id uint) error
}

type PostRepoImpl struct{}

func NewPostRepo() *PostRepoImpl {
	return &PostRepoImpl{}
}

func (r *PostRepoImpl) Create(post *model.Post) error {
	return db.DB.Create(post).Error
}

func (r *PostRepoImpl) FindByID(id uint) (*model.Post, error) {
	var post model.Post
	if err := db.DB.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepoImpl) FindAll(post *model.Post) ([]*model.Post, error) {
	var posts []*model.Post
	if err := db.DB.Where("title like ?", "%"+post.Title+"%").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepoImpl) Update(post *model.Post) error {
	return db.DB.Save(post).Error
}

func (r *PostRepoImpl) Delete(id uint) error {
	return db.DB.Delete(&model.Post{}, id).Error
}
