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
	invoices.Get("/:id", handler.GetInvoiceByID)
	invoices.Post("/", handler.CreateInvoice)
	invoices.Put("/:id", handler.UpdateInvoice)
	invoices.Put("/:id/restore", handler.RestoreInvoice)
	invoices.Delete("/:id", handler.DeleteInvoice)
}
