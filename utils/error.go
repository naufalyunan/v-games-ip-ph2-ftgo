package utils

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// Implement the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Detail: %s", e.Code, e.Message, e.Detail)
}

func NewNotFoundError(message string) *APIError {
	return &APIError{
		Code:    http.StatusNotFound,
		Message: message,
		Detail:  "Resource not found",
	}
}

func NewBadRequestError(message string) *APIError {
	return &APIError{
		Code:    http.StatusBadRequest,
		Message: message,
		Detail:  "Invalid request data",
	}
}

func NewInternalError(message string) *APIError {
	return &APIError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Detail:  "Internal server error",
	}
}

// Fungsi untuk menangani Unauthorized error
func NewUnauthorizedError(message string) *APIError {
	return &APIError{
		Code:    http.StatusUnauthorized,
		Message: message,
		Detail:  "Unauthorized access",
	}
}

func HandleError(c echo.Context, err *APIError) error {
	return c.JSON(err.Code, err)
}
