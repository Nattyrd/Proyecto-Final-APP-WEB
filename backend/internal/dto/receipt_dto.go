package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateReceiptItemRequest struct {
	ProductID uint `json:"productId" binding:"required,gt=0"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

type CreateReceiptRequest struct {
	UserID uint                       `json:"userId" binding:"required,gt=0"`
	Items  []CreateReceiptItemRequest `json:"items" binding:"required,min=1,dive"`
}

type ReceiptItemResponse struct {
	ID        uint            `json:"id"`
	ProductID uint            `json:"productId"`
	Quantity  int             `json:"quantity"`
	UnitPrice decimal.Decimal `json:"unitPrice"`
	Subtotal  decimal.Decimal `json:"subtotal"`
}

type ReceiptResponse struct {
	ID        uint                  `json:"id"`
	UserID    uint                  `json:"userId"`
	Total     decimal.Decimal       `json:"total"`
	Items     []ReceiptItemResponse `json:"items"`
	CreatedAt time.Time             `json:"createdAt"`
}
