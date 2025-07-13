package routes

import (
	"dailytrackr/shared/config"
	"dailytrackr/shared/constants"
	"dailytrackr/shared/utils"
	"dailytrackr/stat-service/handlers"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens for Gin (stat service specific)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader(constants.AuthorizationHeader)
		if authHeader == "" {
			utils.SendUnauthorizedResponse(c.Writer, constants.ErrMissingToken)
			c.Abort()
			return
		}

		// Check if header has Bearer prefix
		if !strings.HasPrefix(authHeader, constants.BearerPrefix) {
			utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, constants.BearerPrefix)
		if token == "" {
			utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
			c.Abort()
			return
		}

		// Load config and validate token
		cfg := config.LoadConfig()
		claims, err := utils.ValidateJWT(token, cfg.JWTSecret)
		if err != nil {
			utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// SetupStatRoutes sets up all statistics-related routes
func SetupStatRoutes(r *gin.Engine, statHandlers *handlers.StatHandlers) {
	// API v1 routes with authentication
	api := r.Group("/api/v1")
	api.Use(AuthMiddleware())

	// Statistics routes
	stats := api.Group("/stats")
	{
		// Dashboard overview
		stats.GET("/dashboard", statHandlers.GetDashboard)

		// Activity statistics
		activities := stats.Group("/activities")
		{
			activities.GET("/summary", statHandlers.GetActivitySummary)
			activities.GET("/chart", statHandlers.GetActivityChart)
		}

		// Habit statistics
		habits := stats.Group("/habits")
		{
			habits.GET("/progress", statHandlers.GetHabitProgress)
		}

		// Expense statistics
		expenses := stats.Group("/expenses")
		{
			expenses.GET("/report", statHandlers.GetExpenseReport)
		}
	}
}
