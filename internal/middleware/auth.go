package middleware

import (
	"myapp/internal/utils"
	"myapp/pkg/helper"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return helper.Fail(c, 401, "Unauthorized", "Authorization header required")
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return helper.Fail(c, 401, "Unauthorized", "Bearer token required")
		}

		// Extract token (remove "Bearer " prefix)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return helper.Fail(c, 401, "Unauthorized", "Token is empty")
		}

		// Validate token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			return helper.Fail(c, 401, "Invalid token", err.Error())
		}

		// Store user info in context
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}
