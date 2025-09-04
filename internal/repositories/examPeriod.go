package repositories

import (
	"database/sql"
	"julia/internal/models"
)

type ExamPeriodRepository interface {
	GetExamPeriodByClassID(classID int64) ([]*models.ExamPeriod, error)
	UpsertExamPeriod(examPeriod *models.ExamPeriod) error
	DeleteExamPeriod(examPeriodID int64) error
}

type examPeriodRepository struct {
	db *sql.DB
}

func NewExamPeriodRepository(db *sql.DB) ExamPeriodRepository {
	return &examPeriodRepository{db: db}
}

func (r *examPeriodRepository) GetExamPeriodByClassID(classID int64) ([]*models.ExamPeriod, error) {
	query := `
		SELECT exam_period_id, class_id, name, start_date, end_date
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
		err := rows.Scan(&examPeriod.ExamPeriodID, &examPeriod.ClassID, &examPeriod.Name, &examPeriod.StartDate, &examPeriod.EndDate)
		if err != nil {
			return nil, err
		}
		examPeriods = append(examPeriods, &examPeriod)
	}
	return examPeriods, nil
}

func (r *examPeriodRepository) UpsertExamPeriod(examPeriod *models.ExamPeriod) error {
	query := `
		INSERT INTO exam_periods (class_id, name, start_date, end_date)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (class_id) DO UPDATE SET
			name = EXCLUDED.name,
			start_date = EXCLUDED.start_date,
			end_date = EXCLUDED.end_date,
			updated_at = now()
	`
	_, err := r.db.Exec(query, examPeriod.ClassID, examPeriod.Name, examPeriod.StartDate, examPeriod.EndDate)
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
