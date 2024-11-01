package service

import (
	"errors"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strings"
)

var (
	ErrAdvertNotFound      = errors.New("advert not found")
	ErrAdvertBadRequest    = errors.New("bad request: invalid advert data")
	ErrAdvertAlreadyExists = errors.New("advert already exists")
)

type AdvertService struct {
	AdvertRepo repository.AdvertRepository
	StaticRepo repository.StaticRepository
	logger     *zap.Logger
}

func NewAdvertService(advertRepo repository.AdvertRepository,
	staticRepo repository.StaticRepository,
	logger *zap.Logger) *AdvertService {
	return &AdvertService{
		AdvertRepo: advertRepo,
		StaticRepo: staticRepo,
	}
}

func (s *AdvertService) advertEntityToDTO(advert *entity.Advert) (*dto.Advert, error) {
	var posterURL string

	if !advert.ImageURL.Valid {
		posterURL = ""
	} else {
		var err error
		posterURL, err = s.StaticRepo.GetStatic(advert.ImageURL.UUID)
		if err != nil {
			posterURL = ""
		}
	}

	advertDTO := dto.Advert{
		ID:          advert.ID,
		SellerId:    advert.SellerId,
		CategoryId:  advert.CategoryId,
		Title:       advert.Title,
		Description: advert.Description,
		Price:       advert.Price,
		ImageURL:    posterURL,
		Status:      dto.AdvertStatus(advert.Status),
		HasDelivery: advert.HasDelivery,
		Location:    advert.Location,
	}

	return &advertDTO, nil
}

func (s *AdvertService) advertEntitiesToDTO(adverts []*entity.Advert) ([]*dto.Advert, error) {
	dtoAdverts := make([]*dto.Advert, 0, len(adverts))
	for _, advert := range adverts {
		advertDTO, err := s.advertEntityToDTO(advert)
		if err != nil {
			return nil, entity.UsecaseWrap(err, err)
		}
		dtoAdverts = append(dtoAdverts, advertDTO)
	}
	return dtoAdverts, nil
}

func (s *AdvertService) GetAdverts(limit, offset int) ([]*dto.Advert, error) {
	adverts, err := s.AdvertRepo.GetAdverts(limit, offset)
	if err != nil {
		return nil, entity.UsecaseWrap(err, err)
	}
	return s.advertEntitiesToDTO(adverts)
}

func (s *AdvertService) GetAdvertsBySellerId(sellerId uuid.UUID) ([]*dto.Advert, error) {
	adverts, err := s.AdvertRepo.GetAdvertsBySellerId(sellerId)
	if err != nil {
		return nil, entity.UsecaseWrap(err, err)
	}
	return s.advertEntitiesToDTO(adverts)
}

func (s *AdvertService) GetSavedAdvertsByUserId(userId uuid.UUID) ([]*dto.Advert, error) {
	savedAdverts, err := s.AdvertRepo.GetSavedAdvertsByUserId(userId)
	if err != nil {
		return nil, entity.UsecaseWrap(err, err)
	}
	return s.advertEntitiesToDTO(savedAdverts)
}

func (s *AdvertService) GetAdvertsByCartId(cartId uuid.UUID) ([]*dto.Advert, error) {
	adverts, err := s.AdvertRepo.GetAdvertsByCartId(cartId)
	if err != nil {
		return nil, entity.UsecaseWrap(err, err)
	}
	return s.advertEntitiesToDTO(adverts)
}

func (s *AdvertService) GetAdvertById(advertId uuid.UUID) (*dto.Advert, error) {
	advert, err := s.AdvertRepo.GetAdvertById(advertId)

	if err != nil {
		if errors.Is(err, repository.ErrAdvertNotFound) {
			return nil, entity.UsecaseWrap(ErrAdvertNotFound, ErrAdvertNotFound)
		}
		return nil, entity.UsecaseWrap(err, err)
	}

	return s.advertEntityToDTO(advert)
}

func (s *AdvertService) AddAdvert(advert *dto.Advert) (*dto.Advert, error) {
	if err := entity.ValidateAdvert(advert.Title, 
        advert.Description, 
        advert.Location, 
        string(advert.Status), 
        int(advert.Price)); err != nil {
		return nil, entity.UsecaseWrap(ErrAdvertBadRequest, ErrAdvertBadRequest)
	}

	entityAdvert, err := s.AdvertRepo.AddAdvert(&entity.Advert{
		SellerId:    advert.SellerId,
		CategoryId:  advert.CategoryId,
		Title:       strings.TrimSpace(advert.Title),
		Description: strings.TrimSpace(advert.Description),
		Price:       advert.Price,
		ImageURL:    uuid.NullUUID{UUID: uuid.MustParse(advert.ImageURL), Valid: true},
		Status:      entity.AdvertStatus(advert.Status),
		HasDelivery: advert.HasDelivery,
		Location:    advert.Location,
	})
	if err != nil {
		return nil, entity.UsecaseWrap(ErrAdvertBadRequest, ErrAdvertBadRequest)
	}

	return s.advertEntityToDTO(entityAdvert)
}

func (s *AdvertService) UpdateAdvert(advert *dto.Advert) error {
	if err := entity.ValidateAdvert(advert.Title, 
        advert.Description, 
        advert.Location, 
        string(advert.Status), 
        int(advert.Price)); err != nil {
		return entity.UsecaseWrap(ErrAdvertBadRequest, ErrAdvertBadRequest)
	}

	err := s.AdvertRepo.UpdateAdvert(&entity.Advert{
		ID:          advert.ID,
		SellerId:    advert.SellerId,
		CategoryId:  advert.CategoryId,
		Title:       strings.TrimSpace(advert.Title),
		Description: strings.TrimSpace(advert.Description),
		Price:       advert.Price,
		ImageURL:    uuid.NullUUID{UUID: uuid.MustParse(advert.ImageURL), Valid: true},
		Status:      entity.AdvertStatus(advert.Status),
		HasDelivery: advert.HasDelivery,
		Location:    advert.Location,
	})
	if err != nil {
		return entity.UsecaseWrap(ErrAdvertBadRequest, ErrAdvertBadRequest)
	}

	return nil
}

func (s *AdvertService) DeleteAdvertById(advertId uuid.UUID) error {
	if err := s.AdvertRepo.DeleteAdvertById(advertId); err != nil {
		if errors.Is(err, repository.ErrAdvertNotFound) {
			return entity.UsecaseWrap(ErrAdvertNotFound, ErrAdvertNotFound)
		}
		return entity.UsecaseWrap(err, err)
	}

	return nil
}

func (s *AdvertService) UpdateAdvertStatus(advertId uuid.UUID, status string) error {
	if err := s.AdvertRepo.UpdateAdvertStatus(advertId, status); err != nil {
		if errors.Is(err, repository.ErrAdvertNotFound) {
			return entity.UsecaseWrap(ErrAdvertNotFound, ErrAdvertNotFound)
		}
		return entity.UsecaseWrap(err, err)
	}

	return nil
}

func (s *AdvertService) GetAdvertsByCategoryId(categoryId uuid.UUID) ([]*dto.Advert, error) {
	adverts, err := s.AdvertRepo.GetAdvertsByCategoryId(categoryId)
	if err != nil {
		return nil, entity.UsecaseWrap(err, err)
	}

	return s.advertEntitiesToDTO(adverts)
}

func (s *AdvertService) UploadImage(advertId uuid.UUID, imageId uuid.UUID) error {
	if err := s.AdvertRepo.UploadImage(advertId, imageId); err != nil {
		return entity.UsecaseWrap(ErrAdvertBadRequest, ErrAdvertBadRequest)
	}

	return nil
}
