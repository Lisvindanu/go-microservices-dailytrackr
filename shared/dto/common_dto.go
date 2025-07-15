package dto

// File: shared/dto/common_dto.go
// Common DTOs yang digunakan oleh multiple services

// Photo Upload Response - digunakan oleh user service dan activity service
type PhotoUploadResponse struct {
	URL       string `json:"url"`
	PublicID  string `json:"public_id,omitempty"`
	SecureURL string `json:"secure_url"`
}

// File Upload Request (untuk multipart form)
type FileUploadRequest struct {
	FileField string `json:"file_field"` // nama field dalam form
	MaxSize   int64  `json:"max_size"`   // maksimum size dalam bytes
}

// Pagination Request
type PaginationRequest struct {
	Page  int `json:"page" validate:"min=1"`
	Limit int `json:"limit" validate:"min=1,max=100"`
}

// Pagination Response
type PaginationResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Date Range Request
type DateRangeRequest struct {
	StartDate string `json:"start_date" validate:"required"` // Format: "2006-01-02"
	EndDate   string `json:"end_date" validate:"required"`   // Format: "2006-01-02"
}

// Health Check Response
type HealthResponse struct {
	Service   string            `json:"service"`
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Timestamp string            `json:"timestamp"`
	Features  []string          `json:"features,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
}
