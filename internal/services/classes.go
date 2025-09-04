package services

import (
	"julia/internal/models"
	"julia/internal/repositories"
)

type ClassService interface {
	GetAllClasses() ([]*models.Class, error)
	GetClassByID(classID int64) (*models.Class, error)
	CreateClass(class *models.Class) error
	UpdateClass(class *models.Class) error
	DeleteClass(classID int64) error
}

type classService struct {
	classRepo repositories.ClassRepository
}

func NewClassService(classRepo repositories.ClassRepository) ClassService {
	return &classService{classRepo: classRepo}
}

func (s *classService) GetAllClasses() ([]*models.Class, error) {
	return s.classRepo.GetAllClasses()
}

func (s *classService) GetClassByID(classID int64) (*models.Class, error) {
	return s.classRepo.GetClassByID(classID)
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
