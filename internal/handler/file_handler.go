package handler

import (
	"io"
	"log"
	"myapp/internal/service"
	"myapp/pkg/helper"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// getFileService creates a new FileService instance
// This ensures fresh S3Client initialization
func getFileService() *service.FileService {
	return service.NewFileService()
}

// Request structs
type UploadFileRequest struct {
	ModelType string `json:"modelType" form:"modelType" form:"ModelType"`
	ModelID   uint   `json:"modelId" form:"modelId" form:"ModelId" form:"modelID"`
}

type UpdateFileRequest struct {
	ModelType string `json:"modelType"`
	ModelID   uint   `json:"modelId"`
	Ext       string `json:"ext"`
}

// handleFileError converts database errors to user-friendly messages for file operations
func handleFileError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}

	errMsg := err.Error()

	// Handle specific application errors
	if strings.Contains(errMsg, "model not found") {
		return 404, "Related model not found"
	}
	if strings.Contains(errMsg, "file not found") {
		return 404, "File not found"
	}
	if strings.Contains(errMsg, "invalid model type") {
		return 400, "Invalid model type. Allowed: product, user, category, brand, location"
	}
	if strings.Contains(errMsg, "invalid file extension") {
		return 400, "Invalid file extension. Allowed: jpg, jpeg, png, gif, pdf, doc, docx, xls, xlsx"
	}
	if strings.Contains(errMsg, "failed to upload file to S3") {
		return 500, "Failed to upload file to storage"
	}

	// Handle validation errors
	if strings.Contains(errMsg, "is required") {
		return 400, errMsg
	}

	// Default to internal server error
	return 500, "Internal server error"
}

func GetFiles(c *fiber.Ctx) error {
	log.Printf("[FILE] Get all files request from IP: %s", c.IP())

	files, err := getFileService().GetAllFiles()
	if err != nil {
		log.Printf("[FILE] Get all files failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch files", err.Error())
	}

	log.Printf("[FILE] Get all files successful")
	return helper.Success(c, 200, "Success", files)
}

func GetFilesByModel(c *fiber.Ctx) error {
	modelType := c.Params("modelType")
	modelID := c.Params("modelId")

	log.Printf("[FILE] Get files by model request - Model Type: %s, Model ID: %s from IP: %s", modelType, modelID, c.IP())

	modelIDUint, err := strconv.ParseUint(modelID, 10, 32)
	if err != nil {
		log.Printf("[FILE] Get files by model failed - Invalid Model ID: %s, error: %v", modelID, err)
		return helper.Fail(c, 400, "Invalid model ID", err.Error())
	}

	files, err := getFileService().GetFilesByModel(modelType, uint(modelIDUint))
	if err != nil {
		log.Printf("[FILE] Get files by model failed - Model Type: %s, Model ID: %d, error: %v", modelType, modelIDUint, err)
		statusCode, message := handleFileError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[FILE] Get files by model successful - Model Type: %s, Model ID: %d", modelType, modelIDUint)
	return helper.Success(c, 200, "Success", files)
}

func GetFileByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[FILE] Get file by ID request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[FILE] Get file by ID failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid file ID", err.Error())
	}

	file, err := getFileService().GetFileByID(uint(idUint))
	if err != nil {
		log.Printf("[FILE] Get file by ID failed - File ID: %d, error: %v", idUint, err)
		statusCode, message := handleFileError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[FILE] Get file by ID successful - File ID: %d", idUint)
	return helper.Success(c, 200, "Success", file)
}

func UploadFile(c *fiber.Ctx) error {
	log.Printf("[FILE] Upload file request from IP: %s", c.IP())

	// Parse form data - handle both camelCase and PascalCase
	modelType := c.FormValue("modelType")
	if modelType == "" {
		modelType = c.FormValue("ModelType") // Try PascalCase
	}

	modelIDStr := c.FormValue("modelId")
	if modelIDStr == "" {
		modelIDStr = c.FormValue("ModelId") // Try PascalCase
		if modelIDStr == "" {
			modelIDStr = c.FormValue("modelID") // Try with uppercase ID
		}
	}

	log.Printf("[FILE] Upload file - Model Type: %s, Model ID: %s", modelType, modelIDStr)

	// Validate required fields
	if modelType == "" {
		log.Printf("[FILE] Upload file failed - Model type is required")
		return helper.Fail(c, 400, "Model type is required", "")
	}

	if modelIDStr == "" {
		log.Printf("[FILE] Upload file failed - Model ID is required")
		return helper.Fail(c, 400, "Model ID is required", "")
	}

	modelID, err := strconv.ParseUint(modelIDStr, 10, 32)
	if err != nil {
		log.Printf("[FILE] Upload file failed - Invalid Model ID: %s, error: %v", modelIDStr, err)
		return helper.Fail(c, 400, "Invalid model ID", err.Error())
	}

	// Get uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Printf("[FILE] Upload file failed - No file uploaded, error: %v", err)
		return helper.Fail(c, 400, "No file uploaded", err.Error())
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("[FILE] Upload file failed - Cannot open file, error: %v", err)
		return helper.Fail(c, 500, "Cannot open uploaded file", err.Error())
	}
	defer file.Close()

	// Read file data
	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Printf("[FILE] Upload file failed - Cannot read file data, error: %v", err)
		return helper.Fail(c, 500, "Cannot read file data", err.Error())
	}

	// Get file extension
	ext := strings.TrimPrefix(filepath.Ext(fileHeader.Filename), ".")
	if ext == "" {
		log.Printf("[FILE] Upload file failed - File has no extension")
		return helper.Fail(c, 400, "File must have an extension", "")
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[FILE] Upload file failed - User not authenticated")
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[FILE] Uploading file with audit - Model Type: %s, Model ID: %d, Extension: %s, User ID: %d, File Size: %d bytes",
		modelType, modelID, ext, userID, len(fileData))

	uploadedFile, err := getFileService().UploadFile(modelType, uint(modelID), ext, fileData, fileHeader.Filename, userID)
	if err != nil {
		log.Printf("[FILE] Upload file failed - Model Type: %s, Model ID: %d, User ID: %d, error: %v", modelType, modelID, userID, err)
		statusCode, message := handleFileError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[FILE] Upload file successful - Model Type: %s, Model ID: %d, Uploaded by User ID: %d", modelType, modelID, userID)
	return helper.Success(c, 201, "File uploaded successfully", uploadedFile)
}

func UpdateFile(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[FILE] Update file request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[FILE] Update file failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid file ID", err.Error())
	}

	var req UpdateFileRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[FILE] Update file failed - Invalid request body for ID: %d, error: %v", idUint, err)
		return helper.Fail(c, 400, "Invalid request body", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[FILE] Update file failed - User not authenticated for File ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	log.Printf("[FILE] Updating file with audit - File ID: %d, User ID: %d", idUint, userID)

	file, err := getFileService().UpdateFile(uint(idUint), req.ModelType, req.ModelID, req.Ext, userID)
	if err != nil {
		log.Printf("[FILE] Update file failed - File ID: %d, User ID: %d, error: %v", idUint, userID, err)
		statusCode, message := handleFileError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[FILE] Update file successful - File ID: %d, Updated by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "File updated successfully", file)
}

func DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[FILE] Delete file request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[FILE] Delete file failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid file ID", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[FILE] Delete file failed - User not authenticated for File ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	err = getFileService().DeleteFile(uint(idUint), userID)
	if err != nil {
		log.Printf("[FILE] Delete file failed - File ID: %d, error: %v", idUint, err)
		statusCode, message := handleFileError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[FILE] Delete file successful - File ID: %d, Deleted by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "File deleted successfully", nil)
}

func GetDeletedFiles(c *fiber.Ctx) error {
	log.Printf("[FILE] Get deleted files request from IP: %s", c.IP())

	files, err := getFileService().GetDeletedFiles()
	if err != nil {
		log.Printf("[FILE] Get deleted files failed - error: %v", err)
		return helper.Fail(c, 500, "Failed to fetch deleted files", err.Error())
	}

	log.Printf("[FILE] Get deleted files successful")
	return helper.Success(c, 200, "Success", files)
}

func RestoreFile(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[FILE] Restore file request - ID: %s from IP: %s", id, c.IP())

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("[FILE] Restore file failed - Invalid ID: %s, error: %v", id, err)
		return helper.Fail(c, 400, "Invalid file ID", err.Error())
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Printf("[FILE] Restore file failed - User not authenticated for File ID: %d", idUint)
		return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
	}

	file, err := getFileService().RestoreFile(uint(idUint), userID)
	if err != nil {
		log.Printf("[FILE] Restore file failed - File ID: %d, error: %v", idUint, err)
		statusCode, message := handleFileError(err)
		return helper.Fail(c, statusCode, message, err.Error())
	}

	log.Printf("[FILE] Restore file successful - File ID: %d, Restored by User ID: %d", idUint, userID)
	return helper.Success(c, 200, "File restored successfully", file)
}
