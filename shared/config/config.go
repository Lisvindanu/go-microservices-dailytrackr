package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Service Ports
	GatewayPort      string
	UserServicePort  string
	ActivityPort     string
	HabitPort        string
	NotificationPort string
	StatPort         string
	AIPort           string

	// JWT
	JWTSecret      string
	JWTExpireHours int

	// External APIs
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
	GeminiAPIKey        string

	// Email
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string

	// WhatsApp (Optional)
	WhatsAppAPIURL   string
	WhatsAppAPIToken string

	// Redis (Optional)
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// Environment
	Environment string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Force set GO111MODULE if specified in .env
	if goModule := os.Getenv("GO111MODULE"); goModule != "" {
		os.Setenv("GO111MODULE", goModule)
	}

	config := &Config{
		// Database - DEFAULT TO MYSQL
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),         // MySQL default port
		DBUser:     getEnv("DB_USER", "root"),         // MySQL default user
		DBPassword: getEnv("DB_PASSWORD", "password"), // MySQL default (no password)
		DBName:     getEnv("DB_NAME", "dailytrackr"),

		// Service Ports
		GatewayPort:      getEnv("GATEWAY_PORT", "3000"),
		UserServicePort:  getEnv("USER_SERVICE_PORT", "3001"),
		ActivityPort:     getEnv("ACTIVITY_SERVICE_PORT", "3002"),
		HabitPort:        getEnv("HABIT_SERVICE_PORT", "3003"),
		NotificationPort: getEnv("NOTIFICATION_SERVICE_PORT", "3004"),
		StatPort:         getEnv("STAT_SERVICE_PORT", "3005"),
		AIPort:           getEnv("AI_SERVICE_PORT", "3006"),

		// JWT
		JWTSecret:      getEnv("JWT_SECRET", "your-super-secret-jwt-key-here"),
		JWTExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", 24),

		// External APIs
		CloudinaryCloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:    getEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret: getEnv("CLOUDINARY_API_SECRET", ""),
		GeminiAPIKey:        getEnv("GEMINI_API_KEY", ""),

		// Email
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),

		// WhatsApp
		WhatsAppAPIURL:   getEnv("WHATSAPP_API_URL", ""),
		WhatsAppAPIToken: getEnv("WHATSAPP_API_TOKEN", ""),

		// Redis
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		// Environment
		Environment: getEnv("ENV", "development"),
	}

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetMySQLDSN returns the MySQL connection string
func (c *Config) GetMySQLDSN() string {
	// MySQL DSN format: username:password@tcp(host:port)/database
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

// GetDatabaseURL returns the PostgreSQL connection string (for compatibility)
func (c *Config) GetDatabaseURL() string {
	return "host=" + c.DBHost + " port=" + c.DBPort + " user=" + c.DBUser +
		" password=" + c.DBPassword + " dbname=" + c.DBName + " sslmode=disable"
}

// IsDevelopment checks if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction checks if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
