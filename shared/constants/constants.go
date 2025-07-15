package constants

// Service Names
const (
	GatewayService      = "gateway"
	UserService         = "user-service"
	ActivityService     = "activity-service"
	HabitService        = "habit-service"
	NotificationService = "notification-service"
	StatService         = "stat-service"
	AIService           = "ai-service"
)

// Default Service Ports
const (
	DefaultGatewayPort      = "3000"
	DefaultUserServicePort  = "3001"
	DefaultActivityPort     = "3002"
	DefaultHabitPort        = "3003"
	DefaultNotificationPort = "3004"
	DefaultStatPort         = "3005"
	DefaultAIPort           = "3006"
)

// Database Tables
const (
	UsersTable      = "users"
	ActivitiesTable = "activities"
	HabitsTable     = "habits"
	HabitLogsTable  = "habit_logs"
)

// Habit Log Status
const (
	HabitStatusDone    = "DONE"
	HabitStatusSkipped = "SKIPPED"
	HabitStatusFailed  = "FAILED"
)

// Environment Values
const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
	EnvTesting     = "testing"
)

// JWT Claims
const (
	JWTIssuer  = "dailytrackr"
	JWTSubject = "user-auth"
)

// Default Values
const (
	DefaultJWTExpireHours = 24
	DefaultPageLimit      = 20
	DefaultPageOffset     = 0
)

// HTTP Headers
const (
	AuthorizationHeader = "Authorization"
	ContentTypeHeader   = "Content-Type"
	BearerPrefix        = "Bearer "
)

// Error Messages - General
const (
	ErrInvalidToken       = "invalid or expired token"
	ErrMissingToken       = "authorization token required"
	ErrInvalidCredentials = "invalid email or password"
	ErrUnauthorizedAccess = "unauthorized access to resource"
	ErrInvalidRequestBody = "invalid request body"
	ErrDatabaseConnection = "database connection error"
)

// Error Messages - User Related
const (
	ErrUserNotFound            = "user not found"
	ErrEmailAlreadyExists      = "email already exists"
	ErrUsernameExists          = "username already exists"
	ErrInvalidCurrentPassword  = "current password is incorrect"
	ErrSamePassword            = "new password must be different from current password"
	ErrWeakPassword            = "password must be at least 6 characters long"
	ErrInvalidFileType         = "invalid file type for profile photo"
	ErrFileTooLarge            = "file size too large"
	ErrPhotoUploadFailed       = "failed to upload profile photo"
	ErrPhotoServiceUnavailable = "photo upload service is not available"
	ErrInvalidBio              = "bio contains invalid content"
	ErrReservedUsername        = "username is reserved and cannot be used"
)

// Error Messages - Activity Related
const (
	ErrActivityNotFound = "activity not found"
)

// Error Messages - Habit Related
const (
	ErrHabitNotFound    = "habit not found"
	ErrHabitLogNotFound = "habit log not found"
)

// Success Messages - User Related
const (
	MsgUserCreated          = "user created successfully"
	MsgLoginSuccess         = "login successful"
	MsgProfileUpdated       = "profile updated successfully"
	MsgPasswordChanged      = "password changed successfully"
	MsgProfilePhotoUploaded = "profile photo uploaded successfully"
	MsgAccountDeleted       = "account deleted successfully"
	MsgProfileRetrieved     = "profile retrieved successfully"
)

// Success Messages - Activity Related
const (
	MsgActivityCreated = "activity created successfully"
	MsgActivityUpdated = "activity updated successfully"
	MsgActivityDeleted = "activity deleted successfully"
)

// Success Messages - Habit Related
const (
	MsgHabitCreated    = "habit created successfully"
	MsgHabitUpdated    = "habit updated successfully"
	MsgHabitDeleted    = "habit deleted successfully"
	MsgHabitLogCreated = "habit log created successfully"
	MsgHabitLogUpdated = "habit log updated successfully"
)

// Success Messages - General
const (
	MsgPhotoUploaded = "photo uploaded successfully"
)

// File Upload Settings
const (
	MaxFileSize         = 10 << 20 // 10 MB for activities
	MaxProfilePhotoSize = 5 << 20  // 5 MB for profile photos
	AllowedImageTypes   = "jpg,jpeg,png,gif,webp"
	AllowedProfileTypes = "jpg,jpeg,png,webp"
	UploadPath          = "uploads/dailytrackr/"
	ProfilePhotoPath    = "uploads/dailytrackr/profiles/"
)

// Notification Types
const (
	NotificationHabitReminder = "habit_reminder"
	NotificationDailySummary  = "daily_summary"
	NotificationWeeklyReport  = "weekly_report"
)

// AI Service Types
const (
	AITypeDailySummary        = "daily_summary"
	AITypeHabitRecommendation = "habit_recommendation"
)

// Time Formats
const (
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04"
	DateTimeFormat = "2006-01-02T15:04:05Z"
)

// User Profile Validation
const (
	MinUsernameLength = 3
	MaxUsernameLength = 50
	MaxEmailLength    = 255
	MinPasswordLength = 6
	MaxPasswordLength = 128
	MaxBioLength      = 500
)
