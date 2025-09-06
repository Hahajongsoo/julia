package services

import (
	"julia/internal/models"
	"julia/internal/repositories"
)

type ExamPeriodService interface {
	GetAllExamPeriods() ([]*models.ExamPeriod, error)
	GetExamPeriodByClassID(classID int64) ([]*models.ExamPeriod, error)
	CreateExamPeriod(examPeriod *models.ExamPeriod) error
	UpdateExamPeriod(examPeriod *models.ExamPeriod) error
	DeleteExamPeriod(examPeriodID int64) error
}

type examPeriodService struct {
	examPeriodRepo repositories.ExamPeriodRepository
}

func NewExamPeriodService(examPeriodRepo repositories.ExamPeriodRepository) ExamPeriodService {
	return &examPeriodService{examPeriodRepo: examPeriodRepo}
}

func (s *examPeriodService) GetAllExamPeriods() ([]*models.ExamPeriod, error) {
	return s.examPeriodRepo.GetAllExamPeriods()
}

func (s *examPeriodService) GetExamPeriodByClassID(classID int64) ([]*models.ExamPeriod, error) {
	return s.examPeriodRepo.GetExamPeriodByClassID(classID)
}

func (s *examPeriodService) CreateExamPeriod(examPeriod *models.ExamPeriod) error {
	return s.examPeriodRepo.CreateExamPeriod(examPeriod)
}

func (s *examPeriodService) UpdateExamPeriod(examPeriod *models.ExamPeriod) error {
	return s.examPeriodRepo.UpdateExamPeriod(examPeriod)
}

func (s *examPeriodService) DeleteExamPeriod(examPeriodID int64) error {
	return s.examPeriodRepo.DeleteExamPeriod(examPeriodID)
}
