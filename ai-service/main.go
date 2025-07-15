package main

import (
	"log"
	"net/http"
	"time"

	"dailytrackr/ai-service/handlers"
	"dailytrackr/ai-service/routes"
	"dailytrackr/shared/config"
	"dailytrackr/shared/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Check Gemini API key
	if cfg.GeminiAPIKey == "" {
		log.Fatal("❌ GEMINI_API_KEY is required for AI service. Please set it in .env file")
	}

	// Initialize database connection
	db, err := database.GetMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Setup Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// ================================================
	// FIXED CORS MIDDLEWARE - No AllowCredentials with wildcard
	// ================================================
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Untuk development
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // 🔧 FIXED: Set to false when using wildcard
		MaxAge:           12 * time.Hour,
	}))

	// Additional manual CORS handling untuk edge cases
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
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
			"service":    "ai-service",
			"status":     "healthy",
			"version":    "1.0.0",
			"cors":       "enabled",
			"ai_enabled": cfg.GeminiAPIKey != "",
		})
	})

	// Initialize handlers
	aiHandlers := handlers.NewAIHandlers(db, cfg)

	// Setup routes
	routes.SetupAIRoutes(r, aiHandlers)

	// Start server
	port := ":" + cfg.AIPort
	log.Printf("🚀 AI Service starting on port %s", port)
	log.Fatal(r.Run(port))
}
