package productunittrack

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProductUnitTrackRoutes(router fiber.Router) {
	productUnitTrack := router.Group("/product-unit-tracks")

	// Apply JWT middleware
	productUnitTrack.Use(middleware.JWTMiddleware())

	// GET /api/v1/product-unit-tracks - Get all product unit tracks
	productUnitTrack.Get("/", handler.GetProductUnitTracks)

	// GET /api/v1/product-unit-tracks/product/:productId - Get product unit tracks by product ID
	productUnitTrack.Get("/product/:productId", handler.GetProductUnitTracksByProduct)

	// GET /api/v1/product-unit-tracks/product-unit/:productUnitId - Get product unit tracks by product unit ID
	productUnitTrack.Get("/product-unit/:productUnitId", handler.GetProductUnitTracksByProductUnit)

	// GET /api/v1/product-unit-tracks/:id - Get product unit track by ID
	productUnitTrack.Get("/:id", handler.GetProductUnitTrackByID)

	// POST /api/v1/product-unit-tracks - Create new product unit track
	productUnitTrack.Post("/", handler.CreateProductUnitTrack)

	// DELETE /api/v1/product-unit-tracks/:id - Delete product unit track by ID
	productUnitTrack.Delete("/:id", handler.DeleteProductUnitTrack)
}
