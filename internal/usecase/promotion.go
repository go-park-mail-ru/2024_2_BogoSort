package usecase

import "github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"

type PromotionUseCase interface {
	// GetPromotionInfo возвращает информацию о продвижении
	GetPromotionInfo() (*entity.Promotion, error)
}
