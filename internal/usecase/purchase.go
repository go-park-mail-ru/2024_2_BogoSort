package usecase

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/google/uuid"
)

type Purchase interface {
	// Add добавляет покупку в базу данных
	Add(purchaseRequest dto.PurchaseRequest, userId uuid.UUID) (*dto.PurchaseResponse, error)

	// GetByUserId получает покупки по UserID
	GetByUserId(userID uuid.UUID) ([]*dto.PurchaseResponse, error)
}
