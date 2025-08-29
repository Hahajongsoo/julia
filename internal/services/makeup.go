package services

import (
	"julia/internal/models"
	"julia/internal/repositories"
)

type MakeupService interface {
	GetMakeupsByDate(date string) ([]*models.Makeup, error)
	GetMakeupsByIDandDate(userID, date string) ([]*models.Makeup, error)
	CreateMakeup(makeup *models.Makeup) error
	UpdateMakeup(userID, date, time string, makeup *models.Makeup) error
	DeleteMakeup(userID, date, time string) error
}

type makeupService struct {
	makeupRepo repositories.MakeupRepository
}

func NewMakeupService(makeupRepo repositories.MakeupRepository) MakeupService {
	return &makeupService{makeupRepo: makeupRepo}
}

func (s *makeupService) GetMakeupsByDate(date string) ([]*models.Makeup, error) {
	return s.makeupRepo.GetMakeupsByDate(date)
}

func (s *makeupService) GetMakeupsByIDandDate(userID, date string) ([]*models.Makeup, error) {
	return s.makeupRepo.GetMakeupsByIDandDate(userID, date)
}

func (s *makeupService) CreateMakeup(makeup *models.Makeup) error {
	return s.makeupRepo.CreateMakeup(makeup)
}

func (s *makeupService) UpdateMakeup(userID, date, time string, makeup *models.Makeup) error {
	return s.makeupRepo.UpdateMakeup(userID, date, time, makeup)
}

func (s *makeupService) DeleteMakeup(userID, date, time string) error {
	return s.makeupRepo.DeleteMakeup(userID, date, time)
}
