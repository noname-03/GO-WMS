package purchaseorderitem

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterPurchaseOrderItemRoutes(app fiber.Router) {
	purchaseOrderItemRoutes := app.Group("/purchase-order-items")

	// All routes require authentication
	purchaseOrderItemRoutes.Use(middleware.JWTMiddleware())

	// CRUD operations for purchase order items
	purchaseOrderItemRoutes.Get("/", handler.GetAllPurchaseOrderItems)                                            // GET /api/v1/purchase-order-items
	purchaseOrderItemRoutes.Get("/deleted", handler.GetDeletedPurchaseOrderItems)                                 // GET /api/v1/purchase-order-items/deleted
	purchaseOrderItemRoutes.Get("/purchase-order/:purchaseOrderId", handler.GetPurchaseOrderItemsByPurchaseOrder) // GET /api/v1/purchase-order-items/purchase-order/:purchaseOrderId
	purchaseOrderItemRoutes.Get("/:id", handler.GetPurchaseOrderItemByID)                                         // GET /api/v1/purchase-order-items/:id
	purchaseOrderItemRoutes.Post("/", handler.CreatePurchaseOrderItem)                                            // POST /api/v1/purchase-order-items
	purchaseOrderItemRoutes.Put("/:id", handler.UpdatePurchaseOrderItem)                                          // PUT /api/v1/purchase-order-items/:id
	purchaseOrderItemRoutes.Put("/:id/restore", handler.RestorePurchaseOrderItem)                                 // PUT /api/v1/purchase-order-items/:id/restore
	purchaseOrderItemRoutes.Delete("/:id", handler.DeletePurchaseOrderItem)                                       // DELETE /api/v1/purchase-order-items/:id
}
