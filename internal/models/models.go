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
	}
}
