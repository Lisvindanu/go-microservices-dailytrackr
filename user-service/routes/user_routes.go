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
			// Profile management
			users.GET("/profile", userHandlers.GetProfile)
			users.PUT("/profile", userHandlers.UpdateProfile)
			users.PATCH("/profile", userHandlers.UpdateProfile) // Support both PUT and PATCH

			// Password management
			users.PUT("/password", userHandlers.ChangePassword)
			users.PATCH("/password", userHandlers.ChangePassword) // Support both PUT and PATCH

			// Profile photo management
			users.POST("/profile/photo", userHandlers.UploadProfilePhoto)
			users.PUT("/profile/photo", userHandlers.UploadProfilePhoto) // Alternative endpoint

			// Account management
			users.DELETE("/account", userHandlers.DeleteAccount)

			// For other services to get user info
			users.GET("/:id", userHandlers.GetUserByID)
		}
	}
}
