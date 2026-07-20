package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/grupo5/ecommerce-api/internal/apperrors"
	"github.com/grupo5/ecommerce-api/internal/auth"
	"github.com/grupo5/ecommerce-api/internal/models"
)

const (
	UserIDKey   = "userId"
	UsernameKey = "username"
	RoleKey     = "role" // clave para recuperar el rol desde el contexto Gin
)

// JWTAuth valida el token Bearer y almacena userId, username y role en el contexto.
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
		c.Set(RoleKey, claims.Role)
		c.Next()
	}
}

// RequireAdmin es un middleware adicional (debe ejecutarse después de JWTAuth).
// Rechaza con 403 si el rol del usuario autenticado no es "ADMIN".
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(RoleKey)
		if !exists || role != models.RoleAdmin {
			apperrors.Abort(c, apperrors.Forbidden("Se requiere rol ADMIN para esta operacion"))
			return
		}
		c.Next()
	}
}

// GetUserID extrae el userID del contexto Gin de forma segura.
func GetUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint)
	return userID, ok
}

// GetRole extrae el rol del contexto Gin de forma segura.
func GetRole(c *gin.Context) (string, bool) {
	value, exists := c.Get(RoleKey)
	if !exists {
		return "", false
	}
	role, ok := value.(string)
	return role, ok
}
