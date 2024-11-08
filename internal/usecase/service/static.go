package service

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
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

func (s *StaticService) UploadFile(data []byte) (uuid.UUID, error) {
	contentType := http.DetectContentType(data)

	if contentType != "image/jpeg" && contentType != "image/png" {
		s.logger.Error("file is not an image", zap.String("content_type", contentType))
		return uuid.Nil, entity.UsecaseWrap(
			errors.New("file is not an image"),
		)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		s.logger.Error("error decoding image", zap.Error(err))
		return uuid.Nil, entity.UsecaseWrap(
			errors.New("file is not an image"),
		)
	}

	const minImageWidth, minImageHeight = 100, 100
	if img.Bounds().Dx() < minImageWidth || img.Bounds().Dy() < minImageHeight {
		s.logger.Error("image size is less than required", zap.Int("width", img.Bounds().Dx()), zap.Int("height", img.Bounds().Dy()))
		return uuid.Nil, entity.UsecaseWrap(
			errors.New(fmt.Sprintf("image dimensions are less than %d x %d", minImageWidth, minImageHeight)),
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
	var opts jpeg.Options
	opts.Quality = 60
	err = jpeg.Encode(&out, squareImage, &opts)
	if err != nil {
		s.logger.Error("error encoding image", zap.Error(err))
		return uuid.Nil, entity.UsecaseWrap(
			errors.New("error processing image"),
		)
	}

	id, err := s.staticRepo.UploadStatic("", uuid.New().String()+".jpg", out.Bytes())
	if err != nil {
		s.logger.Error("error uploading static", zap.Error(err))
		return uuid.Nil, err
	}

	return id, nil
}

func (s *StaticService) GetStaticURL(id uuid.UUID) (string, error) {
	return s.staticRepo.GetStatic(id)
}
