package routes

import (
	"dailytrackr/user-service/handlers"
	"dailytrackr/user-service/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes sets up all user-related routes
func SetupUserRoutes(r *gin.Engine, userHandlers *handlers.UserHandlers) {
	// Public routes (no authentication required)
	auth := r.Group("/auth")
	{
		auth.POST("/register", userHandlers.Register)
		auth.POST("/login", userHandlers.Login)
	}

	// Protected routes (authentication required)
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		// User profile routes
		users := api.Group("/users")
		{
			users.GET("/profile", userHandlers.GetProfile)
			users.GET("/:id", userHandlers.GetUserByID) // For other services
		}
	}
}
