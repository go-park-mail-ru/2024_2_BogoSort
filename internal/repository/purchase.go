package repository

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/jackc/pgx/v5"
)


type PurchaseRepository interface {
	// BeginTransaction начинает транзакцию
	BeginTransaction() (pgx.Tx, error)

	// AddPurchase создает запись о покупке
	AddPurchase(tx pgx.Tx, purchase *entity.Purchase) (*entity.Purchase, error)
}
