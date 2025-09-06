package services

import (
	"julia/internal/models"
	"julia/internal/repositories"
)

type MakeupService interface {
	GetAllMakeups() ([]*models.Makeup, error)
	GetMakeupsByDate(date string) ([]*models.Makeup, error)
	GetMakeupsByMonth(yearMonth string) ([]*models.Makeup, error)
	GetMakeupsByUser(userID string) ([]*models.Makeup, error)
	GetMakeupsByUserAndDate(userID, date string) ([]*models.Makeup, error)
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

func (s *makeupService) GetAllMakeups() ([]*models.Makeup, error) {
	return s.makeupRepo.GetAllMakeups()
}

func (s *makeupService) GetMakeupsByDate(date string) ([]*models.Makeup, error) {
	return s.makeupRepo.GetMakeupsByDate(date)
}

func (s *makeupService) GetMakeupsByMonth(yearMonth string) ([]*models.Makeup, error) {
	return s.makeupRepo.GetMakeupsByMonth(yearMonth)
}

func (s *makeupService) GetMakeupsByUser(userID string) ([]*models.Makeup, error) {
	return s.makeupRepo.GetMakeupsByUser(userID)
}

func (s *makeupService) GetMakeupsByUserAndDate(userID, date string) ([]*models.Makeup, error) {
	return s.makeupRepo.GetMakeupsByUserAndDate(userID, date)
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
