package handler

import (
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var deliveryOrderItemService = service.NewDeliveryOrderItemService()

// handleDeliveryOrderItemError converts database errors to user-friendly messages for delivery order item operations
func handleDeliveryOrderItemError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors first
	if errMsg == "delivery order item not found" {
		return 404, "Delivery order item not found"
	}

	if errMsg == "delivery order not found" {
		return 404, "Delivery order not found"
	}

	if errMsg == "product not found" {
		return 404, "Product not found"
	}

	// Handle PostgreSQL constraint errors as backup
	if strings.Contains(errMsg, "foreign key constraint") {
		return 400, "Invalid delivery order or product ID"
	}

	// Default to 500 for other errors
	return 500, "Internal server error"
}

// Request DTOs
type CreateDeliveryOrderItemRequest struct {
	DeliveryOrderID uint    `json:"deliveryOrderId" validate:"required"`
	ProductID       uint    `json:"productId" validate:"required"`
	QtyDelivered    float64 `json:"qtyDelivered" validate:"required"`
	Description     *string `json:"description"`
}

type UpdateDeliveryOrderItemRequest struct {
	DeliveryOrderID uint    `json:"deliveryOrderId"`
	ProductID       uint    `json:"productId"`
	QtyDelivered    float64 `json:"qtyDelivered"`
	Description     *string `json:"description"`
}

func GetAllDeliveryOrderItems(c *fiber.Ctx) error {
	log.Printf("[DELIVERY_ORDER_ITEM] Get all delivery order items request from IP: %s", c.IP())

	items, err := deliveryOrderItemService.GetAllDeliveryOrderItems()
	if err != nil {
		log.Printf("[DELIVERY_ORDER_ITEM] Get all delivery order items failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch delivery order items", err.Error())
	}

	log.Printf("[DELIVERY_ORDER_ITEM] Get all delivery order items successful")
	return helper.Success(c, 200, "Success", items)
}

func GetDeliveryOrderItemsByDeliveryOrder(c *fiber.Ctx) error {
	deliveryOrderID := c.Params("deliveryOrderId")
	log.Printf("[DELIVERY_ORDER_ITEM] Get delivery order items by delivery order request - Delivery Order ID: %s from IP: %s", deliveryOrderID, c.IP())

	doID, err := strconv.ParseUint(deliveryOrderID, 10, 32)
	if err != nil {
		log.Printf("[DELIVERY_ORDER_ITEM] Get delivery order items by delivery order failed - Invalid Delivery Order ID: %s, error: %v", deliveryOrderID, err)
		return helper.Fail(c, 400, "Invalid delivery order ID", err.Error())
	}

	items, err := deliveryOrderItemService.GetDeliveryOrderItemsByDeliveryOrder(uint(doID))
	if err != nil {
		log.Printf("[DELIVERY_ORDER_ITEM] Get delivery order items by delivery order failed - Delivery Order ID: %d, error: %v", doID, err)
		statusCode, message := handleDeliveryOrderItemError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[DELIVERY_ORDER_ITEM] Get delivery order items by delivery order successful - Delivery Order ID: %d", doID)
	return helper.Success(c, 200, "Success", items)
}

func GetDeliveryOrderItemByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[DELIVERY_ORDER_ITEM] Get delivery order item by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[DELIVERY_ORDER_ITEM] Get delivery order item by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid delivery order item ID", err.Error())
	}

	item, err := deliveryOrderItemService.GetDeliveryOrderItemByID(uint(idUint))
	if err != nil {
		log.Printf("[DELIVERY_ORDER_ITEM] Get delivery order item by ID failed - Item ID: %d not found, error: %v", idUint, err)
		return helper.Fail(c, 404, "Delivery order item not found", err.Error())
	}

	log.Printf("[DELIVERY_ORDER_ITEM] Get delivery order item by ID successful - Item ID: %d", idUint)
	return helper.Success(c, 200, "Success", item)
}

func CreateDeliveryOrderItem(c *fiber.Ctx) error {
	log.Printf("[DELIVERY_ORDER_ITEM] Create delivery order item request from IP: %s", c.IP())

	var req CreateDeliveryOrderItemRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[DELIVERY_ORDER_ITEM] Create delivery order item failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[DELIVERY_ORDER_ITEM] Create delivery order item failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[DELIVERY_ORDER_ITEM] Creating delivery order item with audit - User ID: %d, Delivery Order ID: %d", userID, req.DeliveryOrderID)

	item, err := deliveryOrderItemService.CreateDeliveryOrderItem(
		req.DeliveryOrderID,
		req.ProductID,
		&req.QtyDelivered,
		req.Description,
		userID,
	)

	if err != nil {
		log.Printf("[DELIVERY_ORDER_ITEM] Create delivery order item failed - Delivery Order ID: %d, User ID: %d, error: %v", req.DeliveryOrderID, userID, err)
		statusCode, message := handleDeliveryOrderItemError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[DELIVERY_ORDER_ITEM] Create delivery order item successful - Delivery Order ID: %d, Created by User ID: %d", req.DeliveryOrderID, userID)
	return helper.Success(c, 201, "Delivery order item created successfully", item)
}

func UpdateDeliveryOrderItem(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[DELIVERY_ORDER_ITEM] Update delivery order item request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[DELIVERY_ORDER_ITEM] Update delivery order item failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid delivery order item ID", err.Error())
	}

	var req UpdateDeliveryOrderItemRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[DELIVERY_ORDER_ITEM] Update delivery order item failed - Invalid request body for ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[DELIVERY_ORDER_ITEM] Update delivery order item failed - User not authenticated for Item ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[DELIVERY_ORDER_ITEM] Updating delivery order item with audit - Item ID: %d, User ID: %d", idUint, userID)

	item, err := deliveryOrderItemService.UpdateDeliveryOrderItem(
		uint(idUint),
		req.DeliveryOrderID,
		req.ProductID,
		&req.QtyDelivered,
		req.Description,
		userID,
	)

	if err != nil {
		log.Printf("[DELIVERY_ORDER_ITEM] Update delivery order item failed - Item ID: %d, User ID: %d, error: %v", idUint, userID, err)
		statusCode, message := handleDeliveryOrderItemError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[DELIVERY_ORDER_ITEM] Update delivery order item successful - Item ID: %d, Updated by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Delivery order item updated successfully", item)
}

func DeleteDeliveryOrderItem(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[DELIVERY_ORDER_ITEM] Delete delivery order item request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[DELIVERY_ORDER_ITEM] Delete delivery order item failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid delivery order item ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[DELIVERY_ORDER_ITEM] Delete delivery order item failed - User not authenticated for Item ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = deliveryOrderItemService.DeleteDeliveryOrderItem(uint(idUint), userID)
	if err != nil {
		log.Printf("[DELIVERY_ORDER_ITEM] Delete delivery order item failed - Item ID: %d, error: %v", idUint, err)
		statusCode, message := handleDeliveryOrderItemError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[DELIVERY_ORDER_ITEM] Delete delivery order item successful - Item ID: %d, Deleted by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Delivery order item deleted successfully", nil)
}

func GetDeletedDeliveryOrderItems(c *fiber.Ctx) error {
	log.Printf("[DELIVERY_ORDER_ITEM] Get deleted delivery order items request from IP: %s", c.IP())

	items, err := deliveryOrderItemService.GetDeletedDeliveryOrderItems()
	if err != nil {
		log.Printf("[DELIVERY_ORDER_ITEM] Get deleted delivery order items failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch deleted delivery order items", err.Error())
	}

	log.Printf("[DELIVERY_ORDER_ITEM] Get deleted delivery order items successful")
	return helper.Success(c, 200, "Success", items)
}

func RestoreDeliveryOrderItem(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[DELIVERY_ORDER_ITEM] Restore delivery order item request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[DELIVERY_ORDER_ITEM] Restore delivery order item failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid delivery order item ID", err.Error())
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[DELIVERY_ORDER_ITEM] Restore delivery order item failed - User not authenticated for Item ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	item, err := deliveryOrderItemService.RestoreDeliveryOrderItem(uint(idUint), userID)
	if err != nil {
		log.Printf("[DELIVERY_ORDER_ITEM] Restore delivery order item failed - Item ID: %d, error: %v", idUint, err)
		statusCode, message := handleDeliveryOrderItemError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[DELIVERY_ORDER_ITEM] Restore delivery order item successful - Item ID: %d, Restored by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Delivery order item restored successfully", item)
}
