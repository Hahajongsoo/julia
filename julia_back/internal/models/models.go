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

func parseDateString(dateStr string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		loc = time.UTC
	}

	// 날짜 문자열을 명시적으로 로컬 시간대로 파싱
	date, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return time.Time{}, err
	}
	// 시간을 00:00:00으로 설정하여 하루의 시작으로 처리
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc), nil
}

func parseTimeString(timeStr string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		loc = time.UTC
	}

	startTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(0, 1, 1, startTime.Hour(), startTime.Minute(), 0, 0, loc), nil
}

type Makeup struct {
	MakeupID int64     `json:"makeup_id"`
	UserID   string    `json:"user_id" binding:"required"`
	Date     time.Time `json:"makeup_date" binding:"required"`
	Time     time.Time `json:"start_time" binding:"required"`
	Reason   string    `json:"reason" binding:"required" `
	Status   string    `json:"status"`
}

type MakeupDTO struct {
	MakeupID int64  `json:"makeup_id"`
	UserID   string `json:"user_id"`
	Date     string `json:"makeup_date"`
	Time     string `json:"start_time"`
	Reason   string `json:"reason"`
	Status   string `json:"status"`
}

func (m *Makeup) ToMakeupDTO() *MakeupDTO {
	return &MakeupDTO{
		MakeupID: m.MakeupID,
		UserID:   m.UserID,
		Date:     m.Date.Format("2006-01-02"),
		Time:     m.Time.Format("15:04"),
		Reason:   m.Reason,
		Status:   m.Status,
	}
}

func (m *MakeupDTO) ToMakeup() *Makeup {
	date, err := parseDateString(m.Date)
	if err != nil {
		return nil
	}

	startTime, err := parseTimeString(m.Time)
	if err != nil {
		return nil
	}

	// Set default status if empty
	status := m.Status
	if status == "" {
		status = "pending"
	}

	return &Makeup{
		MakeupID: m.MakeupID,
		UserID:   m.UserID,
		Date:     date,
		Time:     startTime,
		Reason:   m.Reason,
		Status:   status,
	}
}

type Subscription struct {
	UserID    string    `json:"user_id"`
	Endpoint  string    `json:"endpoint"`
	P256dh    string    `json:"p256dh"`
	Auth      string    `json:"auth"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SubscriptionKeys struct {
	P256dh string `json:"p256dh"`
	Auth   string `json:"auth"`
}

type SubscriptionRequest struct {
	UserID   string           `json:"user_id" binding:"required"`
	Endpoint string           `json:"endpoint" binding:"required"`
	Keys     SubscriptionKeys `json:"keys"`
}

type DeleteSubscriptionRequest struct {
	Endpoint string `json:"endpoint" binding:"required"`
}

type Class struct {
	ClassID     int64  `json:"class_id"`
	ClassName   string `json:"class_name" binding:"required"`
	Description string `json:"description"`
}

type ExamPeriod struct {
	ExamPeriodID int64     `json:"exam_period_id"`
	ClassID      int64     `json:"class_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	EnglishDate  time.Time `json:"english_date"`
}

type ExamPeriodDTO struct {
	ExamPeriodID int64     `json:"exam_period_id"`
	ClassID      int64     `json:"class_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	StartDate    string    `json:"start_date"`
	EndDate      string    `json:"end_date"`
	EnglishDate  string    `json:"english_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (e *ExamPeriod) ToExamPeriodDTO() *ExamPeriodDTO {
	var englishDateStr string
	if !e.EnglishDate.IsZero() {
		englishDateStr = e.EnglishDate.Format("2006-01-02")
	}

	return &ExamPeriodDTO{
		ExamPeriodID: e.ExamPeriodID,
		ClassID:      e.ClassID,
		Name:         e.Name,
		Description:  e.Description,
		StartDate:    e.StartDate.Format("2006-01-02"),
		EndDate:      e.EndDate.Format("2006-01-02"),
		EnglishDate:  englishDateStr,
	}
}

func (e *ExamPeriodDTO) ToExamPeriod() *ExamPeriod {
	startDate, err := parseDateString(e.StartDate)
	if err != nil {
		return nil
	}

	endDate, err := parseDateString(e.EndDate)
	if err != nil {
		return nil
	}

	var englishDate time.Time
	if e.EnglishDate != "" {
		var err error
		englishDate, err = parseDateString(e.EnglishDate)
		if err != nil {
			return nil
		}
	}

	return &ExamPeriod{
		ExamPeriodID: e.ExamPeriodID,
		ClassID:      e.ClassID,
		Name:         e.Name,
		Description:  e.Description,
		StartDate:    startDate,
		EndDate:      endDate,
		EnglishDate:  englishDate,
	}
}

type AdminMemo struct {
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Assignment struct {
	AssignmentID int64     `json:"assignment_id"`
	UserID       string    `json:"user_id"`
	MakeupID     int64     `json:"makeup_id"`
	Content      string    `json:"content"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AssignmentRow struct {
	ClassID      int64     `json:"class_id"`
	ClassName    string    `json:"class_name"`
	AssignmentID int64     `json:"assignment_id"`
	UserID       string    `json:"user_id"`
	Content      string    `json:"content"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ClassAssignments struct {
	Classes []ClassBlock `json:"classes"`
}

type ClassBlock struct {
	ClassID   int64               `json:"class_id"`
	ClassName string              `json:"class_name"`
	Students  []StudentAssignment `json:"students"`
}

type StudentAssignment struct {
	UserID      string       `json:"user_id"`
	Assignments []Assignment `json:"assignments"`
}
