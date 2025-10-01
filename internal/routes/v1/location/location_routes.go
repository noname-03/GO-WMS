package location

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func LocationRoutes(router fiber.Router) {
	location := router.Group("/locations")

	// Apply JWT middleware
	location.Use(middleware.JWTMiddleware())

	// GET /api/v1/locations - Get all locations
	location.Get("/", handler.GetLocations)

	// GET /api/v1/locations/deleted - Get deleted locations
	location.Get("/deleted", handler.GetDeletedLocations)

	// GET /api/v1/locations/user/:userId - Get locations by user ID
	location.Get("/user/:userId", handler.GetLocationsByUser)

	// GET /api/v1/locations/type/:type - Get locations by type (gudang/reseller)
	location.Get("/type/:type", handler.GetLocationsByType)

	// GET /api/v1/locations/:id - Get location by ID
	location.Get("/:id", handler.GetLocationByID)

	// POST /api/v1/locations - Create new location
	location.Post("/", handler.CreateLocation)

	// PUT /api/v1/locations/:id - Update location by ID
	location.Put("/:id", handler.UpdateLocation)

	// PUT /api/v1/locations/:id/restore - Restore deleted location
	location.Put("/:id/restore", handler.RestoreLocation)

	// DELETE /api/v1/locations/:id - Delete location by ID
	location.Delete("/:id", handler.DeleteLocation)
}
