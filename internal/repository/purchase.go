package repository

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/google/uuid"
)


type PurchaseRepository interface {
	// BeginTransaction начинает транзакцию
	BeginTransaction() (pgx.Tx, error)

	// AddPurchase создает запись о покупке
	Add(tx pgx.Tx, purchase *entity.Purchase) (*entity.Purchase, error)

	// GetPurchasesByUserID получает покупки по UserID
	GetByUserId(userID uuid.UUID) ([]*entity.Purchase, error)
}
