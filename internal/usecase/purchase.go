package usecase

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/google/uuid"
)

type Purchase interface {
	// AddPurchase добавляет покупку в базу данных
	AddPurchase(purchaseRequest dto.PurchaseRequest) (*dto.PurchaseResponse, error)
	
	// GetPurchasesByUserID получает покупки по UserID
	GetPurchasesByUserID(userID uuid.UUID) ([]*dto.PurchaseResponse, error)
}
