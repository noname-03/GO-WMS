package handler

import (
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var roleService = service.NewRoleService()

// handleRoleError converts database errors to user-friendly messages for role operations
func handleRoleError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors first
	if errMsg == "role already exists" || errMsg == "role name already in use" {
		return 409, "Role name already exists"
	}

	if errMsg == "role not found" {
		return 404, "Role not found"
	}

	// Handle PostgreSQL constraint errors as backup
	if strings.Contains(errMsg, "duplicate key value violates unique constraint") &&
		strings.Contains(errMsg, "uni_roles_name") {
		return 409, "Role name already exists"
	}

	// Default to 500 for other errors
	return 500, "Internal server error"
}

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func GetRoles(c *fiber.Ctx) error {
	log.Printf("[ROLE] Get all roles request from IP: %s", c.IP())

	roles, err := roleService.GetAllRoles()
	if err != nil {
		log.Printf("[ROLE] Get all roles failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch roles", err.Error())
	}

	log.Printf("[ROLE] Get all roles successful - Found %d roles", len(roles))
	return helper.Success(c, 200, "Success", roles)
}

func GetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[ROLE] Get role by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[ROLE] Get role by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid role ID", err.Error())
	}

	role, err := roleService.GetRoleByID(uint(idUint))
	if err != nil {
		log.Printf("[ROLE] Get role by ID failed - Role ID: %d not found, error: %v", idUint, err)
		return helper.Fail(c, 404, "Role not found", err.Error())
	}

	log.Printf("[ROLE] Get role by ID successful - Role ID: %d, Name: %s", role.ID, role.Name)
	return helper.Success(c, 200, "Success", role)
}

func CreateRole(c *fiber.Ctx) error {
	log.Printf("[ROLE] Create role request from IP: %s", c.IP())

	var req CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ROLE] Create role failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[ROLE] Create role failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[ROLE] Creating role with audit - User ID: %d", userID)

	role, err := roleService.CreateRole(req.Name, req.Description, userID)
	if err != nil {
		log.Printf("[ROLE] Create role failed - Name: %s, User ID: %d, error: %v", req.Name, userID, err)
		statusCode, message := handleRoleError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[ROLE] Create role successful - Role ID: %d, Name: %s, Created by User ID: %d", role.ID, role.Name, userID)
	return helper.Success(c, 201, "Role created successfully", role)
}

func UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[ROLE] Update role request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[ROLE] Update role failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid role ID", err.Error())
	}

	var req UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ROLE] Update role failed - Invalid request body for ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[ROLE] Update role failed - User not authenticated for Role ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[ROLE] Updating role with audit - Role ID: %d, User ID: %d", idUint, userID)

	role, err := roleService.UpdateRole(uint(idUint), req.Name, req.Description, userID)
	if err != nil {
		log.Printf("[ROLE] Update role failed - Role ID: %d, User ID: %d, error: %v", idUint, userID, err)
		statusCode, message := handleRoleError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[ROLE] Update role successful - Role ID: %d, Name: %s, Updated by User ID: %d", role.ID, role.Name, userID)
	return helper.Success(c, 200, "Role updated successfully", role)
}

func DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[ROLE] Delete role request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[ROLE] Delete role failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid role ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[ROLE] Delete role failed - User not authenticated for Role ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = roleService.DeleteRole(uint(idUint), userID)
	if err != nil {
		log.Printf("[ROLE] Delete role failed - Role ID: %d, error: %v", idUint, err)
		statusCode, message := handleRoleError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[ROLE] Delete role successful - Role ID: %d, Deleted by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "Role deleted successfully", nil)
}
