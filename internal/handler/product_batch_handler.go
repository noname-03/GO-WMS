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

var productBatchService = service.NewProductBatchService()

// handleProductBatchError converts database errors to user-friendly messages for product batch operations
func handleProductBatchError(err error) (int, string) {
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

type CreateProductBatchRequest struct {
	ProductID   uint     `json:"productId" validate:"required"`
	CodeBatch   *string  `json:"codeBatch"` // Nullable field
	UnitPrice   *float64 `json:"unitPrice"` // Nullable field
	ExpDate     string   `json:"expDate" validate:"required"`
	Description *string  `json:"description"` // Nullable field
}

type UpdateProductBatchRequest struct {
	ProductID   uint     `json:"productId"`
	CodeBatch   *string  `json:"codeBatch"` // Nullable field
	UnitPrice   *float64 `json:"unitPrice"` // Nullable field
	ExpDate     string   `json:"expDate"`
	Description *string  `json:"description"` // Nullable field
}

func GetProductBatches(c *fiber.Ctx) error {
	log.Printf("[PRODUCT_BATCH] Get all product batches request from IP: %s", c.IP())

	batches, err := productBatchService.GetAllProductBatches()
	if err != nil {
		log.Printf("[PRODUCT_BATCH] Get all product batches failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch product batches", err.Error())
	}

	log.Printf("[PRODUCT_BATCH] Get all product batches successful - Found %d batches", len(batches))
	return helper.Success(c, 200, "Success", batches)
}

func GetProductBatchesByProduct(c *fiber.Ctx) error {
	productID := c.Params("productId")
	log.Printf("[PRODUCT_BATCH] Get product batches by product request - Product ID: %s from IP: %s", productID, c.IP())

	productIDUint, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_BATCH] Get product batches by product failed - Invalid Product ID: %s, error: %v", productID, err)
		return helper.Fail(c, 400, "Invalid product ID", err.Error())
	}

	batches, err := productBatchService.GetProductBatchesByProduct(uint(productIDUint))
	if err != nil {
		log.Printf("[PRODUCT_BATCH] Get product batches by product failed - Product ID: %d, error: %v", productIDUint, err)
		statusCode, message := handleProductBatchError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT_BATCH] Get product batches by product successful - Product ID: %d, Found %d batches", productIDUint, len(batches))
	return helper.Success(c, 200, "Success", batches)
}

func GetProductBatchByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_BATCH] Get product batch by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_BATCH] Get product batch by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid product batch ID", err.Error())
	}

	batch, err := productBatchService.GetProductBatchByID(uint(idUint))
	if err != nil {
		log.Printf("[PRODUCT_BATCH] Get product batch by ID failed - Batch ID: %d not found, error: %v", idUint, err)
		return helper.Fail(c, 404, "Product batch not found", err.Error())
	}

	log.Printf("[PRODUCT_BATCH] Get product batch by ID successful - Batch ID: %d, Product ID: %d", batch.ID, batch.ProductID)
	return helper.Success(c, 200, "Success", batch)
}

func CreateProductBatch(c *fiber.Ctx) error {
	log.Printf("[PRODUCT_BATCH] Create product batch request from IP: %s", c.IP())

	var req CreateProductBatchRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PRODUCT_BATCH] Create product batch failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Parse expiry date
	expDate, err := time.Parse("2006-01-02", req.ExpDate)
	if err != nil {
		log.Printf("[PRODUCT_BATCH] Create product batch failed - Invalid exp_date format: %s, error: %v", req.ExpDate, err)
		return helper.Fail(c, 400, "Invalid exp_date format, use YYYY-MM-DD", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_BATCH] Create product batch failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[PRODUCT_BATCH] Creating product batch with audit - User ID: %d, Product ID: %d", userID, req.ProductID)

	batch, err := productBatchService.CreateProductBatch(req.ProductID, req.CodeBatch, req.UnitPrice, expDate, req.Description, userID)
	if err != nil {
		log.Printf("[PRODUCT_BATCH] Create product batch failed - Product ID: %d, User ID: %d, error: %v", req.ProductID, userID, err)
		statusCode, message := handleProductBatchError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT_BATCH] Create product batch successful - Batch ID: %d, Product ID: %d, Created by User ID: %d", batch.ID, batch.ProductID, userID)
	return helper.Success(c, 201, "Product batch created successfully", batch)
}

func UpdateProductBatch(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_BATCH] Update product batch request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_BATCH] Update product batch failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid product batch ID", err.Error())
	}

	var req UpdateProductBatchRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PRODUCT_BATCH] Update product batch failed - Invalid request body for ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Parse expiry date if provided
	var expDate time.Time
	if req.ExpDate != "" {
		expDate, err = time.Parse("2006-01-02", req.ExpDate)
		if err != nil {
			log.Printf("[PRODUCT_BATCH] Update product batch failed - Invalid exp_date format: %s, error: %v", req.ExpDate, err)
			return helper.Fail(c, 400, "Invalid exp_date format, use YYYY-MM-DD", err.Error())
		}
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_BATCH] Update product batch failed - User not authenticated for Batch ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[PRODUCT_BATCH] Updating product batch with audit - Batch ID: %d, User ID: %d", idUint, userID)

	batch, err := productBatchService.UpdateProductBatch(uint(idUint), req.ProductID, req.CodeBatch, req.UnitPrice, expDate, req.Description, userID)
	if err != nil {
		log.Printf("[PRODUCT_BATCH] Update product batch failed - Batch ID: %d, User ID: %d, error: %v", idUint, userID, err)
		statusCode, message := handleProductBatchError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT_BATCH] Update product batch successful - Batch ID: %d, Product ID: %d, Updated by User ID: %d", batch.ID, batch.ProductID, userID)
	return helper.Success(c, 200, "Product batch updated successfully", batch)
}

func DeleteProductBatch(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_BATCH] Delete product batch request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_BATCH] Delete product batch failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid product batch ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_BATCH] Delete product batch failed - User not authenticated for Batch ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = productBatchService.DeleteProductBatch(uint(idUint), userID)
	if err != nil {
		log.Printf("[PRODUCT_BATCH] Delete product batch failed - Batch ID: %d, error: %v", idUint, err)
		statusCode, message := handleProductBatchError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT_BATCH] Delete product batch successful - Batch ID: %d, Deleted by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Product batch deleted successfully", nil)
}
