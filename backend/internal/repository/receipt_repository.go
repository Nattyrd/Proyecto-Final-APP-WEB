package repository

import (
	"context"

	"github.com/grupo5/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type ReceiptRepository struct {
	db *gorm.DB
}

func NewReceiptRepository(db *gorm.DB) *ReceiptRepository {
	return &ReceiptRepository{db: db}
}

func (r *ReceiptRepository) CreateWithItems(ctx context.Context, tx *gorm.DB, receipt *models.Receipt) error {
	return tx.WithContext(ctx).Create(receipt).Error
}

func (r *ReceiptRepository) FindAll(ctx context.Context) ([]models.Receipt, error) {
	var receipts []models.Receipt
	err := r.db.WithContext(ctx).Preload("Items").Order("id desc").Find(&receipts).Error
	return receipts, err
}

func (r *ReceiptRepository) FindByID(ctx context.Context, id uint) (*models.Receipt, error) {
	var receipt models.Receipt
	err := r.db.WithContext(ctx).Preload("Items").First(&receipt, id).Error
	if err != nil {
		return nil, err
	}
	return &receipt, nil
}

func (r *ReceiptRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Receipt, error) {
	var receipts []models.Receipt
	err := r.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).Order("id desc").Find(&receipts).Error
	return receipts, err
}

func (r *ReceiptRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Receipt{}, id).Error
}

func (r *ReceiptRepository) DB() *gorm.DB {
	return r.db
}
