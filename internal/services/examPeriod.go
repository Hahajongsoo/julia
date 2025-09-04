package services

import (
	"julia/internal/models"
	"julia/internal/repositories"
)

type ExamPeriodService interface {
	GetExamPeriodByClassID(classID int64) ([]*models.ExamPeriod, error)
	UpsertExamPeriod(examPeriod *models.ExamPeriod) error
	DeleteExamPeriod(examPeriodID int64) error
}

type examPeriodService struct {
	examPeriodRepo repositories.ExamPeriodRepository
}

func NewExamPeriodService(examPeriodRepo repositories.ExamPeriodRepository) ExamPeriodService {
	return &examPeriodService{examPeriodRepo: examPeriodRepo}
}

func (s *examPeriodService) GetExamPeriodByClassID(classID int64) ([]*models.ExamPeriod, error) {
	return s.examPeriodRepo.GetExamPeriodByClassID(classID)
}

func (s *examPeriodService) UpsertExamPeriod(examPeriod *models.ExamPeriod) error {
	return s.examPeriodRepo.UpsertExamPeriod(examPeriod)
}

func (s *examPeriodService) DeleteExamPeriod(examPeriodID int64) error {
	return s.examPeriodRepo.DeleteExamPeriod(examPeriodID)
}
