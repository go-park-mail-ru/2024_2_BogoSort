package service

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PurchaseService struct {
	purchaseRepo repository.PurchaseRepository
	cartRepo     repository.Cart
	advertRepo   repository.AdvertRepository
}

func NewPurchaseService(purchaseRepo repository.PurchaseRepository, advertRepo repository.AdvertRepository, cartRepo repository.Cart) *PurchaseService {
	return &PurchaseService{purchaseRepo: purchaseRepo, advertRepo: advertRepo, cartRepo: cartRepo}
}

func (s *PurchaseService) purchaseEntityToDTO(purchase *entity.Purchase) (*dto.PurchaseResponse, error) {
	return &dto.PurchaseResponse{
		ID:             purchase.ID,
		CartID:         purchase.CartID,
		Address:        purchase.Address,
		Status:         dto.PurchaseStatus(purchase.Status),
		PaymentMethod:  dto.PaymentMethod(purchase.PaymentMethod),
		DeliveryMethod: dto.DeliveryMethod(purchase.DeliveryMethod),
	}, nil
}

func (s *PurchaseService) Add(purchaseRequest dto.PurchaseRequest, userId uuid.UUID) (*dto.PurchaseResponse, error) {
	ctx := context.Background()
	tx, err := s.purchaseRepo.BeginTransaction()
	if err != nil {
		logger := middleware.GetLogger(ctx)
		logger.Error("failed to begin transaction", zap.Error(err), zap.String("userId", userId.String()))
		return nil, entity.UsecaseWrap(errors.New("failed to begin transaction"), err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	purchase, err := s.purchaseRepo.Add(tx, &entity.Purchase{
		CartID:         purchaseRequest.CartID,
		Address:        purchaseRequest.Address,
		Status:         entity.StatusPending,
		PaymentMethod:  entity.PaymentMethod(purchaseRequest.PaymentMethod),
		DeliveryMethod: entity.DeliveryMethod(purchaseRequest.DeliveryMethod),
	})
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("failed to add purchase"), err)
	}

	err = s.cartRepo.UpdateStatus(tx, purchase.CartID, entity.CartStatusInactive)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("failed to update cart status"), err)
	}

	adverts, err := s.advertRepo.GetByCartId(purchase.CartID, userId)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("failed to get adverts"), err)
	}

	for _, advert := range adverts {
		err = s.advertRepo.UpdateStatus(tx, advert.ID, entity.AdvertStatusReserved)
		if err != nil {
			return nil, entity.UsecaseWrap(errors.New("failed to update advert status"), err)
		}
	}

	return s.purchaseEntityToDTO(purchase)
}

func (s *PurchaseService) GetByUserId(userID uuid.UUID) ([]*dto.PurchaseResponse, error) {
	purchases, err := s.purchaseRepo.GetByUserId(userID)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("failed to get purchases"), err)
	}

	return s.purchaseEntitiesToDTO(purchases)
}

func (s *PurchaseService) purchaseEntitiesToDTO(purchases []*entity.Purchase) ([]*dto.PurchaseResponse, error) {
	ctx := context.Background()
	var purchaseDTOs []*dto.PurchaseResponse

	for _, purchase := range purchases {
		dto, err := s.purchaseEntityToDTO(purchase)
		if err != nil {
			logger := middleware.GetLogger(ctx)
			logger.Error("failed to convert purchase entity to DTO", zap.Error(err))
			return nil, err
		}
		purchaseDTOs = append(purchaseDTOs, dto)
	}

	return purchaseDTOs, nil
}
