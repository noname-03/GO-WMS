package handler

import (
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var locationService = service.NewLocationService()

// Helper functions for logging
func safeStringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// handleLocationError converts database errors to user-friendly messages for location operations
func handleLocationError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors first
	if errMsg == "location name already exists for this user" {
		return 409, "Location name already exists for this user"
	}

	if errMsg == "location not found" {
		return 404, "Location not found"
	}

	if errMsg == "user not found" {
		return 404, "User not found"
	}

	if errMsg == "invalid location type. Must be 'gudang' or 'reseller'" {
		return 400, "Invalid location type. Must be 'gudang' or 'reseller'"
	}

	// Handle PostgreSQL constraint errors as backup
	if strings.Contains(errMsg, "duplicate key value violates unique constraint") {
		return 409, "Location already exists"
	}

	if strings.Contains(errMsg, "foreign key constraint") {
		return 400, "Invalid user ID"
	}

	// Default to 500 for other errors
	return 500, "Internal server error"
}

type CreateLocationRequest struct {
	UserID      uint    `json:"userId" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Address     *string `json:"address"`
	PhoneNumber *string `json:"phoneNumber"`
	Type        string  `json:"type" validate:"required,oneof=gudang reseller"`
}

type UpdateLocationRequest struct {
	UserID      uint    `json:"userId"`
	Name        *string `json:"name"`
	Address     *string `json:"address"`
	PhoneNumber *string `json:"phoneNumber"`
	Type        *string `json:"type"`
}

func GetLocations(c *fiber.Ctx) error {
	log.Printf("[LOCATION] Get all locations request from IP: %s", c.IP())

	locations, err := locationService.GetAllLocations()
	if err != nil {
		log.Printf("[LOCATION] Get all locations failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch locations", err.Error())
	}

	log.Printf("[LOCATION] Get all locations successful")
	return helper.Success(c, 200, "Success", locations)
}

func GetLocationsByUser(c *fiber.Ctx) error {
	userID := c.Params("userId")
	log.Printf("[LOCATION] Get locations by user request - User ID: %s from IP: %s", userID, c.IP())

	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		log.Printf("[LOCATION] Get locations by user failed - Invalid User ID: %s, error: %v", userID, err)
		return helper.Fail(c, 400, "Invalid user ID", err.Error())
	}

	locations, err := locationService.GetLocationsByUser(uint(userIDUint))
	if err != nil {
		log.Printf("[LOCATION] Get locations by user failed - User ID: %d, error: %v", userIDUint, err)
		statusCode, message := handleLocationError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[LOCATION] Get locations by user successful - User ID: %d", userIDUint)
	return helper.Success(c, 200, "Success", locations)
}

func GetLocationsByType(c *fiber.Ctx) error {
	locationType := c.Params("type")
	log.Printf("[LOCATION] Get locations by type request - Type: %s from IP: %s", locationType, c.IP())

	locations, err := locationService.GetLocationsByType(locationType)
	if err != nil {
		log.Printf("[LOCATION] Get locations by type failed - Type: %s, error: %v", locationType, err)
		statusCode, message := handleLocationError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[LOCATION] Get locations by type successful - Type: %s", locationType)
	return helper.Success(c, 200, "Success", locations)
}

func GetLocationByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[LOCATION] Get location by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[LOCATION] Get location by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid location ID", err.Error())
	}

	location, err := locationService.GetLocationByID(uint(idUint))
	if err != nil {
		log.Printf("[LOCATION] Get location by ID failed - Location ID: %d, error: %v", idUint, err)
		statusCode, message := handleLocationError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[LOCATION] Get location by ID successful - Location ID: %d", idUint)
	return helper.Success(c, 200, "Success", location)
}

func CreateLocation(c *fiber.Ctx) error {
	log.Printf("[LOCATION] Create location request from IP: %s", c.IP())

	var req CreateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[LOCATION] Create location failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token for audit trail
	createdByUserID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[LOCATION] Create location failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[LOCATION] Creating location with audit - User ID: %d, Name: %s, Phone: %s, Type: %s", createdByUserID, req.Name, safeStringPtr(req.PhoneNumber), req.Type)

	location, err := locationService.CreateLocation(req.UserID, req.Name, req.Address, req.PhoneNumber, req.Type, createdByUserID)
	if err != nil {
		log.Printf("[LOCATION] Create location failed - User ID: %d, Created by User ID: %d, error: %v", req.UserID, createdByUserID, err)
		statusCode, message := handleLocationError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[LOCATION] Create location successful - Created by User ID: %d", createdByUserID)
	return helper.Success(c, 201, "Location created successfully", location)
}

func UpdateLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[LOCATION] Update location request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[LOCATION] Update location failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid location ID", err.Error())
	}

	var req UpdateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[LOCATION] Update location failed - Invalid request body for ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token for audit trail
	updatedByUserID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[LOCATION] Update location failed - User not authenticated for Location ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[LOCATION] Updating location - ID: %d, Name: %s, Phone: %s, Type: %s, Updated by User ID: %d",
		idUint,
		safeStringPtr(req.Name),
		safeStringPtr(req.PhoneNumber),
		safeStringPtr(req.Type),
		updatedByUserID)

	location, err := locationService.UpdateLocation(uint(idUint), req.UserID, req.Name, req.Address, req.PhoneNumber, req.Type, updatedByUserID)
	if err != nil {
		log.Printf("[LOCATION] Update location failed - Location ID: %d, Updated by User ID: %d, error: %v", idUint, updatedByUserID, err)
		statusCode, message := handleLocationError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[LOCATION] Update location successful - Updated by User ID: %d", updatedByUserID)
	return helper.Success(c, 200, "Location updated successfully", location)
}

func DeleteLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[LOCATION] Delete location request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[LOCATION] Delete location failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid location ID", err.Error())
	}

	// Get user ID from JWT token for audit trail
	deletedByUserID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[LOCATION] Delete location failed - User not authenticated for Location ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = locationService.DeleteLocation(uint(idUint), deletedByUserID)
	if err != nil {
		log.Printf("[LOCATION] Delete location failed - Location ID: %d, error: %v", idUint, err)
		statusCode, message := handleLocationError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[LOCATION] Delete location successful - Location ID: %d", idUint)
	return helper.Success(c, 200, "Location deleted successfully", nil)
}

func GetDeletedLocations(c *fiber.Ctx) error {
	log.Printf("[LOCATION] Get deleted locations request from IP: %s", c.IP())

	locations, err := locationService.GetDeletedLocations()
	if err != nil {
		log.Printf("[LOCATION] Get deleted locations failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch deleted locations", err.Error())
	}

	log.Printf("[LOCATION] Get deleted locations successful")
	return helper.Success(c, 200, "Success", locations)
}

func RestoreLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[LOCATION] Restore location request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[LOCATION] Restore location failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid location ID", err.Error())
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[LOCATION] Restore location failed - User not authenticated for Location ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	location, err := locationService.RestoreLocation(uint(idUint), userID)
	if err != nil {
		log.Printf("[LOCATION] Restore location failed - Location ID: %d, error: %v", idUint, err)
		statusCode, message := handleLocationError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[LOCATION] Restore location successful - Location ID: %d, Restored by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Location restored successfully", location)
}
