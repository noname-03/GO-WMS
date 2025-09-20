package brand

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupBrandRoutes(router fiber.Router) {
	brands := router.Group("/brands")

	// Protected brand routes (require JWT)
	brands.Use(middleware.JWTMiddleware())

	// IMPORTANT: Specific routes MUST come BEFORE parameterized routes
	// Specific routes (no parameters)
	brands.Get("/", handler.GetBrands) // GET /api/v1/brands

	// Parameterized routes (MUST be at the end)
	brands.Get("/:id", handler.GetBrandByID)   // GET /api/v1/brands/:id
	brands.Post("/", handler.CreateBrand)      // POST /api/v1/brands
	brands.Put("/:id", handler.UpdateBrand)    // PUT /api/v1/brands/:id
	brands.Delete("/:id", handler.DeleteBrand) // DELETE /api/v1/brands/:id
}
