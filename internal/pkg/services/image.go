package services

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
)

type ImageService struct {
	images map[uint]string
	mu     sync.RWMutex
}

func NewImageService() *ImageService {
	return &ImageService{
		images: make(map[uint]string),
	}
}

func (s *ImageService) GetImageURL(advertID uint) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if url, ok := s.images[advertID]; ok {
		return url, nil
	}

	return "", errors.New("изображение не найдено")
}

func (s *ImageService) SetImageURL(advertID uint, url string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.images[advertID] = url
}

func (s *ImageService) ValidateImage(imagePath string) error {
	fullPath := filepath.Join("static", imagePath)
	_, err := os.Stat(fullPath)

	return errors.Wrap(err, "failed to validate image")
}
