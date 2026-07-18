package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/grupo5/ecommerce-api/internal/apperrors"
)

func parseUintParam(c *gin.Context, name string) (uint, error) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || value == 0 {
		return 0, apperrors.BadRequest("Identificador invalido")
	}
	return uint(value), nil
}
