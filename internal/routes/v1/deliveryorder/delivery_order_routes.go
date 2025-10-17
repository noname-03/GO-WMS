package deliveryorder

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterDeliveryOrderRoutes(api fiber.Router) {
	// Create a delivery order group
	deliveryOrders := api.Group("/delivery-orders")

	// All routes require JWT authentication
	deliveryOrders.Use(middleware.JWTMiddleware())

	// Public delivery order routes (requires authentication)
	deliveryOrders.Get("/", handler.GetDeliveryOrders)
	deliveryOrders.Get("/filter", handler.GetFilteredDeliveryOrders)
	deliveryOrders.Get("/deleted", handler.GetDeletedDeliveryOrders)
	deliveryOrders.Get("/:id", handler.GetDeliveryOrderByID)
	deliveryOrders.Post("/", handler.CreateDeliveryOrder)
	deliveryOrders.Post("/with-items", handler.CreateDeliveryOrderWithItems)
	deliveryOrders.Put("/:id", handler.UpdateDeliveryOrder)
	deliveryOrders.Put("/:id/restore", handler.RestoreDeliveryOrder)
	deliveryOrders.Delete("/:id", handler.DeleteDeliveryOrder)
}
