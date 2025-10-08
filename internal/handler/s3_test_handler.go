package handler

import (
	"log"
	"myapp/pkg/helper/s3"

	"github.com/gofiber/fiber/v2"
)

type S3TestHandler struct {
	s3Client *s3.S3Client
}

func NewS3TestHandler() *S3TestHandler {
	s3Client := s3.NewS3Client()
	if s3Client == nil {
		log.Printf("❌ Failed to initialize S3 client")
		return &S3TestHandler{s3Client: nil}
	}
	return &S3TestHandler{s3Client: s3Client}
}

// TestS3Connection tests S3 connection by listing buckets or checking bucket
func (h *S3TestHandler) TestS3Connection(c *fiber.Ctx) error {
	if h.s3Client == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "S3 client not initialized",
			"error":   "S3 configuration error",
		})
	}

	// Test connection by checking if we can access the bucket
	bucketName := h.s3Client.GetBucketName()

	// Try to list objects in bucket (this will test authentication and bucket access)
	err := h.s3Client.TestConnection()
	if err != nil {
		log.Printf("❌ S3 Connection Test Failed: %v", err)
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"success":    false,
			"message":    "S3 connection failed",
			"bucket":     bucketName,
			"endpoint":   h.s3Client.GetEndpoint(),
			"error":      err.Error(),
			"suggestion": "Check your S3 credentials, endpoint, and bucket permissions",
		})
	}

	log.Printf("✅ S3 Connection Test Successful!")
	return c.JSON(fiber.Map{
		"success":  true,
		"message":  "S3 connection successful",
		"bucket":   bucketName,
		"endpoint": h.s3Client.GetEndpoint(),
		"config": fiber.Map{
			"use_ssl":    h.s3Client.IsSSLEnabled(),
			"path_style": h.s3Client.IsPathStyleEnabled(),
		},
	})
}
