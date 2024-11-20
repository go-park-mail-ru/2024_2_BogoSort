package repository

import (
	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Seller interface {
	// Add добавляет нового продавца в бд в рамках транзакции
	Add(tx pgx.Tx, userID uuid.UUID) (uuid.UUID, error)

	// GetById возвращает продавца по его ID
	GetById(sellerID uuid.UUID) (*entity.Seller, error)

	// GetByUserId возвращает продавца по ID пользователя
	GetByUserId(userID uuid.UUID) (*entity.Seller, error)
}

var (
	ErrSellerNotFound      = errors.New("seller not found")
	ErrSellerAlreadyExists = errors.New("seller already exists")
)
