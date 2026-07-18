package service

import (
	"context"
	"errors"

	"github.com/grupo5/ecommerce-api/internal/apperrors"
	"github.com/grupo5/ecommerce-api/internal/dto"
	"github.com/grupo5/ecommerce-api/internal/models"
	"github.com/grupo5/ecommerce-api/internal/repository"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

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
	if req.Stock >= 0 {
		product.Stock = req.Stock
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
