package handler

import (
	"log"
	"strings"

	"myapp/internal/service"
	"myapp/internal/utils"
	"myapp/pkg/helper"

	"github.com/gofiber/fiber/v2"
)

var authUserService = service.NewUserService()

// handleAuthError converts database errors to user-friendly messages for auth operations
func handleAuthError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors first
	if errMsg == "email already exists" || errMsg == "user already exists" {
		return 409, "Email already exists"
	}

	if errMsg == "invalid credentials" {
		return 401, "Invalid email or password"
	}

	if errMsg == "user not found" {
		return 404, "User not found"
	}

	// Handle PostgreSQL constraint errors as backup
	if strings.Contains(errMsg, "duplicate key value violates unique constraint") &&
		(strings.Contains(errMsg, "users_email_key") || strings.Contains(errMsg, "uni_users_email")) {
		return 409, "Email already exists"
	}

	// Default to 500 for other errors
	return 500, "Internal server error"
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RegisterResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

func Login(c *fiber.Ctx) error {
	log.Printf("[AUTH] Login request received from IP: %s", c.IP())

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[AUTH] Login failed - Invalid request body for email: %s, error: %v", req.Email, err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Use service instead of direct database access
	user, err := authUserService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		log.Printf("[AUTH] Login failed - Authentication failed for email: %s, error: %v", req.Email, err)
		return helper.Fail(c, 401, "Invalid credentials", err.Error())
	}

	log.Printf("[AUTH] User authenticated successfully - ID: %d, Email: %s", user.ID, user.Email)

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Printf("[AUTH] Login failed - JWT generation error for user ID: %d, error: %v", user.ID, err)
		return helper.Fail(c, 500, "Failed to generate token", err.Error())
	}

	log.Printf("[AUTH] Login successful - User ID: %d, Email: %s, Token generated", user.ID, user.Email)

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
	log.Printf("[AUTH] Register request received from IP: %s", c.IP())

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[AUTH] Register failed - Invalid request body, error: %v", err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Use service instead of direct database access
	user, err := authUserService.CreateUser(req.Name, req.Email, req.Password)
	if err != nil {
		log.Printf("[AUTH] Register failed - User creation failed for email: %s, error: %v", req.Email, err)
		statusCode, message := handleAuthError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[AUTH] User created successfully - ID: %d, Email: %s, Name: %s", user.ID, user.Email, user.Name)

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Printf("[AUTH] Register failed - JWT generation error for user ID: %d, error: %v", user.ID, err)
		return helper.Fail(c, 500, "Failed to generate token", err.Error())
	}

	log.Printf("[AUTH] Registration successful - User ID: %d, Email: %s, Token generated", user.ID, user.Email)

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
	log.Printf("[AUTH] Profile request for user ID: %d from IP: %s", userID, c.IP())

	user, err := authUserService.GetUserByID(userID)
	if err != nil {
		log.Printf("[AUTH] Profile failed - User not found for ID: %d, error: %v", userID, err)
		return helper.Fail(c, 404, "User not found", err.Error())
	}

	log.Printf("[AUTH] Profile retrieved successfully for user ID: %d, Email: %s", user.ID, user.Email)

	return helper.Success(c, 200, "Success", fiber.Map{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// Future auth handlers
func UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	log.Printf("[AUTH] Update profile request for user ID: %d from IP: %s", userID, c.IP())
	log.Printf("[AUTH] Update profile - endpoint not yet implemented")

	// Implement update profile
	return helper.Success(c, 200, "Update profile endpoint", "Coming soon")
}

func Logout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	log.Printf("[AUTH] Logout request for user ID: %d from IP: %s", userID, c.IP())
	log.Printf("[AUTH] Logout successful for user ID: %d", userID)

	// Implement logout (blacklist token)
	return helper.Success(c, 200, "Logout successful", "Token invalidated")
}
