package middleware

import (
	"net/http"
	"strings"

	"dailytrackr/shared/config"
	"dailytrackr/shared/constants"
	"dailytrackr/shared/utils"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware validates JWT tokens for Echo
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get authorization header
			authHeader := c.Request().Header.Get(constants.AuthorizationHeader)
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"message": constants.ErrMissingToken,
				})
			}

			// Check if header has Bearer prefix
			if !strings.HasPrefix(authHeader, constants.BearerPrefix) {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"message": constants.ErrInvalidToken,
				})
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, constants.BearerPrefix)
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"message": constants.ErrInvalidToken,
				})
			}

			// Load config and validate token
			cfg := config.LoadConfig()
			claims, err := utils.ValidateJWT(token, cfg.JWTSecret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"message": constants.ErrInvalidToken,
				})
			}

			// Set user information in context
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}
