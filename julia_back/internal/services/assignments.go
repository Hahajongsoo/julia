package services

import (
	"julia/internal/models"
	"julia/internal/repositories"
)

type AssignmentService interface {
	GetAllAssignments() ([]*models.Assignment, error)
	GetAssignmentsByUserID(userID string) ([]*models.Assignment, error)
	GetAssignmentsByMakeupID(makeupID int64) ([]*models.Assignment, error)
	GetAssignmentWithClassID() (*models.ClassAssignments, error)
	UpsertAssignment(assignment *models.Assignment) error
	DeleteAssignment(assignmentID int64) error
}

type assignmentService struct {
	assignmentRepo repositories.AssignmentRepository
}

func NewAssignmentService(assignmentRepo repositories.AssignmentRepository) AssignmentService {
	return &assignmentService{assignmentRepo: assignmentRepo}
}

func (s *assignmentService) GetAllAssignments() ([]*models.Assignment, error) {
	return s.assignmentRepo.GetAllAssignments()
}

func (s *assignmentService) GetAssignmentsByUserID(userID string) ([]*models.Assignment, error) {
	return s.assignmentRepo.GetAssignmentsByUserID(userID)
}

func (s *assignmentService) GetAssignmentsByMakeupID(makeupID int64) ([]*models.Assignment, error) {
	return s.assignmentRepo.GetAssignmentsByMakeupID(makeupID)
}

func (s *assignmentService) GetAssignmentWithClassID() (*models.ClassAssignments, error) {
	rows, err := s.assignmentRepo.GetAssignmentWithClassID()
	if err != nil {
		return nil, err
	}
	return groupByClassAndUser(rows), nil
}

func groupByClassAndUser(rows []*models.AssignmentRow) *models.ClassAssignments {
	result := &models.ClassAssignments{Classes: []models.ClassBlock{}}
	var (
		curUserID  string
		curClassID int64
		userIdx    int = -1
		classIdx   int = -1
	)
	for _, row := range rows {
		if row.ClassID != curClassID {
			classIdx++
			curClassID = row.ClassID
			result.Classes = append(result.Classes, models.ClassBlock{
				ClassID:   row.ClassID,
				ClassName: row.ClassName,
				Students:  []models.StudentAssignment{},
			})

			userIdx = -1
			curUserID = ""
		}
		if row.UserID != curUserID {
			userIdx++
			curUserID = row.UserID
			result.Classes[classIdx].Students = append(result.Classes[classIdx].Students, models.StudentAssignment{
				UserID:      row.UserID,
				Assignments: []models.Assignment{},
			})
			result.Classes[classIdx].Students[userIdx].Assignments = append(
				result.Classes[classIdx].Students[userIdx].Assignments,
				toAssignment(row))
		} else {
			result.Classes[classIdx].Students[userIdx].Assignments = append(
				result.Classes[classIdx].Students[userIdx].Assignments,
				toAssignment(row))
		}
	}
	return result
}

func toAssignment(row *models.AssignmentRow) models.Assignment {
	return models.Assignment{
		AssignmentID: &row.AssignmentID,
		UserID:       row.UserID,
		Content:      row.Content,
		Status:       row.Status,
		CreatedAt:    row.CreatedAt,
	}
}

func (s *assignmentService) UpsertAssignment(assignment *models.Assignment) error {
	return s.assignmentRepo.UpsertAssignment(assignment)
}

func (s *assignmentService) DeleteAssignment(assignmentID int64) error {
	return s.assignmentRepo.DeleteAssignment(assignmentID)
}
