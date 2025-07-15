package dto

import "time"

// User DTOs
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Profile Management DTOs
type UpdateProfileRequest struct {
	Username string  `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    string  `json:"email,omitempty" validate:"omitempty,email"`
	Bio      *string `json:"bio,omitempty" validate:"omitempty,max=500"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

type DeleteAccountRequest struct {
	Password string `json:"password" validate:"required"`
}

// Response DTOs
type UserResponse struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Bio          string    `json:"bio,omitempty"`
	ProfilePhoto string    `json:"profile_photo,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// Generic Response DTOs (shared across all services)
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
