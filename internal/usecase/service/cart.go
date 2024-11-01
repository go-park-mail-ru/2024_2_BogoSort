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
	cartRepo repository.Cart
	logger   *zap.Logger
}

func NewCartService(cartRepo repository.Cart, logger *zap.Logger) *CartService {
	return &CartService{
		cartRepo: cartRepo,
		logger:   logger,
	}
}

func (c *CartService) AddAdvertToUserCart(userID uuid.UUID, AdvertID uuid.UUID) error {
	cart, err := c.cartRepo.GetCartByUserID(userID)
	switch {
	case errors.Is(err, repository.ErrCartNotFound):
		cartID, err := c.cartRepo.CreateCart(userID)
		if err != nil {
			return entity.PSQLWrap(errors.New("error creating cart"), err)
		}
		return c.cartRepo.AddAdvertToCart(cartID, AdvertID)
	case err != nil:
		return entity.PSQLWrap(errors.New("error getting cart by user id"), err)
	}
	return c.cartRepo.AddAdvertToCart(cart.ID, AdvertID)
}

func (c *CartService) GetCartByID(cartID uuid.UUID) (dto.Cart, error) {
	adverts, err := c.cartRepo.GetAdvertsByCartID(cartID)
	if err != nil {
		return dto.Cart{}, entity.PSQLWrap(errors.New("error getting adverts by cart id"), err)
	}
	cart, err := c.cartRepo.GetCartByID(cartID)
	if err != nil {
		return dto.Cart{}, entity.PSQLWrap(errors.New("error getting cart by id"), err)
	}
	return dto.Cart{Adverts: adverts, ID: cart.ID, UserID: cart.UserID, Status: cart.Status}, nil
}

func (c *CartService) GetCartByUserID(userID uuid.UUID) (dto.Cart, error) {
	cart, err := c.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return dto.Cart{}, entity.PSQLWrap(errors.New("error getting cart by user id"), err)
	}
	return c.GetCartByID(cart.ID)
}
