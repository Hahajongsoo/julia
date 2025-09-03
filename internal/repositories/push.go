package repositories

import (
	"database/sql"
	"fmt"
	"julia/internal/models"
	"log"
)

type PushRepository interface {
	CreateSubscription(subscription *models.Subscription) error
	GetSubscription(userID string) ([]*models.Subscription, error)
	DeleteSubscription(userID string, endpoint string) error
}

type pushRepository struct {
	db *sql.DB
}

func NewPushRepository(db *sql.DB) PushRepository {
	return &pushRepository{db: db}
}

func (r *pushRepository) CreateSubscription(subscription *models.Subscription) error {
	if subscription == nil {
		return fmt.Errorf("subscription cannot be nil")
	}

	query := `
		INSERT INTO push_subscriptions (user_id, endpoint, p256dh, auth)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (endpoint)
		DO UPDATE SET 
			user_id = EXCLUDED.user_id, 
			p256dh = EXCLUDED.p256dh, 
			auth = EXCLUDED.auth, 
			updated_at = now()
	`

	_, err := r.db.Exec(query,
		subscription.UserID,
		subscription.Endpoint,
		subscription.P256dh,
		subscription.Auth,
	)
	if err != nil {
		log.Printf("Database error creating subscription: %v", err)
		return fmt.Errorf("failed to create subscription in database: %w", err)
	}

	return nil
}

func (r *pushRepository) GetSubscription(userID string) ([]*models.Subscription, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	query := `
		SELECT user_id, endpoint, p256dh, auth, created_at, updated_at
		FROM push_subscriptions 
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		log.Printf("Database error getting subscriptions for user %s: %v", userID, err)
		return nil, fmt.Errorf("failed to get subscriptions from database: %w", err)
	}
	defer rows.Close()

	subscriptions := make([]*models.Subscription, 0)
	for rows.Next() {
		var subscription models.Subscription
		err := rows.Scan(
			&subscription.UserID,
			&subscription.Endpoint,
			&subscription.P256dh,
			&subscription.Auth,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning subscription row: %v", err)
			return nil, fmt.Errorf("failed to scan subscription row: %w", err)
		}
		subscriptions = append(subscriptions, &subscription)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating subscription rows: %v", err)
		return nil, fmt.Errorf("error iterating subscription rows: %w", err)
	}

	return subscriptions, nil
}



func (r *pushRepository) DeleteSubscription(userID string, endpoint string) error {
	if userID == "" {
		return fmt.Errorf("user_id is required")
	}

	if endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}

	query := `
		DELETE FROM push_subscriptions 
		WHERE user_id = $1 AND endpoint = $2
	`

	result, err := r.db.Exec(query, userID, endpoint)
	if err != nil {
		log.Printf("Database error deleting subscription for user %s with endpoint %s: %v", userID, endpoint, err)
		return fmt.Errorf("failed to delete subscription from database: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for subscription deletion: %v", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found for user %s with endpoint %s", userID, endpoint)
	}

	log.Printf("Deleted subscription for user %s with endpoint %s", userID, endpoint)
	return nil
}
