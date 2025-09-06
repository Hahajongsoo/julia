package services

import (
	"encoding/json"
	"fmt"
	"julia/internal/models"
	"julia/internal/repositories"
	"log"
	"os"

	"github.com/SherClockHolmes/webpush-go"
)

type PushService interface {
	CreateSubscription(subscription *models.Subscription) error
	GetSubscription(userID string) ([]*models.Subscription, error)
	SendNotification(userID string) error
	DeleteSubscription(userID string, endpoint string) error
}

type pushService struct {
	pushRepo repositories.PushRepository
}

func NewPushService(pushRepo repositories.PushRepository) PushService {
	return &pushService{pushRepo: pushRepo}
}

func (s *pushService) CreateSubscription(subscription *models.Subscription) error {
	if subscription == nil {
		return fmt.Errorf("subscription cannot be nil")
	}

	if subscription.UserID == "" || subscription.Endpoint == "" {
		return fmt.Errorf("user_id and endpoint are required")
	}

	err := s.pushRepo.CreateSubscription(subscription)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

func (s *pushService) GetSubscription(userID string) ([]*models.Subscription, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	subscriptions, err := s.pushRepo.GetSubscription(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}

	return subscriptions, nil
}



func (s *pushService) DeleteSubscription(userID string, endpoint string) error {
	if userID == "" {
		return fmt.Errorf("user_id is required")
	}

	if endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}

	err := s.pushRepo.DeleteSubscription(userID, endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

func (s *pushService) SendNotification(userID string) error {
	if userID == "" {
		return fmt.Errorf("user_id is required")
	}

	subscriptions, err := s.pushRepo.GetSubscription(userID)
	if err != nil {
		return fmt.Errorf("failed to get subscriptions: %w", err)
	}

	if len(subscriptions) == 0 {
		return fmt.Errorf("no subscriptions found for user")
	}

	payload := map[string]any{
		"title": "보강 일정 알림",
		"body":  "오늘 18:00, 김민수 학생 보강이 있어요.",
		"url":   "/calendar",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %w", err)
	}

	successCount := 0
	for _, subscription := range subscriptions {
		resp, err := webpush.SendNotification(body, &webpush.Subscription{
			Endpoint: subscription.Endpoint,
			Keys: webpush.Keys{
				P256dh: subscription.P256dh,
				Auth:   subscription.Auth,
			},
		}, &webpush.Options{
			Subscriber:      "mailto:admin@example.com",
			VAPIDPublicKey:  os.Getenv("VAPID_PUBLIC"),
			VAPIDPrivateKey: os.Getenv("VAPID_PRIVATE"),
			TTL:             60,
		})
		if err != nil {
			log.Printf("Failed to send notification to endpoint %s: %v", subscription.Endpoint, err)
			continue
		}

		log.Printf("Successfully sent notification to endpoint %s, status: %d", subscription.Endpoint, resp.StatusCode)
		successCount++
	}

	if successCount == 0 {
		return fmt.Errorf("failed to send notification to any subscription")
	}

	log.Printf("Successfully sent notifications to %d/%d subscriptions for user: %s", successCount, len(subscriptions), userID)
	return nil
}
