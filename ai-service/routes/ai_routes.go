package routes

import (
	"dailytrackr/ai-service/handlers"
	"dailytrackr/shared/config"
	"dailytrackr/shared/constants"
	"dailytrackr/shared/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens for Gin (AI service specific)
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

// SetupAIRoutes sets up all AI-related routes
func SetupAIRoutes(r *gin.Engine, aiHandlers *handlers.AIHandlers) {
	// API v1 routes with authentication
	api := r.Group("/api/v1")
	api.Use(AuthMiddleware())

	// AI routes
	ai := api.Group("/ai")
	{
		// Daily summary generation
		ai.POST("/daily-summary", aiHandlers.GenerateDailySummary)
		ai.GET("/daily-summary", aiHandlers.GenerateDailySummary) // Support GET with query params

		// Habit recommendations
		ai.POST("/habit-recommendation", aiHandlers.GenerateHabitRecommendation)
		ai.GET("/habit-recommendation", aiHandlers.GenerateHabitRecommendation) // Support GET with query params

		// User insights
		ai.GET("/insights", aiHandlers.GetInsights)

		// Activity analysis
		ai.POST("/analyze-activities", aiHandlers.AnalyzeActivities)
		ai.GET("/analyze-activities", aiHandlers.AnalyzeActivities) // Support GET with query params

		// Productivity tips
		ai.GET("/productivity-tips", aiHandlers.GetProductivityTips)
	}
}
