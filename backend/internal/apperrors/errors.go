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

// ErrorResponse es la estructura base de error devuelta en todos los endpoints.
// Swagger la utiliza como referencia para los tipos de error.
type ErrorResponse struct {
	Code    string `json:"code"    example:"NOT_FOUND"`
	Message string `json:"message" example:"Recurso no encontrado"`
}

// ValidationErrorResponse extiende ErrorResponse con una lista de detalles
// de validación para los errores 400 originados por binding/validación de campos.
type ValidationErrorResponse struct {
	Code    string   `json:"code"    example:"BAD_REQUEST"`
	Message string   `json:"message" example:"Validation failed"`
	Details []string `json:"details" example:"Username: es obligatorio;Email: debe ser un correo valido"`
}
