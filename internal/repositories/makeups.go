package repositories

import (
	"errors"
	"julia/internal/models"

	"database/sql"
)

type MakeupRepository interface {
	GetAllMakeups() ([]*models.Makeup, error)
	GetMakeupsByDate(date string) ([]*models.Makeup, error)
	GetMakeupsByMonth(yearMonth string) ([]*models.Makeup, error)
	GetMakeupsByUser(userID string) ([]*models.Makeup, error)
	GetMakeupsByUserAndDate(userID, date string) ([]*models.Makeup, error)
	CreateMakeup(makeup *models.Makeup) error
	UpdateMakeup(userID, date, time string, makeup *models.Makeup) error
	DeleteMakeup(userID, date, time string) error
}

type makeupRepository struct {
	db *sql.DB
}

func NewMakeupRepository(db *sql.DB) MakeupRepository {
	return &makeupRepository{db: db}
}

func (r *makeupRepository) GetAllMakeups() ([]*models.Makeup, error) {
	query := `
		SELECT user_id, makeup_date, start_time, reason
		FROM makeups
		ORDER BY makeup_date, start_time
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	makeups := make([]*models.Makeup, 0)
	for rows.Next() {
		var makeup models.Makeup
		err := rows.Scan(&makeup.UserID, &makeup.Date, &makeup.Time, &makeup.Reason)
		if err != nil {
			return nil, err
		}
		makeups = append(makeups, &makeup)
	}
	return makeups, nil
}

func (r *makeupRepository) GetMakeupsByDate(date string) ([]*models.Makeup, error) {
	query := `
		SELECT user_id, makeup_date, start_time, reason
		FROM makeups
		WHERE makeup_date = $1
	`
	rows, err := r.db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	makeups := make([]*models.Makeup, 0)
	for rows.Next() {
		var makeup models.Makeup
		err := rows.Scan(&makeup.UserID, &makeup.Date, &makeup.Time, &makeup.Reason)
		if err != nil {
			return nil, err
		}
		makeups = append(makeups, &makeup)
	}
	return makeups, nil
}

func (r *makeupRepository) GetMakeupsByMonth(yearMonth string) ([]*models.Makeup, error) {
	query := `
		SELECT user_id, makeup_date, start_time, reason
		FROM makeups
		WHERE makeup_date >= $1::date AND makeup_date < ($1::date + INTERVAL '1 month')
		ORDER BY makeup_date, start_time
	`
	// YYYY-MM 형식을 YYYY-MM-01로 변환
	startDate := yearMonth + "-01"
	rows, err := r.db.Query(query, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	makeups := make([]*models.Makeup, 0)
	for rows.Next() {
		var makeup models.Makeup
		err := rows.Scan(&makeup.UserID, &makeup.Date, &makeup.Time, &makeup.Reason)
		if err != nil {
			return nil, err
		}
		makeups = append(makeups, &makeup)
	}
	return makeups, nil
}

func (r *makeupRepository) GetMakeupsByUser(userID string) ([]*models.Makeup, error) {
	query := `
		SELECT user_id, makeup_date, start_time, reason
		FROM makeups
		WHERE user_id = $1
		ORDER BY makeup_date, start_time
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	makeups := make([]*models.Makeup, 0)
	for rows.Next() {
		var makeup models.Makeup
		err := rows.Scan(&makeup.UserID, &makeup.Date, &makeup.Time, &makeup.Reason)
		if err != nil {
			return nil, err
		}
		makeups = append(makeups, &makeup)
	}
	return makeups, nil
}

func (r *makeupRepository) GetMakeupsByUserAndDate(userID, date string) ([]*models.Makeup, error) {
	query := `
		SELECT user_id, makeup_date, start_time, reason
		FROM makeups
		WHERE user_id = $1 AND makeup_date = $2
		ORDER BY start_time
	`
	rows, err := r.db.Query(query, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	makeups := make([]*models.Makeup, 0)
	for rows.Next() {
		var makeup models.Makeup
		err := rows.Scan(&makeup.UserID, &makeup.Date, &makeup.Time, &makeup.Reason)
		if err != nil {
			return nil, err
		}
		makeups = append(makeups, &makeup)
	}
	return makeups, nil
}



func (r *makeupRepository) CreateMakeup(makeup *models.Makeup) error {
	query := `
		INSERT INTO makeups (user_id, makeup_date, start_time, reason)
		VALUES ($1, $2, $3, $4)
		RETURNING makeup_id
	`
	err := r.db.QueryRow(query, makeup.UserID, makeup.Date, makeup.Time, makeup.Reason).Scan(&makeup.MakeupID)
	if err != nil {
		return err
	}
	return nil
}

func (r *makeupRepository) UpdateMakeup(userID, date, time string, makeup *models.Makeup) error {
	query := `
		UPDATE makeups
		SET user_id = $1, makeup_date = $2, start_time = $3, reason = $4
		WHERE user_id = $5 AND makeup_date = $6 AND start_time = $7
	`
	result, err := r.db.Exec(query, makeup.UserID, makeup.Date, makeup.Time, makeup.Reason, userID, date, time)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("makeup not found")
	}
	return nil
}

func (r *makeupRepository) DeleteMakeup(userID, date, time string) error {
	query := `
		DELETE FROM makeups
		WHERE user_id = $1 AND makeup_date = $2 AND start_time = $3
	`
	result, err := r.db.Exec(query, userID, date, time)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("makeup not found")
	}
	return nil
}
