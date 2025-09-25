package handler

import (
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var productStockTrackService = service.NewProductStockTrackService()

// handleProductStockTrackError converts database errors to user-friendly messages for product stock track operations
func handleProductStockTrackError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors first
	if errMsg == "product stock track not found" {
		return 404, "Product stock track not found"
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

func GetAllProductStockTracks(c *fiber.Ctx) error {
	log.Printf("[PRODUCT_STOCK_TRACK] Get all product stock tracks request from IP: %s", c.IP())

	result, err := productStockTrackService.GetAllProductStockTracks()
	if err != nil {
		log.Printf("[PRODUCT_STOCK_TRACK] Get all failed, error: %v", err)
		return helper.Fail(c, 400, "Failed to retrieve product stock tracks", err.Error())
	}

	log.Printf("[PRODUCT_STOCK_TRACK] Get all successful")
	return helper.Success(c, 200, "Product stock tracks retrieved successfully", result)
}

func GetProductStockTracksByStock(c *fiber.Ctx) error {
	stockID := c.Params("stockId")
	log.Printf("[PRODUCT_STOCK_TRACK] Get tracks by stock request - Stock ID: %s from IP: %s", stockID, c.IP())

	stockIDUint, err := strconv.ParseUint(stockID, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_STOCK_TRACK] Get tracks by stock failed - Invalid stock ID: %s, error: %v", stockID, err)
		return helper.Fail(c, 400, "Invalid stock ID", err.Error())
	}

	result, err := productStockTrackService.GetProductStockTracksByStock(uint(stockIDUint))
	if err != nil {
		log.Printf("[PRODUCT_STOCK_TRACK] Get tracks by stock failed - Stock ID: %d, error: %v", stockIDUint, err)
		return helper.Fail(c, 400, "Failed to retrieve product stock tracks", err.Error())
	}

	log.Printf("[PRODUCT_STOCK_TRACK] Get tracks by stock successful")
	return helper.Success(c, 200, "Product stock tracks retrieved successfully", result)
}

func GetProductStockTracksByProduct(c *fiber.Ctx) error {
	productID := c.Params("productId")
	log.Printf("[PRODUCT_STOCK_TRACK] Get tracks by product request - Product ID: %s from IP: %s", productID, c.IP())

	productIDUint, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_STOCK_TRACK] Get tracks by product failed - Invalid product ID: %s, error: %v", productID, err)
		return helper.Fail(c, 400, "Invalid product ID", err.Error())
	}

	result, err := productStockTrackService.GetProductStockTracksByProduct(uint(productIDUint))
	if err != nil {
		log.Printf("[PRODUCT_STOCK_TRACK] Get tracks by product failed - Product ID: %d, error: %v", productIDUint, err)
		return helper.Fail(c, 400, "Failed to retrieve product stock tracks", err.Error())
	}

	log.Printf("[PRODUCT_STOCK_TRACK] Get tracks by product successful")
	return helper.Success(c, 200, "Product stock tracks retrieved successfully", result)
}

func GetProductStockTrackByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_STOCK_TRACK] Get track by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_STOCK_TRACK] Get track by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid track ID", err.Error())
	}

	result, err := productStockTrackService.GetProductStockTrackByID(uint(idUint))
	if err != nil {
		log.Printf("[PRODUCT_STOCK_TRACK] Get track by ID failed - ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 404, "Product stock track not found", err.Error())
	}

	log.Printf("[PRODUCT_STOCK_TRACK] Get track by ID successful")
	return helper.Success(c, 200, "Product stock track retrieved successfully", result)
}

func CreateProductStockTrack(c *fiber.Ctx) error {
	log.Printf("[PRODUCT_STOCK_TRACK] Create product stock track request from IP: %s", c.IP())

	var req service.CreateProductStockTrackRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PRODUCT_STOCK_TRACK] Create failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_STOCK_TRACK] Create failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	result, err := productStockTrackService.CreateProductStockTrack(req, userID)
	if err != nil {
		log.Printf("[PRODUCT_STOCK_TRACK] Create failed, error: %v", err)
		return helper.Fail(c, 400, "Failed to create product stock track", err.Error())
	}

	log.Printf("[PRODUCT_STOCK_TRACK] Create successful")
	return helper.Success(c, 201, "Product stock track created successfully", result)
}

func UpdateProductStockTrack(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_STOCK_TRACK] Update track request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_STOCK_TRACK] Update failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid track ID", err.Error())
	}

	var req service.UpdateProductStockTrackRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PRODUCT_STOCK_TRACK] Update failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_STOCK_TRACK] Update failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	result, err := productStockTrackService.UpdateProductStockTrack(uint(idUint), req, userID)
	if err != nil {
		log.Printf("[PRODUCT_STOCK_TRACK] Update failed - Track ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Failed to update product stock track", err.Error())
	}

	log.Printf("[PRODUCT_STOCK_TRACK] Update successful")
	return helper.Success(c, 200, "Product stock track updated successfully", result)
}

func DeleteProductStockTrack(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_STOCK_TRACK] Delete track request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_STOCK_TRACK] Delete failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid track ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_STOCK_TRACK] Delete failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = productStockTrackService.DeleteProductStockTrack(uint(idUint), userID)
	if err != nil {
		log.Printf("[PRODUCT_STOCK_TRACK] Delete failed - Track ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Failed to delete product stock track", err.Error())
	}

	log.Printf("[PRODUCT_STOCK_TRACK] Delete successful")
	return helper.Success(c, 200, "Product stock track deleted successfully", nil)
}
