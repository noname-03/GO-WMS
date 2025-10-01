package handler

import (
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var userService = service.NewUserService()

func GetUsers(c *fiber.Ctx) error {
	log.Printf("[USER] Get all users request from IP: %s", c.IP())

	users, err := userService.GetAllUsers()
	if err != nil {
		log.Printf("[USER] Get all users failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch users", err.Error())
	}

	log.Printf("[USER] Get all users successful - Found %d users", len(users))
	return helper.Success(c, 200, "Success", users)
}

func GetUsersMinimal(c *fiber.Ctx) error {
	log.Printf("[USER] Get users minimal request from IP: %s", c.IP())

	users, err := userService.GetUsersMinimal()
	if err != nil {
		log.Printf("[USER] Get users minimal failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch users", err.Error())
	}

	log.Printf("[USER] Get users minimal successful - Found %d users", len(users))
	return helper.Success(c, 200, "Success", users)
}

func GetUsersRaw(c *fiber.Ctx) error {
	log.Printf("[USER] Get users with stats request from IP: %s", c.IP())

	users, err := userService.GetUsersWithStats()
	if err != nil {
		log.Printf("[USER] Get users with stats failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch users", err.Error())
	}

	log.Printf("[USER] Get users with stats successful - Found %d users", len(users))
	return helper.Success(c, 200, "Success", users)
}

func GetUserByIDRaw(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[USER] Get user by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[USER] Get user by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid user ID", err.Error())
	}

	user, err := userService.GetUserByID(uint(idUint))
	if err != nil {
		log.Printf("[USER] Get user by ID failed - User ID: %d not found, error: %v", idUint, err)
		return helper.Fail(c, 404, "User not found", err.Error())
	}

	log.Printf("[USER] Get user by ID successful - User ID: %d, Email: %s", user.ID, user.Email)
	return helper.Success(c, 200, "Success", user)
}

func GetUsersFromRepository(c *fiber.Ctx) error {
	log.Printf("[USER] Get users from repository request from IP: %s", c.IP())

	users, err := userService.GetUsersWithRawSQL()
	if err != nil {
		log.Printf("[USER] Get users from repository failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch users", err.Error())
	}

	log.Printf("[USER] Get users from repository successful - Found %d users", len(users))
	return helper.Success(c, 200, "Success", users)
}

func SearchUsers(c *fiber.Ctx) error {
	keyword := c.Query("q", "")
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	log.Printf("[USER] Search users request - keyword: '%s', limit: %s, offset: %s from IP: %s", keyword, limitStr, offsetStr, c.IP())

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	users, err := userService.SearchUsers(keyword, limit, offset)
	if err != nil {
		log.Printf("[USER] Search users failed - keyword: '%s', error: %v", keyword, err)
		return helper.Fail(c, 500, "Search failed", err.Error())
	}

	log.Printf("[USER] Search users successful - keyword: '%s', found %d users", keyword, len(users))
	return helper.Success(c, 200, "Success", users)
}

func GetUserStats(c *fiber.Ctx) error {
	log.Printf("[USER] Get user stats request from IP: %s", c.IP())

	stats, err := userService.GetUsersStats()
	if err != nil {
		log.Printf("[USER] Get user stats failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch stats", err.Error())
	}

	log.Printf("[USER] Get user stats successful - Total: %d, Active: %d, Deleted: %d", stats.TotalUsers, stats.ActiveUsers, stats.DeletedUsers)
	return helper.Success(c, 200, "Success", stats)
}

func GetDeletedUsers(c *fiber.Ctx) error {
	log.Printf("[USER] Get deleted users request from IP: %s", c.IP())

	users, err := userService.GetDeletedUsers()
	if err != nil {
		log.Printf("[USER] Get deleted users failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch deleted users", err.Error())
	}

	log.Printf("[USER] Get deleted users successful - Found %d deleted users", len(users))
	return helper.Success(c, 200, "Success", users)
}

func RestoreUser(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[USER] Restore user request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[USER] Restore user failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid user ID", err.Error())
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[USER] Restore user failed - User not authenticated for User ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	user, err := userService.RestoreUser(uint(idUint), userID)
	if err != nil {
		log.Printf("[USER] Restore user failed - User ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 500, "Failed to restore user", err.Error())
	}

	log.Printf("[USER] Restore user successful - User ID: %d, Restored by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "User restored successfully", user)
}
