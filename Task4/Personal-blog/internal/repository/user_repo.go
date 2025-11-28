package repository

import (
	"Personal-blog/internal/model"
	"Personal-blog/internal/pkg/db"
)

type UserRepo interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
}

type UserRepoImpl struct{}

func NewUserRepo() *UserRepoImpl {
	return &UserRepoImpl{}
}

func (r *UserRepoImpl) Create(user *model.User) error {
	return db.DB.Create(user).Error
}

func (r *UserRepoImpl) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := db.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepoImpl) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := db.DB.Where("username = ? ", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepoImpl) Update(user *model.User) error {
	return db.DB.Save(user).Error
}

func (r *UserRepoImpl) Delete(id uint) error {
	return db.DB.Delete(&model.User{}, id).Error
}
