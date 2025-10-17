package handler

import (
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var invoiceItemService = service.NewInvoiceItemService()

// handleInvoiceItemError converts database errors to user-friendly messages for invoice item operations
func handleInvoiceItemError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors first
	if errMsg == "invoice item not found" {
		return 404, "Invoice item not found"
	}

	if errMsg == "invoice not found" {
		return 404, "Invoice not found"
	}

	if errMsg == "product not found" {
		return 404, "Product not found"
	}

	// Handle PostgreSQL constraint errors as backup
	if strings.Contains(errMsg, "foreign key constraint") {
		return 400, "Invalid invoice or product ID"
	}

	// Default to 500 for other errors
	return 500, "Internal server error"
}

// Request DTOs
type CreateInvoiceItemRequest struct {
	InvoiceID   uint    `json:"invoiceId" validate:"required"`
	ProductID   uint    `json:"productId" validate:"required"`
	QtyInvoiced float64 `json:"qtyInvoiced" validate:"required"`
	UnitPrice   float64 `json:"unitPrice" validate:"required"`
	TotalPrice  float64 `json:"totalPrice" validate:"required"`
	Description *string `json:"description"`
}

type UpdateInvoiceItemRequest struct {
	InvoiceID   uint    `json:"invoiceId"`
	ProductID   uint    `json:"productId"`
	QtyInvoiced float64 `json:"qtyInvoiced"`
	UnitPrice   float64 `json:"unitPrice"`
	TotalPrice  float64 `json:"totalPrice"`
	Description *string `json:"description"`
}

func GetAllInvoiceItems(c *fiber.Ctx) error {
	log.Printf("[INVOICE_ITEM] Get all invoice items request from IP: %s", c.IP())

	items, err := invoiceItemService.GetAllInvoiceItems()
	if err != nil {
		log.Printf("[INVOICE_ITEM] Get all invoice items failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch invoice items", err.Error())
	}

	log.Printf("[INVOICE_ITEM] Get all invoice items successful")
	return helper.Success(c, 200, "Success", items)
}

func GetInvoiceItemsByInvoice(c *fiber.Ctx) error {
	invoiceID := c.Params("invoiceId")
	log.Printf("[INVOICE_ITEM] Get invoice items by invoice request - Invoice ID: %s from IP: %s", invoiceID, c.IP())

	invID, err := strconv.ParseUint(invoiceID, 10, 32)
	if err != nil {
		log.Printf("[INVOICE_ITEM] Get invoice items by invoice failed - Invalid Invoice ID: %s, error: %v", invoiceID, err)
		return helper.Fail(c, 400, "Invalid invoice ID", err.Error())
	}

	items, err := invoiceItemService.GetInvoiceItemsByInvoice(uint(invID))
	if err != nil {
		log.Printf("[INVOICE_ITEM] Get invoice items by invoice failed - Invoice ID: %d, error: %v", invID, err)
		statusCode, message := handleInvoiceItemError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[INVOICE_ITEM] Get invoice items by invoice successful - Invoice ID: %d", invID)
	return helper.Success(c, 200, "Success", items)
}

func GetInvoiceItemByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INVOICE_ITEM] Get invoice item by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[INVOICE_ITEM] Get invoice item by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid invoice item ID", err.Error())
	}

	item, err := invoiceItemService.GetInvoiceItemByID(uint(idUint))
	if err != nil {
		log.Printf("[INVOICE_ITEM] Get invoice item by ID failed - Item ID: %d not found, error: %v", idUint, err)
		return helper.Fail(c, 404, "Invoice item not found", err.Error())
	}

	log.Printf("[INVOICE_ITEM] Get invoice item by ID successful - Item ID: %d", idUint)
	return helper.Success(c, 200, "Success", item)
}

func CreateInvoiceItem(c *fiber.Ctx) error {
	log.Printf("[INVOICE_ITEM] Create invoice item request from IP: %s", c.IP())

	var req CreateInvoiceItemRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[INVOICE_ITEM] Create invoice item failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[INVOICE_ITEM] Create invoice item failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[INVOICE_ITEM] Creating invoice item with audit - User ID: %d, Invoice ID: %d", userID, req.InvoiceID)

	item, err := invoiceItemService.CreateInvoiceItem(
		req.InvoiceID,
		req.ProductID,
		req.QtyInvoiced,
		req.UnitPrice,
		req.TotalPrice,
		req.Description,
		userID,
	)

	if err != nil {
		log.Printf("[INVOICE_ITEM] Create invoice item failed - Invoice ID: %d, User ID: %d, error: %v", req.InvoiceID, userID, err)
		statusCode, message := handleInvoiceItemError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[INVOICE_ITEM] Create invoice item successful - Invoice ID: %d, Created by User ID: %d", req.InvoiceID, userID)
	return helper.Success(c, 201, "Invoice item created successfully", item)
}

func UpdateInvoiceItem(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INVOICE_ITEM] Update invoice item request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[INVOICE_ITEM] Update invoice item failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid invoice item ID", err.Error())
	}

	var req UpdateInvoiceItemRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[INVOICE_ITEM] Update invoice item failed - Invalid request body for ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[INVOICE_ITEM] Update invoice item failed - User not authenticated for Item ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[INVOICE_ITEM] Updating invoice item with audit - Item ID: %d, User ID: %d", idUint, userID)

	item, err := invoiceItemService.UpdateInvoiceItem(
		uint(idUint),
		req.InvoiceID,
		req.ProductID,
		req.QtyInvoiced,
		req.UnitPrice,
		req.TotalPrice,
		req.Description,
		userID,
	)

	if err != nil {
		log.Printf("[INVOICE_ITEM] Update invoice item failed - Item ID: %d, User ID: %d, error: %v", idUint, userID, err)
		statusCode, message := handleInvoiceItemError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[INVOICE_ITEM] Update invoice item successful - Item ID: %d, Updated by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Invoice item updated successfully", item)
}

func DeleteInvoiceItem(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INVOICE_ITEM] Delete invoice item request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[INVOICE_ITEM] Delete invoice item failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid invoice item ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[INVOICE_ITEM] Delete invoice item failed - User not authenticated for Item ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = invoiceItemService.DeleteInvoiceItem(uint(idUint), userID)
	if err != nil {
		log.Printf("[INVOICE_ITEM] Delete invoice item failed - Item ID: %d, error: %v", idUint, err)
		statusCode, message := handleInvoiceItemError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[INVOICE_ITEM] Delete invoice item successful - Item ID: %d, Deleted by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Invoice item deleted successfully", nil)
}

func GetDeletedInvoiceItems(c *fiber.Ctx) error {
	log.Printf("[INVOICE_ITEM] Get deleted invoice items request from IP: %s", c.IP())

	items, err := invoiceItemService.GetDeletedInvoiceItems()
	if err != nil {
		log.Printf("[INVOICE_ITEM] Get deleted invoice items failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch deleted invoice items", err.Error())
	}

	log.Printf("[INVOICE_ITEM] Get deleted invoice items successful")
	return helper.Success(c, 200, "Success", items)
}

func RestoreInvoiceItem(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INVOICE_ITEM] Restore invoice item request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[INVOICE_ITEM] Restore invoice item failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid invoice item ID", err.Error())
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[INVOICE_ITEM] Restore invoice item failed - User not authenticated for Item ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	item, err := invoiceItemService.RestoreInvoiceItem(uint(idUint), userID)
	if err != nil {
		log.Printf("[INVOICE_ITEM] Restore invoice item failed - Item ID: %d, error: %v", idUint, err)
		statusCode, message := handleInvoiceItemError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[INVOICE_ITEM] Restore invoice item successful - Item ID: %d, Restored by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Invoice item restored successfully", item)
}
