package repositories

import (
	"database/sql"
	"julia/internal/models"
)

type ClassRepository interface {
	GetAllClasses() ([]*models.Class, error)
	CreateClass(class *models.Class) error
	UpdateClass(class *models.Class) error
	DeleteClass(classID int64) error
}

type classRepository struct {
	db *sql.DB
}

func NewClassRepository(db *sql.DB) ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) GetAllClasses() ([]*models.Class, error) {
	query := `
		SELECT class_id, class_name
		FROM classes
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	classes := make([]*models.Class, 0)
	for rows.Next() {
		var class models.Class
		err := rows.Scan(&class.ClassID, &class.ClassName)
		if err != nil {
			return nil, err
		}
		classes = append(classes, &class)
	}
	return classes, nil
}

func (r *classRepository) CreateClass(class *models.Class) error {
	query := `
		INSERT INTO classes (class_id, class_name)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, class.ClassID, class.ClassName)
	if err != nil {
		return err
	}
	return nil
}

func (r *classRepository) UpdateClass(class *models.Class) error {
	query := `
		UPDATE classes
		SET class_name = $1
		WHERE class_id = $3
	`
	_, err := r.db.Exec(query, class.ClassName, class.ClassID)
	if err != nil {
		return err
	}
	return nil
}

func (r *classRepository) DeleteClass(classID int64) error {
	query := `
		DELETE FROM classes
		WHERE class_id = $1
	`
	_, err := r.db.Exec(query, classID)
	if err != nil {
		return err
	}
	return nil
}
