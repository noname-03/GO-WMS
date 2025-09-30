package handler

import (
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var productUnitService = service.NewProductUnitService()

// handleProductUnitError converts database errors to user-friendly messages for product unit operations
func handleProductUnitError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors first
	if errMsg == "product unit with this name already exists for this product" {
		return 409, "Product unit name already exists for this product"
	}

	if errMsg == "barcode already exists for this product and location" {
		return 409, "Barcode already exists for this product and location"
	}

	if errMsg == "barcode already exists for another product" {
		return 409, "barcode already exists for another product"
	}

	if errMsg == "product unit with same name and location already exists" {
		return 409, "product unit with same name and location already exists"
	}

	if errMsg == "product unit name already in use for this product and location" {
		return 409, "product unit name already in use for this product and location"
	}

	if errMsg == "product unit with this name already exists for this product and location" {
		return 409, "product unit with this name already exists for this product and location"
	}

	if errMsg == "product unit barcode already in use for this product" {
		return 409, "Product unit barcode already in use for this product"
	}

	if errMsg == "product unit not found" {
		return 404, "Product unit not found"
	}

	if errMsg == "product not found" {
		return 404, "Product not found"
	}

	if errMsg == "invalid location ID" {
		return 400, "Invalid location ID"
	}
	if errMsg == "location ID is required" {
		return 400, "Invalid location ID"
	}

	// Handle PostgreSQL constraint errors as backup
	if strings.Contains(errMsg, "duplicate key value violates unique constraint") {
		return 409, "Product unit already exists"
	}

	if strings.Contains(errMsg, "foreign key constraint") {
		return 400, "Invalid product ID"
	}

	// Default to 500 for other errors
	return 500, "Internal server error"
}

type CreateProductUnitRequest struct {
	ProductID      uint     `json:"productId" validate:"required"`
	LocationID     uint     `json:"locationId" validate:"required"`
	ProductBatchID uint     `json:"productBatchId" validate:"required"`
	Name           *string  `json:"name"`
	Quantity       *float64 `json:"quantity"`
	UnitPrice      *float64 `json:"unitPrice"`
	Barcode        *string  `json:"barcode"`
	Description    *string  `json:"description"`
}

type UpdateProductUnitRequest struct {
	ProductID      uint     `json:"productId" validate:"required"`
	LocationID     uint     `json:"locationId" validate:"required"`
	ProductBatchID uint     `json:"productBatchId" validate:"required"`
	Name           *string  `json:"name"`
	Quantity       *float64 `json:"quantity"`
	UnitPrice      *float64 `json:"unitPrice"`
	Barcode        *string  `json:"barcode"`
	Description    *string  `json:"description"`
}

func GetProductUnits(c *fiber.Ctx) error {
	log.Printf("[PRODUCT_UNIT] Get all product units request from IP: %s", c.IP())

	productUnits, err := productUnitService.GetAllProductUnits()
	if err != nil {
		log.Printf("[PRODUCT_UNIT] Get all product units failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch product units", err.Error())
	}

	log.Printf("[PRODUCT_UNIT] Get all product units successful")
	return helper.Success(c, 200, "Success", productUnits)
}

func GetProductUnitsByProduct(c *fiber.Ctx) error {
	productID := c.Params("productId")
	log.Printf("[PRODUCT_UNIT] Get product units by product request - Product ID: %s from IP: %s", productID, c.IP())

	productIDUint, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_UNIT] Get product units by product failed - Invalid Product ID: %s, error: %v", productID, err)
		return helper.Fail(c, 400, "Invalid product ID", err.Error())
	}

	productUnits, err := productUnitService.GetProductUnitsByProduct(uint(productIDUint))
	if err != nil {
		log.Printf("[PRODUCT_UNIT] Get product units by product failed - Product ID: %d, error: %v", productIDUint, err)
		statusCode, message := handleProductUnitError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT_UNIT] Get product units by product successful - Product ID: %d", productIDUint)
	return helper.Success(c, 200, "Success", productUnits)
}

func GetProductUnitByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_UNIT] Get product unit by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_UNIT] Get product unit by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid product unit ID", err.Error())
	}

	productUnit, err := productUnitService.GetProductUnitByID(uint(idUint))
	if err != nil {
		log.Printf("[PRODUCT_UNIT] Get product unit by ID failed - Product Unit ID: %d not found, error: %v", idUint, err)
		return helper.Fail(c, 404, "Product unit not found", err.Error())
	}

	log.Printf("[PRODUCT_UNIT] Get product unit by ID successful - Product Unit ID: %d", idUint)
	return helper.Success(c, 200, "Success", productUnit)
}

func CreateProductUnit(c *fiber.Ctx) error {
	log.Printf("[PRODUCT_UNIT] Create product unit request from IP: %s", c.IP())

	var req CreateProductUnitRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PRODUCT_UNIT] Create product unit failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_UNIT] Create product unit failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[PRODUCT_UNIT] Creating product unit with audit - User ID: %d, Product ID: %d, Product Batch ID: %d", userID, req.ProductID, req.ProductBatchID)

	productUnit, err := productUnitService.CreateProductUnit(req.ProductID, req.LocationID, req.ProductBatchID, req.Name, req.Quantity, req.UnitPrice, req.Barcode, req.Description, userID)
	if err != nil {
		log.Printf("[PRODUCT_UNIT] Create product unit failed - Product ID: %d, User ID: %d, error: %v", req.ProductID, userID, err)
		statusCode, message := handleProductUnitError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT_UNIT] Create product unit successful - Product Unit: %v, Created by User ID: %d", productUnit, userID)
	return helper.Success(c, 201, "Product unit created successfully", productUnit)
}

func UpdateProductUnit(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_UNIT] Update product unit request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_UNIT] Update product unit failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid product unit ID", err.Error())
	}

	var req UpdateProductUnitRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PRODUCT_UNIT] Update product unit failed - Invalid request body for ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_UNIT] Update product unit failed - User not authenticated for Product Unit ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[PRODUCT_UNIT] Updating product unit with audit - Product Unit ID: %d, User ID: %d, Product Batch ID: %d", idUint, userID, req.ProductBatchID)

	productUnit, err := productUnitService.UpdateProductUnit(uint(idUint), req.ProductID, req.LocationID, req.ProductBatchID, req.Name, req.Quantity, req.UnitPrice, req.Barcode, req.Description, userID)
	if err != nil {
		log.Printf("[PRODUCT_UNIT] Update product unit failed - Product Unit ID: %d, User ID: %d, error: %v", idUint, userID, err)
		statusCode, message := handleProductUnitError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT_UNIT] Update product unit successful - Product Unit: %v, Updated by User ID: %d", productUnit, userID)
	return helper.Success(c, 200, "Product unit updated successfully", productUnit)
}

func DeleteProductUnit(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_UNIT] Delete product unit request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_UNIT] Delete product unit failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid product unit ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_UNIT] Delete product unit failed - User not authenticated for Product Unit ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = productUnitService.DeleteProductUnit(uint(idUint), userID)
	if err != nil {
		log.Printf("[PRODUCT_UNIT] Delete product unit failed - Product Unit ID: %d, error: %v", idUint, err)
		statusCode, message := handleProductUnitError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT_UNIT] Delete product unit successful - Product Unit ID: %d", idUint)
	return helper.Success(c, 200, "Product unit deleted successfully", nil)
}
