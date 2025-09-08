package services

import (
	"julia/internal/models"
	"julia/internal/repositories"
)

type ClassService interface {
	GetAllClasses() ([]*models.Class, error)
	GetClassByID(classID int64) (*models.Class, error)
	GetUsersByClassID(classID int64) ([]*models.User, error)
	CreateClass(class *models.Class) error
	UpdateClass(class *models.Class) error
	DeleteClass(classID int64) error
}

type classService struct {
	classRepo repositories.ClassRepository
	userRepo  repositories.UserRepository
}

func NewClassService(classRepo repositories.ClassRepository, userRepo repositories.UserRepository) ClassService {
	return &classService{classRepo: classRepo, userRepo: userRepo}
}

func (s *classService) GetAllClasses() ([]*models.Class, error) {
	return s.classRepo.GetAllClasses()
}

func (s *classService) GetClassByID(classID int64) (*models.Class, error) {
	return s.classRepo.GetClassByID(classID)
}

func (s *classService) GetUsersByClassID(classID int64) ([]*models.User, error) {
	return s.userRepo.GetUsersByClassID(classID)
}

func (s *classService) CreateClass(class *models.Class) error {
	return s.classRepo.CreateClass(class)
}

func (s *classService) UpdateClass(class *models.Class) error {
	return s.classRepo.UpdateClass(class)
}

func (s *classService) DeleteClass(classID int64) error {
	return s.classRepo.DeleteClass(classID)
}
