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

// Error Messages
const (
	ErrInvalidToken       = "invalid or expired token"
	ErrMissingToken       = "authorization token required"
	ErrInvalidCredentials = "invalid email or password"
	ErrUserNotFound       = "user not found"
	ErrEmailAlreadyExists = "email already exists"
	ErrUsernameExists     = "username already exists"
	ErrActivityNotFound   = "activity not found"
	ErrHabitNotFound      = "habit not found"
	ErrHabitLogNotFound   = "habit log not found"
	ErrUnauthorizedAccess = "unauthorized access to resource"
	ErrInvalidRequestBody = "invalid request body"
	ErrDatabaseConnection = "database connection error"
)

// Success Messages
const (
	MsgUserCreated     = "user created successfully"
	MsgLoginSuccess    = "login successful"
	MsgActivityCreated = "activity created successfully"
	MsgActivityUpdated = "activity updated successfully"
	MsgActivityDeleted = "activity deleted successfully"
	MsgHabitCreated    = "habit created successfully"
	MsgHabitUpdated    = "habit updated successfully"
	MsgHabitDeleted    = "habit deleted successfully"
	MsgHabitLogCreated = "habit log created successfully"
	MsgHabitLogUpdated = "habit log updated successfully"
	MsgPhotoUploaded   = "photo uploaded successfully"
)

// File Upload
const (
	MaxFileSize       = 10 << 20 // 10 MB
	AllowedImageTypes = "jpg,jpeg,png,gif,webp"
	UploadPath        = "uploads/dailytrackr/"
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
