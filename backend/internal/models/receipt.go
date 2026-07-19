package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Receipt struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	UserID    uint            `gorm:"not null;index" json:"userId"`
	User      User            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	Total     decimal.Decimal `gorm:"type:numeric(14,2);not null" json:"total"`
	Items     []ReceiptItem   `gorm:"foreignKey:ReceiptID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"items"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

type ReceiptItem struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	ReceiptID uint            `gorm:"not null;index" json:"receiptId"`
	ProductID uint            `gorm:"not null;index" json:"productId"`
	Product   Product         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	Quantity  int             `gorm:"not null" json:"quantity"`
	UnitPrice decimal.Decimal `gorm:"type:numeric(12,2);not null" json:"unitPrice"`
	Subtotal  decimal.Decimal `gorm:"type:numeric(14,2);not null" json:"subtotal"`
	CreatedAt time.Time       `json:"createdAt"`
}
