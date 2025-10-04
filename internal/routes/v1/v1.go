package v1

import (
	"myapp/internal/routes/v1/auth"
	"myapp/internal/routes/v1/brand"
	"myapp/internal/routes/v1/category"
	"myapp/internal/routes/v1/file"
	"myapp/internal/routes/v1/location"
	"myapp/internal/routes/v1/product"
	"myapp/internal/routes/v1/productbatch"
	"myapp/internal/routes/v1/productitem"
	"myapp/internal/routes/v1/productitemtrack"
	"myapp/internal/routes/v1/productstock"
	"myapp/internal/routes/v1/productstocktrack"
	"myapp/internal/routes/v1/productunit"
	"myapp/internal/routes/v1/productunittrack"
	"myapp/internal/routes/v1/role"
	"myapp/internal/routes/v1/s3test"
	"myapp/internal/routes/v1/user"

	"github.com/gofiber/fiber/v2"
	// Import modules lain di sini untuk future development
	// "myapp/internal/routes/v1/order"
	// "myapp/internal/routes/v1/warehouse"
)

func SetupV1Routes(app *fiber.App) {
	// Create v1 API group
	v1 := app.Group("/api/v1")

	// Setup module routes
	auth.SetupAuthRoutes(v1)
	user.SetupUserRoutes(v1)
	role.SetupRoleRoutes(v1)
	brand.SetupBrandRoutes(v1)
	category.SetupCategoryRoutes(v1)
	file.SetupFileRoutes(v1)
	product.RegisterProductRoutes(v1)
	productbatch.RegisterProductBatchRoutes(v1)
	productunit.ProductUnitRoutes(v1)
	productunittrack.ProductUnitTrackRoutes(v1)
	location.LocationRoutes(v1)
	productstock.ProductStockRoutes(v1)
	productstocktrack.ProductStockTrackRoutes(v1)
	productitem.ProductItemRoutes(v1)
	productitemtrack.ProductItemTrackRoutes(v1)

	// S3 test routes (for debugging)
	s3test.SetupS3TestRoutes(v1)

	// Future modules
	// order.SetupOrderRoutes(v1)
	// warehouse.SetupWarehouseRoutes(v1)

	// Health check endpoint
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "OK",
			"version": "v1",
			"message": "GO-WMS API v1 is running",
		})
	})
}
