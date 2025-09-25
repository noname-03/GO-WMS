package handler

import (
	"log"
	"strconv"
	"time"

	"myapp/internal/service"
	"myapp/pkg/helper"

	"github.com/gofiber/fiber/v2"
)

var productItemTrackService = service.NewProductItemTrackService()

type CreateProductItemTrackRequest struct {
	ProductItemID  uint     `json:"product_item_id" validate:"required"`
	ProductStockID *uint    `json:"product_stock_id,omitempty"`
	ProductID      *uint    `json:"product_id,omitempty"`
	ProductBatchID *uint    `json:"product_batch_id,omitempty"`
	UnitPrice      *string  `json:"unit_price"`
	StockIn        *float64 `json:"stock_in" validate:"omitempty,gte=0"`
	StockOut       *float64 `json:"stock_out" validate:"omitempty,gte=0"`
	Quantity       *float64 `json:"quantity" validate:"omitempty,gt=0"`
	Operation      *string  `json:"operation" validate:"omitempty,oneof=In Out Plus Minus"`
	Stock          *float64 `json:"stock" validate:"omitempty,gte=0"`
	Action         string   `json:"action,omitempty"` // CREATE, UPDATE, DELETE, STOCK_IN, STOCK_OUT
}

type UpdateProductItemTrackRequest struct {
	UnitPrice *string  `json:"unit_price,omitempty"`
	Quantity  *float64 `json:"quantity,omitempty" validate:"omitempty,gt=0"`
	Operation *string  `json:"operation,omitempty" validate:"omitempty,oneof=In Out Plus Minus"`
	Stock     *float64 `json:"stock,omitempty" validate:"omitempty,gte=0"`
}

type DateRangeRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func GetAllProductItemTracks(c *fiber.Ctx) error {
	log.Printf("[PRODUCT_ITEM_TRACK] Get all product item tracks request from IP: %s", c.IP())

	result, err := productItemTrackService.GetAllProductItemTracks()
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Get all failed, error: %v", err)
		return helper.Fail(c, 400, "Failed to retrieve product item tracks", err.Error())
	}

	log.Printf("[PRODUCT_ITEM_TRACK] Get all successful")
	return helper.Success(c, 200, "Product item tracks retrieved successfully", result)
}

func GetProductItemTracksByItem(c *fiber.Ctx) error {
	itemID := c.Params("itemId")
	log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by item request - Item ID: %s from IP: %s", itemID, c.IP())

	itemIDUint, err := strconv.ParseUint(itemID, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by item failed - Invalid item ID: %s, error: %v", itemID, err)
		return helper.Fail(c, 400, "Invalid item ID", err.Error())
	}

	result, err := productItemTrackService.GetProductItemTracksByItem(uint(itemIDUint))
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by item failed - Item ID: %d, error: %v", itemIDUint, err)
		return helper.Fail(c, 400, "Failed to retrieve product item tracks", err.Error())
	}

	log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by item successful")
	return helper.Success(c, 200, "Product item tracks retrieved successfully", result)
}

func GetProductItemTracksByStock(c *fiber.Ctx) error {
	stockID := c.Params("stockId")
	log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by stock request - Stock ID: %s from IP: %s", stockID, c.IP())

	stockIDUint, err := strconv.ParseUint(stockID, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by stock failed - Invalid stock ID: %s, error: %v", stockID, err)
		return helper.Fail(c, 400, "Invalid stock ID", err.Error())
	}

	result, err := productItemTrackService.GetProductItemTracksByStock(uint(stockIDUint))
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by stock failed - Stock ID: %d, error: %v", stockIDUint, err)
		return helper.Fail(c, 400, "Failed to retrieve product item tracks", err.Error())
	}

	log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by stock successful")
	return helper.Success(c, 200, "Product item tracks retrieved successfully", result)
}

func GetProductItemTracksByProduct(c *fiber.Ctx) error {
	productID := c.Params("productId")
	log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by product request - Product ID: %s from IP: %s", productID, c.IP())

	productIDUint, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by product failed - Invalid product ID: %s, error: %v", productID, err)
		return helper.Fail(c, 400, "Invalid product ID", err.Error())
	}

	result, err := productItemTrackService.GetProductItemTracksByProduct(uint(productIDUint))
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by product failed - Product ID: %d, error: %v", productIDUint, err)
		return helper.Fail(c, 400, "Failed to retrieve product item tracks", err.Error())
	}

	log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by product successful")
	return helper.Success(c, 200, "Product item tracks retrieved successfully", result)
}

func GetProductItemTracksByDateRange(c *fiber.Ctx) error {
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by date range request - Start: %s, End: %s from IP: %s", startDateStr, endDateStr, c.IP())

	if startDateStr == "" || endDateStr == "" {
		log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by date range failed - Missing date parameters")
		return helper.Fail(c, 400, "Invalid parameters", "startDate and endDate query parameters are required")
	}

	_, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by date range failed - Invalid start date: %s, error: %v", startDateStr, err)
		return helper.Fail(c, 400, "Invalid date format", "Invalid startDate format. Use YYYY-MM-DD")
	}

	_, err = time.Parse("2006-01-02", endDateStr)
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by date range failed - Invalid end date: %s, error: %v", endDateStr, err)
		return helper.Fail(c, 400, "Invalid date format", "Invalid endDate format. Use YYYY-MM-DD")
	}

	// For now, return all tracks (date range filtering can be implemented later)
	result, err := productItemTrackService.GetAllProductItemTracks()
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by date range failed, error: %v", err)
		return helper.Fail(c, 400, "Failed to retrieve product item tracks", err.Error())
	}

	log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by date range successful")
	return helper.Success(c, 200, "Product item tracks retrieved successfully", result)
}

func GetTracksByOperation(c *fiber.Ctx) error {
	operation := c.Params("operation")
	log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by operation request - Operation: %s from IP: %s", operation, c.IP())

	// For now, return all tracks (operation filtering can be implemented later)
	result, err := productItemTrackService.GetAllProductItemTracks()
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by operation failed - Operation: %s, error: %v", operation, err)
		return helper.Fail(c, 400, "Failed to retrieve product item tracks", err.Error())
	}

	log.Printf("[PRODUCT_ITEM_TRACK] Get tracks by operation successful")
	return helper.Success(c, 200, "Product item tracks retrieved successfully", result)
}

func GetValueReportByProduct(c *fiber.Ctx) error {
	log.Printf("[PRODUCT_ITEM_TRACK] Get value report request from IP: %s", c.IP())

	// For now, return all tracks (value report can be implemented later)
	result, err := productItemTrackService.GetAllProductItemTracks()
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Get value report failed, error: %v", err)
		return helper.Fail(c, 400, "Failed to retrieve value report", err.Error())
	}

	log.Printf("[PRODUCT_ITEM_TRACK] Get value report successful")
	return helper.Success(c, 200, "Value report retrieved successfully", result)
}

func GetProductItemTrackByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_ITEM_TRACK] Get track by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Get track by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid track ID", err.Error())
	}

	result, err := productItemTrackService.GetProductItemTrackByID(uint(idUint))
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Get track by ID failed - ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 404, "Product item track not found", err.Error())
	}

	log.Printf("[PRODUCT_ITEM_TRACK] Get track by ID successful")
	return helper.Success(c, 200, "Product item track retrieved successfully", result)
}

func CreateProductItemTrack(c *fiber.Ctx) error {
	log.Printf("[PRODUCT_ITEM_TRACK] Create product item track request from IP: %s", c.IP())

	var req CreateProductItemTrackRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Create failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_ITEM_TRACK] Create failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	result, err := productItemTrackService.CreateProductItemTrack(req.ProductItemID, req.ProductStockID, req.ProductID, req.ProductBatchID, req.UnitPrice, req.StockIn, req.StockOut, req.Quantity, req.Operation, req.Stock, req.Action, userID)
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Create failed, error: %v", err)
		return helper.Fail(c, 400, "Failed to create product item track", err.Error())
	}

	log.Printf("[PRODUCT_ITEM_TRACK] Create successful")
	return helper.Success(c, 201, "Product item track created successfully", result)
}

func UpdateProductItemTrack(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_ITEM_TRACK] Update track request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Update failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid track ID", err.Error())
	}

	var req UpdateProductItemTrackRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Update failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_ITEM_TRACK] Update failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	result, err := productItemTrackService.UpdateProductItemTrack(uint(idUint), req.UnitPrice, req.Quantity, req.Operation, req.Stock, userID)
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Update failed - Track ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Failed to update product item track", err.Error())
	}

	log.Printf("[PRODUCT_ITEM_TRACK] Update successful")
	return helper.Success(c, 200, "Product item track updated successfully", result)
}

func DeleteProductItemTrack(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_ITEM_TRACK] Delete track request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Delete failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid track ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_ITEM_TRACK] Delete failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = productItemTrackService.DeleteProductItemTrack(uint(idUint), userID)
	if err != nil {
		log.Printf("[PRODUCT_ITEM_TRACK] Delete failed - Track ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Failed to delete product item track", err.Error())
	}

	log.Printf("[PRODUCT_ITEM_TRACK] Delete successful")
	return helper.Success(c, 200, "Product item track deleted successfully", nil)
}
