package handler

import (
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var brandService = service.NewBrandService()

// handleBrandError converts database errors to user-friendly messages for brand operations
func handleBrandError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors first
	if errMsg == "brand already exists" || errMsg == "brand name already in use" {
		return 409, "Brand name already exists"
	}

	if errMsg == "brand not found" {
		return 404, "Brand not found"
	}

	// Handle PostgreSQL constraint errors as backup
	if strings.Contains(errMsg, "duplicate key value violates unique constraint") &&
		strings.Contains(errMsg, "uni_brands_name") {
		return 409, "Brand name already exists"
	}

	// Default to 500 for other errors
	return 500, "Internal server error"
}

type CreateBrandRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"` // Nullable field
}

type UpdateBrandRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"` // Nullable field
}

func GetBrands(c *fiber.Ctx) error {
	log.Printf("[BRAND] Get all brands request from IP: %s", c.IP())

	brands, err := brandService.GetAllBrands()
	if err != nil {
		log.Printf("[BRAND] Get all brands failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch brands", err.Error())
	}

	log.Printf("[BRAND] Get all brands successful - Found %d brands", len(brands))
	return helper.Success(c, 200, "Success", brands)
}

func GetBrandByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[BRAND] Get brand by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[BRAND] Get brand by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid brand ID", err.Error())
	}

	brand, err := brandService.GetBrandByID(uint(idUint))
	if err != nil {
		log.Printf("[BRAND] Get brand by ID failed - Brand ID: %d not found, error: %v", idUint, err)
		return helper.Fail(c, 404, "Brand not found", err.Error())
	}

	log.Printf("[BRAND] Get brand by ID successful - Brand ID: %d, Name: %s", brand.ID, brand.Name)
	return helper.Success(c, 200, "Success", brand)
}

func CreateBrand(c *fiber.Ctx) error {
	log.Printf("[BRAND] Create brand request from IP: %s", c.IP())

	var req CreateBrandRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[BRAND] Create brand failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[BRAND] Create brand failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[BRAND] Creating brand with audit - User ID: %d", userID)

	brand, err := brandService.CreateBrand(req.Name, req.Description, userID)
	if err != nil {
		log.Printf("[BRAND] Create brand failed - Name: %s, User ID: %d, error: %v", req.Name, userID, err)
		statusCode, message := handleBrandError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[BRAND] Create brand successful - Brand ID: %d, Name: %s, Created by User ID: %d", brand.ID, brand.Name, userID)
	return helper.Success(c, 200, "Brand created successfully", brand)
}

func UpdateBrand(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[BRAND] Update brand request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[BRAND] Update brand failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid brand ID", err.Error())
	}

	var req UpdateBrandRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[BRAND] Update brand failed - Invalid request body for ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[BRAND] Update brand failed - User not authenticated for Brand ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[BRAND] Updating brand with audit - Brand ID: %d, User ID: %d", idUint, userID)

	brand, err := brandService.UpdateBrand(uint(idUint), req.Name, req.Description, userID)
	if err != nil {
		log.Printf("[BRAND] Update brand failed - Brand ID: %d, User ID: %d, error: %v", idUint, userID, err)
		statusCode, message := handleBrandError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[BRAND] Update brand successful - Brand ID: %d, Name: %s, Updated by User ID: %d", brand.ID, brand.Name, userID)
	return helper.Success(c, 200, "Brand updated successfully", brand)
}

func DeleteBrand(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[BRAND] Delete brand request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[BRAND] Delete brand failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid brand ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[BRAND] Delete brand failed - User not authenticated for Brand ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = brandService.DeleteBrand(uint(idUint), userID)
	if err != nil {
		log.Printf("[BRAND] Delete brand failed - Brand ID: %d, error: %v", idUint, err)
		statusCode, message := handleBrandError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[BRAND] Delete brand successful - Brand ID: %d, Deleted by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Brand deleted successfully", nil)
}
