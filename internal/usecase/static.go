package usecase

import (
	"errors"
	"io"

	"github.com/google/uuid"
)

type StaticUseCase interface {
	// GetAvatar возвращает url аватара по id
	GetAvatar(staticID uuid.UUID) (string, error)

	// UploadStatic загружает файл и возвращает id загруженного файла
	UploadStatic(data io.ReadSeeker) (uuid.UUID, error)

	// GetStatic возвращает url статики по id
	GetStatic(id uuid.UUID) (string, error)

	// GetStaticFile возвращает файл по uri
	GetStaticFile(uri string) (io.ReadSeeker, error)
}

var (
	ErrStaticFileNotFound    = errors.New("static file not found")
	ErrStaticTooBigFile      = errors.New("static file too big")
	ErrStaticNotImage        = errors.New("static file is not image")
	ErrStaticImageDimensions = errors.New("static image dimensions are invalid")
	ErrStaticNotFound        = errors.New("static not found")
)
