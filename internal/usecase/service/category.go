package service

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
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

func (s *CategoryService) categoryEntityToDTO(category *entity.Category) (*dto.Category, error) {
	categoryDTO := dto.Category{
		ID:   category.ID,
		Title: category.Title,
	}

	return &categoryDTO, nil
}

func (s *CategoryService) categoryEntitiesToDTO(categories []*entity.Category) ([]*dto.Category, error) {
	categoryDTOs := make([]*dto.Category, 0, len(categories))
	for _, category := range categories {
		categoryDTO, err := s.categoryEntityToDTO(category)
		if err != nil {
			s.logger.Error("error converting category entity to dto", zap.Error(err))
			return nil, entity.UsecaseWrap(err, err)
		}
		categoryDTOs = append(categoryDTOs, categoryDTO)
	}

	return categoryDTOs, nil
}

func (s *CategoryService) GetCategories() ([]*dto.Category, error) {
	categories, err := s.categoryRepo.GetCategories()
	if err != nil {
		s.logger.Error("error getting categories", zap.Error(err))
		return nil, entity.UsecaseWrap(err, err)
	}

	return s.categoryEntitiesToDTO(categories)
}
