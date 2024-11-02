package usecase

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/jackc/pgx/v5"
)

type Purchase interface {
	// AddPurchase добавляет покупку в базу данных
	AddPurchase(tx pgx.Tx, purchaseRequest dto.PurchaseRequest) (*dto.PurchaseResponse, error)
}
