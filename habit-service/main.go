package main

import (
	"log"
	"net/http"
	"time"

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

	// ================================================
	// FIXED CORS MIDDLEWARE - Proper configuration
	// ================================================
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS, echo.PATCH},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		ExposeHeaders:    []string{echo.HeaderContentLength},
		AllowCredentials: false,
		MaxAge:           int(12 * time.Hour / time.Second),
	}))

	// Additional manual CORS handling for complex preflight requests
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Set CORS headers for all requests
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With")
			c.Response().Header().Set("Access-Control-Max-Age", "43200")

			// Handle preflight requests
			if c.Request().Method == "OPTIONS" {
				return c.NoContent(http.StatusNoContent)
			}

			return next(c)
		}
	})

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"service": "habit-service",
			"status":  "healthy",
			"version": "1.0.0",
			"cors":    "enabled",
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
