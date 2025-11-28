package service

import (
	"Personal-blog/internal/model"
	"Personal-blog/internal/pkg/logger"
	"Personal-blog/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepoImpl
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepo(),
	}
}

func (s *UserService) CreateUser(user *model.User) error {
	err := s.userRepo.Create(user)
	if err != nil {
		logger.Error("Failed to create user:", err)
		return err
	}
	return nil
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		logger.Error("Failed to get user by ID:", err)
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		logger.Error("Failed to get user by username:", err)
		return nil, err
	}
	return user, nil
}
func (s *UserService) UpdateUser(user *model.User) error {
	err := s.userRepo.Update(user)
	if err != nil {
		logger.Error("Failed to update user:", err)
		return err
	}
	return nil
}

func (s *UserService) DeleteUser(id uint) error {
	err := s.userRepo.Delete(id)
	if err != nil {
		logger.Error("Failed to delete user:", err)
		return err
	}
	return nil
}
