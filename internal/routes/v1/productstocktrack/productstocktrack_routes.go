package productstocktrack

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProductStockTrackRoutes(router fiber.Router) {
	tracks := router.Group("/product-stock-tracks")
	tracks.Use(middleware.JWTMiddleware()) // All routes require authentication
	{
		// GET /api/v1/product-stock-tracks - Get all stock tracks
		tracks.Get("", handler.GetAllProductStockTracks)

		// GET /api/v1/product-stock-tracks/:id - Get stock track by ID
		tracks.Get("/:id", handler.GetProductStockTrackByID)

		// GET /api/v1/product-stock-tracks/stock/:stockId - Get tracks by stock ID
		tracks.Get("/stock/:stockId", handler.GetProductStockTracksByStock)

		// GET /api/v1/product-stock-tracks/product/:productId - Get tracks by product ID
		tracks.Get("/product/:productId", handler.GetProductStockTracksByProduct)

		// GET /api/v1/product-stock-tracks/date-range?startDate=YYYY-MM-DD&endDate=YYYY-MM-DD - Get tracks by date range
		// tracks.Get("/date-range", handler.GetProductStockTracksByDateRange)

		// POST /api/v1/product-stock-tracks - Create new stock track
		tracks.Post("", handler.CreateProductStockTrack)

		// PUT /api/v1/product-stock-tracks/:id - Update stock track
		tracks.Put("/:id", handler.UpdateProductStockTrack)

		// DELETE /api/v1/product-stock-tracks/:id - Delete stock track
		tracks.Delete("/:id", handler.DeleteProductStockTrack)
	}
}
