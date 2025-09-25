package handler

import (
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var productStockService = service.NewProductStockService()

// handleProductStockError converts database errors to user-friendly messages for product stock operations
func handleProductStockError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors first
	if errMsg == "product batch not found" {
		return 404, "Product batch not found"
	}

	if errMsg == "product not found" {
		return 404, "Product not found"
	}
	// Handle PostgreSQL constraint errors as backup
	if strings.Contains(errMsg, "foreign key constraint") {
		return 400, "Invalid product ID"
	}

	// Default to 500 for other errors
	return 500, "Internal server error"
}

type CreateProductStockRequest struct {
	ProductBatchID uint     `json:"productBatchId" validate:"required"`
	ProductID      uint     `json:"productId" validate:"required"`
	LocationID     uint     `json:"locationId" validate:"required"`
	Quantity       *float64 `json:"quantity" validate:"omitempty,gte=0"`
}

type UpdateProductStockRequest struct {
	ProductBatchID uint     `json:"productBatchId" validate:"required"`
	ProductID      uint     `json:"productId,omitempty" validate:"omitempty,min=1"`
	LocationID     uint     `json:"locationId,omitempty" validate:"omitempty,min=1"`
	Quantity       *float64 `json:"quantity,omitempty" validate:"omitempty,gte=0"`
}

func GetAllProductStocks(c *fiber.Ctx) error {
	log.Printf("[PRODUCT_STOCK] Get all product stocks request from IP: %s", c.IP())

	result, err := productStockService.GetAllProductStocks()
	if err != nil {
		log.Printf("[PRODUCT_STOCK] Get all failed, error: %v", err)
		return helper.Fail(c, 400, "Failed to retrieve product stocks", err.Error())
	}

	log.Printf("[PRODUCT_STOCK] Get all successful")
	return helper.Success(c, 200, "Product stocks retrieved successfully", result)
}

func GetProductStocksByProduct(c *fiber.Ctx) error {
	productID := c.Params("productId")
	log.Printf("[PRODUCT_STOCK] Get stocks by product request - Product ID: %s from IP: %s", productID, c.IP())

	productIDUint, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_STOCK] Get stocks by product failed - Invalid product ID: %s, error: %v", productID, err)
		return helper.Fail(c, 400, "Invalid product ID", err.Error())
	}

	result, err := productStockService.GetProductStocksByProduct(uint(productIDUint))
	if err != nil {
		log.Printf("[PRODUCT_STOCK] Get stocks by product failed - Product ID: %d, error: %v", productIDUint, err)
		return helper.Fail(c, 400, "Failed to retrieve product stocks", err.Error())
	}

	log.Printf("[PRODUCT_STOCK] Get stocks by product successful")
	return helper.Success(c, 200, "Product stocks retrieved successfully", result)
}

func GetProductStockByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_STOCK] Get stock by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_STOCK] Get stock by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid stock ID", err.Error())
	}

	result, err := productStockService.GetProductStockByID(uint(idUint))
	if err != nil {
		log.Printf("[PRODUCT_STOCK] Get stock by ID failed - Stock ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 404, "Product stock not found", err.Error())
	}

	log.Printf("[PRODUCT_STOCK] Get stock by ID successful")
	return helper.Success(c, 200, "Product stock retrieved successfully", result)
}

func CreateProductStock(c *fiber.Ctx) error {
	log.Printf("[PRODUCT_STOCK] Create product stock request from IP: %s", c.IP())

	var req CreateProductStockRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PRODUCT_STOCK] Create failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_STOCK] Create failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	// result, err := productStockService.CreateProductStock(req.ProductBatchID, req.ProductID, req.LocationID, req.Quantity, userID)
	result, err := productStockService.CreateProductStock(req.ProductBatchID, req.ProductID, req.LocationID, req.Quantity, userID)
	if err != nil {
		log.Printf("[PRODUCT_STOCK] Create failed, error: %v", err)
		return helper.Fail(c, 400, "Failed to create product stock", err.Error())
	}

	log.Printf("[PRODUCT_STOCK] Create successful")
	return helper.Success(c, 201, "Product stock created successfully", result)
}

func UpdateProductStock(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_STOCK] Update stock request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_STOCK] Update failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid stock ID", err.Error())
	}

	var req UpdateProductStockRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PRODUCT_STOCK] Update failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_STOCK] Update failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	// result, err := productStockService.UpdateProductStock(uint(idUint), req.ProductBatchID, req.ProductID, req.LocationID, req.Quantity, userID)
	result, err := productStockService.UpdateProductStock(uint(idUint), req.ProductBatchID, req.ProductID, req.LocationID, req.Quantity, userID)
	if err != nil {
		log.Printf("[PRODUCT_STOCK] Update failed - Stock ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Failed to update product stock", err.Error())
	}

	log.Printf("[PRODUCT_STOCK] Update successful")
	return helper.Success(c, 200, "Product stock updated successfully", result)
}

func DeleteProductStock(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_STOCK] Delete stock request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_STOCK] Delete failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid stock ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_STOCK] Delete failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = productStockService.DeleteProductStock(uint(idUint), userID)
	if err != nil {
		log.Printf("[PRODUCT_STOCK] Delete failed - Stock ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Failed to delete product stock", err.Error())
	}

	log.Printf("[PRODUCT_STOCK] Delete successful")
	return helper.Success(c, 200, "Product stock deleted successfully", nil)
}
