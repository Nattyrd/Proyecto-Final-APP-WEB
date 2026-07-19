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

// FindAll devuelve todos los productos sin paginación (uso interno).
func (r *ProductRepository) FindAll(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).Order("id asc").Find(&products).Error
	return products, err
}

// FindPaginated devuelve productos con paginación. Página indexada desde 1.
// Retorna también el total de registros para calcular metadatos en el service.
func (r *ProductRepository) FindPaginated(ctx context.Context, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Order("id asc").
		Limit(pageSize).
		Offset(offset).
		Find(&products).Error

	return products, total, err
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
