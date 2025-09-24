package handler

import (
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var productUnitTrackService = service.NewProductUnitTrackService()

// handleProductUnitTrackError converts database errors to user-friendly messages for product unit track operations
func handleProductUnitTrackError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors first
	if errMsg == "product unit track not found" {
		return 404, "Product unit track not found"
	}

	if errMsg == "product unit not found" {
		return 404, "Product unit not found"
	}

	if errMsg == "product not found" {
		return 404, "Product not found"
	}

	// Handle PostgreSQL constraint errors as backup
	if strings.Contains(errMsg, "duplicate key value violates unique constraint") {
		return 409, "Product unit track already exists"
	}

	if strings.Contains(errMsg, "foreign key constraint") {
		return 400, "Invalid product unit ID"
	}

	// Default to 500 for other errors
	return 500, "Internal server error"
}

type CreateProductUnitTrackRequest struct {
	ProductUnitID uint    `json:"productUnitId" validate:"required"`
	Description   *string `json:"description"`
}

type UpdateProductUnitTrackRequest struct {
	ProductUnitID uint    `json:"productUnitId"`
	Description   *string `json:"description"`
}

func GetProductUnitTracks(c *fiber.Ctx) error {
	log.Printf("[PRODUCT_UNIT_TRACK] Get all product unit tracks request from IP: %s", c.IP())

	tracks, err := productUnitTrackService.GetAllProductUnitTracks()
	if err != nil {
		log.Printf("[PRODUCT_UNIT_TRACK] Get all product unit tracks failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch product unit tracks", err.Error())
	}

	log.Printf("[PRODUCT_UNIT_TRACK] Get all product unit tracks successful")
	return helper.Success(c, 200, "Success", tracks)
}

func GetProductUnitTracksByProduct(c *fiber.Ctx) error {
	productID := c.Params("productId")
	log.Printf("[PRODUCT_UNIT_TRACK] Get product unit tracks by product request - Product ID: %s from IP: %s", productID, c.IP())

	productIDUint, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_UNIT_TRACK] Get product unit tracks by product failed - Invalid Product ID: %s, error: %v", productID, err)
		return helper.Fail(c, 400, "Invalid product ID", err.Error())
	}

	tracks, err := productUnitTrackService.GetProductUnitTracksByProduct(uint(productIDUint))
	if err != nil {
		log.Printf("[PRODUCT_UNIT_TRACK] Get product unit tracks by product failed - Product ID: %d, error: %v", productIDUint, err)
		statusCode, message := handleProductUnitTrackError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT_UNIT_TRACK] Get product unit tracks by product successful - Product ID: %d", productIDUint)
	return helper.Success(c, 200, "Success", tracks)
}

func GetProductUnitTracksByProductUnit(c *fiber.Ctx) error {
	productUnitID := c.Params("productUnitId")
	log.Printf("[PRODUCT_UNIT_TRACK] Get product unit tracks by product unit request - Product Unit ID: %s from IP: %s", productUnitID, c.IP())

	productUnitIDUint, err := strconv.ParseUint(productUnitID, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_UNIT_TRACK] Get product unit tracks by product unit failed - Invalid Product Unit ID: %s, error: %v", productUnitID, err)
		return helper.Fail(c, 400, "Invalid product unit ID", err.Error())
	}

	tracks, err := productUnitTrackService.GetProductUnitTracksByProductUnit(uint(productUnitIDUint))
	if err != nil {
		log.Printf("[PRODUCT_UNIT_TRACK] Get product unit tracks by product unit failed - Product Unit ID: %d, error: %v", productUnitIDUint, err)
		statusCode, message := handleProductUnitTrackError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT_UNIT_TRACK] Get product unit tracks by product unit successful - Product Unit ID: %d", productUnitIDUint)
	return helper.Success(c, 200, "Success", tracks)
}

func GetProductUnitTrackByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_UNIT_TRACK] Get product unit track by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_UNIT_TRACK] Get product unit track by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid product unit track ID", err.Error())
	}

	track, err := productUnitTrackService.GetProductUnitTrackByID(uint(idUint))
	if err != nil {
		log.Printf("[PRODUCT_UNIT_TRACK] Get product unit track by ID failed - Product Unit Track ID: %d, error: %v", idUint, err)
		statusCode, message := handleProductUnitTrackError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT_UNIT_TRACK] Get product unit track by ID successful - Product Unit Track ID: %d", idUint)
	return helper.Success(c, 200, "Success", track)
}

func CreateProductUnitTrack(c *fiber.Ctx) error {
	log.Printf("[PRODUCT_UNIT_TRACK] Create product unit track request from IP: %s", c.IP())

	var req CreateProductUnitTrackRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[PRODUCT_UNIT_TRACK] Create product unit track failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_UNIT_TRACK] Create product unit track failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[PRODUCT_UNIT_TRACK] Creating product unit track with audit - User ID: %d, Product Unit ID: %d", userID, req.ProductUnitID)

	// Convert *string to string for description
	description := ""
	if req.Description != nil {
		description = *req.Description
	}

	track, err := productUnitTrackService.CreateProductUnitTrack(
		req.ProductUnitID,
		description,
		userID,
	)
	if err != nil {
		log.Printf("[PRODUCT_UNIT_TRACK] Create product unit track failed - Product Unit ID: %d, User ID: %d, error: %v", req.ProductUnitID, userID, err)
		statusCode, message := handleProductUnitTrackError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT_UNIT_TRACK] Create product unit track successful - Track: %v, Created by User ID: %d", track, userID)
	return helper.Success(c, 201, "Product unit track created successfully", track)
}

func DeleteProductUnitTrack(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[PRODUCT_UNIT_TRACK] Delete product unit track request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[PRODUCT_UNIT_TRACK] Delete product unit track failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid product unit track ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[PRODUCT_UNIT_TRACK] Delete product unit track failed - User not authenticated for Track ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = productUnitTrackService.DeleteProductUnitTrack(uint(idUint), userID)
	if err != nil {
		log.Printf("[PRODUCT_UNIT_TRACK] Delete product unit track failed - Track ID: %d, error: %v", idUint, err)
		statusCode, message := handleProductUnitTrackError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[PRODUCT_UNIT_TRACK] Delete product unit track successful - Track ID: %d", idUint)
	return helper.Success(c, 200, "Product unit track deleted successfully", nil)
}
