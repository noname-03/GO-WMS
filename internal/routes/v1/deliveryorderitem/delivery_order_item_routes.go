package deliveryorderitem

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterDeliveryOrderItemRoutes(api fiber.Router) {
	// Create a delivery order item group
	deliveryOrderItems := api.Group("/delivery-order-items")

	// All routes require JWT authentication
	deliveryOrderItems.Use(middleware.JWTMiddleware())

	// Public delivery order item routes (requires authentication)
	deliveryOrderItems.Get("/", handler.GetAllDeliveryOrderItems)
	deliveryOrderItems.Get("/deleted", handler.GetDeletedDeliveryOrderItems)
	deliveryOrderItems.Get("/:id", handler.GetDeliveryOrderItemByID)
	deliveryOrderItems.Get("/delivery-order/:deliveryOrderId", handler.GetDeliveryOrderItemsByDeliveryOrder)
	deliveryOrderItems.Post("/", handler.CreateDeliveryOrderItem)
	deliveryOrderItems.Put("/:id", handler.UpdateDeliveryOrderItem)
	deliveryOrderItems.Put("/:id/restore", handler.RestoreDeliveryOrderItem)
	deliveryOrderItems.Delete("/:id", handler.DeleteDeliveryOrderItem)
}
