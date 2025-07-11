package middleware

import (
	"strings"

	"dailytrackr/shared/config"
	"dailytrackr/shared/constants"
	"dailytrackr/shared/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens
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
