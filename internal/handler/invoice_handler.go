package handler

import (
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var invoiceService = service.NewInvoiceService()

// handleInvoiceError converts database errors to user-friendly messages for invoice operations
func handleInvoiceError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors first
	if errMsg == "invoice not found" {
		return 404, "Invoice not found"
	}

	if errMsg == "user not found" {
		return 404, "User not found"
	}

	if errMsg == "purchase order not found" {
		return 404, "Purchase order not found"
	}

	if errMsg == "delivery order not found" {
		return 404, "Delivery order not found"
	}

	if errMsg == "invoice number already exists" {
		return 400, "Invoice number already exists"
	}

	// Handle PostgreSQL constraint errors as backup
	if strings.Contains(errMsg, "foreign key constraint") {
		return 400, "Invalid user, purchase order, or delivery order ID"
	}

	// Default to 500 for other errors
	return 500, "Internal server error"
}

// Request DTOs
type CreateInvoiceRequest struct {
	InvoiceNumber   string  `json:"invoiceNumber" validate:"required"`
	UserID          uint    `json:"userId" validate:"required"`
	PurchaseOrderID *uint   `json:"purchaseOrderId"`
	DeliveryOrderID *uint   `json:"deliveryOrderId"`
	InvoiceDate     string  `json:"invoiceDate" validate:"required"`
	Status          string  `json:"status"` // draft, sent, paid, closed
	TotalAmount     float64 `json:"totalAmount" validate:"required"`
	Description     *string `json:"description"`
}

type UpdateInvoiceRequest struct {
	InvoiceNumber   string  `json:"invoiceNumber"`
	UserID          uint    `json:"userId"`
	PurchaseOrderID *uint   `json:"purchaseOrderId"`
	DeliveryOrderID *uint   `json:"deliveryOrderId"`
	InvoiceDate     string  `json:"invoiceDate"`
	Status          string  `json:"status"`
	TotalAmount     float64 `json:"totalAmount"`
	Description     *string `json:"description"`
}

func GetInvoices(c *fiber.Ctx) error {
	log.Printf("[INVOICE] Get all invoices request from IP: %s", c.IP())

	invoices, err := invoiceService.GetAllInvoices()
	if err != nil {
		log.Printf("[INVOICE] Get all invoices failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch invoices", err.Error())
	}

	log.Printf("[INVOICE] Get all invoices successful")
	return helper.Success(c, 200, "Success", invoices)
}

func GetFilteredInvoices(c *fiber.Ctx) error {
	log.Printf("[INVOICE] Get filtered invoices request from IP: %s", c.IP())

	// Parse query parameters
	var userID *uint
	var purchaseOrderID *uint
	var deliveryOrderID *uint
	var invoiceDateFrom *time.Time
	var invoiceDateTo *time.Time
	var status *string

	// Parse user_id
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		parsedUserID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			log.Printf("[INVOICE] Invalid user_id parameter: %s - error: %v", userIDStr, err)
			return helper.Fail(c, 400, "Invalid user_id parameter", err.Error())
		}
		userIDUint := uint(parsedUserID)
		userID = &userIDUint
	}

	// Parse purchase_order_id
	if poIDStr := c.Query("purchase_order_id"); poIDStr != "" {
		parsedPOID, err := strconv.ParseUint(poIDStr, 10, 32)
		if err != nil {
			log.Printf("[INVOICE] Invalid purchase_order_id parameter: %s - error: %v", poIDStr, err)
			return helper.Fail(c, 400, "Invalid purchase_order_id parameter", err.Error())
		}
		poIDUint := uint(parsedPOID)
		purchaseOrderID = &poIDUint
	}

	// Parse delivery_order_id
	if doIDStr := c.Query("delivery_order_id"); doIDStr != "" {
		parsedDOID, err := strconv.ParseUint(doIDStr, 10, 32)
		if err != nil {
			log.Printf("[INVOICE] Invalid delivery_order_id parameter: %s - error: %v", doIDStr, err)
			return helper.Fail(c, 400, "Invalid delivery_order_id parameter", err.Error())
		}
		doIDUint := uint(parsedDOID)
		deliveryOrderID = &doIDUint
	}

	// Parse invoice_date_from
	if dateFromStr := c.Query("invoice_date_from"); dateFromStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateFromStr)
		if err != nil {
			log.Printf("[INVOICE] Invalid invoice_date_from parameter: %s - error: %v", dateFromStr, err)
			return helper.Fail(c, 400, "Invalid invoice_date_from parameter, expected format: YYYY-MM-DD", err.Error())
		}
		invoiceDateFrom = &parsedDate
	}

	// Parse invoice_date_to
	if dateToStr := c.Query("invoice_date_to"); dateToStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateToStr)
		if err != nil {
			log.Printf("[INVOICE] Invalid invoice_date_to parameter: %s - error: %v", dateToStr, err)
			return helper.Fail(c, 400, "Invalid invoice_date_to parameter, expected format: YYYY-MM-DD", err.Error())
		}
		invoiceDateTo = &parsedDate
	}

	// Parse status
	if statusStr := c.Query("status"); statusStr != "" {
		status = &statusStr
	}

	log.Printf("[INVOICE] Filter parameters - user_id: %v, purchase_order_id: %v, delivery_order_id: %v, invoice_date_from: %v, invoice_date_to: %v, status: %v",
		userID, purchaseOrderID, deliveryOrderID, invoiceDateFrom, invoiceDateTo, status)

	invoices, err := invoiceService.GetFilteredInvoices(userID, purchaseOrderID, deliveryOrderID, invoiceDateFrom, invoiceDateTo, status)
	if err != nil {
		log.Printf("[INVOICE] Get filtered invoices failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch filtered invoices", err.Error())
	}

	log.Printf("[INVOICE] Get filtered invoices successful")
	return helper.Success(c, 200, "Success", invoices)
}

func GetInvoiceByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INVOICE] Get invoice by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[INVOICE] Get invoice by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid invoice ID", err.Error())
	}

	invoice, err := invoiceService.GetInvoiceByID(uint(idUint))
	if err != nil {
		log.Printf("[INVOICE] Get invoice by ID failed - Invoice ID: %d not found, error: %v", idUint, err)
		return helper.Fail(c, 404, "Invoice not found", err.Error())
	}

	log.Printf("[INVOICE] Get invoice by ID successful - Invoice ID: %d", idUint)
	return helper.Success(c, 200, "Success", invoice)
}

func CreateInvoice(c *fiber.Ctx) error {
	log.Printf("[INVOICE] Create invoice request from IP: %s", c.IP())

	var req CreateInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[INVOICE] Create invoice failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Parse invoice date
	invoiceDate, err := time.Parse("2006-01-02", req.InvoiceDate)
	if err != nil {
		log.Printf("[INVOICE] Create invoice failed - Invalid invoice_date format: %s, error: %v", req.InvoiceDate, err)
		return helper.Fail(c, 400, "Invalid invoice_date format, use YYYY-MM-DD", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[INVOICE] Create invoice failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[INVOICE] Creating invoice with audit - User ID: %d, Invoice Number: %s", userID, req.InvoiceNumber)

	invoice, err := invoiceService.CreateInvoice(
		req.InvoiceNumber,
		req.UserID,
		req.PurchaseOrderID,
		req.DeliveryOrderID,
		invoiceDate,
		req.Status,
		req.TotalAmount,
		req.Description,
		userID,
	)

	if err != nil {
		log.Printf("[INVOICE] Create invoice failed - Invoice Number: %s, User ID: %d, error: %v", req.InvoiceNumber, userID, err)
		statusCode, message := handleInvoiceError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[INVOICE] Create invoice successful - Invoice Number: %s, Created by User ID: %d", req.InvoiceNumber, userID)
	return helper.Success(c, 201, "Invoice created successfully", invoice)
}

func UpdateInvoice(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INVOICE] Update invoice request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[INVOICE] Update invoice failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid invoice ID", err.Error())
	}

	var req UpdateInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[INVOICE] Update invoice failed - Invalid request body for ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Parse invoice date if provided
	var invoiceDate time.Time
	if req.InvoiceDate != "" {
		invoiceDate, err = time.Parse("2006-01-02", req.InvoiceDate)
		if err != nil {
			log.Printf("[INVOICE] Update invoice failed - Invalid invoice_date format: %s, error: %v", req.InvoiceDate, err)
			return helper.Fail(c, 400, "Invalid invoice_date format, use YYYY-MM-DD", err.Error())
		}
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[INVOICE] Update invoice failed - User not authenticated for Invoice ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[INVOICE] Updating invoice with audit - Invoice ID: %d, User ID: %d", idUint, userID)

	invoice, err := invoiceService.UpdateInvoice(
		uint(idUint),
		req.InvoiceNumber,
		req.UserID,
		req.PurchaseOrderID,
		req.DeliveryOrderID,
		invoiceDate,
		req.Status,
		req.TotalAmount,
		req.Description,
		userID,
	)

	if err != nil {
		log.Printf("[INVOICE] Update invoice failed - Invoice ID: %d, User ID: %d, error: %v", idUint, userID, err)
		statusCode, message := handleInvoiceError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[INVOICE] Update invoice successful - Invoice ID: %d, Updated by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Invoice updated successfully", invoice)
}

func DeleteInvoice(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INVOICE] Delete invoice request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[INVOICE] Delete invoice failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid invoice ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[INVOICE] Delete invoice failed - User not authenticated for Invoice ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = invoiceService.DeleteInvoice(uint(idUint), userID)
	if err != nil {
		log.Printf("[INVOICE] Delete invoice failed - Invoice ID: %d, error: %v", idUint, err)
		statusCode, message := handleInvoiceError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[INVOICE] Delete invoice successful - Invoice ID: %d, Deleted by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Invoice deleted successfully", nil)
}

func GetDeletedInvoices(c *fiber.Ctx) error {
	log.Printf("[INVOICE] Get deleted invoices request from IP: %s", c.IP())

	invoices, err := invoiceService.GetDeletedInvoices()
	if err != nil {
		log.Printf("[INVOICE] Get deleted invoices failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch deleted invoices", err.Error())
	}

	log.Printf("[INVOICE] Get deleted invoices successful")
	return helper.Success(c, 200, "Success", invoices)
}

func RestoreInvoice(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INVOICE] Restore invoice request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[INVOICE] Restore invoice failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid invoice ID", err.Error())
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[INVOICE] Restore invoice failed - User not authenticated for Invoice ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	invoice, err := invoiceService.RestoreInvoice(uint(idUint), userID)
	if err != nil {
		log.Printf("[INVOICE] Restore invoice failed - Invoice ID: %d, error: %v", idUint, err)
		statusCode, message := handleInvoiceError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[INVOICE] Restore invoice successful - Invoice ID: %d, Restored by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Invoice restored successfully", invoice)
}
