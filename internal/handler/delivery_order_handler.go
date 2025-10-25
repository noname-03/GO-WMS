package handler

import (
	"log"
	"myapp/internal/model"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var deliveryOrderService = service.NewDeliveryOrderService()

// handleDeliveryOrderError converts database errors to user-friendly messages for delivery order operations
func handleDeliveryOrderError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors first
	if errMsg == "delivery order not found" {
		return 404, "Delivery order not found"
	}

	if errMsg == "purchase order not found" {
		return 404, "Purchase order not found"
	}

	if errMsg == "DO number already exists" {
		return 400, "DO number already exists"
	}

	// Handle PostgreSQL constraint errors as backup
	if strings.Contains(errMsg, "foreign key constraint") {
		return 400, "Invalid purchase order ID"
	}

	// Default to 500 for other errors
	return 500, "Internal server error"
}

// Request DTOs
type CreateDeliveryOrderRequest struct {
	DONumber        string  `json:"doNumber"`
	PurchaseOrderID uint    `json:"purchaseOrderId"`
	DeliveryDate    string  `json:"deliveryDate"` // YYYY-MM-DD format
	Status          string  `json:"status"`
	Description     *string `json:"description"`
}

type UpdateDeliveryOrderRequest struct {
	DONumber        string  `json:"doNumber"`
	PurchaseOrderID uint    `json:"purchaseOrderId"`
	DeliveryDate    string  `json:"deliveryDate"` // YYYY-MM-DD format
	Status          string  `json:"status"`
	Description     *string `json:"description"`
}

// UpdateDeliveryOrderWithItemsRequest is the request body for updating a delivery order with items
type UpdateDeliveryOrderWithItemsRequest struct {
	DONumber        string                     `json:"doNumber"`
	PurchaseOrderID uint                       `json:"purchaseOrderId"`
	DeliveryDate    string                     `json:"deliveryDate"` // YYYY-MM-DD format
	Status          string                     `json:"status"`
	Description     *string                    `json:"description"`
	Items           []DeliveryOrderItemRequest `json:"items" validate:"required"`
}

// CreateDeliveryOrderWithItemsRequest is the request body for creating a delivery order with items
type CreateDeliveryOrderWithItemsRequest struct {
	DONumber        string                     `json:"doNumber" validate:"required"`
	PurchaseOrderID uint                       `json:"purchaseOrderId" validate:"required"`
	DeliveryDate    string                     `json:"deliveryDate" validate:"required"` // YYYY-MM-DD format
	Status          string                     `json:"status"`
	Description     *string                    `json:"description"`
	Items           []DeliveryOrderItemRequest `json:"items" validate:"required"`
}

type DeliveryOrderItemRequest struct {
	ProductID    uint     `json:"productId" validate:"required"`
	QtyDelivered float64  `json:"qtyDelivered" validate:"required"`
	ReceivedQty  *float64 `json:"receivedQty"`
	DamagedQty   *float64 `json:"damagedQty"`
	Description  *string  `json:"description"`
}

func GetDeliveryOrders(c *fiber.Ctx) error {
	log.Printf("[DELIVERY_ORDER] Get all delivery orders request from IP: %s", c.IP())

	orders, err := deliveryOrderService.GetAllDeliveryOrders()
	if err != nil {
		log.Printf("[DELIVERY_ORDER] Get all delivery orders failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch delivery orders", err.Error())
	}

	log.Printf("[DELIVERY_ORDER] Get all delivery orders successful")
	return helper.Success(c, 200, "Success", orders)
}

// GetFilteredDeliveryOrders returns delivery orders with optional query parameters
// Query params: purchase_order_id, delivery_date_from, delivery_date_to, status
func GetFilteredDeliveryOrders(c *fiber.Ctx) error {
	log.Printf("[DELIVERY_ORDER] Get filtered delivery orders request from IP: %s", c.IP())

	// Parse query parameters
	var purchaseOrderID *uint
	if poIDStr := c.Query("purchase_order_id"); poIDStr != "" {
		id, err := strconv.ParseUint(poIDStr, 10, 32)
		if err != nil {
			log.Printf("[DELIVERY_ORDER] Invalid purchase_order_id parameter: %s", poIDStr)
			return helper.Fail(c, 400, "Invalid purchase_order_id parameter", err.Error())
		}
		poid := uint(id)
		purchaseOrderID = &poid
	}

	var deliveryDateFrom *time.Time
	if dateStr := c.Query("delivery_date_from"); dateStr != "" {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("[DELIVERY_ORDER] Invalid delivery_date_from parameter: %s", dateStr)
			return helper.Fail(c, 400, "Invalid delivery_date_from format, use YYYY-MM-DD", err.Error())
		}
		deliveryDateFrom = &date
	}

	var deliveryDateTo *time.Time
	if dateStr := c.Query("delivery_date_to"); dateStr != "" {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("[DELIVERY_ORDER] Invalid delivery_date_to parameter: %s", dateStr)
			return helper.Fail(c, 400, "Invalid delivery_date_to format, use YYYY-MM-DD", err.Error())
		}
		deliveryDateTo = &date
	}

	var status *string
	if statusStr := c.Query("status"); statusStr != "" {
		status = &statusStr
	}

	log.Printf("[DELIVERY_ORDER] Filtering with - PurchaseOrderID: %v, DateFrom: %v, DateTo: %v, Status: %v", purchaseOrderID, deliveryDateFrom, deliveryDateTo, status)

	orders, err := deliveryOrderService.GetFilteredDeliveryOrders(purchaseOrderID, deliveryDateFrom, deliveryDateTo, status)
	if err != nil {
		log.Printf("[DELIVERY_ORDER] Get filtered delivery orders failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch filtered delivery orders", err.Error())
	}

	log.Printf("[DELIVERY_ORDER] Get filtered delivery orders successful")
	return helper.Success(c, 200, "Success", orders)
}

func GetDeliveryOrderByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[DELIVERY_ORDER] Get delivery order by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[DELIVERY_ORDER] Get delivery order by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid delivery order ID", err.Error())
	}

	order, err := deliveryOrderService.GetDeliveryOrderByID(uint(idUint))
	if err != nil {
		log.Printf("[DELIVERY_ORDER] Get delivery order by ID failed - Order ID: %d not found, error: %v", idUint, err)
		return helper.Fail(c, 404, "Delivery order not found", err.Error())
	}

	log.Printf("[DELIVERY_ORDER] Get delivery order by ID successful - Order ID: %d", idUint)
	return helper.Success(c, 200, "Success", order)
}

// GetDeliveryOrderWithItems returns delivery order with all its items
func GetDeliveryOrderWithItems(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[DELIVERY_ORDER] Get delivery order with items request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[DELIVERY_ORDER] Get delivery order with items failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid delivery order ID", err.Error())
	}

	order, err := deliveryOrderService.GetDeliveryOrderWithItems(uint(idUint))
	if err != nil {
		log.Printf("[DELIVERY_ORDER] Get delivery order with items failed - Order ID: %d not found, error: %v", idUint, err)
		return helper.Fail(c, 404, "Delivery order not found", err.Error())
	}

	log.Printf("[DELIVERY_ORDER] Get delivery order with items successful - Order ID: %d", idUint)
	return helper.Success(c, 200, "Success", order)
}

func CreateDeliveryOrder(c *fiber.Ctx) error {
	log.Printf("[DELIVERY_ORDER] Create delivery order request from IP: %s", c.IP())

	var req CreateDeliveryOrderRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[DELIVERY_ORDER] Create delivery order failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Parse delivery date
	deliveryDate, err := time.Parse("2006-01-02", req.DeliveryDate)
	if err != nil {
		log.Printf("[DELIVERY_ORDER] Create delivery order failed - Invalid delivery_date format: %s, error: %v", req.DeliveryDate, err)
		return helper.Fail(c, 400, "Invalid delivery_date format, use YYYY-MM-DD", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[DELIVERY_ORDER] Create delivery order failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[DELIVERY_ORDER] Creating delivery order with audit - User ID: %d, DO Number: %s", userID, req.DONumber)

	order, err := deliveryOrderService.CreateDeliveryOrder(
		req.DONumber,
		req.PurchaseOrderID,
		deliveryDate,
		req.Status,
		req.Description,
		userID,
	)

	if err != nil {
		log.Printf("[DELIVERY_ORDER] Create delivery order failed - DO Number: %s, User ID: %d, error: %v", req.DONumber, userID, err)
		statusCode, message := handleDeliveryOrderError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[DELIVERY_ORDER] Create delivery order successful - DO Number: %s, Created by User ID: %d", req.DONumber, userID)
	return helper.Success(c, 201, "Delivery order created successfully", order)
}

// CreateDeliveryOrderWithItems creates a delivery order with multiple items in a transaction
func CreateDeliveryOrderWithItems(c *fiber.Ctx) error {
	log.Printf("[DELIVERY_ORDER] Create delivery order with items request from IP: %s", c.IP())

	var req CreateDeliveryOrderWithItemsRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[DELIVERY_ORDER] Create delivery order with items failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Validate items array
	if len(req.Items) == 0 {
		log.Printf("[DELIVERY_ORDER] Create delivery order with items failed - No items provided")
		return helper.Fail(c, 400, "At least one item is required", "Items array is empty")
	}

	// Parse delivery date
	deliveryDate, err := time.Parse("2006-01-02", req.DeliveryDate)
	if err != nil {
		log.Printf("[DELIVERY_ORDER] Create delivery order with items failed - Invalid delivery_date format: %s, error: %v", req.DeliveryDate, err)
		return helper.Fail(c, 400, "Invalid delivery_date format, use YYYY-MM-DD", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[DELIVERY_ORDER] Create delivery order with items failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	// Convert request items to model items
	items := make([]model.DeliveryOrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = model.DeliveryOrderItem{
			ProductID:    item.ProductID,
			QtyDelivered: &item.QtyDelivered,
			Description:  item.Description,
		}
	}

	log.Printf("[DELIVERY_ORDER] Creating delivery order with %d items - User ID: %d, DO Number: %s", len(items), userID, req.DONumber)

	order, err := deliveryOrderService.CreateDeliveryOrderWithItems(
		req.DONumber,
		req.PurchaseOrderID,
		deliveryDate,
		req.Status,
		req.Description,
		items,
		userID,
	)

	if err != nil {
		log.Printf("[DELIVERY_ORDER] Create delivery order with items failed - DO Number: %s, User ID: %d, error: %v", req.DONumber, userID, err)
		statusCode, message := handleDeliveryOrderError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[DELIVERY_ORDER] Create delivery order with items successful - DO Number: %s, Created by User ID: %d", req.DONumber, userID)
	return helper.Success(c, 201, "Delivery order with items created successfully", order)
}

func UpdateDeliveryOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[DELIVERY_ORDER] Update delivery order request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[DELIVERY_ORDER] Update delivery order failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid delivery order ID", err.Error())
	}

	var req UpdateDeliveryOrderRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[DELIVERY_ORDER] Update delivery order failed - Invalid request body for ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Parse delivery date if provided
	var deliveryDate time.Time
	if req.DeliveryDate != "" {
		deliveryDate, err = time.Parse("2006-01-02", req.DeliveryDate)
		if err != nil {
			log.Printf("[DELIVERY_ORDER] Update delivery order failed - Invalid delivery_date format: %s, error: %v", req.DeliveryDate, err)
			return helper.Fail(c, 400, "Invalid delivery_date format, use YYYY-MM-DD", err.Error())
		}
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[DELIVERY_ORDER] Update delivery order failed - User not authenticated for Order ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[DELIVERY_ORDER] Updating delivery order with audit - Order ID: %d, User ID: %d", idUint, userID)

	order, err := deliveryOrderService.UpdateDeliveryOrder(
		uint(idUint),
		req.DONumber,
		req.PurchaseOrderID,
		deliveryDate,
		req.Status,
		req.Description,
		userID,
	)

	if err != nil {
		log.Printf("[DELIVERY_ORDER] Update delivery order failed - Order ID: %d, User ID: %d, error: %v", idUint, userID, err)
		statusCode, message := handleDeliveryOrderError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[DELIVERY_ORDER] Update delivery order successful - Order ID: %d, Updated by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Delivery order updated successfully", order)
}

// UpdateDeliveryOrderWithItems updates a delivery order and replaces all its items in a transaction
func UpdateDeliveryOrderWithItems(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[DELIVERY_ORDER] Update delivery order with items request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[DELIVERY_ORDER] Update delivery order with items failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid delivery order ID", err.Error())
	}

	var req UpdateDeliveryOrderWithItemsRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[DELIVERY_ORDER] Update delivery order with items failed - Invalid request body for ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Validate items array
	if len(req.Items) == 0 {
		log.Printf("[DELIVERY_ORDER] Update delivery order with items failed - No items provided for ID: %d", idUint)
		return helper.Fail(c, 400, "At least one item is required", "Items array is empty")
	}

	// Parse delivery date if provided
	var deliveryDate time.Time
	if req.DeliveryDate != "" {
		deliveryDate, err = time.Parse("2006-01-02", req.DeliveryDate)
		if err != nil {
			log.Printf("[DELIVERY_ORDER] Update delivery order with items failed - Invalid delivery_date format: %s, error: %v", req.DeliveryDate, err)
			return helper.Fail(c, 400, "Invalid delivery_date format, use YYYY-MM-DD", err.Error())
		}
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[DELIVERY_ORDER] Update delivery order with items failed - User not authenticated for Order ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	// Convert request items to model items
	items := make([]model.DeliveryOrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = model.DeliveryOrderItem{
			ProductID:    item.ProductID,
			QtyDelivered: &item.QtyDelivered,
			Description:  item.Description,
		}
	}

	log.Printf("[DELIVERY_ORDER] Updating delivery order with %d items - Order ID: %d, User ID: %d", len(items), idUint, userID)

	order, err := deliveryOrderService.UpdateDeliveryOrderWithItems(
		uint(idUint),
		req.DONumber,
		req.PurchaseOrderID,
		deliveryDate,
		req.Status,
		req.Description,
		items,
		userID,
	)

	if err != nil {
		log.Printf("[DELIVERY_ORDER] Update delivery order with items failed - Order ID: %d, User ID: %d, error: %v", idUint, userID, err)
		statusCode, message := handleDeliveryOrderError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[DELIVERY_ORDER] Update delivery order with items successful - Order ID: %d, Updated by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Delivery order with items updated successfully", order)
}

func DeleteDeliveryOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[DELIVERY_ORDER] Delete delivery order request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[DELIVERY_ORDER] Delete delivery order failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid delivery order ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[DELIVERY_ORDER] Delete delivery order failed - User not authenticated for Order ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = deliveryOrderService.DeleteDeliveryOrder(uint(idUint), userID)
	if err != nil {
		log.Printf("[DELIVERY_ORDER] Delete delivery order failed - Order ID: %d, error: %v", idUint, err)
		statusCode, message := handleDeliveryOrderError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[DELIVERY_ORDER] Delete delivery order successful - Order ID: %d, Deleted by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Delivery order deleted successfully", nil)
}

func GetDeletedDeliveryOrders(c *fiber.Ctx) error {
	log.Printf("[DELIVERY_ORDER] Get deleted delivery orders request from IP: %s", c.IP())

	orders, err := deliveryOrderService.GetDeletedDeliveryOrders()
	if err != nil {
		log.Printf("[DELIVERY_ORDER] Get deleted delivery orders failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch deleted delivery orders", err.Error())
	}

	log.Printf("[DELIVERY_ORDER] Get deleted delivery orders successful")
	return helper.Success(c, 200, "Success", orders)
}

func RestoreDeliveryOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[DELIVERY_ORDER] Restore delivery order request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[DELIVERY_ORDER] Restore delivery order failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid delivery order ID", err.Error())
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[DELIVERY_ORDER] Restore delivery order failed - User not authenticated for Order ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	order, err := deliveryOrderService.RestoreDeliveryOrder(uint(idUint), userID)
	if err != nil {
		log.Printf("[DELIVERY_ORDER] Restore delivery order failed - Order ID: %d, error: %v", idUint, err)
		statusCode, message := handleDeliveryOrderError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[DELIVERY_ORDER] Restore delivery order successful - Order ID: %d, Restored by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Delivery order restored successfully", order)
}
