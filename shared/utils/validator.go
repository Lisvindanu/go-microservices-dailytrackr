package utils

import (
	"errors"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct using struct tags
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// GetValidationErrors returns formatted validation errors
func GetValidationErrors(err error) []string {
	var errors []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errors = append(errors, getFieldErrorMessage(fieldError))
		}
	}

	return errors
}

// getFieldErrorMessage returns a user-friendly error message for field validation
func getFieldErrorMessage(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())

	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + fe.Param() + " characters"
	case "max":
		return field + " must not exceed " + fe.Param() + " characters"
	case "oneof":
		return field + " must be one of: " + fe.Param()
	default:
		return field + " is invalid"
	}
}

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	// Check for at least one letter
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	if !hasLetter {
		return errors.New("password must contain at least one letter")
	}

	return nil
}

// ValidateUsername validates username format
func ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	if len(username) > 50 {
		return errors.New("username must not exceed 50 characters")
	}

	// Check for valid characters (letters, numbers, underscore, dash)
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(username)
	if !validUsername {
		return errors.New("username can only contain letters, numbers, underscore, and dash")
	}

	return nil
}
