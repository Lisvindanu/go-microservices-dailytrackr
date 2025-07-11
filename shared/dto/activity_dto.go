package dto

import "time"

// Activity DTOs
type CreateActivityRequest struct {
	Title        string `json:"title" validate:"required,min=3,max=200"`
	StartTime    string `json:"start_time" validate:"required"` // Format: "2006-01-02T15:04:05Z"
	DurationMins int    `json:"duration_mins" validate:"required,min=1"`
	Cost         *int   `json:"cost,omitempty"` // Nullable
	Note         string `json:"note,omitempty"`
}

type UpdateActivityRequest struct {
	Title        string `json:"title,omitempty"`
	StartTime    string `json:"start_time,omitempty"`
	DurationMins int    `json:"duration_mins,omitempty"`
	Cost         *int   `json:"cost,omitempty"`
	Note         string `json:"note,omitempty"`
}

type ActivityResponse struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Title        string    `json:"title"`
	StartTime    time.Time `json:"start_time"`
	DurationMins int       `json:"duration_mins"`
	Cost         *int      `json:"cost,omitempty"`
	PhotoURL     string    `json:"photo_url,omitempty"`
	Note         string    `json:"note,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ActivityListResponse struct {
	Activities []ActivityResponse `json:"activities"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
}

// Photo Upload
type PhotoUploadResponse struct {
	URL       string `json:"url"`
	PublicID  string `json:"public_id"`
	SecureURL string `json:"secure_url"`
}
