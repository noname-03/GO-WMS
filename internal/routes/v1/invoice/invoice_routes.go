package invoice

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterInvoiceRoutes(api fiber.Router) {
	// Create an invoice group
	invoices := api.Group("/invoices")

	// All routes require JWT authentication
	invoices.Use(middleware.JWTMiddleware())

	// Public invoice routes (requires authentication)
	invoices.Get("/", handler.GetInvoices)
	invoices.Get("/deleted", handler.GetDeletedInvoices)
	invoices.Get("/filter", handler.GetFilteredInvoices) // Query params: user_id, purchase_order_id, delivery_order_id, invoice_date_from, invoice_date_to, status
	invoices.Get("/:id", handler.GetInvoiceByID)
	invoices.Post("/", handler.CreateInvoice)
	invoices.Put("/:id", handler.UpdateInvoice)
	invoices.Put("/:id/restore", handler.RestoreInvoice)
	invoices.Delete("/:id", handler.DeleteInvoice)
}
