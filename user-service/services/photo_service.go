package services

import (
	"context"
	"fmt"
	"log"
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
	log.Printf("Initializing PhotoService for User Service...")
	log.Printf("Cloudinary Config - Name: '%s', Key: '%s', Secret Set: %v",
		cfg.CloudinaryCloudName,
		cfg.CloudinaryAPIKey,
		cfg.CloudinaryAPISecret != "")

	// Initialize Cloudinary if credentials are available
	if cfg.CloudinaryCloudName == "" || cfg.CloudinaryAPIKey == "" || cfg.CloudinaryAPISecret == "" {
		log.Printf("⚠️  Warning: Cloudinary credentials not configured. Photo upload will be disabled.")
		return &PhotoService{
			cloudinary: nil,
			config:     cfg,
		}
	}

	cld, err := cloudinary.NewFromParams(
		cfg.CloudinaryCloudName,
		cfg.CloudinaryAPIKey,
		cfg.CloudinaryAPISecret,
	)
	if err != nil {
		log.Printf("❌ Failed to initialize Cloudinary client: %v", err)
		return &PhotoService{
			cloudinary: nil,
			config:     cfg,
		}
	}

	log.Println("✅ Cloudinary client initialized successfully for User Service!")

	return &PhotoService{
		cloudinary: cld,
		config:     cfg,
	}
}

// UploadPhoto uploads a profile photo to Cloudinary
func (s *PhotoService) UploadPhoto(file *multipart.FileHeader) (string, error) {
	// Check if Cloudinary is available
	if s.cloudinary == nil {
		return "", fmt.Errorf("photo upload service is not available. Please configure Cloudinary credentials")
	}

	// Validate file type
	if !s.isValidImageType(file.Filename) {
		return "", fmt.Errorf("invalid file type. Allowed types: %s", constants.AllowedProfileTypes)
	}

	// Check file size (profile photos should be smaller)
	if file.Size > constants.MaxProfilePhotoSize {
		return "", fmt.Errorf("file size too large. Maximum size for profile photos: 5 MB")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Generate unique filename for profile photo
	filename := s.generateProfileFilename()

	// Upload to Cloudinary with profile-specific settings
	log.Printf("Uploading profile photo '%s' to Cloudinary...", filename)

	// Create bool pointers for Cloudinary parameters
	uniqueFilename := true
	useFilename := false

	uploadResult, err := s.cloudinary.Upload.Upload(
		context.Background(),
		src,
		uploader.UploadParams{
			PublicID:       filename,
			Folder:         constants.ProfilePhotoPath,
			ResourceType:   "image",
			Transformation: "c_fill,w_300,h_300,q_auto,f_auto", // Optimize for profile photos
			AllowedFormats: []string{"jpg", "png", "webp"},
			UniqueFilename: &uniqueFilename, // Use pointer to bool
			UseFilename:    &useFilename,    // Use pointer to bool
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %v", err)
	}

	log.Printf("✅ Profile photo uploaded successfully. Secure URL: %s", uploadResult.SecureURL)
	return uploadResult.SecureURL, nil
}

// isValidImageType checks if the file type is allowed for profile photos
func (s *PhotoService) isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	if len(ext) > 0 {
		ext = ext[1:] // Remove the dot
	}

	// Only allow common image types for profile photos
	allowedTypes := []string{"jpg", "jpeg", "png", "webp"}
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			return true
		}
	}
	return false
}

// generateProfileFilename generates a unique filename for profile photos
func (s *PhotoService) generateProfileFilename() string {
	timestamp := time.Now().Unix()
	randomStr := s.generateRandomString(8)
	return fmt.Sprintf("profile_%d_%s", timestamp, randomStr)
}

// generateRandomString generates a random string of specified length
func (s *PhotoService) generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		seed := time.Now().UnixNano() + int64(i)
		result[i] = charset[seed%int64(len(charset))]
	}
	return string(result)
}

// DeletePhoto removes a photo from Cloudinary (optional feature)
func (s *PhotoService) DeletePhoto(publicID string) error {
	if s.cloudinary == nil {
		return fmt.Errorf("photo service is not available")
	}

	if publicID == "" {
		return fmt.Errorf("public ID is required")
	}

	// Extract public ID from URL if full URL is provided
	if strings.Contains(publicID, "cloudinary.com") {
		parts := strings.Split(publicID, "/")
		for i, part := range parts {
			if part == "upload" && i+1 < len(parts) {
				// Skip version if present
				if strings.HasPrefix(parts[i+1], "v") && i+2 < len(parts) {
					publicID = strings.Join(parts[i+2:], "/")
				} else {
					publicID = strings.Join(parts[i+1:], "/")
				}
				break
			}
		}

		// Remove file extension
		if dotIndex := strings.LastIndex(publicID, "."); dotIndex != -1 {
			publicID = publicID[:dotIndex]
		}
	}

	log.Printf("Deleting photo with public ID: %s", publicID)

	_, err := s.cloudinary.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: publicID,
	})

	if err != nil {
		log.Printf("Failed to delete photo: %v", err)
		return fmt.Errorf("failed to delete photo: %v", err)
	}

	log.Printf("✅ Photo deleted successfully: %s", publicID)
	return nil
}

// ValidateFileUpload validates file before upload
func (s *PhotoService) ValidateFileUpload(file *multipart.FileHeader) error {
	// Check file size
	if file.Size > constants.MaxProfilePhotoSize {
		return fmt.Errorf("file size too large. Maximum size: %d MB", constants.MaxProfilePhotoSize/(1024*1024))
	}

	if file.Size == 0 {
		return fmt.Errorf("file is empty")
	}

	// Check file type
	if !s.isValidImageType(file.Filename) {
		return fmt.Errorf("invalid file type. Allowed types: %s", constants.AllowedProfileTypes)
	}

	return nil
}
