package service

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/services"
)

type AdvertService struct {
	AdvertRepo repository.AdvertRepository
	ImageService *services.ImageService
}

func NewAdvertService(advertRepo repository.AdvertRepository) *AdvertService {
	return &AdvertService{
		AdvertRepo: advertRepo,
	}
}

func (s *AdvertService) convertAdvertToDTO(advert *entity.Advert) (*dto.Advert, error) {
	return &dto.Advert{
		ID:          advert.ID,
		SellerId:    advert.SellerId,
		CategoryId:  advert.CategoryId,
		Title:       advert.Title,
		Description: advert.Description,
		Price:       advert.Price,
		ImageURL:    advert.ImageURL,
		Status:      dto.AdvertStatus(advert.Status),
		HasDelivery: advert.HasDelivery,
		Location:    advert.Location,
	}, nil
}

func (s *AdvertService) convertAdvertsToDTO(adverts []*entity.Advert) (*dto.AdvertResponseList, error) {
	reviews := make([]dto.Advert, 0)
	for _, advert := range adverts {
		toDTO, err := s.convertAdvertToDTO(advert)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, *toDTO)
	}
	return &dto.AdvertResponseList{Adverts: reviews}, nil
}

func (s *AdvertService) GetAdverts(limit, offset int) (*dto.AdvertResponseList, error) {
	adverts, err := s.AdvertRepo.GetAdverts(limit, offset)
	if err != nil {
		return nil, err
	}

	return s.convertAdvertsToDTO(adverts)
}
