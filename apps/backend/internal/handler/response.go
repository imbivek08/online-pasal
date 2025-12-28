package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// StandardResponse represents a standard API response
type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SendSuccess sends a successful response
func SendSuccess(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendError sends an error response
func SendError(c echo.Context, statusCode int, err error, message string) error {
	return c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   err.Error(),
		Message: message,
	})
}

// SendValidationError sends a validation error response
func SendValidationError(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{
		Success: false,
		Error:   "validation_error",
		Message: message,
	})
}

// SendNotFound sends a not found error response
func SendNotFound(c echo.Context, message string) error {
	return c.JSON(http.StatusNotFound, ErrorResponse{
		Success: false,
		Error:   "not_found",
		Message: message,
	})
}

// SendInternalError sends an internal server error response
func SendInternalError(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Error:   "internal_server_error",
		Message: err.Error(),
	})
}
