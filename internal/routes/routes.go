package routes

import (
    "github.com/gofiber/fiber/v2"
    v1 "myapp/internal/routes/v1"
    // Import v2 untuk future development
    // v2 "myapp/internal/routes/v2"
)

func SetupRoutes(app *fiber.App) {
    // Setup API versioning
    v1.SetupV1Routes(app)
    
    // Future API versions
    // v2.SetupV2Routes(app)
    
    // Global routes (not versioned)
    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Welcome to GO-WMS API",
            "versions": []string{
                "/api/v1",
                // "/api/v2", // Future
            },
            "docs": "/api/v1/health",
        })
    })
    
    // Global health check
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "OK",
            "service": "GO-WMS",
            "message": "Service is running",
        })
    })
}