package reseller

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterResellerRoutes(app fiber.Router) {
	resellerRoutes := app.Group("/resellers")

	// All routes require authentication
	resellerRoutes.Use(middleware.JWTMiddleware())

	// Get all resellers
	resellerRoutes.Get("/", handler.GetAllResellers) // GET /api/v1/resellers
}	
