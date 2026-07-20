package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grupo5/ecommerce-api/internal/apperrors"
	"github.com/grupo5/ecommerce-api/internal/dto"
	"github.com/grupo5/ecommerce-api/internal/middleware"
	"github.com/grupo5/ecommerce-api/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register godoc
// @Summary      Registrar usuario
// @Description  Crea un nuevo usuario y devuelve un token JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Datos de registro"
// @Success      201 {object} dto.AuthResponse
// @Failure      400 {object} apperrors.ErrorResponse
// @Failure      409 {object} apperrors.ErrorResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(formatValidationError(err)))
		return
	}

	response, err := h.userService.Register(c.Request.Context(), req)
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary      Iniciar sesion
// @Description  Autentica un usuario y devuelve un token JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Credenciales"
// @Success      200 {object} dto.AuthResponse
// @Failure      400 {object} apperrors.ErrorResponse
// @Failure      401 {object} apperrors.ErrorResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(formatValidationError(err)))
		return
	}

	response, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetByID godoc
// @Summary      Obtener usuario por ID
// @Description  Consulta los datos de un usuario (sin contrasena)
// @Tags         users
// @Produce      json
// @Param        id path int true "ID del usuario"
// @Success      200 {object} dto.UserResponse
// @Failure      400 {object} apperrors.ErrorResponse
// @Failure      404 {object} apperrors.ErrorResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	response, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAll godoc
// @Summary      Obtener todos los usuarios
// @Description  Consulta la lista de todos los usuarios (solo ADMIN)
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} dto.UserResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	response, err := h.userService.GetAll(c.Request.Context())
	if err != nil {
		apperrors.Abort(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

// Update godoc
// @Summary      Actualizar usuario
// @Description  Actualiza los datos de un usuario autenticado
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del usuario"
// @Param        request body dto.UpdateUserRequest true "Datos a actualizar"
// @Success      200 {object} dto.UserResponse
// @Failure      400 {object} apperrors.ErrorResponse
// @Failure      401 {object} apperrors.ErrorResponse
// @Failure      403 {object} apperrors.ErrorResponse
// @Failure      404 {object} apperrors.ErrorResponse
// @Failure      409 {object} apperrors.ErrorResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	authUserID, ok := middleware.GetUserID(c)
	authRole, _ := middleware.GetRole(c)
	if (!ok || authUserID != id) && authRole != "ADMIN" {
		apperrors.Abort(c, apperrors.Forbidden("No puede modificar otro usuario"))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(formatValidationError(err)))
		return
	}

	response, err := h.userService.Update(c.Request.Context(), id, req)
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Delete godoc
// @Summary      Eliminar usuario
// @Description  Elimina un usuario autenticado
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del usuario"
// @Success      204
// @Failure      400 {object} apperrors.ErrorResponse
// @Failure      401 {object} apperrors.ErrorResponse
// @Failure      404 {object} apperrors.ErrorResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	authUserID, ok := middleware.GetUserID(c)
	authRole, _ := middleware.GetRole(c)
	if (!ok || authUserID != id) && authRole != "ADMIN" {
		apperrors.Abort(c, apperrors.Forbidden("No puede eliminar otro usuario"))
		return
	}

	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		apperrors.Abort(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
