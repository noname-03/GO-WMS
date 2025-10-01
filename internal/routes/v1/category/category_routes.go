package category

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupCategoryRoutes(router fiber.Router) {
	categories := router.Group("/categories")

	// Protected category routes (require JWT)
	categories.Use(middleware.JWTMiddleware())

	// IMPORTANT: Specific routes MUST come BEFORE parameterized routes
	// Specific routes (no parameters)
	categories.Get("/", handler.GetCategories)                      // GET /api/v1/categories
	categories.Get("/deleted", handler.GetDeletedCategories)        // GET /api/v1/categories/deleted
	categories.Get("/brand/:brandId", handler.GetCategoriesByBrand) // GET /api/v1/categories/brand/:brandId

	// Parameterized routes (MUST be at the end)
	categories.Get("/:id", handler.GetCategoryByID)         // GET /api/v1/categories/:id
	categories.Post("/", handler.CreateCategory)            // POST /api/v1/categories
	categories.Put("/:id", handler.UpdateCategory)          // PUT /api/v1/categories/:id
	categories.Put("/:id/restore", handler.RestoreCategory) // PUT /api/v1/categories/:id/restore
	categories.Delete("/:id", handler.DeleteCategory)       // DELETE /api/v1/categories/:id
}
