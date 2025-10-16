package handler

import (
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"

	"github.com/gofiber/fiber/v2"
)

var resellerService = service.NewResellerService()

func GetAllResellers(c *fiber.Ctx) error {
	log.Printf("[RESELLER] Get all resellers request from IP: %s", c.IP())

	resellers, err := resellerService.GetAllResellers()
	if err != nil {
		log.Printf("[RESELLER] Get all resellers failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch resellers", err.Error())
	}

	log.Printf("[RESELLER] Get all resellers successful")
	return helper.Success(c, 200, "Success", resellers)
}
