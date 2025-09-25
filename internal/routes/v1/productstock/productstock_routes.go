package productstock

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProductStockRoutes(router fiber.Router) {
	stocks := router.Group("/product-stocks")
	stocks.Use(middleware.JWTMiddleware()) // All routes require authentication
	{
		// GET /api/v1/product-stocks - Get all product stocks
		stocks.Get("", handler.GetAllProductStocks)

		// GET /api/v1/product-stocks/:id - Get product stock by ID
		stocks.Get("/:id", handler.GetProductStockByID)

		// GET /api/v1/product-stocks/product/:productId - Get stocks by product ID
		stocks.Get("/product/:productId", handler.GetProductStocksByProduct)

		// POST /api/v1/product-stocks - Create new product stock
		stocks.Post("", handler.CreateProductStock)

		// PUT /api/v1/product-stocks/:id - Update product stock
		stocks.Put("/:id", handler.UpdateProductStock)

		// DELETE /api/v1/product-stocks/:id - Delete product stock
		stocks.Delete("/:id", handler.DeleteProductStock)
	}
}
