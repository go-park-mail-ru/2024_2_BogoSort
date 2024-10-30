package repository

import (
	"errors"
	"github.com/google/uuid"
)

type StaticRepository interface {
	GetStatic(staticID uuid.UUID) (string, error)
	UploadStatic(path, filename string, data []byte) (uuid.UUID, error)
}

var (
	ErrStaticNotFound = errors.New("статика не найдена")
	ErrStaticTooLarge = errors.New("статика слишком большая")
)
