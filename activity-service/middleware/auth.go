package middleware

import (
	"strings"

	"dailytrackr/shared/config"
	"dailytrackr/shared/constants"
	"dailytrackr/shared/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT tokens for Fiber
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get(constants.AuthorizationHeader)
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": constants.ErrMissingToken,
			})
		}

		// Check if header has Bearer prefix
		if !strings.HasPrefix(authHeader, constants.BearerPrefix) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": constants.ErrInvalidToken,
			})
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, constants.BearerPrefix)
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": constants.ErrInvalidToken,
			})
		}

		// Load config and validate token
		cfg := config.LoadConfig()
		claims, err := utils.ValidateJWT(token, cfg.JWTSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": constants.ErrInvalidToken,
			})
		}

		// Set user information in context
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}
