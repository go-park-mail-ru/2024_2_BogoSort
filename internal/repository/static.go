package repository

import (
	"errors"

	"github.com/google/uuid"
)

type StaticRepository interface {
	// Get возвращает путь к статическому файлу по его ID
	Get(staticID uuid.UUID) (string, error)

	// Upload загружает статический файл и возвращает его ID
	Upload(path, filename string, data []byte) (uuid.UUID, error)

	// GetMaxSize возвращает максимальный размер файла
	GetMaxSize() int
}

var (
	ErrStaticNotFound = errors.New("статика не найдена")
	ErrStaticTooLarge = errors.New("статика слишком большая")
)
