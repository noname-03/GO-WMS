package handler

import (
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var productService = service.NewProductService()

// handleProductError converts database errors to user-friendly messages for product operations
func handleProductError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors first
	if errMsg == "product already exists for this category" {
		return 409, "Product name already exists for this category"
	}

	if errMsg == "product not found" {
		return 404, "Product not found"
	}

	if errMsg == "category not found" {
		return 404, "Category not found"
	}

	// Handle PostgreSQL constraint errors as backup
	if strings.Contains(errMsg, "duplicate key value violates unique constraint") {
		return 409, "Product name already exists for this category"
	}

	if strings.Contains(errMsg, "foreign key constraint") {
		return 400, "Invalid category ID"
	}

	// Default to 500 for other errors
	return 500, "Internal server error"
}

type CreateProductRequest struct {
	CategoryID  uint    `json:"categoryId" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"` // Nullable field
}

type UpdateProductRequest struct {
	CategoryID  uint    `json:"categoryId"`
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"` // Nullable field
}

func GetProducts(c *fiber.Ctx) error {
	log.Printf("[PRODUCT] Get all products request from IP: %s", c.IP())

	products, err := productService.GetAllProducts()
	if err != nil {
		log.Printf("[PRODUCT] Get all products failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch products", err.Error())
	}

	log.Printf("[PRODUCT] Get all products successful")
	return helper.Success(c, 200, "Success", products)
}

func GetProductsByCategory(c *fiber.Ctx) error {
	categoryID := c.Params("categoryId")
	log.Printf("[PRODUCT] Get products by category request - Category ID: %s from IP: %s", categoryID, c.IP())

	categoryIDUint, err := strconv.ParseUint(categoryID, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT] Get products by category failed - Invalid Category ID: %s, error: %v", categoryID, err)
		return helper.Fail(c, 400, "Invalid category ID", err.Error())
	}

	products, err := productService.GetProductsByCategory(uint(categoryIDUint))
	if err != nil {
		log.Printf("[PRODUCT] Get products by category failed - Category ID: %d, error: %v", categoryIDUint, err)
		statusCode, message := handleProductError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT] Get products by category successful")
	return helper.Success(c, 200, "Success", products)
}

func GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT] Get product by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT] Get product by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid product ID", err.Error())
	}

	product, err := productService.GetProductByID(uint(idUint))
	if err != nil {
		log.Printf("[PRODUCT] Get product by ID failed - Product ID: %d not found, error: %v", idUint, err)
		return helper.Fail(c, 404, "Product not found", err.Error())
	}

	log.Printf("[PRODUCT] Get product by ID successful")
	return helper.Success(c, 200, "Success", product)
}

func CreateProduct(c *fiber.Ctx) error {
	log.Printf("[PRODUCT] Create product request from IP: %s", c.IP())

	var req CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PRODUCT] Create product failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT] Create product failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[PRODUCT] Creating product with audit - User ID: %d, Category ID: %d", userID, req.CategoryID)

	product, err := productService.CreateProduct(req.CategoryID, req.Name, req.Description, userID)
	if err != nil {
		log.Printf("[PRODUCT] Create product failed - Name: %s, Category ID: %d, User ID: %d, error: %v", req.Name, req.CategoryID, userID, err)
		statusCode, message := handleProductError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT] Create product successful - Name: %s, Category ID: %d, Created by User ID: %d", req.Name, req.CategoryID, userID)
	return helper.Success(c, 201, "Product created successfully", product)
}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT] Update product request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT] Update product failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid product ID", err.Error())
	}

	var req UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PRODUCT] Update product failed - Invalid request body for ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT] Update product failed - User not authenticated for Product ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[PRODUCT] Updating product with audit - Product ID: %d, User ID: %d", idUint, userID)

	product, err := productService.UpdateProduct(uint(idUint), req.CategoryID, req.Name, req.Description, userID)
	if err != nil {
		log.Printf("[PRODUCT] Update product failed - Product ID: %d, User ID: %d, error: %v", idUint, userID, err)
		statusCode, message := handleProductError(err)
		log.Printf("[PRODUCT] Update product failed - Error: %s", err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT] Update product successful - Product ID: %d, Updated by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Product updated successfully", product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT] Delete product request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT] Delete product failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid product ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT] Delete product failed - User not authenticated for Product ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = productService.DeleteProduct(uint(idUint), userID)
	if err != nil {
		log.Printf("[PRODUCT] Delete product failed - Product ID: %d, error: %v", idUint, err)
		statusCode, message := handleProductError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT] Delete product successful")
	return helper.Success(c, 200, "Product deleted successfully", nil)
}
