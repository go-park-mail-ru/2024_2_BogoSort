package service

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"net/http"
	"github.com/chai2010/webp"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type StaticService struct {
	logger     *zap.Logger
	staticRepo repository.StaticRepository
}

func NewStaticService(staticRepo repository.StaticRepository, logger *zap.Logger) *StaticService {
	return &StaticService{
		logger:     logger,
		staticRepo: staticRepo,
	}
}

func (s *StaticService) GetAvatar(staticID uuid.UUID) (string, error) {
	path, err := s.staticRepo.GetStatic(staticID)
	if err != nil {
		s.logger.Error("failed to get static", zap.Error(err), zap.String("static_id", staticID.String()))
		return "", err
	}
	return path, nil
}

func (s *StaticService) UploadStatic(reader io.ReadSeeker) (uuid.UUID, error) {
	size, err := reader.Seek(0, io.SeekEnd)
	if err != nil {
		return uuid.Nil, entity.UsecaseWrap(err, errors.New("ошибка при определении размера файла"))
	}
	if size > int64(s.staticRepo.GetMaxSize()) {
		return uuid.Nil, usecase.ErrStaticTooBigFile
	}
	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return uuid.Nil, entity.UsecaseWrap(err, errors.New("ошибка при возвращении io.ReadSeeker в начало файла"))
	}

	headerBytes := make([]byte, 512)
	_, err = reader.Read(headerBytes)
	if err != nil {
		return uuid.Nil, entity.UsecaseWrap(err, errors.New("ошибка при чтении заголовка файла"))
	}
	contentType := http.DetectContentType(headerBytes)
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		return uuid.Nil, usecase.ErrStaticNotImage
	}
	_, err = reader.Seek(0, io.SeekStart) 
	if err != nil {
		return uuid.Nil, entity.UsecaseWrap(err, errors.New("ошибка при возвращении io.ReadSeeker в начало файла"))
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		return uuid.Nil, usecase.ErrStaticNotImage
	}

	const minImageWidth, minImageHeight = 100, 100
	if img.Bounds().Dx() < minImageWidth || img.Bounds().Dy() < minImageHeight {
		return uuid.Nil, entity.UsecaseWrap(
			usecase.ErrStaticImageDimensions,
			fmt.Errorf(
				"изображение имеет размеры %dx%d, а должно быть как минимум %dx%d",
				img.Bounds().Dx(), img.Bounds().Dy(), minImageWidth, minImageHeight,
			),
		)
	}

	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	var squareImage *image.RGBA
	var start image.Point
	var squareSize int
	if width > height {
		start.X = (width - height) / 2
		squareSize = height
	} else {
		start.Y = (height - width) / 2
		squareSize = width
	}
	squareImage = image.NewRGBA(image.Rect(0, 0, squareSize, squareSize))
	draw.Draw(squareImage, squareImage.Bounds(), img, start, draw.Src)

	var out bytes.Buffer
	var opts webp.Options
	opts.Lossless = false
	opts.Quality = 60
	if err = webp.Encode(&out, squareImage, &opts); err != nil {
		return uuid.Nil, errors.Wrap(err, "ошибка при конвертации изображения в формат WEBP")
	}

	id, err := s.staticRepo.UploadStatic("avatars", uuid.New().String()+".webp", out.Bytes())
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *StaticService) GetStatic(id uuid.UUID) (string, error) {
	return s.staticRepo.GetStatic(id)
}
