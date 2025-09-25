package productitem

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProductItemRoutes(router fiber.Router) {
	items := router.Group("/product-items")
	items.Use(middleware.JWTMiddleware()) // All routes require authentication
	{
		// GET /api/v1/product-items - Get all product items
		items.Get("", handler.GetAllProductItems)

		// GET /api/v1/product-items/:id - Get product item by ID
		items.Get("/:id", handler.GetProductItemByID)

		// GET /api/v1/product-items/stock/:stockId - Get items by stock ID
		items.Get("/stock/:stockId", handler.GetProductItemsByStock)

		// GET /api/v1/product-items/product/:productId - Get items by product ID
		items.Get("/product/:productId", handler.GetProductItemsByProduct)

		// GET /api/v1/product-items/location/:locationId - Get items by location ID
		items.Get("/location/:locationId", handler.GetProductItemsByLocation)

		// GET /api/v1/product-items/summary/by-product - Get items summary grouped by product
		items.Get("/summary/by-product", handler.GetItemsSummaryByProduct)

		// POST /api/v1/product-items - Create new product item
		items.Post("", handler.CreateProductItem)

		// PUT /api/v1/product-items/:id - Update product item
		items.Put("/:id", handler.UpdateProductItem)

		// DELETE /api/v1/product-items/:id - Delete product item
		items.Delete("/:id", handler.DeleteProductItem)
	}
}
