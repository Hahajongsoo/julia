package services

import (
	"julia/internal/models"
	"julia/internal/repositories"
	"time"
)

type NotificationService interface {
	EnqueueMakeupNotifications(m *models.Makeup) error
	CancelMakeupNotifications(userID, date, time string) error
}

type notificationService struct {
	notifRepo repositories.NotificationRepository
}

func NewNotificationService(notifRepo repositories.NotificationRepository) NotificationService {
	return &notificationService{notifRepo: notifRepo}
}

func makeupStartAtKST(m *models.Makeup) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return time.Time{}, err
	}

	d := m.Date
	if d.Location() == nil {
		d = d.UTC()
	}
	t := m.Time
	if t.Location() == nil {
		t = t.UTC()
	}

	d = d.In(loc)
	t = t.In(loc)

	start := time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), 0, 0, loc)
	return start, nil
}

func (s *notificationService) EnqueueMakeupNotifications(m *models.Makeup) error {
	startAt, err := makeupStartAtKST(m)
	if err != nil {
		return err
	}

	oneDayBefore := startAt.Add(-24 * time.Hour)
	thirtyMinBefore := startAt.Add(-30 * time.Minute)

	now := time.Now()
	if oneDayBefore.Before(now) {
		oneDayBefore = now
	}
	if thirtyMinBefore.Before(now) {
		thirtyMinBefore = now
	}

	if err := s.notifRepo.UpsertNotification(m.MakeupID, m.UserID, "makeup-1d", oneDayBefore.UTC()); err != nil {
		return err
	}
	if err := s.notifRepo.UpsertNotification(m.MakeupID, m.UserID, "makeup-30m", thirtyMinBefore.UTC()); err != nil {
		return err
	}
	return nil
}

func (s *notificationService) CancelMakeupNotifications(userID, date, time string) error {
	return s.notifRepo.CancelNotificationByMakeup(userID, date, time)
}
