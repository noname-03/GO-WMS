package handler

import (
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var inventoryService = service.NewInventoryService()

// GetInventoryStock returns inventory stock information with optional filters
func GetInventoryStock(c *fiber.Ctx) error {
	log.Printf("[INVENTORY] Get inventory stock request from IP: %s", c.IP())

	// Parse query parameters
	brandIDStr := c.Query("brandId")
	categoryIDStr := c.Query("categoryId")
	productIDStr := c.Query("productId")
	productBatchIDStr := c.Query("productBatchId")
	locationIDStr := c.Query("locationId")

	var brandID, categoryID, productID, productBatchID, locationID *uint

	// Parse brandId if provided
	if brandIDStr != "" {
		id, err := strconv.ParseUint(brandIDStr, 10, 32)
		if err != nil {
			log.Printf("[INVENTORY] Get inventory stock failed - Invalid brandId: %s, error: %v", brandIDStr, err)
			return helper.Fail(c, 400, "Invalid brandId", err.Error())
		}
		brandIDUint := uint(id)
		brandID = &brandIDUint
	}

	// Parse categoryId if provided
	if categoryIDStr != "" {
		id, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err != nil {
			log.Printf("[INVENTORY] Get inventory stock failed - Invalid categoryId: %s, error: %v", categoryIDStr, err)
			return helper.Fail(c, 400, "Invalid categoryId", err.Error())
		}
		categoryIDUint := uint(id)
		categoryID = &categoryIDUint
	}

	// Parse productId if provided
	if productIDStr != "" {
		id, err := strconv.ParseUint(productIDStr, 10, 32)
		if err != nil {
			log.Printf("[INVENTORY] Get inventory stock failed - Invalid productId: %s, error: %v", productIDStr, err)
			return helper.Fail(c, 400, "Invalid productId", err.Error())
		}
		productIDUint := uint(id)
		productID = &productIDUint
	}

	// Parse productBatchId if provided
	if productBatchIDStr != "" {
		id, err := strconv.ParseUint(productBatchIDStr, 10, 32)
		if err != nil {
			log.Printf("[INVENTORY] Get inventory stock failed - Invalid productBatchId: %s, error: %v", productBatchIDStr, err)
			return helper.Fail(c, 400, "Invalid productBatchId", err.Error())
		}
		productBatchIDUint := uint(id)
		productBatchID = &productBatchIDUint
	}

	// Parse locationId if provided
	if locationIDStr != "" {
		id, err := strconv.ParseUint(locationIDStr, 10, 32)
		if err != nil {
			log.Printf("[INVENTORY] Get inventory stock failed - Invalid locationId: %s, error: %v", locationIDStr, err)
			return helper.Fail(c, 400, "Invalid locationId", err.Error())
		}
		locationIDUint := uint(id)
		locationID = &locationIDUint
	}

	log.Printf("[INVENTORY] Fetching inventory stock - Brand ID: %v, Category ID: %v, Product ID: %v, Product Batch ID: %v, Location ID: %v",
		brandID, categoryID, productID, productBatchID, locationID)

	results, err := inventoryService.GetInventoryStock(brandID, categoryID, productID, productBatchID, locationID)
	if err != nil {
		log.Printf("[INVENTORY] Get inventory stock failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch inventory stock", err.Error())
	}

	log.Printf("[INVENTORY] Get inventory stock successful")
	return helper.Success(c, 200, "Success", results)
}
