package handlers

import (
	"database/sql"
	"strconv"

	"dailytrackr/shared/config"
	"dailytrackr/shared/constants"
	"dailytrackr/shared/dto"
	"dailytrackr/shared/utils"
	"dailytrackr/user-service/models"
	"dailytrackr/user-service/services"
	"dailytrackr/user-service/validators"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandlers struct {
	userRepo     *models.UserRepository
	photoService *services.PhotoService
	validator    *validators.UserValidator
	config       *config.Config
}

// NewUserHandlers creates a new user handlers instance
func NewUserHandlers(db *sql.DB, cfg *config.Config) *UserHandlers {
	return &UserHandlers{
		userRepo:     models.NewUserRepository(db),
		photoService: services.NewPhotoService(cfg),
		validator:    validators.NewUserValidator(),
		config:       cfg,
	}
}

// Register handles user registration
func (h *UserHandlers) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestResponse(c.Writer, "Invalid request body", err)
		return
	}

	// Advanced validation
	if err := h.validator.ValidateRegistration(req); err != nil {
		utils.SendBadRequestResponse(c.Writer, "Validation failed", err)
		return
	}

	// Check if email already exists
	emailExists, err := h.userRepo.EmailExists(req.Email)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Database error", err)
		return
	}
	if emailExists {
		utils.SendBadRequestResponse(c.Writer, constants.ErrEmailAlreadyExists, nil)
		return
	}

	// Check if username already exists
	usernameExists, err := h.userRepo.UsernameExists(req.Username)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Database error", err)
		return
	}
	if usernameExists {
		utils.SendBadRequestResponse(c.Writer, constants.ErrUsernameExists, nil)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to hash password", err)
		return
	}

	// Create user
	user := &models.User{
		Username:     h.validator.SanitizeInput(req.Username),
		Email:        h.validator.SanitizeInput(req.Email),
		PasswordHash: string(hashedPassword),
	}

	if err := h.userRepo.Create(user); err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to create user", err)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(
		user.ID,
		user.Username,
		user.Email,
		h.config.JWTSecret,
		h.config.JWTExpireHours,
	)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to generate token", err)
		return
	}

	// Prepare response
	userResponse := h.convertToUserResponse(user)
	authResponse := dto.AuthResponse{
		Token: token,
		User:  userResponse,
	}

	utils.SendCreatedResponse(c.Writer, constants.MsgUserCreated, authResponse)
}

// Login handles user login
func (h *UserHandlers) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestResponse(c.Writer, "Invalid request body", err)
		return
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		utils.SendBadRequestResponse(c.Writer, "Validation failed", err)
		return
	}

	// Sanitize input
	email := h.validator.SanitizeInput(req.Email)

	// Get user by email
	user, err := h.userRepo.GetByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidCredentials)
			return
		}
		utils.SendInternalServerErrorResponse(c.Writer, "Database error", err)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidCredentials)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(
		user.ID,
		user.Username,
		user.Email,
		h.config.JWTSecret,
		h.config.JWTExpireHours,
	)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to generate token", err)
		return
	}

	// Prepare response
	userResponse := h.convertToUserResponse(user)
	authResponse := dto.AuthResponse{
		Token: token,
		User:  userResponse,
	}

	utils.SendSuccessResponse(c.Writer, constants.MsgLoginSuccess, authResponse)
}

// GetProfile handles getting user profile
func (h *UserHandlers) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
		return
	}

	// Get user from database
	user, err := h.userRepo.GetByID(userID.(int64))
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendNotFoundResponse(c.Writer, constants.ErrUserNotFound)
			return
		}
		utils.SendInternalServerErrorResponse(c.Writer, "Database error", err)
		return
	}

	// Prepare response
	userResponse := h.convertToUserResponse(user)
	utils.SendSuccessResponse(c.Writer, constants.MsgProfileRetrieved, userResponse)
}

// UpdateProfile handles updating user profile
func (h *UserHandlers) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestResponse(c.Writer, "Invalid request body", err)
		return
	}

	// Advanced validation
	if err := h.validator.ValidateProfileUpdate(req); err != nil {
		utils.SendBadRequestResponse(c.Writer, "Validation failed", err)
		return
	}

	// Get existing user
	user, err := h.userRepo.GetByID(userID.(int64))
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendNotFoundResponse(c.Writer, constants.ErrUserNotFound)
			return
		}
		utils.SendInternalServerErrorResponse(c.Writer, "Database error", err)
		return
	}

	// Check if new username already exists (if changed)
	if req.Username != "" && req.Username != user.Username {
		usernameExists, err := h.userRepo.UsernameExists(req.Username)
		if err != nil {
			utils.SendInternalServerErrorResponse(c.Writer, "Database error", err)
			return
		}
		if usernameExists {
			utils.SendBadRequestResponse(c.Writer, constants.ErrUsernameExists, nil)
			return
		}
		user.Username = h.validator.SanitizeInput(req.Username)
	}

	// Check if new email already exists (if changed)
	if req.Email != "" && req.Email != user.Email {
		emailExists, err := h.userRepo.EmailExists(req.Email)
		if err != nil {
			utils.SendInternalServerErrorResponse(c.Writer, "Database error", err)
			return
		}
		if emailExists {
			utils.SendBadRequestResponse(c.Writer, constants.ErrEmailAlreadyExists, nil)
			return
		}
		user.Email = h.validator.SanitizeInput(req.Email)
	}

	// Update bio if provided
	if req.Bio != nil {
		user.Bio = h.validator.SanitizeInput(*req.Bio)
	}

	// Update user in database
	if err := h.userRepo.Update(user); err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to update profile", err)
		return
	}

	// Get updated user
	updatedUser, err := h.userRepo.GetByID(userID.(int64))
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to get updated profile", err)
		return
	}

	userResponse := h.convertToUserResponse(updatedUser)
	utils.SendSuccessResponse(c.Writer, constants.MsgProfileUpdated, userResponse)
}

// ChangePassword handles password change
func (h *UserHandlers) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestResponse(c.Writer, "Invalid request body", err)
		return
	}

	// Advanced validation
	if err := h.validator.ValidatePasswordChange(req); err != nil {
		utils.SendBadRequestResponse(c.Writer, "Validation failed", err)
		return
	}

	// Get user from database
	user, err := h.userRepo.GetByID(userID.(int64))
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendNotFoundResponse(c.Writer, constants.ErrUserNotFound)
			return
		}
		utils.SendInternalServerErrorResponse(c.Writer, "Database error", err)
		return
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidCurrentPassword)
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to hash password", err)
		return
	}

	// Update password
	if err := h.userRepo.UpdatePassword(user.ID, string(hashedPassword)); err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to update password", err)
		return
	}

	utils.SendSuccessResponse(c.Writer, constants.MsgPasswordChanged, nil)
}

// UploadProfilePhoto handles profile photo upload
func (h *UserHandlers) UploadProfilePhoto(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
		return
	}

	// Get file from form
	file, err := c.FormFile("photo")
	if err != nil {
		utils.SendBadRequestResponse(c.Writer, "No photo file provided", err)
		return
	}

	// Validate file before upload
	if err := h.photoService.ValidateFileUpload(file); err != nil {
		utils.SendBadRequestResponse(c.Writer, "File validation failed", err)
		return
	}

	// Get current user to check if there's an existing photo
	user, err := h.userRepo.GetByID(userID.(int64))
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to get user info", err)
		return
	}

	// Upload new photo
	photoURL, err := h.photoService.UploadPhoto(file)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to upload photo", err)
		return
	}

	// Update user profile photo in database
	if err := h.userRepo.UpdateProfilePhoto(userID.(int64), photoURL); err != nil {
		// If database update fails, try to delete the uploaded photo
		if deleteErr := h.photoService.DeletePhoto(photoURL); deleteErr != nil {
			// Log the error but don't fail the request
			utils.SendInternalServerErrorResponse(c.Writer, "Failed to update profile photo", err)
			return
		}
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to update profile photo", err)
		return
	}

	// Optionally delete old photo if it exists
	if user.ProfilePhoto != "" && user.ProfilePhoto != photoURL {
		if deleteErr := h.photoService.DeletePhoto(user.ProfilePhoto); deleteErr != nil {
			// Log but don't fail - old photo cleanup is not critical
			// log.Printf("Warning: Failed to delete old profile photo: %v", deleteErr)
		}
	}

	response := dto.PhotoUploadResponse{
		URL:       photoURL,
		SecureURL: photoURL,
	}

	utils.SendSuccessResponse(c.Writer, constants.MsgProfilePhotoUploaded, response)
}

// DeleteAccount handles account deletion
func (h *UserHandlers) DeleteAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
		return
	}

	var req dto.DeleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestResponse(c.Writer, "Invalid request body", err)
		return
	}

	// Get user from database
	user, err := h.userRepo.GetByID(userID.(int64))
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendNotFoundResponse(c.Writer, constants.ErrUserNotFound)
			return
		}
		utils.SendInternalServerErrorResponse(c.Writer, "Database error", err)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		utils.SendUnauthorizedResponse(c.Writer, "Password is incorrect")
		return
	}

	// Delete profile photo if exists
	if user.ProfilePhoto != "" {
		if deleteErr := h.photoService.DeletePhoto(user.ProfilePhoto); deleteErr != nil {
			// Log but don't fail - photo cleanup is not critical for account deletion
			// log.Printf("Warning: Failed to delete profile photo during account deletion: %v", deleteErr)
		}
	}

	// Delete user account
	if err := h.userRepo.Delete(userID.(int64)); err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to delete account", err)
		return
	}

	utils.SendSuccessResponse(c.Writer, constants.MsgAccountDeleted, nil)
}

// GetUserByID handles getting user by ID (for other services)
func (h *UserHandlers) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestResponse(c.Writer, "Invalid user ID", err)
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendNotFoundResponse(c.Writer, constants.ErrUserNotFound)
			return
		}
		utils.SendInternalServerErrorResponse(c.Writer, "Database error", err)
		return
	}

	userResponse := h.convertToUserResponse(user)
	utils.SendSuccessResponse(c.Writer, "User retrieved successfully", userResponse)
}

// convertToUserResponse converts User model to UserResponse DTO
func (h *UserHandlers) convertToUserResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Bio:          user.Bio,
		ProfilePhoto: user.ProfilePhoto,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
