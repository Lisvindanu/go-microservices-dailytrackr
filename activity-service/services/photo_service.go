package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"dailytrackr/shared/config"
	"dailytrackr/shared/constants"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type PhotoService struct {
	cloudinary *cloudinary.Cloudinary
	config     *config.Config
}

// NewPhotoService creates a new photo service instance
func NewPhotoService(cfg *config.Config) *PhotoService {
	// Initialize Cloudinary only if credentials are provided
	var cld *cloudinary.Cloudinary
	if cfg.CloudinaryCloudName != "" && cfg.CloudinaryAPIKey != "" && cfg.CloudinaryAPISecret != "" {
		cld, _ = cloudinary.NewFromParams(
			cfg.CloudinaryCloudName,
			cfg.CloudinaryAPIKey,
			cfg.CloudinaryAPISecret,
		)
	}

	return &PhotoService{
		cloudinary: cld,
		config:     cfg,
	}
}

// UploadPhoto uploads a photo to Cloudinary
func (s *PhotoService) UploadPhoto(file *multipart.FileHeader) (string, error) {
	// Check if Cloudinary is configured
	if s.cloudinary == nil {
		return s.UploadPhotoSimple(file)
	}

	// Validate file type
	if !s.isValidImageType(file.Filename) {
		return "", fmt.Errorf("invalid file type. Allowed types: %s", constants.AllowedImageTypes)
	}

	// Check file size
	if file.Size > constants.MaxFileSize {
		return "", fmt.Errorf("file size too large. Maximum size: %d bytes", constants.MaxFileSize)
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Generate unique filename
	filename := s.generateFilename(file.Filename)

	// Upload to Cloudinary with corrected parameters
	uploadResult, err := s.cloudinary.Upload.Upload(
		context.Background(),
		src,
		uploader.UploadParams{
			PublicID:     filename,
			Folder:       constants.UploadPath,
			ResourceType: "image",
			Format:       "auto",
			// Remove problematic Transformation and Quality fields for now
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %v", err)
	}

	return uploadResult.SecureURL, nil
}

// UploadPhotoSimple uploads photo with basic validation (fallback if Cloudinary not configured)
func (s *PhotoService) UploadPhotoSimple(file *multipart.FileHeader) (string, error) {
	// Validate file type
	if !s.isValidImageType(file.Filename) {
		return "", fmt.Errorf("invalid file type. Allowed types: %s", constants.AllowedImageTypes)
	}

	// Check file size
	if file.Size > constants.MaxFileSize {
		return "", fmt.Errorf("file size too large. Maximum size: %d bytes", constants.MaxFileSize)
	}

	// Generate mock URL (for development/testing)
	filename := s.generateFilename(file.Filename)
	mockURL := fmt.Sprintf("https://placeholder.dailytrackr.com/photos/%s", filename)

	// In a real scenario, you would:
	// 1. Save file to local storage or
	// 2. Upload to another cloud storage service
	// For now, return mock URL
	return mockURL, nil
}

// isValidImageType checks if the file type is allowed
func (s *PhotoService) isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	if len(ext) > 0 {
		ext = ext[1:] // Remove the dot
	}

	allowedTypes := strings.Split(constants.AllowedImageTypes, ",")
	for _, allowedType := range allowedTypes {
		if ext == strings.TrimSpace(allowedType) {
			return true
		}
	}

	return false
}

// generateFilename generates a unique filename
func (s *PhotoService) generateFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("activity_%d_%s%s", timestamp, generateRandomString(8), ext)
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		// Use timestamp-based seeding for better randomness
		seed := time.Now().UnixNano() + int64(i)
		result[i] = charset[seed%int64(len(charset))]
	}
	return string(result)
}
