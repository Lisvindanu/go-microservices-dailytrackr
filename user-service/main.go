package main

import (
	"log"
	"net/http"
	"time"

	"dailytrackr/shared/config"
	"dailytrackr/shared/database"
	"dailytrackr/user-service/handlers"
	"dailytrackr/user-service/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	db, err := database.GetMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	log.Println("✅ Database connection established")

	// Initialize handlers with database
	userHandlers := handlers.NewUserHandlers(db, cfg)

	// Setup Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Untuk development
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // Set to false when using wildcard
		MaxAge:           12 * time.Hour,
	}))

	// Additional manual CORS handling untuk edge cases
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept")
		c.Header("Access-Control-Max-Age", "43200")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":  "user-service",
			"status":   "healthy",
			"version":  "1.0.0",
			"features": []string{"authentication", "profile_management", "photo_upload"},
			"cors":     "enabled",
		})
	})

	// Setup routes
	routes.SetupUserRoutes(r, userHandlers)

	// Start server
	port := ":" + cfg.UserServicePort
	log.Printf("🚀 User Service starting on port %s", port)
	log.Printf("📋 Available endpoints:")
	log.Printf("   Authentication:")
	log.Printf("   - POST /auth/register")
	log.Printf("   - POST /auth/login")
	log.Printf("   Profile Management:")
	log.Printf("   - GET  /api/v1/users/profile")
	log.Printf("   - PUT  /api/v1/users/profile")
	log.Printf("   - PATCH /api/v1/users/profile")
	log.Printf("   Security:")
	log.Printf("   - PUT  /api/v1/users/password")
	log.Printf("   - PATCH /api/v1/users/password")
	log.Printf("   Photo:")
	log.Printf("   - POST /api/v1/users/profile/photo")
	log.Printf("   - PUT  /api/v1/users/profile/photo")
	log.Printf("   Account:")
	log.Printf("   - DELETE /api/v1/users/account")
	log.Printf("   Service:")
	log.Printf("   - GET  /api/v1/users/:id")
	log.Printf("   - GET  /health")
	log.Printf("🌐 Service URL: http://localhost:%s", cfg.UserServicePort)

	// Warn if Cloudinary not configured
	if cfg.CloudinaryCloudName == "" || cfg.CloudinaryAPIKey == "" || cfg.CloudinaryAPISecret == "" {
		log.Printf("⚠️  Warning: Cloudinary not configured. Photo upload will be disabled.")
		log.Printf("   To enable photo upload, set CLOUDINARY_CLOUD_NAME, CLOUDINARY_API_KEY, and CLOUDINARY_API_SECRET in .env")
	} else {
		log.Printf("📸 Photo upload service enabled with Cloudinary")
	}

	log.Fatal(r.Run(port))
}
