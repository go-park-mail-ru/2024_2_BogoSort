package service

import (
	"errors"
	"strings"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrAdvertNotFound      = errors.New("advert not found")
	ErrAdvertBadRequest    = errors.New("bad request: invalid advert data")
	ErrAdvertAlreadyExists = errors.New("advert already exists")
	ErrForbidden           = errors.New("forbidden: you cannot modify this source")
)

type AdvertService struct {
	advertRepo repository.AdvertRepository
	staticRepo repository.StaticRepository
	sellerRepo repository.Seller
	logger     *zap.Logger
}

func NewAdvertService(advertRepo repository.AdvertRepository,
	staticRepo repository.StaticRepository,
	sellerRepo repository.Seller,
	logger *zap.Logger) *AdvertService {
	return &AdvertService{
		advertRepo: advertRepo,
		staticRepo: staticRepo,
		sellerRepo: sellerRepo,
	}
}

func (s *AdvertService) advertEntityToDTO(advert *entity.Advert) (*dto.AdvertResponse, error) {
	var posterURL string

	if !advert.ImageURL.Valid {
		posterURL = ""
	} else {
		var err error
		posterURL, err = s.staticRepo.GetStatic(advert.ImageURL.UUID)
		if err != nil {
			posterURL = ""
		}
	}

	advertDTO := dto.AdvertResponse{
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
		CreatedAt:   advert.CreatedAt,
		UpdatedAt:   advert.UpdatedAt,
	}

	return &advertDTO, nil
}

func (s *AdvertService) advertEntitiesToDTO(adverts []*entity.Advert) ([]*dto.AdvertResponse, error) {
	dtoAdverts := make([]*dto.AdvertResponse, 0, len(adverts))
	for _, advert := range adverts {
		advertDTO, err := s.advertEntityToDTO(advert)
		if err != nil {
			return nil, entity.UsecaseWrap(err, err)
		}
		dtoAdverts = append(dtoAdverts, advertDTO)
	}
	return dtoAdverts, nil
}

func (s *AdvertService) GetAdverts(limit, offset int) ([]*dto.AdvertResponse, error) {
	adverts, err := s.advertRepo.GetAdverts(limit, offset)
	if err != nil {
		return nil, entity.UsecaseWrap(err, err)
	}
	return s.advertEntitiesToDTO(adverts)
}

func (s *AdvertService) GetAdvertsByUserId(userId uuid.UUID) ([]*dto.AdvertResponse, error) {
	seller, err := s.sellerRepo.GetSellerByUserID(userId)
	if err != nil {
		return nil, entity.UsecaseWrap(err, repository.ErrSellerNotFound)
	}

	adverts, err := s.advertRepo.GetAdvertsBySellerId(seller.ID)
	if err != nil {
		return nil, entity.UsecaseWrap(err, err)
	}
	return s.advertEntitiesToDTO(adverts)
}

func (s *AdvertService) GetAdvertsByCartId(cartId uuid.UUID) ([]*dto.AdvertResponse, error) {
	adverts, err := s.advertRepo.GetAdvertsByCartId(cartId)
	if err != nil {
		return nil, entity.UsecaseWrap(err, err)
	}
	return s.advertEntitiesToDTO(adverts)
}

func (s *AdvertService) GetAdvertById(advertId uuid.UUID) (*dto.AdvertResponse, error) {
	advert, err := s.advertRepo.GetAdvertById(advertId)

	if err != nil {
		if errors.Is(err, repository.ErrAdvertNotFound) {
			return nil, entity.UsecaseWrap(ErrAdvertNotFound, ErrAdvertNotFound)
		}
		return nil, entity.UsecaseWrap(err, err)
	}

	return s.advertEntityToDTO(advert)
}

func (s *AdvertService) AddAdvert(advert *dto.AdvertRequest, userId uuid.UUID) (*dto.AdvertResponse, error) {
	if err := entity.ValidateAdvert(advert.Title,
		advert.Description,
		advert.Location,
		string(advert.Status),
		int(advert.Price)); err != nil {
		return nil, entity.UsecaseWrap(ErrAdvertBadRequest, ErrAdvertBadRequest)
	}

	seller, err := s.sellerRepo.GetSellerByUserID(userId)
	if err != nil {
		return nil, entity.UsecaseWrap(err, repository.ErrSellerNotFound)
	}

	entityAdvert, err := s.advertRepo.AddAdvert(&entity.Advert{
		SellerId:    seller.ID,
		CategoryId:  advert.CategoryId,
		Title:       strings.TrimSpace(advert.Title),
		Description: strings.TrimSpace(advert.Description),
		Price:       advert.Price,
		Status:      entity.AdvertStatus(advert.Status),
		HasDelivery: advert.HasDelivery,
		Location:    advert.Location,
	})
	if err != nil {
		return nil, entity.UsecaseWrap(ErrAdvertBadRequest, ErrAdvertBadRequest)
	}

	return s.advertEntityToDTO(entityAdvert)
}

func (s *AdvertService) UpdateAdvert(advert *dto.AdvertRequest, userId uuid.UUID, advertId uuid.UUID) error {
	if err := entity.ValidateAdvert(advert.Title,
		advert.Description,
		advert.Location,
		string(advert.Status),
		int(advert.Price)); err != nil {
		return entity.UsecaseWrap(ErrAdvertBadRequest, ErrAdvertBadRequest)
	}

	seller, err := s.sellerRepo.GetSellerByUserID(userId)
	if err != nil {
		return entity.UsecaseWrap(err, repository.ErrSellerNotFound)
	}

	existingAdvert, err := s.advertRepo.GetAdvertById(advertId)
	if err != nil {
		if errors.Is(err, repository.ErrAdvertNotFound) {
			return entity.UsecaseWrap(ErrAdvertNotFound, ErrAdvertNotFound)
		}
		return entity.UsecaseWrap(err, err)
	}
	if existingAdvert.SellerId != seller.ID {
		return entity.UsecaseWrap(ErrForbidden, ErrForbidden)
	}

	err = s.advertRepo.UpdateAdvert(&entity.Advert{
		ID:          advertId,
		SellerId:    seller.ID,
		CategoryId:  advert.CategoryId,
		Title:       strings.TrimSpace(advert.Title),
		Description: strings.TrimSpace(advert.Description),
		Price:       advert.Price,
		Status:      entity.AdvertStatus(advert.Status),
		HasDelivery: advert.HasDelivery,
		Location:    advert.Location,
	})
	if err != nil {
		return entity.UsecaseWrap(ErrAdvertBadRequest, ErrAdvertBadRequest)
	}

	return nil
}

func (s *AdvertService) DeleteAdvertById(advertId uuid.UUID, userId uuid.UUID) error {
	seller, err := s.sellerRepo.GetSellerByUserID(userId)
	if err != nil {
		return entity.UsecaseWrap(err, repository.ErrSellerNotFound)
	}

	existingAdvert, err := s.advertRepo.GetAdvertById(advertId)
	if err != nil {
		return entity.UsecaseWrap(err, repository.ErrAdvertNotFound)
	}
	if existingAdvert.SellerId != seller.ID {
		return entity.UsecaseWrap(ErrForbidden, ErrForbidden)
	}

	if err := s.advertRepo.DeleteAdvertById(advertId); err != nil {
		if errors.Is(err, repository.ErrAdvertNotFound) {
			return entity.UsecaseWrap(ErrAdvertNotFound, ErrAdvertNotFound)
		}
		return entity.UsecaseWrap(err, err)
	}

	return nil
}

func (s *AdvertService) UpdateAdvertStatus(advertId uuid.UUID, status string, userId uuid.UUID) error {
	seller, err := s.sellerRepo.GetSellerByUserID(userId)
	if err != nil {
		return entity.UsecaseWrap(err, repository.ErrSellerNotFound)
	}

	existingAdvert, err := s.advertRepo.GetAdvertById(advertId)
	if err != nil {
		return entity.UsecaseWrap(err, repository.ErrAdvertNotFound)
	}
	if existingAdvert.SellerId != seller.ID {
		return entity.UsecaseWrap(ErrForbidden, ErrForbidden)
	}

	if err := s.advertRepo.UpdateAdvertStatus(advertId, status); err != nil {
		if errors.Is(err, repository.ErrAdvertNotFound) {
			return entity.UsecaseWrap(ErrAdvertNotFound, ErrAdvertNotFound)
		}
		return entity.UsecaseWrap(err, err)
	}

	return nil
}

func (s *AdvertService) GetAdvertsByCategoryId(categoryId uuid.UUID) ([]*dto.AdvertResponse, error) {
	adverts, err := s.advertRepo.GetAdvertsByCategoryId(categoryId)
	if err != nil {
		return nil, entity.UsecaseWrap(err, err)
	}

	return s.advertEntitiesToDTO(adverts)
}

func (s *AdvertService) UploadImage(advertId uuid.UUID, imageId uuid.UUID, userId uuid.UUID) error {
	seller, err := s.sellerRepo.GetSellerByUserID(userId)
	if err != nil {
		return entity.UsecaseWrap(err, repository.ErrSellerNotFound)
	}

	existingAdvert, err := s.advertRepo.GetAdvertById(advertId)
	if err != nil {
		return entity.UsecaseWrap(err, repository.ErrAdvertNotFound)
	}
	if existingAdvert.SellerId != seller.ID {
		return entity.UsecaseWrap(ErrForbidden, ErrForbidden)
	}

	if err := s.advertRepo.UploadImage(advertId, imageId); err != nil {
		return entity.UsecaseWrap(ErrAdvertBadRequest, ErrAdvertBadRequest)
	}

	return nil
}
