package repository

import (
	"errors"
	"github.com/google/uuid"
)

type StaticRepository interface {
	// GetStatic возвращает путь к статическому файлу по его ID
	GetStatic(staticID uuid.UUID) (string, error)
	
	// UploadStatic загружает статический файл и возвращает его ID
	UploadStatic(path, filename string, data []byte) (uuid.UUID, error)

	// GetMaxSize возвращает максимальный размер файла
	GetMaxSize() int
}

var (
	ErrStaticNotFound = errors.New("статика не найдена")
	ErrStaticTooLarge = errors.New("статика слишком большая")
)
