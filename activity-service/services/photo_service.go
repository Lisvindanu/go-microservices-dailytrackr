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

	"github.com/cloudinary/cloudinary-go/v2" // <-- PERBAIKAN: Menghapus "com" yang berlebih
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type PhotoService struct {
	cloudinary *cloudinary.Cloudinary
	config     *config.Config
}

// NewPhotoService creates a new photo service instance.
// It will now log a fatal error if Cloudinary credentials are not properly configured.
func NewPhotoService(cfg *config.Config) *PhotoService {
	log.Printf("Cloudinary Config Check - Name: '%s', Key: '%s', Secret Is Set: %v",
		cfg.CloudinaryCloudName,
		cfg.CloudinaryAPIKey,
		cfg.CloudinaryAPISecret != "")

	// Initialize Cloudinary. It's now mandatory.
	if cfg.CloudinaryCloudName == "" || cfg.CloudinaryAPIKey == "" || cfg.CloudinaryAPISecret == "" {
		log.Fatalf("FATAL: Cloudinary credentials are not configured. Please check your .env file.")
	}

	cld, err := cloudinary.NewFromParams(
		cfg.CloudinaryCloudName,
		cfg.CloudinaryAPIKey,
		cfg.CloudinaryAPISecret,
	)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize Cloudinary client: %v", err)
	}

	log.Println("✅ Cloudinary client initialized successfully!")

	return &PhotoService{
		cloudinary: cld,
		config:     cfg,
	}
}

// UploadPhoto uploads a photo to Cloudinary. The fallback logic has been removed.
func (s *PhotoService) UploadPhoto(file *multipart.FileHeader) (string, error) {
	// Validate file type
	if !s.isValidImageType(file.Filename) {
		return "", fmt.Errorf("invalid file type. Allowed types: %s", constants.AllowedImageTypes)
	}

	// Check file size
	if file.Size > constants.MaxFileSize {
		return "", fmt.Errorf("file size too large. Maximum size: %d MB", constants.MaxFileSize/1024/1024)
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Generate unique filename for Cloudinary
	filename := s.generateFilename(file.Filename)

	// Upload to Cloudinary
	log.Printf("Uploading '%s' to Cloudinary folder '%s'", filename, constants.UploadPath)
	uploadResult, err := s.cloudinary.Upload.Upload(
		context.Background(),
		src,
		uploader.UploadParams{
			PublicID:     filename,
			Folder:       constants.UploadPath,
			ResourceType: "image",
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %v", err)
	}

	log.Printf("Successfully uploaded. Secure URL: %s", uploadResult.SecureURL)
	return uploadResult.SecureURL, nil
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
	// Variabel 'ext' dihapus karena tidak digunakan.
	timestamp := time.Now().Unix()
	randomStr := generateRandomString(8)
	// Return filename without extension, Cloudinary will handle it.
	return fmt.Sprintf("activity_%d_%s", timestamp, randomStr)
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	// Use a more robust random source if available, but this is fine for filenames.
	for i := range result {
		seed := time.Now().UnixNano() + int64(i)
		result[i] = charset[seed%int64(len(charset))]
	}
	return string(result)
}
