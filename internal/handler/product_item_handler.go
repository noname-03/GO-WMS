package handler

import (
	"log"
	"strconv"

	"myapp/internal/service"
	"myapp/pkg/helper"

	"github.com/gofiber/fiber/v2"
)

var productItemService = service.NewProductItemService()

type CreateProductItemRequest struct {
	ProductStockID uint     `json:"productStockId" validate:"required"`
	ProductID      uint     `json:"productId" validate:"required"`
	ProductBatchID uint     `json:"productBatchId" validate:"required"`
	LocationID     *uint    `json:"locationId"`
	StockIn        *float64 `json:"stockIn" validate:"omitempty,gte=0"`
	StockOut       *float64 `json:"stockOut" validate:"omitempty,gte=0"`
	Quantity       *float64 `json:"quantity" validate:"omitempty,gte=0"`
}

type UpdateProductItemRequest struct {
	ProductStockID *uint    `json:"productStockId,omitempty" validate:"omitempty,min=1"`
	ProductID      *uint    `json:"productId,omitempty" validate:"omitempty,min=1"`
	LocationID     *uint    `json:"locationId,omitempty"`
	StockIn        *float64 `json:"stockIn,omitempty" validate:"omitempty,gte=0"`
	StockOut       *float64 `json:"stockOut,omitempty" validate:"omitempty,gte=0"`
	Quantity       *float64 `json:"quantity,omitempty" validate:"omitempty,gte=0"`
}

func GetAllProductItems(c *fiber.Ctx) error {
	log.Printf("[PRODUCT_ITEM] Get all product items request from IP: %s", c.IP())

	result, err := productItemService.GetAllProductItems()
	if err != nil {
		log.Printf("[PRODUCT_ITEM] Get all failed, error: %v", err)
		return helper.Fail(c, 400, "Failed to retrieve product items", err.Error())
	}

	log.Printf("[PRODUCT_ITEM] Get all successful")
	return helper.Success(c, 200, "Product items retrieved successfully", result)
}

func GetProductItemsByStock(c *fiber.Ctx) error {
	stockID := c.Params("stockId")
	log.Printf("[PRODUCT_ITEM] Get items by stock request - Stock ID: %s from IP: %s", stockID, c.IP())

	stockIDUint, err := strconv.ParseUint(stockID, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_ITEM] Get items by stock failed - Invalid stock ID: %s, error: %v", stockID, err)
		return helper.Fail(c, 400, "Invalid stock ID", err.Error())
	}

	result, err := productItemService.GetProductItemsByStock(uint(stockIDUint))
	if err != nil {
		log.Printf("[PRODUCT_ITEM] Get items by stock failed - Stock ID: %d, error: %v", stockIDUint, err)
		return helper.Fail(c, 400, "Failed to retrieve product items", err.Error())
	}

	log.Printf("[PRODUCT_ITEM] Get items by stock successful")
	return helper.Success(c, 200, "Product items retrieved successfully", result)
}

func GetProductItemsByProduct(c *fiber.Ctx) error {
	productID := c.Params("productId")
	log.Printf("[PRODUCT_ITEM] Get items by product request - Product ID: %s from IP: %s", productID, c.IP())

	productIDUint, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_ITEM] Get items by product failed - Invalid product ID: %s, error: %v", productID, err)
		return helper.Fail(c, 400, "Invalid product ID", err.Error())
	}

	result, err := productItemService.GetProductItemsByProduct(uint(productIDUint))
	if err != nil {
		log.Printf("[PRODUCT_ITEM] Get items by product failed - Product ID: %d, error: %v", productIDUint, err)
		return helper.Fail(c, 400, "Failed to retrieve product items", err.Error())
	}

	log.Printf("[PRODUCT_ITEM] Get items by product successful")
	return helper.Success(c, 200, "Product items retrieved successfully", result)
}

func GetProductItemsByLocation(c *fiber.Ctx) error {
	locationID := c.Params("locationId")
	log.Printf("[PRODUCT_ITEM] Get items by location request - Location ID: %s from IP: %s", locationID, c.IP())

	locationIDUint, err := strconv.ParseUint(locationID, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_ITEM] Get items by location failed - Invalid location ID: %s, error: %v", locationID, err)
		return helper.Fail(c, 400, "Invalid location ID", err.Error())
	}

	result, err := productItemService.GetProductItemsByLocation(uint(locationIDUint))
	if err != nil {
		log.Printf("[PRODUCT_ITEM] Get items by location failed - Location ID: %d, error: %v", locationIDUint, err)
		return helper.Fail(c, 400, "Failed to retrieve product items", err.Error())
	}

	log.Printf("[PRODUCT_ITEM] Get items by location successful")
	return helper.Success(c, 200, "Product items retrieved successfully", result)
}

func GetItemsSummaryByProduct(c *fiber.Ctx) error {
	log.Printf("[PRODUCT_ITEM] Get items summary request from IP: %s", c.IP())

	result, err := productItemService.GetAllProductItems()
	if err != nil {
		log.Printf("[PRODUCT_ITEM] Get items summary failed, error: %v", err)
		return helper.Fail(c, 400, "Failed to retrieve items summary", err.Error())
	}

	log.Printf("[PRODUCT_ITEM] Get items summary successful")
	return helper.Success(c, 200, "Items summary retrieved successfully", result)
}

func GetProductItemByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_ITEM] Get item by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_ITEM] Get item by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid item ID", err.Error())
	}

	result, err := productItemService.GetProductItemByID(uint(idUint))
	if err != nil {
		log.Printf("[PRODUCT_ITEM] Get item by ID failed - ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 404, "Product item not found", err.Error())
	}

	log.Printf("[PRODUCT_ITEM] Get item by ID successful")
	return helper.Success(c, 200, "Product item retrieved successfully", result)
}

func CreateProductItem(c *fiber.Ctx) error {
	log.Printf("[PRODUCT_ITEM] Create product item request from IP: %s", c.IP())

	var req CreateProductItemRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PRODUCT_ITEM] Create failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_ITEM] Create failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	result, err := productItemService.CreateProductItem(req.ProductStockID, req.ProductID, req.ProductBatchID, req.LocationID, req.StockIn, req.StockOut, req.Quantity, userID)
	if err != nil {
		log.Printf("[PRODUCT_ITEM] Create failed, error: %v", err)
		return helper.Fail(c, 400, "Failed to create product item", err.Error())
	}

	log.Printf("[PRODUCT_ITEM] Create successful")
	return helper.Success(c, 201, "Product item created successfully", result)
}

func UpdateProductItem(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_ITEM] Update item request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_ITEM] Update failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid item ID", err.Error())
	}

	var req UpdateProductItemRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PRODUCT_ITEM] Update failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_ITEM] Update failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	result, err := productItemService.UpdateProductItem(uint(idUint), req.ProductStockID, req.ProductID, req.LocationID, req.StockIn, req.StockOut, req.Quantity, userID)
	if err != nil {
		log.Printf("[PRODUCT_ITEM] Update failed - Item ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Failed to update product item", err.Error())
	}

	log.Printf("[PRODUCT_ITEM] Update successful")
	return helper.Success(c, 200, "Product item updated successfully", result)
}

func DeleteProductItem(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_ITEM] Delete item request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_ITEM] Delete failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid item ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_ITEM] Delete failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = productItemService.DeleteProductItem(uint(idUint), userID)
	if err != nil {
		log.Printf("[PRODUCT_ITEM] Delete failed - Item ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Failed to delete product item", err.Error())
	}

	log.Printf("[PRODUCT_ITEM] Delete successful")
	return helper.Success(c, 200, "Product item deleted successfully", nil)
}
