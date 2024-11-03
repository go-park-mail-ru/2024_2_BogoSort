package repository

import "github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"

type CategoryRepository interface {
	// GetCategories возвращает все категории
	GetCategories() ([]*entity.Category, error)
}
