package main

import (
	"log"
	"time"

	"dailytrackr/activity-service/handlers"
	"dailytrackr/activity-service/routes"
	"dailytrackr/shared/config"
	"dailytrackr/shared/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	db, err := database.GetMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": "Internal server error",
				"error":   err.Error(),
			})
		},
	})

	// ================================================
	// FIXED CORS MIDDLEWARE - No AllowCredentials with wildcard
	// ================================================
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Authorization,Accept",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: false,                             // 🔧 FIXED: Set to false when using wildcard
		MaxAge:           int(12 * time.Hour / time.Second), // 12 hours in seconds
	}))

	// Additional manual CORS handling untuk edge cases
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept")
		c.Set("Access-Control-Max-Age", "43200")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	})

	app.Use(logger.New())

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "activity-service",
			"status":  "healthy",
			"version": "1.0.0",
			"cors":    "enabled",
		})
	})

	// Initialize handlers
	activityHandlers := handlers.NewActivityHandlers(db, cfg)

	// Setup routes
	routes.SetupActivityRoutes(app, activityHandlers)

	// Start server
	port := ":" + cfg.ActivityPort
	log.Printf("🚀 Activity Service starting on port %s", port)
	log.Printf("📋 Available endpoints:")
	log.Printf("   - GET  /health")
	log.Printf("   - POST /api/v1/activities")
	log.Printf("   - GET  /api/v1/activities")
	log.Printf("   - GET  /api/v1/activities/:id")
	log.Printf("   - PUT  /api/v1/activities/:id")
	log.Printf("   - DELETE /api/v1/activities/:id")
	log.Printf("   - POST /api/v1/activities/:id/photo")

	log.Fatal(app.Listen(port))
}
