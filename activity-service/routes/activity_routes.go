package routes

import (
	"dailytrackr/activity-service/handlers"
	"dailytrackr/activity-service/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupActivityRoutes sets up all activity-related routes
func SetupActivityRoutes(app *fiber.App, activityHandlers *handlers.ActivityHandlers) {
	// API v1 routes with authentication
	api := app.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())

	// Activity routes
	activities := api.Group("/activities")
	{
		// CRUD operations
		activities.Post("/", activityHandlers.CreateActivity)      // POST /api/v1/activities
		activities.Get("/", activityHandlers.GetActivities)        // GET /api/v1/activities?page=1&limit=20&start_date=2025-01-01&end_date=2025-01-31
		activities.Get("/:id", activityHandlers.GetActivity)       // GET /api/v1/activities/:id
		activities.Put("/:id", activityHandlers.UpdateActivity)    // PUT /api/v1/activities/:id
		activities.Delete("/:id", activityHandlers.DeleteActivity) // DELETE /api/v1/activities/:id

		// Photo upload
		activities.Post("/:id/photo", activityHandlers.UploadPhoto) // POST /api/v1/activities/:id/photo
	}
}
