package usecase

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
)

type CategoryUseCase interface {
	// GetCategories возвращает все категории
	GetCategories() ([]*dto.Category, error)
}
