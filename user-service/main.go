package main

import (
	"log"
	"net/http"

	"dailytrackr/shared/config"
	"dailytrackr/shared/database"
	"dailytrackr/user-service/handlers"
	"dailytrackr/user-service/routes"

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

	// Initialize handlers with database
	userHandlers := handlers.NewUserHandlers(db, cfg)

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
			"service": "user-service",
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	// Setup routes
	routes.SetupUserRoutes(r, userHandlers)

	// Start server
	port := ":" + cfg.UserServicePort
	log.Printf("User Service starting on port %s", port)
	log.Fatal(r.Run(port))
}
