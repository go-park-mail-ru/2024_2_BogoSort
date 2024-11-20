package usecase

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
)

type CategoryUseCase interface {
	// Get возвращает все категории
	Get() ([]*entity.Category, error)
}
