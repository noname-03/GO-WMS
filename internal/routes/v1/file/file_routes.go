package file

import (
	"myapp/internal/handler"
	"myapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupFileRoutes(router fiber.Router) {
	files := router.Group("/files")

	// Protected file routes (require JWT)
	files.Use(middleware.JWTMiddleware())

	// IMPORTANT: Specific routes MUST come BEFORE parameterized routes
	// Specific routes (no parameters)
	files.Get("/", handler.GetFiles)                                 // GET /api/v1/files
	files.Get("/deleted", handler.GetDeletedFiles)                   // GET /api/v1/files/deleted
	files.Get("/model/:modelType/:modelId", handler.GetFilesByModel) // GET /api/v1/files/model/product/1

	// Parameterized routes (MUST be at the end)
	files.Get("/:id", handler.GetFileByID)         // GET /api/v1/files/:id
	files.Post("/upload", handler.UploadFile)      // POST /api/v1/files/upload
	files.Put("/:id", handler.UpdateFile)          // PUT /api/v1/files/:id
	files.Put("/:id/restore", handler.RestoreFile) // PUT /api/v1/files/:id/restore
	files.Delete("/:id", handler.DeleteFile)       // DELETE /api/v1/files/:id
}
