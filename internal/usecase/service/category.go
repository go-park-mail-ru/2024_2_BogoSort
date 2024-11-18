package service

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"go.uber.org/zap"
)

type CategoryService struct {
	categoryRepo repository.CategoryRepository
	logger *zap.Logger
}

func NewCategoryService(categoryRepo repository.CategoryRepository, logger *zap.Logger) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		logger: logger,
	}
}

func (s *CategoryService) Get() ([]*entity.Category, error) {
	categories, err := s.categoryRepo.Get()
	if err != nil {
		s.logger.Error("error getting categories", zap.Error(err))
		return nil, entity.UsecaseWrap(err, err)
	}

	return categories, nil
}
