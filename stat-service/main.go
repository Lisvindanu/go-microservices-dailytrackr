package main

import (
	"log"
	"net/http"

	"dailytrackr/shared/config"
	"dailytrackr/shared/database"
	"dailytrackr/stat-service/handlers"
	"dailytrackr/stat-service/routes"

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
			"service": "stat-service",
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	// Initialize handlers
	statHandlers := handlers.NewStatHandlers(db, cfg)

	// Setup routes
	routes.SetupStatRoutes(r, statHandlers)

	// Start server
	port := ":" + cfg.StatPort
	log.Printf("🚀 Statistics Service starting on port %s", port)
	log.Printf("📊 Available endpoints:")
	log.Printf("   - GET  /health")
	log.Printf("   - GET  /api/v1/stats/dashboard")
	log.Printf("   - GET  /api/v1/stats/activities/summary")
	log.Printf("   - GET  /api/v1/stats/habits/progress")
	log.Printf("   - GET  /api/v1/stats/activities/chart")
	log.Printf("   - GET  /api/v1/stats/expenses/report")

	log.Fatal(r.Run(port))
}
