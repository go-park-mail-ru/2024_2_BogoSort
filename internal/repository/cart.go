package repository

import (
	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Cart interface {
	GetAdvertsByCartId(cartID uuid.UUID) ([]entity.Advert, error)
	AddAdvert(cartID uuid.UUID, AdvertID uuid.UUID) error
	DeleteAdvert(cartID uuid.UUID, AdvertID uuid.UUID) error
	UpdateStatus(tx pgx.Tx, cartID uuid.UUID, status entity.CartStatus) error
	GetByUserId(userID uuid.UUID) (entity.Cart, error)
	Create(userID uuid.UUID) (uuid.UUID, error)
	GetById(cartID uuid.UUID) (entity.Cart, error)
}

var (
	ErrCartNotFound         = errors.New("cart not found")
	ErrCartOrAdvertNotFound = errors.New("cart or advert not found")
	ErrCartAlreadyExists    = errors.New("cart already exists")
)
