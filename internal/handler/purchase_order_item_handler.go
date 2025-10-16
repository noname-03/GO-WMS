package handler

import (
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var purchaseOrderItemService = service.NewPurchaseOrderItemService()

// handlePurchaseOrderItemError converts database errors to user-friendly messages for purchase order item operations
func handlePurchaseOrderItemError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors first
	if errMsg == "purchase order item not found" {
		return 404, "Purchase order item not found"
	}

	if errMsg == "purchase order not found" {
		return 404, "Purchase order not found"
	}

	if errMsg == "product not found" {
		return 404, "Product not found"
	}

	// Handle PostgreSQL constraint errors as backup
	if strings.Contains(errMsg, "foreign key constraint") {
		return 400, "Invalid purchase order ID or product ID"
	}

	// Default to 500 for other errors
	return 500, "Internal server error"
}

type CreatePurchaseOrderItemRequest struct {
	PurchaseOrderID uint     `json:"purchaseOrderId" validate:"required"`
	ProductID       uint     `json:"productId" validate:"required"`
	QtyOrdered      *float64 `json:"qtyOrdered" validate:"required,gt=0"`
	UnitPrice       *float64 `json:"unitPrice" validate:"required,gte=0"`
	Discount        *float64 `json:"discount" validate:"omitempty,gte=0"`
	TotalPrice      *float64 `json:"totalPrice" validate:"required,gte=0"`
	Description     *string  `json:"description"`
}

type UpdatePurchaseOrderItemRequest struct {
	PurchaseOrderID uint     `json:"purchaseOrderId"`
	ProductID       uint     `json:"productId"`
	QtyOrdered      *float64 `json:"qtyOrdered" validate:"omitempty,gt=0"`
	UnitPrice       *float64 `json:"unitPrice" validate:"omitempty,gte=0"`
	Discount        *float64 `json:"discount" validate:"omitempty,gte=0"`
	TotalPrice      *float64 `json:"totalPrice" validate:"omitempty,gte=0"`
	Description     *string  `json:"description"`
}

func GetAllPurchaseOrderItems(c *fiber.Ctx) error {
	log.Printf("[PURCHASE_ORDER_ITEM] Get all purchase order items request from IP: %s", c.IP())

	items, err := purchaseOrderItemService.GetAllPurchaseOrderItems()
	if err != nil {
		log.Printf("[PURCHASE_ORDER_ITEM] Get all items failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch purchase order items", err.Error())
	}

	log.Printf("[PURCHASE_ORDER_ITEM] Get all items successful")
	return helper.Success(c, 200, "Success", items)
}

func GetPurchaseOrderItemsByPurchaseOrder(c *fiber.Ctx) error {
	purchaseOrderID := c.Params("purchaseOrderId")
	log.Printf("[PURCHASE_ORDER_ITEM] Get items by purchase order request - Purchase Order ID: %s from IP: %s", purchaseOrderID, c.IP())

	purchaseOrderIDUint, err := strconv.ParseUint(purchaseOrderID, 10, 32)
	if err != nil {
		log.Printf("[PURCHASE_ORDER_ITEM] Get items by purchase order failed - Invalid Purchase Order ID: %s, error: %v", purchaseOrderID, err)
		return helper.Fail(c, 400, "Invalid purchase order ID", err.Error())
	}

	items, err := purchaseOrderItemService.GetPurchaseOrderItemsByPurchaseOrder(uint(purchaseOrderIDUint))
	if err != nil {
		log.Printf("[PURCHASE_ORDER_ITEM] Get items by purchase order failed - Purchase Order ID: %d, error: %v", purchaseOrderIDUint, err)
		statusCode, message := handlePurchaseOrderItemError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PURCHASE_ORDER_ITEM] Get items by purchase order successful")
	return helper.Success(c, 200, "Success", items)
}

func GetPurchaseOrderItemByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PURCHASE_ORDER_ITEM] Get item by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PURCHASE_ORDER_ITEM] Get item by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid purchase order item ID", err.Error())
	}

	item, err := purchaseOrderItemService.GetPurchaseOrderItemByID(uint(idUint))
	if err != nil {
		log.Printf("[PURCHASE_ORDER_ITEM] Get item by ID failed - Item ID: %d not found, error: %v", idUint, err)
		return helper.Fail(c, 404, "Purchase order item not found", err.Error())
	}

	log.Printf("[PURCHASE_ORDER_ITEM] Get item by ID successful - Item ID: %d", idUint)
	return helper.Success(c, 200, "Success", item)
}

func CreatePurchaseOrderItem(c *fiber.Ctx) error {
	log.Printf("[PURCHASE_ORDER_ITEM] Create purchase order item request from IP: %s", c.IP())

	var req CreatePurchaseOrderItemRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PURCHASE_ORDER_ITEM] Create item failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PURCHASE_ORDER_ITEM] Create item failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[PURCHASE_ORDER_ITEM] Creating purchase order item with audit - User ID: %d, Purchase Order ID: %d", userID, req.PurchaseOrderID)

	item, err := purchaseOrderItemService.CreatePurchaseOrderItem(req.PurchaseOrderID, req.ProductID, req.QtyOrdered, req.UnitPrice, req.Discount, req.TotalPrice, req.Description, userID)
	if err != nil {
		log.Printf("[PURCHASE_ORDER_ITEM] Create item failed - Purchase Order ID: %d, User ID: %d, error: %v", req.PurchaseOrderID, userID, err)
		statusCode, message := handlePurchaseOrderItemError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PURCHASE_ORDER_ITEM] Create item successful - Purchase Order ID: %d, Created by User ID: %d", req.PurchaseOrderID, userID)
	return helper.Success(c, 201, "Purchase order item created successfully", item)
}

func UpdatePurchaseOrderItem(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PURCHASE_ORDER_ITEM] Update item request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PURCHASE_ORDER_ITEM] Update item failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid purchase order item ID", err.Error())
	}

	var req UpdatePurchaseOrderItemRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PURCHASE_ORDER_ITEM] Update item failed - Invalid request body for ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PURCHASE_ORDER_ITEM] Update item failed - User not authenticated for Item ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[PURCHASE_ORDER_ITEM] Updating purchase order item with audit - Item ID: %d, User ID: %d", idUint, userID)

	item, err := purchaseOrderItemService.UpdatePurchaseOrderItem(uint(idUint), req.PurchaseOrderID, req.ProductID, req.QtyOrdered, req.UnitPrice, req.Discount, req.TotalPrice, req.Description, userID)
	if err != nil {
		log.Printf("[PURCHASE_ORDER_ITEM] Update item failed - Item ID: %d, User ID: %d, error: %v", idUint, userID, err)
		statusCode, message := handlePurchaseOrderItemError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PURCHASE_ORDER_ITEM] Update item successful - Item ID: %d, Updated by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Purchase order item updated successfully", item)
}

func DeletePurchaseOrderItem(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PURCHASE_ORDER_ITEM] Delete item request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PURCHASE_ORDER_ITEM] Delete item failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid purchase order item ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PURCHASE_ORDER_ITEM] Delete item failed - User not authenticated for Item ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = purchaseOrderItemService.DeletePurchaseOrderItem(uint(idUint), userID)
	if err != nil {
		log.Printf("[PURCHASE_ORDER_ITEM] Delete item failed - Item ID: %d, error: %v", idUint, err)
		statusCode, message := handlePurchaseOrderItemError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PURCHASE_ORDER_ITEM] Delete item successful - Item ID: %d, Deleted by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Purchase order item deleted successfully", nil)
}

func GetDeletedPurchaseOrderItems(c *fiber.Ctx) error {
	log.Printf("[PURCHASE_ORDER_ITEM] Get deleted items request from IP: %s", c.IP())

	items, err := purchaseOrderItemService.GetDeletedPurchaseOrderItems()
	if err != nil {
		log.Printf("[PURCHASE_ORDER_ITEM] Get deleted items failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch deleted purchase order items", err.Error())
	}

	log.Printf("[PURCHASE_ORDER_ITEM] Get deleted items successful")
	return helper.Success(c, 200, "Success", items)
}

func RestorePurchaseOrderItem(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PURCHASE_ORDER_ITEM] Restore item request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PURCHASE_ORDER_ITEM] Restore item failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid purchase order item ID", err.Error())
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PURCHASE_ORDER_ITEM] Restore item failed - User not authenticated for Item ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	item, err := purchaseOrderItemService.RestorePurchaseOrderItem(uint(idUint), userID)
	if err != nil {
		log.Printf("[PURCHASE_ORDER_ITEM] Restore item failed - Item ID: %d, error: %v", idUint, err)
		statusCode, message := handlePurchaseOrderItemError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PURCHASE_ORDER_ITEM] Restore item successful - Item ID: %d, Restored by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Purchase order item restored successfully", item)
}
