package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grupo5/ecommerce-api/internal/apperrors"
	"github.com/grupo5/ecommerce-api/internal/dto"
	"github.com/grupo5/ecommerce-api/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// Create godoc
// @Summary      Crear producto
// @Description  Registra un nuevo producto en el catalogo
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateProductRequest true "Datos del producto"
// @Success      201 {object} dto.ProductResponse
// @Failure      400 {object} apperrors.ErrorResponse
// @Failure      401 {object} apperrors.ErrorResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(formatValidationError(err)))
		return
	}

	response, err := h.productService.Create(c.Request.Context(), req)
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetAll godoc
// @Summary      Listar productos
// @Description  Obtiene todos los productos disponibles
// @Tags         products
// @Produce      json
// @Success      200 {array} dto.ProductResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /products [get]
func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.productService.GetAll(c.Request.Context())
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetByID godoc
// @Summary      Obtener producto por ID
// @Description  Consulta un producto por su identificador
// @Tags         products
// @Produce      json
// @Param        id path int true "ID del producto"
// @Success      200 {object} dto.ProductResponse
// @Failure      400 {object} apperrors.ErrorResponse
// @Failure      404 {object} apperrors.ErrorResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /products/{id} [get]
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	response, err := h.productService.GetByID(c.Request.Context(), id)
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Update godoc
// @Summary      Actualizar producto
// @Description  Modifica los datos de un producto existente
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del producto"
// @Param        request body dto.UpdateProductRequest true "Datos a actualizar"
// @Success      200 {object} dto.ProductResponse
// @Failure      400 {object} apperrors.ErrorResponse
// @Failure      401 {object} apperrors.ErrorResponse
// @Failure      404 {object} apperrors.ErrorResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(formatValidationError(err)))
		return
	}

	response, err := h.productService.Update(c.Request.Context(), id, req)
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Delete godoc
// @Summary      Eliminar producto
// @Description  Elimina un producto del catalogo
// @Tags         products
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del producto"
// @Success      204
// @Failure      400 {object} apperrors.ErrorResponse
// @Failure      401 {object} apperrors.ErrorResponse
// @Failure      404 {object} apperrors.ErrorResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	if err := h.productService.Delete(c.Request.Context(), id); err != nil {
		apperrors.Abort(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
