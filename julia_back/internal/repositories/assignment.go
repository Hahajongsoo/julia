package repositories

import (
	"database/sql"
	"julia/internal/models"
)

type AssignmentRepository interface {
	GetAllAssignments() ([]*models.Assignment, error)
	GetAssignmentsByUserID(userID string) ([]*models.Assignment, error)
	GetAssignmentsByMakeupID(makeupID int64) ([]*models.Assignment, error)
	GetAssignmentWithClassID() ([]*models.AssignmentRow, error)
	UpsertAssignment(assignment *models.Assignment) error
	DeleteAssignment(assignmentID int64) error
}

type assignmentRepository struct {
	db *sql.DB
}

func NewAssignmentRepository(db *sql.DB) AssignmentRepository {
	return &assignmentRepository{db: db}
}

func (r *assignmentRepository) GetAllAssignments() ([]*models.Assignment, error) {
	query := `
		SELECT assignment_id, user_id, content, status
		FROM assignments
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	assignments := make([]*models.Assignment, 0)
	for rows.Next() {
		var assignment models.Assignment
		err := rows.Scan(&assignment.AssignmentID, &assignment.UserID, &assignment.Content, &assignment.Status)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, &assignment)
	}
	return assignments, nil
}

func (r *assignmentRepository) GetAssignmentsByUserID(userID string) ([]*models.Assignment, error) {
	query := `
		SELECT assignment_id, user_id, content, status, created_at
		FROM assignments
		WHERE user_id = $1
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	assignments := make([]*models.Assignment, 0)
	for rows.Next() {
		var assignment models.Assignment
		err := rows.Scan(&assignment.AssignmentID, &assignment.UserID, &assignment.Content, &assignment.Status, &assignment.CreatedAt)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, &assignment)
	}
	return assignments, nil
}

func (r *assignmentRepository) GetAssignmentsByMakeupID(makeupID int64) ([]*models.Assignment, error) {
	query := `
		SELECT assignment_id, user_id, content, status, created_at
		FROM assignments
		WHERE makeup_id = $1
	`
	rows, err := r.db.Query(query, makeupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	assignments := make([]*models.Assignment, 0)
	for rows.Next() {
		var assignment models.Assignment
		err := rows.Scan(&assignment.AssignmentID, &assignment.UserID, &assignment.Content, &assignment.Status, &assignment.CreatedAt)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, &assignment)
	}
	return assignments, nil
}

func (r *assignmentRepository) GetAssignmentWithClassID() ([]*models.AssignmentRow, error) {
	query := `
		SELECT c.class_id, c.class_name, a.assignment_id, a.user_id, a.content, a.status, a.created_at
		FROM assignments a
		JOIN users u ON a.user_id = u.id
		JOIN classes c ON u.class_id = c.class_id
		ORDER BY c.class_id, a.user_id, a.created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	assignments := make([]*models.AssignmentRow, 0)
	for rows.Next() {
		var assignment models.AssignmentRow
		err := rows.Scan(&assignment.ClassID, &assignment.AssignmentID, &assignment.UserID, &assignment.Content, &assignment.Status, &assignment.CreatedAt)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, &assignment)
	}
	return assignments, nil
}

func (r *assignmentRepository) UpsertAssignment(assignment *models.Assignment) error {
	query := `
		INSERT INTO assignments (user_id, content, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (assignment_id) DO UPDATE SET
			content = EXCLUDED.content,
			status = EXCLUDED.status
	`
	_, err := r.db.Exec(query, assignment.UserID, assignment.Content, assignment.Status)
	if err != nil {
		return err
	}
	return nil
}

func (r *assignmentRepository) DeleteAssignment(assignmentID int64) error {
	query := `
		DELETE FROM assignments
		WHERE assignment_id = $1
	`
	_, err := r.db.Exec(query, assignmentID)
	if err != nil {
		return err
	}
	return nil
}
