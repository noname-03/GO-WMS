package user

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(router fiber.Router) {
	// User routes group
	users := router.Group("/users")

	// Protected user routes (require JWT)
	users.Use(middleware.JWTMiddleware())

	// IMPORTANT: Specific routes MUST come BEFORE parameterized routes
	// Specific routes (no parameters)
	users.Get("/", handler.GetUsers)                         // GET /api/v1/users
	users.Get("/deleted", handler.GetDeletedUsers)           // GET /api/v1/users/deleted
	users.Get("/minimal", handler.GetUsersMinimal)           // GET /api/v1/users/minimal
	users.Get("/raw", handler.GetUsersRaw)                   // GET /api/v1/users/raw
	users.Get("/repository", handler.GetUsersFromRepository) // GET /api/v1/users/repository
	users.Get("/search", handler.SearchUsers)                // GET /api/v1/users/search?q=keyword
	users.Get("/stats", handler.GetUserStats)                // GET /api/v1/users/stats
	users.Get("/rawQuery", handler.GetUsersFromRepository)   // GET /api/v1/users/rawQuery

	// Parameterized routes (MUST be at the end)
	users.Get("/:id", handler.GetUserByIDRaw)      // GET /api/v1/users/:id
	users.Put("/:id/restore", handler.RestoreUser) // PUT /api/v1/users/:id/restore
}
