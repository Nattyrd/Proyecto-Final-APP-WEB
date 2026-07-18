package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grupo5/ecommerce-api/internal/apperrors"
	"github.com/grupo5/ecommerce-api/internal/dto"
	"github.com/grupo5/ecommerce-api/internal/middleware"
	"github.com/grupo5/ecommerce-api/internal/service"
)

type ReceiptHandler struct {
	receiptService *service.ReceiptService
}

func NewReceiptHandler(receiptService *service.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{receiptService: receiptService}
}

// Create godoc
// @Summary      Crear recibo
// @Description  Procesa una compra, calcula el total en backend y descuenta stock
// @Tags         receipts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateReceiptRequest true "Datos de la compra"
// @Success      201 {object} dto.ReceiptResponse
// @Failure      400 {object} apperrors.ErrorResponse
// @Failure      401 {object} apperrors.ErrorResponse
// @Failure      404 {object} apperrors.ErrorResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /receipts [post]
func (h *ReceiptHandler) Create(c *gin.Context) {
	var req dto.CreateReceiptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(formatValidationError(err)))
		return
	}

	authUserID, ok := middleware.GetUserID(c)
	if !ok || authUserID != req.UserID {
		apperrors.Abort(c, apperrors.Forbidden("Solo puede crear recibos para su propio usuario"))
		return
	}

	response, err := h.receiptService.Create(c.Request.Context(), req)
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetAll godoc
// @Summary      Listar recibos
// @Description  Obtiene todos los recibos registrados
// @Tags         receipts
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} dto.ReceiptResponse
// @Failure      401 {object} apperrors.ErrorResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /receipts [get]
func (h *ReceiptHandler) GetAll(c *gin.Context) {
	receipts, err := h.receiptService.GetAll(c.Request.Context())
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, receipts)
}

// GetByID godoc
// @Summary      Obtener recibo por ID
// @Description  Consulta un recibo con sus items
// @Tags         receipts
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del recibo"
// @Success      200 {object} dto.ReceiptResponse
// @Failure      400 {object} apperrors.ErrorResponse
// @Failure      401 {object} apperrors.ErrorResponse
// @Failure      404 {object} apperrors.ErrorResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /receipts/{id} [get]
func (h *ReceiptHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	response, err := h.receiptService.GetByID(c.Request.Context(), id)
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetByUserID godoc
// @Summary      Listar recibos por usuario
// @Description  Obtiene los recibos asociados a un usuario
// @Tags         receipts
// @Produce      json
// @Security     BearerAuth
// @Param        userId path int true "ID del usuario"
// @Success      200 {array} dto.ReceiptResponse
// @Failure      400 {object} apperrors.ErrorResponse
// @Failure      401 {object} apperrors.ErrorResponse
// @Failure      404 {object} apperrors.ErrorResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /receipts/user/{userId} [get]
func (h *ReceiptHandler) GetByUserID(c *gin.Context) {
	userID, err := parseUintParam(c, "userId")
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	authUserID, ok := middleware.GetUserID(c)
	if !ok || authUserID != userID {
		apperrors.Abort(c, apperrors.Forbidden("Solo puede consultar sus propios recibos"))
		return
	}

	receipts, err := h.receiptService.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, receipts)
}

// Delete godoc
// @Summary      Eliminar recibo
// @Description  Elimina un recibo existente
// @Tags         receipts
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del recibo"
// @Success      204
// @Failure      400 {object} apperrors.ErrorResponse
// @Failure      401 {object} apperrors.ErrorResponse
// @Failure      404 {object} apperrors.ErrorResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /receipts/{id} [delete]
func (h *ReceiptHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	if err := h.receiptService.Delete(c.Request.Context(), id); err != nil {
		apperrors.Abort(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
