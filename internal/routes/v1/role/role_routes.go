package role

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoleRoutes(router fiber.Router) {
	roles := router.Group("/roles")

	// Protected role routes (require JWT)
	roles.Use(middleware.JWTMiddleware())

	// IMPORTANT: Specific routes MUST come BEFORE parameterized routes
	// Specific routes (no parameters)
	roles.Get("/", handler.GetRoles)               // GET /api/v1/roles
	roles.Get("/deleted", handler.GetDeletedRoles) // GET /api/v1/roles/deleted

	// Parameterized routes (MUST be at the end)
	roles.Get("/:id", handler.GetRoleByID)         // GET /api/v1/roles/:id
	roles.Post("/", handler.CreateRole)            // POST /api/v1/roles
	roles.Put("/:id", handler.UpdateRole)          // PUT /api/v1/roles/:id
	roles.Put("/:id/restore", handler.RestoreRole) // PUT /api/v1/roles/:id/restore
	roles.Delete("/:id", handler.DeleteRole)       // DELETE /api/v1/roles/:id
}
