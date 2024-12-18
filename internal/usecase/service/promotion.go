package service

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
)

type PromotionService struct {
	promotionRepo repository.PromotionRepository
}

func NewPromotionService(promotionRepo repository.PromotionRepository) *PromotionService {
	return &PromotionService{promotionRepo: promotionRepo}
}

func (s *PromotionService) GetPromotionInfo() (*entity.Promotion, error) {
	promotion, err := s.promotionRepo.GetPromotionInfo()
	if err != nil {
		return nil, entity.UsecaseWrap(err, err)
	}

	return promotion, nil
}
