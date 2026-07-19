package service

import (
	"context"
	"errors"
	"math"

	"github.com/grupo5/ecommerce-api/internal/apperrors"
	"github.com/grupo5/ecommerce-api/internal/dto"
	"github.com/grupo5/ecommerce-api/internal/models"
	"github.com/grupo5/ecommerce-api/internal/repository"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// DefaultPageSize es el tamaño de página por defecto si el cliente no lo especifica.
const DefaultPageSize = 10

// MaxPageSize es el máximo permitido para evitar consultas masivas.
const MaxPageSize = 100

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) Create(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	if req.Price.LessThanOrEqual(decimal.Zero) {
		return nil, apperrors.BadRequest("El precio debe ser mayor a cero")
	}

	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, apperrors.Internal("No se pudo crear el producto")
	}

	response := toProductResponse(product)
	return &response, nil
}

// GetAllPaginated devuelve los productos con soporte de paginación.
// page y pageSize son opcionales; si son <= 0 se usan los valores por defecto.
func (s *ProductService) GetAllPaginated(ctx context.Context, page, pageSize int) (*dto.PaginatedProductsResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	products, total, err := s.productRepo.FindPaginated(ctx, page, pageSize)
	if err != nil {
		return nil, apperrors.Internal("No se pudo listar los productos")
	}

	responses := make([]dto.ProductResponse, 0, len(products))
	for i := range products {
		responses = append(responses, toProductResponse(&products[i]))
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	return &dto.PaginatedProductsResponse{
		Data:       responses,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

// GetAll devuelve todos los productos sin paginación (conservado por compatibilidad interna).
func (s *ProductService) GetAll(ctx context.Context) ([]dto.ProductResponse, error) {
	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		return nil, apperrors.Internal("No se pudo listar los productos")
	}

	responses := make([]dto.ProductResponse, 0, len(products))
	for i := range products {
		responses = append(responses, toProductResponse(&products[i]))
	}

	return responses, nil
}

func (s *ProductService) GetByID(ctx context.Context, id uint) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("Producto no encontrado")
		}
		return nil, apperrors.Internal("No se pudo consultar el producto")
	}

	response := toProductResponse(product)
	return &response, nil
}

func (s *ProductService) Update(ctx context.Context, id uint, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("Producto no encontrado")
		}
		return nil, apperrors.Internal("No se pudo consultar el producto")
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if !req.Price.IsZero() {
		if req.Price.LessThanOrEqual(decimal.Zero) {
			return nil, apperrors.BadRequest("El precio debe ser mayor a cero")
		}
		product.Price = req.Price
	}
	// FIX: Se usa *int para detectar si el campo fue enviado o no.
	// Si req.Stock == nil, el cliente no envió el campo → no se modifica.
	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, apperrors.Internal("No se pudo actualizar el producto")
	}

	response := toProductResponse(product)
	return &response, nil
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	if _, err := s.productRepo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound("Producto no encontrado")
		}
		return apperrors.Internal("No se pudo consultar el producto")
	}

	if err := s.productRepo.Delete(ctx, id); err != nil {
		return apperrors.Internal("No se pudo eliminar el producto")
	}

	return nil
}

func toProductResponse(product *models.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}
}
