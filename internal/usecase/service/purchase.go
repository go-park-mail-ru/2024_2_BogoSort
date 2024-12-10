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
	var advertCards []dto.PreviewAdvertCard
	for _, advert := range purchase.Adverts {
		advertCards = append(advertCards, convertAdvertToPreviewCard(advert))
	}
	
	return &dto.PurchaseResponse{
		ID:             purchase.ID,
		SellerID:       purchase.SellerID,
		CustomerID:     purchase.CustomerID,
		Address:        purchase.Address,
		Status:         dto.PurchaseStatus(purchase.Status),
		PaymentMethod:  dto.PaymentMethod(purchase.PaymentMethod),
		DeliveryMethod: dto.DeliveryMethod(purchase.DeliveryMethod),
		Adverts:        advertCards,
	}, nil
}

func (s *PurchaseService) Add(purchaseRequest dto.PurchaseRequest, userId uuid.UUID) ([]*dto.PurchaseResponse, error) {
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

	cart, err := s.cartRepo.GetById(purchaseRequest.CartID)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("failed to get cart"), err)
	}

	if cart.UserID != userId {
		return nil, entity.UsecaseWrap(errors.New("cart does not belong to user"), nil)
	}

	var purchases []*dto.PurchaseResponse

	for _, cartPurchase := range cart.CartPurchases {
		purchase, err := s.purchaseRepo.Add(tx, &entity.Purchase{
			SellerID:       cartPurchase.SellerID,
			CustomerID:     userId,
			CartID:         purchaseRequest.CartID,
			Address:        purchaseRequest.Address,
			Status:         entity.StatusPending,
			PaymentMethod:  entity.PaymentMethod(purchaseRequest.PaymentMethod),
			DeliveryMethod: entity.DeliveryMethod(purchaseRequest.DeliveryMethod),
			Adverts:        cartPurchase.Adverts,
		})
		if err != nil {
			return nil, entity.UsecaseWrap(errors.New("failed to add purchase"), err)
		}

		for _, advert := range cartPurchase.Adverts {
			err = s.advertRepo.UpdateStatus(tx, advert.ID, entity.AdvertStatusReserved)
			if err != nil {
				return nil, entity.UsecaseWrap(errors.New("failed to update advert status"), err)
			}
		}

		purchaseDTO, err := s.purchaseEntityToDTO(purchase)
		if err != nil {
			return nil, err
		}
		purchases = append(purchases, purchaseDTO)
	}

	err = s.cartRepo.UpdateStatus(tx, cart.ID, entity.CartStatusInactive)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("failed to update cart status"), err)
	}

	logger := middleware.GetLogger(ctx)
	logger.Info("purchases added", zap.Any("purchases", purchases))

	return purchases, nil
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

func convertAdvertToPreviewCard(advert entity.Advert) dto.PreviewAdvertCard {
	return dto.PreviewAdvertCard{
		Preview: dto.PreviewAdvert{
			ID:          advert.ID,
			SellerId:    advert.SellerId,
			CategoryId:  advert.CategoryId,
			Title:       advert.Title,
			Price:       advert.Price,
			ImageId:     advert.ImageId,
			Status:      dto.AdvertStatus(advert.Status),
			Location:    advert.Location,
			HasDelivery: advert.HasDelivery,
		},
		IsSaved:  advert.IsSaved,
		IsViewed: advert.IsViewed,
	}
}
