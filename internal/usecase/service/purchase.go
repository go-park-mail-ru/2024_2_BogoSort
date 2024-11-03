package service

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"go.uber.org/zap"
)

type PurchaseService struct {
	purchaseRepo repository.PurchaseRepository
	cartRepo repository.Cart
	logger       *zap.Logger
}

func NewPurchaseService(purchaseRepo repository.PurchaseRepository, cartRepo repository.Cart, logger *zap.Logger) *PurchaseService {
	return &PurchaseService{purchaseRepo: purchaseRepo, cartRepo: cartRepo, logger: logger}
}

func (s *PurchaseService) purchaseEntityToDTO(purchase *entity.Purchase) (*dto.PurchaseResponse, error) {
	return &dto.PurchaseResponse{
		ID: purchase.ID,
		CartID: purchase.CartID,
		Address: purchase.Address,
		Status: dto.PurchaseStatus(purchase.Status),
		PaymentMethod: dto.PaymentMethod(purchase.PaymentMethod),
		DeliveryMethod: dto.DeliveryMethod(purchase.DeliveryMethod),
	}, nil
}

func (s *PurchaseService) AddPurchase(purchaseRequest dto.PurchaseRequest) (*dto.PurchaseResponse, error) {
	ctx := context.Background()
	tx, err := s.purchaseRepo.BeginTransaction()
	if err != nil {
		s.logger.Error("failed to begin transaction", zap.Error(err))
		return nil, entity.UsecaseWrap(errors.New("failed to begin transaction"), err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	purchase, err := s.purchaseRepo.AddPurchase(tx, &entity.Purchase{
		CartID: purchaseRequest.CartID,
		Address: purchaseRequest.Address,
		Status: entity.StatusPending,
		PaymentMethod: entity.PaymentMethod(purchaseRequest.PaymentMethod),
		DeliveryMethod: entity.DeliveryMethod(purchaseRequest.DeliveryMethod),
	})
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("failed to add purchase"), err)
	}

	err = s.cartRepo.UpdateCartStatus(tx, purchase.CartID, entity.CartStatusInactive)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("failed to update cart status"), err)
	}

	return s.purchaseEntityToDTO(purchase)
}
