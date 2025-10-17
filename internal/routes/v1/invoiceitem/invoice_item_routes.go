package invoiceitem

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterInvoiceItemRoutes(api fiber.Router) {
	// Create an invoice item group
	invoiceItems := api.Group("/invoice-items")

	// All routes require JWT authentication
	invoiceItems.Use(middleware.JWTMiddleware())

	// Public invoice item routes (requires authentication)
	invoiceItems.Get("/", handler.GetAllInvoiceItems)
	invoiceItems.Get("/deleted", handler.GetDeletedInvoiceItems)
	invoiceItems.Get("/:id", handler.GetInvoiceItemByID)
	invoiceItems.Get("/invoice/:invoiceId", handler.GetInvoiceItemsByInvoice)
	invoiceItems.Post("/", handler.CreateInvoiceItem)
	invoiceItems.Put("/:id", handler.UpdateInvoiceItem)
	invoiceItems.Put("/:id/restore", handler.RestoreInvoiceItem)
	invoiceItems.Delete("/:id", handler.DeleteInvoiceItem)
}
