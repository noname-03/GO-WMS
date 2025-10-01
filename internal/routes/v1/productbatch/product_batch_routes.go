package productbatch

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterProductBatchRoutes(app fiber.Router) {
	productBatchRoutes := app.Group("/product-batches")

	// All routes require authentication
	productBatchRoutes.Use(middleware.JWTMiddleware())

	// CRUD operations for product batches
	productBatchRoutes.Get("/", handler.GetProductBatches)               // GET /api/v1/product-batches
	productBatchRoutes.Get("/deleted", handler.GetDeletedProductBatches) // GET /api/v1/product-batches/deleted
	productBatchRoutes.Get("/:id", handler.GetProductBatchByID)          // GET /api/v1/product-batches/:id
	productBatchRoutes.Post("/", handler.CreateProductBatch)             // POST /api/v1/product-batches
	productBatchRoutes.Put("/:id", handler.UpdateProductBatch)           // PUT /api/v1/product-batches/:id
	productBatchRoutes.Put("/:id/restore", handler.RestoreProductBatch)  // PUT /api/v1/product-batches/:id/restore
	productBatchRoutes.Delete("/:id", handler.DeleteProductBatch)        // DELETE /api/v1/product-batches/:id
}
