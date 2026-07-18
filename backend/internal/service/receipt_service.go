package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/grupo5/ecommerce-api/internal/apperrors"
	"github.com/grupo5/ecommerce-api/internal/dto"
	"github.com/grupo5/ecommerce-api/internal/models"
	"github.com/grupo5/ecommerce-api/internal/repository"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReceiptService struct {
	receiptRepo *repository.ReceiptRepository
	productRepo *repository.ProductRepository
	userRepo    *repository.UserRepository
}

func NewReceiptService(
	receiptRepo *repository.ReceiptRepository,
	productRepo *repository.ProductRepository,
	userRepo *repository.UserRepository,
) *ReceiptService {
	return &ReceiptService{
		receiptRepo: receiptRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
	}
}

func (s *ReceiptService) Create(ctx context.Context, req dto.CreateReceiptRequest) (*dto.ReceiptResponse, error) {
	if _, err := s.userRepo.FindByID(ctx, req.UserID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("Usuario no encontrado")
		}
		return nil, apperrors.Internal("No se pudo validar el usuario")
	}

	var createdReceipt *models.Receipt

	err := s.receiptRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		total := decimal.Zero
		items := make([]models.ReceiptItem, 0, len(req.Items))

		for _, itemReq := range req.Items {
			var product models.Product
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				First(&product, itemReq.ProductID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return apperrors.NotFound(fmt.Sprintf("Producto %d no encontrado", itemReq.ProductID))
				}
				return apperrors.Internal("No se pudo consultar el producto")
			}

			if product.Stock < itemReq.Quantity {
				return apperrors.BadRequest(
					fmt.Sprintf("Stock insuficiente para el producto '%s'. Disponible: %d, solicitado: %d",
						product.Name, product.Stock, itemReq.Quantity),
				)
			}

			unitPrice := product.Price
			subtotal := unitPrice.Mul(decimal.NewFromInt(int64(itemReq.Quantity)))
			total = total.Add(subtotal)

			product.Stock -= itemReq.Quantity
			if err := tx.Save(&product).Error; err != nil {
				return apperrors.Internal("No se pudo actualizar el stock")
			}

			items = append(items, models.ReceiptItem{
				ProductID: itemReq.ProductID,
				Quantity:  itemReq.Quantity,
				UnitPrice: unitPrice,
				Subtotal:  subtotal,
			})
		}

		receipt := &models.Receipt{
			UserID: req.UserID,
			Total:  total,
			Items:  items,
		}

		if err := tx.Create(receipt).Error; err != nil {
			return apperrors.Internal("No se pudo crear el recibo")
		}

		createdReceipt = receipt
		return nil
	})

	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		}
		return nil, apperrors.Internal("No se pudo procesar la compra")
	}

	fullReceipt, err := s.receiptRepo.FindByID(ctx, createdReceipt.ID)
	if err != nil {
		return nil, apperrors.Internal("No se pudo consultar el recibo creado")
	}

	response := toReceiptResponse(fullReceipt)
	return &response, nil
}

func (s *ReceiptService) GetAll(ctx context.Context) ([]dto.ReceiptResponse, error) {
	receipts, err := s.receiptRepo.FindAll(ctx)
	if err != nil {
		return nil, apperrors.Internal("No se pudo listar los recibos")
	}

	responses := make([]dto.ReceiptResponse, 0, len(receipts))
	for i := range receipts {
		responses = append(responses, toReceiptResponse(&receipts[i]))
	}

	return responses, nil
}

func (s *ReceiptService) GetByID(ctx context.Context, id uint) (*dto.ReceiptResponse, error) {
	receipt, err := s.receiptRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("Recibo no encontrado")
		}
		return nil, apperrors.Internal("No se pudo consultar el recibo")
	}

	response := toReceiptResponse(receipt)
	return &response, nil
}

func (s *ReceiptService) GetByUserID(ctx context.Context, userID uint) ([]dto.ReceiptResponse, error) {
	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("Usuario no encontrado")
		}
		return nil, apperrors.Internal("No se pudo validar el usuario")
	}

	receipts, err := s.receiptRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal("No se pudo listar los recibos del usuario")
	}

	responses := make([]dto.ReceiptResponse, 0, len(receipts))
	for i := range receipts {
		responses = append(responses, toReceiptResponse(&receipts[i]))
	}

	return responses, nil
}

func (s *ReceiptService) Delete(ctx context.Context, id uint) error {
	if _, err := s.receiptRepo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound("Recibo no encontrado")
		}
		return apperrors.Internal("No se pudo consultar el recibo")
	}

	if err := s.receiptRepo.Delete(ctx, id); err != nil {
		return apperrors.Internal("No se pudo eliminar el recibo")
	}

	return nil
}

func toReceiptResponse(receipt *models.Receipt) dto.ReceiptResponse {
	items := make([]dto.ReceiptItemResponse, 0, len(receipt.Items))
	for _, item := range receipt.Items {
		items = append(items, dto.ReceiptItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
		})
	}

	return dto.ReceiptResponse{
		ID:        receipt.ID,
		UserID:    receipt.UserID,
		Total:     receipt.Total,
		Items:     items,
		CreatedAt: receipt.CreatedAt,
	}
}
