package repository

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
)

// PromotionRepository интерфейс для работы с продвижением
type PromotionRepository interface {
	// GetPromotionInfo возвращает информацию о продвижении
	GetPromotionInfo() (*entity.Promotion, error)
}
