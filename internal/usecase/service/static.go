package service

import (
	"bytes"
	"fmt"
	"github.com/chai2010/webp"
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"

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
	path, err := s.staticRepo.Get(staticID)
	if err != nil {
		s.logger.Error("failed to get static", zap.Error(err), zap.String("static_id", staticID.String()))
		return "", err
	}
	return path, nil
}

func (s *StaticService) UploadStatic(reader io.ReadSeeker) (uuid.UUID, error) {
	size, err := reader.Seek(0, io.SeekEnd)
	if err != nil {
		return uuid.Nil, entity.UsecaseWrap(err, errors.New("error determining file size"))
	}
	if size > int64(s.staticRepo.GetMaxSize()) {
		return uuid.Nil, usecase.ErrStaticTooBigFile
	}
	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return uuid.Nil, entity.UsecaseWrap(err, errors.New("error returning io.ReadSeeker to the start of the file"))
	}

	headerBytes := make([]byte, 512)
	_, err = reader.Read(headerBytes)
	if err != nil {
		return uuid.Nil, entity.UsecaseWrap(err, errors.New("error reading file header"))
	}
	contentType := http.DetectContentType(headerBytes)
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		return uuid.Nil, usecase.ErrStaticNotImage
	}
	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return uuid.Nil, entity.UsecaseWrap(err, errors.New("error returning io.ReadSeeker to the start of the file"))
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
				"image dimensions are %dx%d, but must be at least %dx%d",
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
		return uuid.Nil, errors.Wrap(err, "error converting image to WEBP format")
	}

	id, err := s.staticRepo.Upload("images", uuid.New().String()+".webp", out.Bytes())
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *StaticService) GetStaticFile(staticURI string) (io.ReadSeeker, error) {
	absolutePath, err := filepath.Abs(staticURI)
	zap.L().Info("Getting static file", zap.String("uri", staticURI))
	zap.L().Info("Absolute path", zap.String("path", absolutePath))
	if err != nil {
		return nil, entity.UsecaseWrap(err, errors.New("error determining absolute path"))
	}

	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return nil, usecase.ErrStaticNotFound
	}

	file, err := os.Open(absolutePath)
	if err != nil {
		return nil, entity.UsecaseWrap(err, errors.New("error opening file"))
	}
	return file, nil
}

func (s *StaticService) GetStatic(id uuid.UUID) (string, error) {
	return s.staticRepo.Get(id)
}
