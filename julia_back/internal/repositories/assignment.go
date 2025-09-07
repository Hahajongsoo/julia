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
		SELECT COALESCE(c.class_id, -1) as class_id, COALESCE(c.class_name, '미지정') as class_name, a.assignment_id, a.user_id, a.content, a.status, a.created_at
		FROM assignments a
		JOIN users u ON a.user_id = u.id
		LEFT JOIN classes c ON u.class_id = c.class_id
		WHERE a.status = 'pending'
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
		err := rows.Scan(&assignment.ClassID, &assignment.ClassName, &assignment.AssignmentID, &assignment.UserID, &assignment.Content, &assignment.Status, &assignment.CreatedAt)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, &assignment)
	}
	return assignments, nil
}

func (r *assignmentRepository) UpsertAssignment(assignment *models.Assignment) error {
	if assignment.AssignmentID == 0 {
		query := `
			INSERT INTO assignments (user_id, makeup_id, content, status)
			VALUES ($1, $2, $3, $4)
			RETURNING assignment_id, created_at, updated_at
		`
		err := r.db.QueryRow(query, assignment.UserID, assignment.MakeupID, assignment.Content, assignment.Status).
			Scan(&assignment.AssignmentID, &assignment.CreatedAt, &assignment.UpdatedAt)
		if err != nil {
			return err
		}
	} else {
		query := `
			UPDATE assignments
			SET user_id = $2, makeup_id = $3, content = $4, status = $5, updated_at = now()
			WHERE assignment_id = $1
			RETURNING updated_at
		`
		err := r.db.QueryRow(query, assignment.AssignmentID, assignment.UserID, assignment.MakeupID, assignment.Content, assignment.Status).
			Scan(&assignment.UpdatedAt)
		if err != nil {
			return err
		}
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
