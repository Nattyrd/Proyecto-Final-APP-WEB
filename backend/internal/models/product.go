package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"size:120;not null" json:"name"`
	Description string          `gorm:"size:500" json:"description"`
	Price       decimal.Decimal `gorm:"type:numeric(12,2);not null" json:"price"`
	Stock       int             `gorm:"not null;default:0" json:"stock"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}
