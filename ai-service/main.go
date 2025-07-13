package main

import (
	"log"
	"net/http"

	"dailytrackr/ai-service/handlers"
	"dailytrackr/ai-service/routes"
	"dailytrackr/shared/config"
	"dailytrackr/shared/database"

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

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

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
			"ai_enabled": cfg.GeminiAPIKey != "",
			"features": []string{
				"daily_summary",
				"habit_recommendation",
				"activity_insights",
				"productivity_analysis",
			},
		})
	})

	// Initialize handlers
	aiHandlers := handlers.NewAIHandlers(db, cfg)

	// Setup routes
	routes.SetupAIRoutes(r, aiHandlers)

	// Start server
	port := ":" + cfg.AIPort
	log.Printf("🚀 AI Service starting on port %s", port)
	log.Printf("🧠 Gemini AI integration: ✅ ENABLED")
	log.Printf("📋 Available endpoints:")
	log.Printf("   - GET  /health")
	log.Printf("   - POST /api/v1/ai/daily-summary")
	log.Printf("   - POST /api/v1/ai/habit-recommendation")
	log.Printf("   - GET  /api/v1/ai/insights")
	log.Printf("   - POST /api/v1/ai/analyze-activities")
	log.Printf("   - GET  /api/v1/ai/productivity-tips")

	log.Fatal(r.Run(port))
}
