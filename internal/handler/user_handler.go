package handler

import (
    "github.com/gofiber/fiber/v2"
    "myapp/internal/service"
    "myapp/pkg/helper"
    "strconv"
    "log"
)

var userService = service.NewUserService()

func GetUsers(c *fiber.Ctx) error {
    users, err := userService.GetAllUsers()
    if err != nil {
        log.Println("Error fetching users:", err)
        return helper.Fail(c, 500, "Failed to fetch users", err.Error())
    }
    
    log.Printf("Found %d users in database\n", len(users))
    return helper.Success(c, 200, "Success", users)
}

func GetUsersMinimal(c *fiber.Ctx) error {
    users, err := userService.GetUsersMinimal()
    if err != nil {
        return helper.Fail(c, 500, "Failed to fetch users", err.Error())
    }
    
    log.Printf("Found %d users in database\n", len(users))
    return helper.Success(c, 200, "Success", users)
}

func GetUsersRaw(c *fiber.Ctx) error {
    users, err := userService.GetUsersWithStats()
    if err != nil {
        return helper.Fail(c, 500, "Failed to fetch users", err.Error())
    }

    log.Printf("Found %d users with raw SQL\n", len(users))
    return helper.Success(c, 200, "Success", users)
}

func GetUserByIDRaw(c *fiber.Ctx) error {
    id := c.Params("id")
    idUint, err := strconv.ParseUint(id, 10, 32)
    if err != nil {
        return helper.Fail(c, 400, "Invalid user ID", err.Error())
    }

    user, err := userService.GetUserByID(uint(idUint))
    if err != nil {
        return helper.Fail(c, 404, "User not found", err.Error())
    }

    return helper.Success(c, 200, "Success", user)
}

func GetUsersFromRepository(c *fiber.Ctx) error {
    users, err := userService.GetUsersWithRawSQL()
    if err != nil {
        return helper.Fail(c, 500, "Failed to fetch users", err.Error())
    }
    
    return helper.Success(c, 200, "Success", users)
}

func SearchUsers(c *fiber.Ctx) error {
    keyword := c.Query("q", "")
    limitStr := c.Query("limit", "10")
    offsetStr := c.Query("offset", "0")
    
    limit, _ := strconv.Atoi(limitStr)
    offset, _ := strconv.Atoi(offsetStr)
    
    users, err := userService.SearchUsers(keyword, limit, offset)
    if err != nil {
        return helper.Fail(c, 500, "Search failed", err.Error())
    }
    
    return helper.Success(c, 200, "Success", users)
}

func GetUserStats(c *fiber.Ctx) error {
    stats, err := userService.GetUsersStats()
    if err != nil {
        return helper.Fail(c, 500, "Failed to fetch stats", err.Error())
    }
    
    return helper.Success(c, 200, "Success", stats)
}