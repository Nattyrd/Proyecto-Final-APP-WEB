package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/grupo5/ecommerce-api/internal/apperrors"
	"github.com/grupo5/ecommerce-api/internal/auth"
)

const (
	UserIDKey   = "userId"
	UsernameKey = "username"
)

func JWTAuth(tokenService *auth.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			apperrors.Abort(c, apperrors.Unauthorized("Token de autenticacion requerido"))
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			apperrors.Abort(c, apperrors.Unauthorized("Formato de token invalido"))
			return
		}

		claims, err := tokenService.Validate(parts[1])
		if err != nil {
			apperrors.Abort(c, apperrors.Unauthorized("Token invalido o expirado"))
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UsernameKey, claims.Username)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint)
	return userID, ok
}
