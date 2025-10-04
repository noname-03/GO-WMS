package s3test

import (
	"myapp/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupS3TestRoutes(api fiber.Router) {
	s3TestHandler := handler.NewS3TestHandler()

	// S3 test routes
	s3TestGroup := api.Group("/s3-test")
	s3TestGroup.Get("/connection", s3TestHandler.TestS3Connection)
}
