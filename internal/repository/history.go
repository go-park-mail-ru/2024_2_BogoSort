package repository

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/google/uuid"
)

type HistoryRepository interface {
	// Получение истории изменения цены объявления
	GetAdvertPriceHistory(advertID uuid.UUID) ([]*entity.PriceHistory, error)

	// Добавление истории изменения цены объявления
	AddAdvertPriceChange(advertID uuid.UUID, oldPrice int, newPrice int) error
}
