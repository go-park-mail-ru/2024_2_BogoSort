package repository

import (
	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Seller interface {
	// AddSeller добавляет нового продавца в бд в рамках транзакции
	AddSeller(tx pgx.Tx, userID uuid.UUID) (uuid.UUID, error)

	// GetSellerByID возвращает продавца по его ID
	GetSellerByID(sellerID uuid.UUID) (*entity.Seller, error)

	// GetSellerByUserID возвращает продавца по ID пользователя
	GetSellerByUserID(userID uuid.UUID) (*entity.Seller, error)
}

var (
	ErrSellerNotFound      = errors.New("seller not found")
	ErrSellerAlreadyExists = errors.New("seller already exists")
)
