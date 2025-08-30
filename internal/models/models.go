package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Password  string    `json:"password"`
	Phone     string    `json:"phone"`
	ClassID   *int64    `json:"class_id"`
	CreatedAt time.Time `json:"created_at"`
	Role      string    `json:"role"`
}

type ResponseUser struct {
	ID        string    `json:"id"`
	Phone     string    `json:"phone"`
	ClassID   *int64    `json:"class_id"`
	CreatedAt time.Time `json:"created_at"`
	Role      string    `json:"role"`
}

type LoginRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

func (u *User) ToResponseUser() *ResponseUser {
	return &ResponseUser{
		ID:        u.ID,
		Phone:     u.Phone,
		ClassID:   u.ClassID,
		CreatedAt: u.CreatedAt,
		Role:      u.Role,
	}
}

type Makeup struct {
	MakeupID string    `json:"makeup_id"`
	UserID   string    `json:"user_id" binding:"required"`
	Date     time.Time `json:"makeup_date" binding:"required"`
	Time     time.Time `json:"start_time" binding:"required"`
	Reason   string    `json:"reason" binding:"required" `
}

type MakeupDTO struct {
	MakeupID string `json:"makeup_id"`
	UserID   string `json:"user_id"`
	Date     string `json:"makeup_date"`
	Time     string `json:"start_time"`
	Reason   string `json:"reason"`
}

func (m *Makeup) ToMakeupDTO() *MakeupDTO {
	return &MakeupDTO{
		MakeupID: m.MakeupID,
		UserID:   m.UserID,
		Date:     m.Date.Format("2006-01-02"),
		Time:     m.Time.Format("15:04"),
		Reason:   m.Reason,
	}
}

func (m *MakeupDTO) ToMakeup() *Makeup {
	date, err := time.Parse("2006-01-02", m.Date)
	if err != nil {
		return nil
	}
	startTime, err := time.Parse("15:04", m.Time)
	if err != nil {
		return nil
	}
	return &Makeup{
		MakeupID: m.MakeupID,
		UserID:   m.UserID,
		Date:     date,
		Time:     startTime,
		Reason:   m.Reason,
	}
}
