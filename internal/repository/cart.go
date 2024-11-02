package repository

import (
	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/google/uuid"
)

type Cart interface {
	GetAdvertsByCartID(cartID uuid.UUID) ([]entity.Advert, error)
	AddAdvertToCart(cartID uuid.UUID, AdvertID uuid.UUID) error
	DeleteAdvertFromCart(cartID uuid.UUID, AdvertID uuid.UUID) error
	UpdateCartStatus(cartID uuid.UUID, status entity.CartStatus) error
	GetCartByUserID(userID uuid.UUID) (entity.Cart, error)
	CreateCart(userID uuid.UUID) (uuid.UUID, error)
	GetCartByID(cartID uuid.UUID) (entity.Cart, error)
}

var (
	ErrCartNotFound         = errors.New("cart not found")
	ErrCartOrAdvertNotFound = errors.New("cart or advert not found")
	ErrCartAlreadyExists    = errors.New("cart already exists")
)
