package repository

import (
	"context"

	"github.com/grupo5/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).Order("id asc").Find(&products).Error
	return products, err
}

func (r *ProductRepository) FindByID(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindByIDForUpdate(ctx context.Context, tx *gorm.DB, id uint) (*models.Product, error) {
	var product models.Product
	err := tx.WithContext(ctx).Clauses().First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *ProductRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

func (r *ProductRepository) DB() *gorm.DB {
	return r.db
}
