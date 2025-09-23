package helper

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SuccessMessage untuk response tanpa data (hanya message)
func Success(c *fiber.Ctx, code int, message string, data interface{}) error {
	return c.Status(code).JSON(APIResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// SuccessOptional dengan data optional menggunakan variadic
func SuccessOptional(c *fiber.Ctx, code int, message string, data ...interface{}) error {
	response := APIResponse{
		Code:    code,
		Message: message,
	}

	// Jika ada data, set ke response
	if len(data) > 0 && data[0] != nil {
		response.Data = data[0]
	}

	return c.Status(code).JSON(response)
}

func Fail(c *fiber.Ctx, code int, message string, err interface{}) error {
	return c.Status(code).JSON(APIResponse{
		Code:    code,
		Message: message,
		Error:   err,
	})
}

// Message Sukses
// "Success"
// "Created successfully"
// "Updated successfully"
// "Deleted successfully"
// "Fetched successfully"
// "Login successful"
// "Logout successful"
// Message Error
// "Bad request"
// "Validation failed"
// "Unauthorized"
// "Forbidden"
// "Not found"
// "Internal server error"
// "Email is required"
// "Password is invalid"
// "Resource not found"
// HTTP Code
// 200: OK
// 201: Created
// 400: Bad Request
// 401: Unauthorized
// 403: Forbidden
// 404: Not Found
// 422: Unprocessable Entity
// 500: Internal Server Error
