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

var purchaseOrderService = service.NewPurchaseOrderService()

// handlePurchaseOrderError converts database errors to user-friendly messages for purchase order operations
func handlePurchaseOrderError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors first
	if errMsg == "purchase order not found" {
		return 404, "Purchase order not found"
	}

	if errMsg == "user not found" {
		return 404, "User not found"
	}

	if errMsg == "PO number already exists" {
		return 400, "PO number already exists"
	}

	// Handle PostgreSQL constraint errors as backup
	if strings.Contains(errMsg, "foreign key constraint") {
		return 400, "Invalid user ID"
	}

	// Default to 500 for other errors
	return 500, "Internal server error"
}

type CreatePurchaseOrderRequest struct {
	PONumber    string   `json:"poNumber" validate:"required"`
	UserID      uint     `json:"userId" validate:"required"`
	OrderDate   string   `json:"orderDate" validate:"required"`
	Status      string   `json:"status"` // draft, submitted, approved, received, closed
	TotalAmount *float64 `json:"totalAmount"`
	Description *string  `json:"description"`
}

type UpdatePurchaseOrderRequest struct {
	PONumber    string   `json:"poNumber"`
	UserID      uint     `json:"userId"`
	OrderDate   string   `json:"orderDate"`
	Status      string   `json:"status"`
	TotalAmount *float64 `json:"totalAmount"`
	Description *string  `json:"description"`
}

func GetPurchaseOrders(c *fiber.Ctx) error {
	log.Printf("[PURCHASE_ORDER] Get all purchase orders request from IP: %s", c.IP())

	orders, err := purchaseOrderService.GetAllPurchaseOrders()
	if err != nil {
		log.Printf("[PURCHASE_ORDER] Get all purchase orders failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch purchase orders", err.Error())
	}

	log.Printf("[PURCHASE_ORDER] Get all purchase orders successful")
	return helper.Success(c, 200, "Success", orders)
}

// GetFilteredPurchaseOrders returns purchase orders with optional query parameters
// Query params: user_id, order_date_from, order_date_to, status
func GetFilteredPurchaseOrders(c *fiber.Ctx) error {
	log.Printf("[PURCHASE_ORDER] Get filtered purchase orders request from IP: %s", c.IP())

	// Parse query parameters
	var userID *uint
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		id, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			log.Printf("[PURCHASE_ORDER] Invalid user_id parameter: %s", userIDStr)
			return helper.Fail(c, 400, "Invalid user_id parameter", err.Error())
		}
		uid := uint(id)
		userID = &uid
	}

	var orderDateFrom *time.Time
	if dateStr := c.Query("order_date_from"); dateStr != "" {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("[PURCHASE_ORDER] Invalid order_date_from parameter: %s", dateStr)
			return helper.Fail(c, 400, "Invalid order_date_from format, use YYYY-MM-DD", err.Error())
		}
		orderDateFrom = &date
	}

	var orderDateTo *time.Time
	if dateStr := c.Query("order_date_to"); dateStr != "" {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("[PURCHASE_ORDER] Invalid order_date_to parameter: %s", dateStr)
			return helper.Fail(c, 400, "Invalid order_date_to format, use YYYY-MM-DD", err.Error())
		}
		orderDateTo = &date
	}

	var status *string
	if statusStr := c.Query("status"); statusStr != "" {
		status = &statusStr
	}

	log.Printf("[PURCHASE_ORDER] Filtering with - UserID: %v, DateFrom: %v, DateTo: %v, Status: %v", userID, orderDateFrom, orderDateTo, status)

	orders, err := purchaseOrderService.GetFilteredPurchaseOrders(userID, orderDateFrom, orderDateTo, status)
	if err != nil {
		log.Printf("[PURCHASE_ORDER] Get filtered purchase orders failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch filtered purchase orders", err.Error())
	}

	log.Printf("[PURCHASE_ORDER] Get filtered purchase orders successful")
	return helper.Success(c, 200, "Success", orders)
}

func GetPurchaseOrderByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PURCHASE_ORDER] Get purchase order by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PURCHASE_ORDER] Get purchase order by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid purchase order ID", err.Error())
	}

	order, err := purchaseOrderService.GetPurchaseOrderByID(uint(idUint))
	if err != nil {
		log.Printf("[PURCHASE_ORDER] Get purchase order by ID failed - Order ID: %d not found, error: %v", idUint, err)
		return helper.Fail(c, 404, "Purchase order not found", err.Error())
	}

	log.Printf("[PURCHASE_ORDER] Get purchase order by ID successful - Order ID: %d", idUint)
	return helper.Success(c, 200, "Success", order)
}

func CreatePurchaseOrder(c *fiber.Ctx) error {
	log.Printf("[PURCHASE_ORDER] Create purchase order request from IP: %s", c.IP())

	var req CreatePurchaseOrderRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PURCHASE_ORDER] Create purchase order failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Parse order date
	orderDate, err := time.Parse("2006-01-02", req.OrderDate)
	if err != nil {
		log.Printf("[PURCHASE_ORDER] Create purchase order failed - Invalid order_date format: %s, error: %v", req.OrderDate, err)
		return helper.Fail(c, 400, "Invalid order_date format, use YYYY-MM-DD", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PURCHASE_ORDER] Create purchase order failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[PURCHASE_ORDER] Creating purchase order with audit - User ID: %d, PO Number: %s", userID, req.PONumber)

	order, err := purchaseOrderService.CreatePurchaseOrder(req.PONumber, req.UserID, orderDate, req.Status, req.TotalAmount, req.Description, userID)
	if err != nil {
		log.Printf("[PURCHASE_ORDER] Create purchase order failed - PO Number: %s, User ID: %d, error: %v", req.PONumber, userID, err)
		statusCode, message := handlePurchaseOrderError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PURCHASE_ORDER] Create purchase order successful - PO Number: %s, Created by User ID: %d", req.PONumber, userID)
	return helper.Success(c, 201, "Purchase order created successfully", order)
}

func UpdatePurchaseOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PURCHASE_ORDER] Update purchase order request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PURCHASE_ORDER] Update purchase order failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid purchase order ID", err.Error())
	}

	var req UpdatePurchaseOrderRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PURCHASE_ORDER] Update purchase order failed - Invalid request body for ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Parse order date if provided
	var orderDate time.Time
	if req.OrderDate != "" {
		orderDate, err = time.Parse("2006-01-02", req.OrderDate)
		if err != nil {
			log.Printf("[PURCHASE_ORDER] Update purchase order failed - Invalid order_date format: %s, error: %v", req.OrderDate, err)
			return helper.Fail(c, 400, "Invalid order_date format, use YYYY-MM-DD", err.Error())
		}
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PURCHASE_ORDER] Update purchase order failed - User not authenticated for Order ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[PURCHASE_ORDER] Updating purchase order with audit - Order ID: %d, User ID: %d", idUint, userID)

	order, err := purchaseOrderService.UpdatePurchaseOrder(uint(idUint), req.PONumber, req.UserID, orderDate, req.Status, req.TotalAmount, req.Description, userID)
	if err != nil {
		log.Printf("[PURCHASE_ORDER] Update purchase order failed - Order ID: %d, User ID: %d, error: %v", idUint, userID, err)
		statusCode, message := handlePurchaseOrderError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PURCHASE_ORDER] Update purchase order successful - Order ID: %d, Updated by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Purchase order updated successfully", order)
}

func DeletePurchaseOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PURCHASE_ORDER] Delete purchase order request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PURCHASE_ORDER] Delete purchase order failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid purchase order ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PURCHASE_ORDER] Delete purchase order failed - User not authenticated for Order ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = purchaseOrderService.DeletePurchaseOrder(uint(idUint), userID)
	if err != nil {
		log.Printf("[PURCHASE_ORDER] Delete purchase order failed - Order ID: %d, error: %v", idUint, err)
		statusCode, message := handlePurchaseOrderError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PURCHASE_ORDER] Delete purchase order successful - Order ID: %d, Deleted by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Purchase order deleted successfully", nil)
}

func GetDeletedPurchaseOrders(c *fiber.Ctx) error {
	log.Printf("[PURCHASE_ORDER] Get deleted purchase orders request from IP: %s", c.IP())

	orders, err := purchaseOrderService.GetDeletedPurchaseOrders()
	if err != nil {
		log.Printf("[PURCHASE_ORDER] Get deleted purchase orders failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch deleted purchase orders", err.Error())
	}

	log.Printf("[PURCHASE_ORDER] Get deleted purchase orders successful")
	return helper.Success(c, 200, "Success", orders)
}

func RestorePurchaseOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PURCHASE_ORDER] Restore purchase order request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PURCHASE_ORDER] Restore purchase order failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid purchase order ID", err.Error())
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PURCHASE_ORDER] Restore purchase order failed - User not authenticated for Order ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	order, err := purchaseOrderService.RestorePurchaseOrder(uint(idUint), userID)
	if err != nil {
		log.Printf("[PURCHASE_ORDER] Restore purchase order failed - Order ID: %d, error: %v", idUint, err)
		statusCode, message := handlePurchaseOrderError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PURCHASE_ORDER] Restore purchase order successful - Order ID: %d, Restored by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Purchase order restored successfully", order)
}
