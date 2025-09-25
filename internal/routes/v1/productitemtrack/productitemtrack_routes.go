package productitemtrack

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProductItemTrackRoutes(router fiber.Router) {
	tracks := router.Group("/product-item-tracks")
	tracks.Use(middleware.JWTMiddleware()) // All routes require authentication
	{
		// GET /api/v1/product-item-tracks - Get all item tracks
		tracks.Get("", handler.GetAllProductItemTracks)

		// GET /api/v1/product-item-tracks/:id - Get item track by ID
		tracks.Get("/:id", handler.GetProductItemTrackByID)

		// GET /api/v1/product-item-tracks/item/:itemId - Get tracks by item ID
		tracks.Get("/item/:itemId", handler.GetProductItemTracksByItem)

		// GET /api/v1/product-item-tracks/stock/:stockId - Get tracks by stock ID
		tracks.Get("/stock/:stockId", handler.GetProductItemTracksByStock)

		// GET /api/v1/product-item-tracks/product/:productId - Get tracks by product ID
		tracks.Get("/product/:productId", handler.GetProductItemTracksByProduct)

		// GET /api/v1/product-item-tracks/date-range?startDate=YYYY-MM-DD&endDate=YYYY-MM-DD - Get tracks by date range
		tracks.Get("/date-range", handler.GetProductItemTracksByDateRange)

		// GET /api/v1/product-item-tracks/operation/:operation - Get tracks by operation (In/Out/Plus/Minus)
		tracks.Get("/operation/:operation", handler.GetTracksByOperation)

		// GET /api/v1/product-item-tracks/reports/value-by-product - Get value report grouped by product
		tracks.Get("/reports/value-by-product", handler.GetValueReportByProduct)

		// POST /api/v1/product-item-tracks - Create new item track
		tracks.Post("", handler.CreateProductItemTrack)

		// PUT /api/v1/product-item-tracks/:id - Update item track
		tracks.Put("/:id", handler.UpdateProductItemTrack)

		// DELETE /api/v1/product-item-tracks/:id - Delete item track
		tracks.Delete("/:id", handler.DeleteProductItemTrack)
	}
}
