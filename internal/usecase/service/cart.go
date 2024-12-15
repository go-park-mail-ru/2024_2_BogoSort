package service

import (
	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
)

type CartService struct {
	cartRepo   repository.Cart
	advertRepo repository.AdvertRepository
}

func NewCartService(cartRepo repository.Cart, advertRepo repository.AdvertRepository) *CartService {
	return &CartService{
		cartRepo:   cartRepo,
		advertRepo: advertRepo,
	}
}

func convertAdvertToResponse(advert entity.Advert) dto.PreviewAdvertCard {
	return dto.PreviewAdvertCard{
		Preview: dto.PreviewAdvert{
			ID:          advert.ID,
			Title:       advert.Title,
			CategoryId:  advert.CategoryId,
			SellerId:    advert.SellerId,
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

func (c *CartService) AddAdvert(userID, advertID uuid.UUID) error {
	cart, err := c.cartRepo.GetByUserId(userID)
	switch {
	case errors.Is(err, repository.ErrCartNotFound):
		cartID, err := c.cartRepo.Create(userID)
		if err != nil {
			return entity.PSQLWrap(errors.New("error creating cart"), err)
		}
		advert, err := c.advertRepo.GetById(advertID, userID)
		if err != nil {
			return entity.PSQLWrap(errors.New("error getting advert by id"), err)
		}
		if advert.Status != entity.AdvertStatusActive {
			return entity.UsecaseWrap(errors.New("advert is not active"), nil)
		}
		return c.cartRepo.AddAdvert(cartID, advertID)
	case err != nil:
		return entity.PSQLWrap(errors.New("error getting cart by user id"), err)
	}
	return c.cartRepo.AddAdvert(cart.ID, advertID)
}

func (c *CartService) DeleteAdvert(cartID uuid.UUID, advertID uuid.UUID) error {
	return c.cartRepo.DeleteAdvert(cartID, advertID)
}

func (c *CartService) GetById(cartID uuid.UUID) (dto.Cart, error) {
	adverts, err := c.cartRepo.GetAdvertsByCartId(cartID)

	switch {
	case errors.Is(err, repository.ErrCartNotFound):
		return dto.Cart{}, entity.UsecaseWrap(errors.New("cart not found"), err)
	case err != nil:
		return dto.Cart{}, entity.PSQLWrap(errors.New("error getting adverts by cart id"), err)
	}
	cart, err := c.cartRepo.GetById(cartID)
	if err != nil {
		return dto.Cart{}, entity.PSQLWrap(errors.New("error getting cart by id"), err)
	}

	advertsBySeller := make(map[uuid.UUID][]dto.PreviewAdvertCard)
	for _, advert := range adverts {
		advertsBySeller[advert.SellerId] = append(
			advertsBySeller[advert.SellerId],
			convertAdvertToResponse(advert),
		)
	}

	cartPurchases := make([]dto.CartPurchase, 0, len(advertsBySeller))
	for sellerID, sellerAdverts := range advertsBySeller {
		cartPurchases = append(cartPurchases, dto.CartPurchase{
			SellerID: sellerID,
			Adverts:  sellerAdverts,
		})
	}

	return dto.Cart{
		ID:            cart.ID,
		UserID:        cart.UserID,
		Status:        cart.Status,
		CartPurchases: cartPurchases,
	}, nil
}

func (c *CartService) GetByUserId(userID uuid.UUID) (dto.Cart, error) {
	cart, err := c.cartRepo.GetByUserId(userID)
	switch {
	case errors.Is(err, repository.ErrCartNotFound):
		return dto.Cart{}, entity.UsecaseWrap(errors.New("cart not found"), err)
	case err != nil:
		return dto.Cart{}, entity.PSQLWrap(errors.New("error getting cart by user id"), err)
	}

	return dto.Cart{
		ID:            cart.ID,
		UserID:        cart.UserID,
		Status:        cart.Status,
		CartPurchases: []dto.CartPurchase{},
	}, nil
}

func (c *CartService) CheckExists(userID uuid.UUID) (uuid.UUID, error) {
	cart, err := c.cartRepo.GetByUserId(userID)
	switch {
	case errors.Is(err, repository.ErrCartNotFound):
		return uuid.Nil, nil
	case err != nil:
		return uuid.Nil, entity.PSQLWrap(errors.New("error checking cart existence"), err)
	}
	return cart.ID, nil
}
