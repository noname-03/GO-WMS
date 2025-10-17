package purchaseorder

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterPurchaseOrderRoutes(app fiber.Router) {
	purchaseOrderRoutes := app.Group("/purchase-orders")

	// All routes require authentication
	purchaseOrderRoutes.Use(middleware.JWTMiddleware())

	// CRUD operations for purchase orders
	purchaseOrderRoutes.Get("/", handler.GetPurchaseOrders)                       // GET /api/v1/purchase-orders
	purchaseOrderRoutes.Get("/filter", handler.GetFilteredPurchaseOrders)         // GET /api/v1/purchase-orders/filter?user_id=&order_date_from=&order_date_to=&status=
	purchaseOrderRoutes.Get("/deleted", handler.GetDeletedPurchaseOrders)         // GET /api/v1/purchase-orders/deleted
	purchaseOrderRoutes.Get("/:id", handler.GetPurchaseOrderByID)                 // GET /api/v1/purchase-orders/:id
	purchaseOrderRoutes.Post("/", handler.CreatePurchaseOrder)                    // POST /api/v1/purchase-orders
	purchaseOrderRoutes.Post("/with-items", handler.CreatePurchaseOrderWithItems) // POST /api/v1/purchase-orders/with-items
	purchaseOrderRoutes.Put("/:id", handler.UpdatePurchaseOrder)                  // PUT /api/v1/purchase-orders/:id
	purchaseOrderRoutes.Put("/:id/restore", handler.RestorePurchaseOrder)         // PUT /api/v1/purchase-orders/:id/restore
	purchaseOrderRoutes.Delete("/:id", handler.DeletePurchaseOrder)               // DELETE /api/v1/purchase-orders/:id
}
