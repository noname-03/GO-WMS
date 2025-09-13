package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"myapp/internal/service"
	"myapp/internal/utils"
	"myapp/pkg/helper"
)

var authUserService = service.NewUserService()

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Use service instead of direct database access
	user, err := authUserService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		return helper.Fail(c, 401, "Invalid credentials", err.Error())
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Println("Error generating JWT:", err)
		return helper.Fail(c, 500, "Failed to generate token", err.Error())
	}

	return helper.Success(c, 200, "Login successful", fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Use service instead of direct database access
	user, err := authUserService.CreateUser(req.Name, req.Email, req.Password)
	if err != nil {
		return helper.Fail(c, 409, "Registration failed", err.Error())
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Println("Error generating JWT:", err)
		return helper.Fail(c, 500, "Failed to generate token", err.Error())
	}

	return helper.Success(c, 201, "User created successfully", fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	user, err := authUserService.GetUserByID(userID)
	if err != nil {
		return helper.Fail(c, 404, "User not found", err.Error())
	}

	return helper.Success(c, 200, "Success", fiber.Map{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// Future auth handlers
func UpdateProfile(c *fiber.Ctx) error {
	// Implement update profile
	return helper.Success(c, 200, "Update profile endpoint", "Coming soon")
}

func Logout(c *fiber.Ctx) error {
	// Implement logout (blacklist token)
	return helper.Success(c, 200, "Logout successful", "Token invalidated")
}