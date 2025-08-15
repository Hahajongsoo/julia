package services

import (
	"julia/internal/models"
	"julia/internal/repositories"
	"julia/utils"
)

type UserService interface {
	GetUserByID(id string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(id string, user *models.User) error
	DeleteUser(id string) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByID(id string) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *userService) CreateUser(user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.userRepo.CreateUser(user)
}

func (s *userService) UpdateUser(id string, user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.userRepo.UpdateUser(id, user)
}

func (s *userService) DeleteUser(id string) error {
	return s.userRepo.DeleteUser(id)
}
