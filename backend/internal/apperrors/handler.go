package apperrors

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var appErr *AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.Status, ErrorResponse{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
			return
		}

		log.Printf("error no controlado: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Ocurrio un error interno en el servidor",
		})
	}
}

func Abort(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}
