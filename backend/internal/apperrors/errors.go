package apperrors

import "net/http"

type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(status int, code, message string) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message)
}

func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func Forbidden(message string) *AppError {
	return New(http.StatusForbidden, "FORBIDDEN", message)
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, "NOT_FOUND", message)
}

func Conflict(message string) *AppError {
	return New(http.StatusConflict, "CONFLICT", message)
}

func Internal(message string) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_ERROR", message)
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
