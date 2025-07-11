package dto

import "time"

// Habit DTOs
type CreateHabitRequest struct {
	Title        string `json:"title" validate:"required,min=3,max=200"`
	StartDate    string `json:"start_date" validate:"required"` // Format: "2006-01-02"
	EndDate      string `json:"end_date" validate:"required"`   // Format: "2006-01-02"
	ReminderTime string `json:"reminder_time,omitempty"`        // Format: "15:04"
}

type UpdateHabitRequest struct {
	Title        string `json:"title,omitempty"`
	ReminderTime string `json:"reminder_time,omitempty"`
}

type HabitResponse struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Title        string    `json:"title"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	ReminderTime string    `json:"reminder_time,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Habit Log DTOs
type CreateHabitLogRequest struct {
	HabitID int64  `json:"habit_id" validate:"required"`
	Date    string `json:"date" validate:"required"` // Format: "2006-01-02"
	Status  string `json:"status" validate:"required,oneof=DONE SKIPPED FAILED"`
	Note    string `json:"note,omitempty"`
}

type UpdateHabitLogRequest struct {
	Status string `json:"status,omitempty"`
	Note   string `json:"note,omitempty"`
}

type HabitLogResponse struct {
	ID        int64     `json:"id"`
	HabitID   int64     `json:"habit_id"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"` // DONE, SKIPPED, FAILED
	PhotoURL  string    `json:"photo_url,omitempty"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HabitWithLogsResponse struct {
	Habit HabitResponse      `json:"habit"`
	Logs  []HabitLogResponse `json:"logs"`
	Stats HabitStatsResponse `json:"stats"`
}

type HabitStatsResponse struct {
	TotalDays     int     `json:"total_days"`
	CompletedDays int     `json:"completed_days"`
	SkippedDays   int     `json:"skipped_days"`
	FailedDays    int     `json:"failed_days"`
	SuccessRate   float64 `json:"success_rate"`
	CurrentStreak int     `json:"current_streak"`
	LongestStreak int     `json:"longest_streak"`
}
