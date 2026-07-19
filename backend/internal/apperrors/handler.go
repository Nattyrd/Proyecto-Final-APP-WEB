package apperrors

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorHandler es el middleware centralizado de manejo de errores.
// Se registra en Gin DESPUÉS de todos los handlers mediante router.Use().
// Procesa los errores adjuntados vía c.Error(err) y los serializa en JSON.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var appErr *AppError
		if errors.As(err, &appErr) {
			// Si el error es BAD_REQUEST y contiene detalles de validación (separados por ";"),
			// devolvemos la estructura extendida con Details.
			if appErr.Status == http.StatusBadRequest && strings.Contains(appErr.Message, ":") {
				details := parseDetails(appErr.Message)
				c.JSON(appErr.Status, ValidationErrorResponse{
					Code:    appErr.Code,
					Message: "Error de validación en los datos enviados",
					Details: details,
				})
				return
			}

			c.JSON(appErr.Status, ErrorResponse{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
			return
		}

		// Errores de validación directos de Gin (no pasados por apperrors.Abort)
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			details := make([]string, 0, len(validationErrors))
			for _, fe := range validationErrors {
				details = append(details, fe.Field()+": valor inválido (tag: "+fe.Tag()+")")
			}
			c.JSON(http.StatusBadRequest, ValidationErrorResponse{
				Code:    "BAD_REQUEST",
				Message: "Error de validación en los datos enviados",
				Details: details,
			})
			return
		}

		log.Printf("[ERROR] error no controlado: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Ocurrió un error interno en el servidor",
		})
	}
}

// Abort agrega el error a Gin y aborta la cadena de handlers.
func Abort(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}

// parseDetails divide el mensaje de error de validación (ej: "Field1: msg; Field2: msg")
// en un slice de strings para la respuesta estructurada.
func parseDetails(message string) []string {
	parts := strings.Split(message, ";")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
