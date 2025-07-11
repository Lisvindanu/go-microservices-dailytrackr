package handlers

import (
	"database/sql"
	"strconv"

	"dailytrackr/shared/config"
	"dailytrackr/shared/constants"
	"dailytrackr/shared/dto"
	"dailytrackr/shared/utils"
	"dailytrackr/user-service/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandlers struct {
	userRepo *models.UserRepository
	config   *config.Config
}

// NewUserHandlers creates a new user handlers instance
func NewUserHandlers(db *sql.DB, cfg *config.Config) *UserHandlers {
	return &UserHandlers{
		userRepo: models.NewUserRepository(db),
		config:   cfg,
	}
}

// Register handles user registration
func (h *UserHandlers) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestResponse(c.Writer, "Invalid request body", err)
		return
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
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
		Username:     req.Username,
		Email:        req.Email,
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
	userResponse := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

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

	// Get user by email
	user, err := h.userRepo.GetByEmail(req.Email)
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
	userResponse := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	authResponse := dto.AuthResponse{
		Token: token,
		User:  userResponse,
	}

	utils.SendSuccessResponse(c.Writer, constants.MsgLoginSuccess, authResponse)
}

// GetProfile handles getting user profile
func (h *UserHandlers) GetProfile(c *gin.Context) {
	// Get user ID from JWT token (will be set by auth middleware)
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
	userResponse := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	utils.SendSuccessResponse(c.Writer, "Profile retrieved successfully", userResponse)
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

	userResponse := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	utils.SendSuccessResponse(c.Writer, "User retrieved successfully", userResponse)
}
