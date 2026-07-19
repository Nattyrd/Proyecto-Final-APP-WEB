package handlers

import (
	"net/http"
	"strconv"

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
// @Description  Registra un nuevo producto en el catálogo
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
// @Summary      Listar productos (paginado)
// @Description  Obtiene los productos con soporte de paginación mediante query params. Si no se especifican, se devuelven los primeros 10 resultados.
// @Tags         products
// @Produce      json
// @Param        page      query int false "Número de página (defecto: 1)"      minimum(1)
// @Param        pageSize  query int false "Resultados por página (defecto: 10, máximo: 100)" minimum(1) maximum(100)
// @Success      200 {object} dto.PaginatedProductsResponse
// @Failure      500 {object} apperrors.ErrorResponse
// @Router       /products [get]
func (h *ProductHandler) GetAll(c *gin.Context) {
	page := parseIntQuery(c, "page", 1)
	pageSize := parseIntQuery(c, "pageSize", 10)

	result, err := h.productService.GetAllPaginated(c.Request.Context(), page, pageSize)
	if err != nil {
		apperrors.Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
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
// @Description  Modifica los datos de un producto existente. Solo los campos enviados son modificados.
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
// @Description  Elimina un producto del catálogo
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

// parseIntQuery lee un query param como entero con un fallback.
func parseIntQuery(c *gin.Context, name string, fallback int) int {
	raw := c.Query(name)
	if raw == "" {
		return fallback
	}
	val, err := strconv.Atoi(raw)
	if err != nil || val <= 0 {
		return fallback
	}
	return val
}
