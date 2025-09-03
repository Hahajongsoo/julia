package repositories

import (
	"database/sql"
	"time"
)

type NotificationRepository interface {
	UpsertNotification(makeupID int64, userID, kind string, at time.Time) error
	CancelNotificationByMakeup(userID, date, time string) error
}

type notifRepository struct{ DB *sql.DB }

func NewNotificationRepository(db *sql.DB) NotificationRepository {
	return &notifRepository{DB: db}
}

func (r *notifRepository) UpsertNotification(makeupID int64, userID, kind string, at time.Time) error {
	_, err := r.DB.Exec(`
	  INSERT INTO scheduled_notifications (makeup_id, user_id, kind, scheduled_at, status, attempts)
	  VALUES ($1,$2,$3,$4,'pending',0)
	  ON CONFLICT (makeup_id, kind) DO UPDATE
	    SET scheduled_at = EXCLUDED.scheduled_at,
	        status = 'pending',
	        attempts = 0,
	        last_error = NULL
	`, makeupID, userID, kind, at)
	return err
}

func (r *notifRepository) CancelNotificationByMakeup(userID, date, time string) error {
	_, err := r.DB.Exec(`
	  DELETE FROM scheduled_notifications
	  WHERE makeup_id = (
	    SELECT makeup_id
		FROM makeups
		WHERE user_id = $1 and makeup_date = $2 and start_time = $3
	  ) 
	  AND status = 'pending'
	`, userID, date, time)
	return err
}
