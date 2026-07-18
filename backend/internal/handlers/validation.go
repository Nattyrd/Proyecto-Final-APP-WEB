package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func formatValidationError(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		messages := make([]string, 0, len(validationErrors))
		for _, fieldErr := range validationErrors {
			messages = append(messages, fmt.Sprintf("%s: %s", fieldErr.Field(), validationMessage(fieldErr)))
		}
		return strings.Join(messages, "; ")
	}

	return err.Error()
}

func validationMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "es obligatorio"
	case "email":
		return "debe ser un correo valido"
	case "min":
		return fmt.Sprintf("debe tener al menos %s caracteres", fieldErr.Param())
	case "max":
		return fmt.Sprintf("debe tener maximo %s caracteres", fieldErr.Param())
	case "gt":
		return fmt.Sprintf("debe ser mayor a %s", fieldErr.Param())
	case "gte":
		return fmt.Sprintf("debe ser mayor o igual a %s", fieldErr.Param())
	default:
		return "valor invalido"
	}
}
