package productunit

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProductUnitRoutes(router fiber.Router) {
	productUnit := router.Group("/product-units")

	// Apply JWT middleware
	productUnit.Use(middleware.JWTMiddleware())

	// GET /api/v1/product-units - Get all product units
	productUnit.Get("/", handler.GetProductUnits)

	// GET /api/v1/product-units/product/:productId - Get product units by product ID
	productUnit.Get("/product/:productId", handler.GetProductUnitsByProduct)

	// GET /api/v1/product-units/:id - Get product unit by ID
	productUnit.Get("/:id", handler.GetProductUnitByID)

	// POST /api/v1/product-units - Create new product unit
	productUnit.Post("/", handler.CreateProductUnit)

	// PUT /api/v1/product-units/:id - Update product unit by ID
	productUnit.Put("/:id", handler.UpdateProductUnit)

	// DELETE /api/v1/product-units/:id - Delete product unit by ID
	productUnit.Delete("/:id", handler.DeleteProductUnit)
}
