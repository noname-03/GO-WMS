package inventory

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInventoryRoutes(api fiber.Router) {
	inventory := api.Group("/inventory")

	// Protected routes - require JWT authentication
	inventory.Get("/stock", middleware.JWTMiddleware(), handler.GetInventoryStock)
}
