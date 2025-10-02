package product

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterProductRoutes(app fiber.Router) {
	productRoutes := app.Group("/products")

	// Public routes (optional, if any products are publicly viewable)
	// productRoutes.Get("/public", handler.GetPublicProducts)

	// Protected routes - all require authentication
	productRoutes.Use(middleware.JWTMiddleware())

	// CRUD operations for products
	productRoutes.Get("/", handler.GetProducts)                                 // GET /api/v1/products
	productRoutes.Get("/categories/:categoryId", handler.GetProductsByCategory) // GET /api/v1/products/categories/:id
	productRoutes.Get("/deleted", handler.GetDeletedProducts)                   // GET /api/v1/products/deleted
	productRoutes.Get("/:id", handler.GetProductByID)                           // GET /api/v1/products/:id
	productRoutes.Post("/", handler.CreateProduct)                              // POST /api/v1/products
	productRoutes.Put("/:id", handler.UpdateProduct)                            // PUT /api/v1/products/:id
	productRoutes.Put("/:id/restore", handler.RestoreProduct)                   // PUT /api/v1/products/:id/restore
	productRoutes.Delete("/:id", handler.DeleteProduct)                         // DELETE /api/v1/products/:id

	// Product batch routes - nested under products
	productRoutes.Get("/:productId/batches", handler.GetProductBatchesByProduct) // GET /api/v1/products/:productId/batches
}
