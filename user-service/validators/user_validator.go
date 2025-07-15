package validators

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	"dailytrackr/shared/dto"
)

// UserValidator provides advanced validation for user-related operations
type UserValidator struct{}

// NewUserValidator creates a new user validator instance
func NewUserValidator() *UserValidator {
	return &UserValidator{}
}

// ValidateRegistration validates user registration data
func (v *UserValidator) ValidateRegistration(req dto.RegisterRequest) error {
	// Validate username
	if err := v.ValidateUsername(req.Username); err != nil {
		return err
	}

	// Validate email
	if err := v.ValidateEmail(req.Email); err != nil {
		return err
	}

	// Validate password
	if err := v.ValidatePassword(req.Password); err != nil {
		return err
	}

	return nil
}

// ValidateProfileUpdate validates profile update data
func (v *UserValidator) ValidateProfileUpdate(req dto.UpdateProfileRequest) error {
	// Validate username if provided
	if req.Username != "" {
		if err := v.ValidateUsername(req.Username); err != nil {
			return err
		}
	}

	// Validate email if provided
	if req.Email != "" {
		if err := v.ValidateEmail(req.Email); err != nil {
			return err
		}
	}

	// Validate bio if provided
	if req.Bio != nil {
		if err := v.ValidateBio(*req.Bio); err != nil {
			return err
		}
	}

	return nil
}

// ValidatePasswordChange validates password change request
func (v *UserValidator) ValidatePasswordChange(req dto.ChangePasswordRequest) error {
	// Validate current password is not empty
	if strings.TrimSpace(req.CurrentPassword) == "" {
		return errors.New("current password is required")
	}

	// Validate new password
	if err := v.ValidatePassword(req.NewPassword); err != nil {
		return err
	}

	// Check if new password is different from current
	if req.CurrentPassword == req.NewPassword {
		return errors.New("new password must be different from current password")
	}

	return nil
}

// ValidateUsername validates username format and rules
func (v *UserValidator) ValidateUsername(username string) error {
	username = strings.TrimSpace(username)

	// Check length
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

	// Check if starts with letter or number (not underscore or dash)
	if !regexp.MustCompile(`^[a-zA-Z0-9]`).MatchString(username) {
		return errors.New("username must start with a letter or number")
	}

	// Check for consecutive special characters
	if regexp.MustCompile(`[_-]{2,}`).MatchString(username) {
		return errors.New("username cannot contain consecutive special characters")
	}

	// Reserved usernames
	reservedUsernames := []string{
		"admin", "administrator", "root", "api", "www", "mail", "ftp",
		"test", "demo", "guest", "user", "null", "undefined", "system",
		"support", "help", "info", "contact", "about", "dailytrackr",
	}

	usernameLower := strings.ToLower(username)
	for _, reserved := range reservedUsernames {
		if usernameLower == reserved {
			return errors.New("username is reserved and cannot be used")
		}
	}

	return nil
}

// ValidateEmail validates email format
func (v *UserValidator) ValidateEmail(email string) error {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return errors.New("email is required")
	}

	// Basic email regex validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	// Check length
	if len(email) > 255 {
		return errors.New("email must not exceed 255 characters")
	}

	// Check for multiple @ symbols
	if strings.Count(email, "@") != 1 {
		return errors.New("email must contain exactly one @ symbol")
	}

	// Split and validate parts
	parts := strings.Split(email, "@")
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return errors.New("email local and domain parts cannot be empty")
	}

	// Check local part (before @)
	localPart := parts[0]
	if len(localPart) > 64 {
		return errors.New("email local part must not exceed 64 characters")
	}

	// Check domain part (after @)
	domainPart := parts[1]
	if len(domainPart) > 253 {
		return errors.New("email domain part must not exceed 253 characters")
	}

	// Domain must contain at least one dot
	if !strings.Contains(domainPart, ".") {
		return errors.New("email domain must contain at least one dot")
	}

	return nil
}

// ValidatePassword validates password strength
func (v *UserValidator) ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	if len(password) > 128 {
		return errors.New("password must not exceed 128 characters")
	}

	// Check for at least one letter
	hasLetter := false
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
			break
		}
	}

	if !hasLetter {
		return errors.New("password must contain at least one letter")
	}

	// Optional: Check for common weak passwords
	weakPasswords := []string{
		"password", "123456", "password123", "admin", "qwerty",
		"letmein", "welcome", "monkey", "dragon", "master",
	}

	passwordLower := strings.ToLower(password)
	for _, weak := range weakPasswords {
		if passwordLower == weak {
			return errors.New("password is too common, please choose a stronger password")
		}
	}

	return nil
}

// ValidateBio validates user bio
func (v *UserValidator) ValidateBio(bio string) error {
	bio = strings.TrimSpace(bio)

	// Check length
	if len(bio) > 500 {
		return errors.New("bio must not exceed 500 characters")
	}

	// Check for excessive whitespace
	if regexp.MustCompile(`\s{3,}`).MatchString(bio) {
		return errors.New("bio contains excessive whitespace")
	}

	// Check for potentially harmful content (basic check)
	suspiciousPatterns := []string{
		"<script", "javascript:", "onclick", "onerror", "onload",
		"eval(", "document.cookie", "alert(",
	}

	bioLower := strings.ToLower(bio)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(bioLower, pattern) {
			return errors.New("bio contains potentially harmful content")
		}
	}

	return nil
}

// ValidateFileUpload validates file upload for profile photos
func (v *UserValidator) ValidateFileUpload(filename string, fileSize int64) error {
	// Check file size (5MB max for profile photos)
	maxSize := int64(5 << 20) // 5 MB
	if fileSize > maxSize {
		return errors.New("file size too large. Maximum size: 5MB")
	}

	if fileSize == 0 {
		return errors.New("file is empty")
	}

	// Check file extension
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}
	filename = strings.ToLower(filename)

	hasValidExtension := false
	for _, ext := range allowedExtensions {
		if strings.HasSuffix(filename, ext) {
			hasValidExtension = true
			break
		}
	}

	if !hasValidExtension {
		return errors.New("invalid file type. Allowed types: JPG, JPEG, PNG, WebP")
	}

	return nil
}

// SanitizeInput removes potentially harmful characters from input
func (v *UserValidator) SanitizeInput(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Trim whitespace
	input = strings.TrimSpace(input)

	// Remove control characters except tab, newline, carriage return
	var result strings.Builder
	for _, char := range input {
		if unicode.IsControl(char) && char != '\t' && char != '\n' && char != '\r' {
			continue
		}
		result.WriteRune(char)
	}

	return result.String()
}
