package handler

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"log"
)

var roleService = service.NewRoleService()

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func GetRoles(c *fiber.Ctx) error {
	roles, err := roleService.GetAllRoles()
	if err != nil {
		log.Println("Error fetching roles:", err)
		return helper.Fail(c, 500, "Failed to fetch roles", err.Error())
	}

	log.Printf("Found %d roles in database\n", len(roles))
	return helper.Success(c, 200, "Success", roles)
}

func GetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return helper.Fail(c, 400, "Invalid role ID", err.Error())
	}

	role, err := roleService.GetRoleByID(uint(idUint))
	if err != nil {
		return helper.Fail(c, 404, "Role not found", err.Error())
	}

	return helper.Success(c, 200, "Success", role)
}

func CreateRole(c *fiber.Ctx) error {
	var req CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	role, err := roleService.CreateRole(req.Name, req.Description)
	if err != nil {
		log.Println("Error creating role:", err)
		return helper.Fail(c, 500, "Failed to create role", err.Error())
	}

	return helper.Success(c, 201, "Role created successfully", role)
}

func UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return helper.Fail(c, 400, "Invalid role ID", err.Error())
	}

	var req UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	role, err := roleService.UpdateRole(uint(idUint), req.Name, req.Description)
	if err != nil {
		log.Println("Error updating role:", err)
		return helper.Fail(c, 400, "Failed to update role", err.Error())
	}

	return helper.Success(c, 200, "Role updated successfully", role)
}

func DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return helper.Fail(c, 400, "Invalid role ID", err.Error())
	}

	err = roleService.DeleteRole(uint(idUint))
	if err != nil {
		log.Println("Error deleting role:", err)
		return helper.Fail(c, 400, "Failed to delete role", err.Error())
	}

	return helper.Success(c, 200, "Role deleted successfully", nil)
}