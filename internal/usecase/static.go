package usecase

import (
	"io"
	"errors"

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

var ErrStaticFileNotFound = errors.New("static file not found")
var ErrStaticTooBigFile = errors.New("static file too big")
var ErrStaticNotImage = errors.New("static file is not image")
var ErrStaticImageDimensions = errors.New("static image dimensions are invalid")
var ErrStaticNotFound = errors.New("static not found")