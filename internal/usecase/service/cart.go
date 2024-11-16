package service

import (
	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CartService struct {
	cartRepo   repository.Cart
	advertRepo repository.AdvertRepository
	logger     *zap.Logger
}

func NewCartService(cartRepo repository.Cart, advertRepo repository.AdvertRepository, logger *zap.Logger) *CartService {
	return &CartService{
		cartRepo:   cartRepo,
		advertRepo: advertRepo,
		logger:     logger,
	}
}

func ConvertAdvertToResponse(advert entity.Advert) dto.AdvertResponse {
	return dto.AdvertResponse{ID: advert.ID, Title: advert.Title, Price: advert.Price, ImageId: advert.ImageId}
}

func (c *CartService) AddAdvertToUserCart(userID uuid.UUID, AdvertID uuid.UUID) error {
	cart, err := c.cartRepo.GetCartByUserID(userID)
	switch {
	case errors.Is(err, repository.ErrCartNotFound):
		cartID, err := c.cartRepo.CreateCart(userID)
		if err != nil {
			return entity.PSQLWrap(errors.New("error creating cart"), err)
		}
		advert, err := c.advertRepo.GetAdvertById(AdvertID)
		if err != nil {
			return entity.PSQLWrap(errors.New("error getting advert by id"), err)
		}
		if advert.Status != entity.AdvertStatusActive {
			return entity.UsecaseWrap(errors.New("advert is not active"), nil)
		}
		return c.cartRepo.AddAdvertToCart(cartID, AdvertID)
	case err != nil:
		return entity.PSQLWrap(errors.New("error getting cart by user id"), err)
	}
	return c.cartRepo.AddAdvertToCart(cart.ID, AdvertID)
}

func (c *CartService) DeleteAdvertFromCart(cartID uuid.UUID, AdvertID uuid.UUID) error {
	return c.cartRepo.DeleteAdvertFromCart(cartID, AdvertID)
}

func (c *CartService) GetCartByID(cartID uuid.UUID) (dto.Cart, error) {
	adverts, err := c.cartRepo.GetAdvertsByCartID(cartID)
	switch {
	case errors.Is(err, repository.ErrCartNotFound):
		return dto.Cart{}, entity.UsecaseWrap(errors.New("cart not found"), err)
	case err != nil:
		return dto.Cart{}, entity.PSQLWrap(errors.New("error getting adverts by cart id"), err)
	}
	cart, err := c.cartRepo.GetCartByID(cartID)
	if err != nil {
		return dto.Cart{}, entity.PSQLWrap(errors.New("error getting cart by id"), err)
	}

	advertsResponse := make([]dto.AdvertResponse, len(adverts))
	for i, advert := range adverts {
		advertsResponse[i] = ConvertAdvertToResponse(advert)
	}
	return dto.Cart{Adverts: advertsResponse, ID: cart.ID, UserID: cart.UserID, Status: cart.Status}, nil
}

func (c *CartService) GetCartByUserID(userID uuid.UUID) (dto.Cart, error) {
	cart, err := c.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return dto.Cart{}, entity.PSQLWrap(errors.New("error getting cart by user id"), err)
	}
	return c.GetCartByID(cart.ID)
}

func (c *CartService) CheckCartExists(userID uuid.UUID) (bool, error) {
	_, err := c.cartRepo.GetCartByUserID(userID)
	switch {
	case errors.Is(err, repository.ErrCartNotFound):
		return false, nil
	case err != nil:
		return false, entity.PSQLWrap(errors.New("error checking cart existence"), err)
	}
	return true, nil
}
