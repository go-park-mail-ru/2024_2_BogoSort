package service

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/services"
	"go.uber.org/zap"
	"errors"
	"github.com/google/uuid"
	"database/sql"
)

var (
    ErrAdvertNotFound     = errors.New("advert not found")
    ErrAdvertBadRequest   = errors.New("bad request: invalid advert data")
    ErrAdvertAlreadyExists = errors.New("advert already exists")
)

type AdvertService struct {
	AdvertRepo   repository.AdvertRepository
	ImageService *services.ImageService
	logger       *zap.Logger
}

func NewAdvertService(advertRepo repository.AdvertRepository, logger *zap.Logger) (*AdvertService, error) {
	return &AdvertService{
		AdvertRepo: advertRepo,
		logger:     logger,
	}, nil
}

func (s *AdvertService) convertAdvertToDTO(advert *entity.Advert) (*dto.Advert, error) {
	return &dto.Advert{
		ID:          advert.ID,
		SellerId:    advert.SellerId,
		CategoryId:  advert.CategoryId,
		Title:       advert.Title,
		Description: advert.Description,
		Price:       advert.Price,
		ImageURL:    advert.ImageURL.String,
		Status:      dto.AdvertStatus(advert.Status),
		HasDelivery: advert.HasDelivery,
		Location:    advert.Location,
	}, nil
}

func (s *AdvertService) convertAdvertToEntity(advert *dto.Advert) (*entity.Advert, error) {
    return &entity.Advert{
        ID:          advert.ID,
        SellerId:    advert.SellerId,
        CategoryId:  advert.CategoryId,
        Title:       advert.Title,
        Description: advert.Description,
        Price:       advert.Price,
        ImageURL:    sql.NullString{String: advert.ImageURL, Valid: advert.ImageURL != ""},
        Status:      entity.AdvertStatus(advert.Status),
        HasDelivery: advert.HasDelivery,
        Location:    advert.Location,
    }, nil
}

func (s *AdvertService) convertAdvertsToDTO(adverts []*entity.Advert) ([]*dto.Advert, error) {
    dtoAdverts := make([]*dto.Advert, 0, len(adverts))
    for _, advert := range adverts {
        toDTO, err := s.convertAdvertToDTO(advert)
        if err != nil {
            return nil, err
        }
        dtoAdverts = append(dtoAdverts, toDTO)
    }
    return dtoAdverts, nil
}

func (s *AdvertService) GetAdverts(limit, offset int) ([]*dto.Advert, error) {
    adverts, err := s.AdvertRepo.GetAdverts(limit, offset)
    if err != nil {
        return nil, err
    }
    return s.convertAdvertsToDTO(adverts)
}

func (s *AdvertService) GetAdvertsByUserId(userId uuid.UUID) ([]*dto.Advert, error) {
    adverts, err := s.AdvertRepo.GetAdvertsByUserId(userId)
    if err != nil {
        return nil, err
    }
    return s.convertAdvertsToDTO(adverts)
}

func (s *AdvertService) GetSavedAdvertsByUserId(userId uuid.UUID) ([]*dto.Advert, error) {
    savedAdverts, err := s.AdvertRepo.GetSavedAdvertsByUserId(userId)
    if err != nil {
        return nil, err
    }
    return s.convertAdvertsToDTO(savedAdverts)
}

func (s *AdvertService) GetAdvertsByCartId(cartId uuid.UUID) ([]*dto.Advert, error) {
    adverts, err := s.AdvertRepo.GetAdvertsByCartId(cartId)
    if err != nil {
        return nil, err
    }
    return s.convertAdvertsToDTO(adverts)
}

func (s *AdvertService) GetAdvertById(advertId uuid.UUID) (*dto.Advert, error) {
    advert, err := s.AdvertRepo.GetAdvertById(advertId)
    if err != nil {
        if errors.Is(err, repository.ErrAdvertNotFound) {
            return nil, ErrAdvertNotFound
        }
        return nil, err
    }
    return s.convertAdvertToDTO(advert)
}

func (s *AdvertService) AddAdvert(advert *dto.Advert) (*dto.Advert, error) {
    entityAdvert, err := s.convertAdvertToEntity(advert)
    if err != nil {
        return nil, ErrAdvertBadRequest
    }

    createdAdvert, err := s.AdvertRepo.AddAdvert(entityAdvert)
    if err != nil {
        if errors.Is(err, repository.ErrAdvertAlreadyExists) {
            return nil, ErrAdvertAlreadyExists
        }
        return nil, err
    }

    return s.convertAdvertToDTO(createdAdvert)
}

func (s *AdvertService) UpdateAdvert(advert *dto.Advert) error {
    entityAdvert, err := s.convertAdvertToEntity(advert)
    if err != nil {
        return ErrAdvertBadRequest
    }

    if err := s.AdvertRepo.UpdateAdvert(entityAdvert); err != nil {
        if errors.Is(err, repository.ErrAdvertNotFound) {
            return ErrAdvertNotFound
        }
        return err
    }

    return nil
}

func (s *AdvertService) DeleteAdvertById(advertId uuid.UUID) error {
    if err := s.AdvertRepo.DeleteAdvertById(advertId); err != nil {
        if errors.Is(err, repository.ErrAdvertNotFound) {
            return ErrAdvertNotFound
        }
        return err
    }
    return nil
}
