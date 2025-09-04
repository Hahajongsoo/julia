package repositories

import (
	"database/sql"
	"julia/internal/models"
)

type ExamPeriodRepository interface {
	GetAllExamPeriods() ([]*models.ExamPeriod, error)
	GetExamPeriodByClassID(classID int64) ([]*models.ExamPeriod, error)
	CreateExamPeriod(examPeriod *models.ExamPeriod) error
	UpdateExamPeriod(examPeriod *models.ExamPeriod) error
	DeleteExamPeriod(examPeriodID int64) error
}

type examPeriodRepository struct {
	db *sql.DB
}

func NewExamPeriodRepository(db *sql.DB) ExamPeriodRepository {
	return &examPeriodRepository{db: db}
}

func (r *examPeriodRepository) GetAllExamPeriods() ([]*models.ExamPeriod, error) {
	query := `
		SELECT exam_period_id, class_id, name, description, start_date, end_date
		FROM exam_periods
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	examPeriods := make([]*models.ExamPeriod, 0)
	for rows.Next() {
		var examPeriod models.ExamPeriod
		err := rows.Scan(&examPeriod.ExamPeriodID, &examPeriod.ClassID, &examPeriod.Name, &examPeriod.Description, &examPeriod.StartDate, &examPeriod.EndDate)
		if err != nil {
			return nil, err
		}
		examPeriods = append(examPeriods, &examPeriod)
	}
	return examPeriods, nil
}

func (r *examPeriodRepository) GetExamPeriodByClassID(classID int64) ([]*models.ExamPeriod, error) {
	query := `
		SELECT exam_period_id, class_id, name, description, start_date, end_date
		FROM exam_periods
		WHERE class_id = $1
	`
	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	examPeriods := make([]*models.ExamPeriod, 0)
	for rows.Next() {
		var examPeriod models.ExamPeriod
		err := rows.Scan(&examPeriod.ExamPeriodID, &examPeriod.ClassID, &examPeriod.Name, &examPeriod.Description, &examPeriod.StartDate, &examPeriod.EndDate)
		if err != nil {
			return nil, err
		}
		examPeriods = append(examPeriods, &examPeriod)
	}
	return examPeriods, nil
}

func (r *examPeriodRepository) CreateExamPeriod(examPeriod *models.ExamPeriod) error {
	query := `
		INSERT INTO exam_periods (class_id, name, description, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, examPeriod.ClassID, examPeriod.Name, examPeriod.Description, examPeriod.StartDate, examPeriod.EndDate)
	if err != nil {
		return err
	}
	return nil
}

func (r *examPeriodRepository) UpdateExamPeriod(examPeriod *models.ExamPeriod) error {
	query := `
		UPDATE exam_periods
		SET name = $1, description = $2, start_date = $3, end_date = $4
		WHERE exam_period_id = $5
	`
	_, err := r.db.Exec(query, examPeriod.Name, examPeriod.Description, examPeriod.StartDate, examPeriod.EndDate, examPeriod.ExamPeriodID)
	if err != nil {
		return err
	}
	return nil
}

func (r *examPeriodRepository) DeleteExamPeriod(examPeriodID int64) error {
	query := `
		DELETE FROM exam_periods
		WHERE exam_period_id = $1
	`
	_, err := r.db.Exec(query, examPeriodID)
	if err != nil {
		return err
	}
	return nil
}
