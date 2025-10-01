package handler

import (
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var categoryService = service.NewCategoryService()

// handleCategoryError converts database errors to user-friendly messages for category operations
func handleCategoryError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors first
	if errMsg == "category already exists for this brand" || errMsg == "category name already in use for this brand" {
		return 409, "Category name already exists for this brand"
	}

	if errMsg == "category not found" {
		return 404, "Category not found"
	}

	if errMsg == "brand not found" {
		return 404, "Brand not found"
	}

	// Handle PostgreSQL constraint errors as backup
	if strings.Contains(errMsg, "duplicate key value violates unique constraint") {
		return 409, "Category name already exists for this brand"
	}

	if strings.Contains(errMsg, "foreign key constraint") {
		return 400, "Invalid brand ID"
	}

	// Default to 500 for other errors
	return 500, "Internal server error"
}

type CreateCategoryRequest struct {
	BrandID     uint    `json:"brandId" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"` // Nullable field
}

type UpdateCategoryRequest struct {
	BrandID     uint    `json:"brandId"`
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"` // Nullable field
}

func GetCategories(c *fiber.Ctx) error {
	log.Printf("[CATEGORY] Get all categories request from IP: %s", c.IP())

	categories, err := categoryService.GetAllCategories()
	if err != nil {
		log.Printf("[CATEGORY] Get all categories failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch categories", err.Error())
	}

	log.Printf("[CATEGORY] Get all categories successful")
	return helper.Success(c, 200, "Success", categories)
}

func GetCategoriesByBrand(c *fiber.Ctx) error {
	brandID := c.Params("brandId")
	log.Printf("[CATEGORY] Get categories by brand request - Brand ID: %s from IP: %s", brandID, c.IP())

	brandIDUint, err := strconv.ParseUint(brandID, 10, 32)
	if err != nil {
		log.Printf("[CATEGORY] Get categories by brand failed - Invalid Brand ID: %s, error: %v", brandID, err)
		return helper.Fail(c, 400, "Invalid brand ID", err.Error())
	}

	categories, err := categoryService.GetCategoriesByBrand(uint(brandIDUint))
	if err != nil {
		log.Printf("[CATEGORY] Get categories by brand failed - Brand ID: %d, error: %v", brandIDUint, err)
		statusCode, message := handleCategoryError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[CATEGORY] Get categories by brand successful - Brand ID: %d", brandIDUint)
	return helper.Success(c, 200, "Success", categories)
}

func GetCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[CATEGORY] Get category by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[CATEGORY] Get category by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid category ID", err.Error())
	}

	category, err := categoryService.GetCategoryByID(uint(idUint))
	if err != nil {
		log.Printf("[CATEGORY] Get category by ID failed - Category ID: %d not found, error: %v", idUint, err)
		return helper.Fail(c, 404, "Category not found", err.Error())
	}

	log.Printf("[CATEGORY] Get category by ID successful")
	return helper.Success(c, 200, "Success", category)
}

func CreateCategory(c *fiber.Ctx) error {
	log.Printf("[CATEGORY] Create category request from IP: %s", c.IP())

	var req CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[CATEGORY] Create category failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[CATEGORY] Create category failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[CATEGORY] Creating category with audit - User ID: %d, Brand ID: %d", userID, req.BrandID)

	category, err := categoryService.CreateCategory(req.BrandID, req.Name, req.Description, userID)
	if err != nil {
		log.Printf("[CATEGORY] Create category failed - Name: %s, Brand ID: %d, User ID: %d, error: %v", req.Name, req.BrandID, userID, err)
		statusCode, message := handleCategoryError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[CATEGORY] Create category successful")
	return helper.Success(c, 201, "Category created successfully", category)
}

func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[CATEGORY] Update category request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[CATEGORY] Update category failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid category ID", err.Error())
	}

	var req UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[CATEGORY] Update category failed - Invalid request body for ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[CATEGORY] Update category failed - User not authenticated for Category ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[CATEGORY] Updating category with audit - Category ID: %d, User ID: %d", idUint, userID)

	category, err := categoryService.UpdateCategory(uint(idUint), req.BrandID, req.Name, req.Description, userID)
	if err != nil {
		log.Printf("[CATEGORY] Update category failed - Category ID: %d, User ID: %d, error: %v", idUint, userID, err)
		statusCode, message := handleCategoryError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[CATEGORY] Update category successful")
	return helper.Success(c, 200, "Category updated successfully", category)
}

func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[CATEGORY] Delete category request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[CATEGORY] Delete category failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid category ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[CATEGORY] Delete category failed - User not authenticated for Category ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = categoryService.DeleteCategory(uint(idUint), userID)
	if err != nil {
		log.Printf("[CATEGORY] Delete category failed - Category ID: %d, error: %v", idUint, err)
		statusCode, message := handleCategoryError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[CATEGORY] Delete category successful - Category ID: %d, Deleted by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Category deleted successfully", nil)
}

func GetDeletedCategories(c *fiber.Ctx) error {
	log.Printf("[CATEGORY] Get deleted categories request from IP: %s", c.IP())

	categories, err := categoryService.GetDeletedCategories()
	if err != nil {
		log.Printf("[CATEGORY] Get deleted categories failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch deleted categories", err.Error())
	}

	log.Printf("[CATEGORY] Get deleted categories successful")
	return helper.Success(c, 200, "Success", categories)
}

func RestoreCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[CATEGORY] Restore category request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[CATEGORY] Restore category failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid category ID", err.Error())
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[CATEGORY] Restore category failed - User not authenticated for Category ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	category, err := categoryService.RestoreCategory(uint(idUint), userID)
	if err != nil {
		log.Printf("[CATEGORY] Restore category failed - Category ID: %d, error: %v", idUint, err)
		statusCode, message := handleCategoryError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[CATEGORY] Restore category successful - Category ID: %d, Restored by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Category restored successfully", category)
}
