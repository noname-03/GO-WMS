package auth

import (
    "github.com/gofiber/fiber/v2"
    "myapp/internal/handler"
    "myapp/internal/middleware"
)

func SetupAuthRoutes(router fiber.Router) {
    // Auth routes group
    auth := router.Group("/auth")
    
    // Public auth routes (no JWT required)
    auth.Post("/login", handler.Login)          // POST /api/v1/auth/login
    auth.Post("/register", handler.Register)    // POST /api/v1/auth/register
    
    // Protected auth routes (require JWT)
    authProtected := auth.Group("/", middleware.JWTMiddleware())
    authProtected.Get("/profile", handler.GetProfile)    // GET /api/v1/auth/profile
    authProtected.Put("/profile", handler.UpdateProfile) // PUT /api/v1/auth/profile (future)
    authProtected.Post("/logout", handler.Logout)        // POST /api/v1/auth/logout (future)
}