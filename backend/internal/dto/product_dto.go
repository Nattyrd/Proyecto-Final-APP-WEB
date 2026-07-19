package dto

import "github.com/shopspring/decimal"

type CreateProductRequest struct {
	Name        string          `json:"name" binding:"required,min=2,max=120"`
	Description string          `json:"description" binding:"max=500"`
	Price       decimal.Decimal `json:"price" binding:"required"`
	Stock       int             `json:"stock" binding:"required,gte=0"`
}

type UpdateProductRequest struct {
	Name        string          `json:"name" binding:"omitempty,min=2,max=120"`
	Description string          `json:"description" binding:"omitempty,max=500"`
	Price       decimal.Decimal `json:"price" binding:"omitempty"`
	Stock       *int            `json:"stock" binding:"omitempty,gte=0"`
}

type ProductResponse struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	Stock       int             `json:"stock"`
}

type PaginatedProductsResponse struct {
	Data       []ProductResponse `json:"data"`
	Page       int               `json:"page"`
	PageSize   int               `json:"pageSize"`
	TotalItems int64             `json:"totalItems"`
	TotalPages int               `json:"totalPages"`
}
