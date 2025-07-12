package main

import (
	"log"

	"dailytrackr/habit-service/handlers"
	"dailytrackr/habit-service/routes"
	"dailytrackr/shared/config"
	"dailytrackr/shared/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"service": "habit-service",
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	// Initialize handlers
	habitHandlers := handlers.NewHabitHandlers(db, cfg)

	// Setup routes
	routes.SetupHabitRoutes(e, habitHandlers)

	// Start server
	port := ":" + cfg.HabitPort
	log.Printf("🚀 Habit Service starting on port %s", port)
	log.Printf("📋 Available endpoints:")
	log.Printf("   - GET  /health")
	log.Printf("   - POST /api/v1/habits")
	log.Printf("   - GET  /api/v1/habits")
	log.Printf("   - GET  /api/v1/habits/:id")
	log.Printf("   - PUT  /api/v1/habits/:id")
	log.Printf("   - DELETE /api/v1/habits/:id")
	log.Printf("   - POST /api/v1/habits/:id/logs")
	log.Printf("   - GET  /api/v1/habits/:id/logs")
	log.Printf("   - PUT  /api/v1/habit-logs/:id")
	log.Printf("   - GET  /api/v1/habits/:id/stats")

	log.Fatal(e.Start(port))
}
