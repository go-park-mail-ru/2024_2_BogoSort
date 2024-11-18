package service

import (
	"errors"
	"strings"
	"context"

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

func (s *AdvertService) Get(limit, offset int, userId uuid.UUID) ([]*dto.PreviewAdvertCard, error) {
	adverts, err := s.advertRepo.Get(limit, offset, userId)
	if err != nil {
		return nil, entity.UsecaseWrap(err, err)
	}

	dtoAdverts := make([]*dto.PreviewAdvertCard, 0, len(adverts))
	for _, advert := range adverts {
		advertDTO := dto.PreviewAdvertCard{
			Preview: dto.PreviewAdvert{
				ID:          advert.ID,
				SellerId:    advert.SellerId,
				CategoryId:  advert.CategoryId,
				Title:       advert.Title,
				Price:       advert.Price,
				ImageURL:    advert.ImageURL.UUID.String(),
				Status:      dto.AdvertStatus(advert.Status),
				HasDelivery: advert.HasDelivery,
				Location:    advert.Location,
			},
			IsSaved: advert.IsSaved,
			IsViewed: advert.IsViewed,
		}
		dtoAdverts = append(dtoAdverts, &advertDTO)
	}

	return dtoAdverts, nil
}

func (s *AdvertService) GetByUserId(userId uuid.UUID) ([]*dto.MyPreviewAdvertCard, error) {
	seller, err := s.sellerRepo.GetByUserId(userId)
	if err != nil {
		return nil, entity.UsecaseWrap(err, repository.ErrSellerNotFound)
	}

	adverts, err := s.advertRepo.GetBySellerId(seller.ID, userId)
	if err != nil {
		return nil, entity.UsecaseWrap(err, err)
	}
	
	dtoAdverts := make([]*dto.MyPreviewAdvertCard, 0, len(adverts))
	for _, advert := range adverts {
		advertDTO := dto.MyPreviewAdvertCard{
			Preview: dto.PreviewAdvert{
				ID:          advert.ID,
				SellerId:    advert.SellerId,
				CategoryId:  advert.CategoryId,
				Title:       advert.Title,
				Price:       advert.Price,
				ImageURL:    advert.ImageURL.UUID.String(),
				Status:      dto.AdvertStatus(advert.Status),
				HasDelivery: advert.HasDelivery,
				Location:    advert.Location,
			},
			ViewsNumber: advert.ViewsNumber,
			SavesNumber: advert.SavesNumber,
		}
		dtoAdverts = append(dtoAdverts, &advertDTO)
	}
	return dtoAdverts, nil
}

func (s *AdvertService) GetByCartId(cartId, userId uuid.UUID) ([]*dto.PreviewAdvertCard, error) {
	adverts, err := s.advertRepo.GetByCartId(cartId, userId)
	if err != nil {
		return nil, entity.UsecaseWrap(err, err)
	}

	dtoAdverts := make([]*dto.PreviewAdvertCard, 0, len(adverts))
	for _, advert := range adverts {
		advertDTO := dto.PreviewAdvertCard{
			Preview: dto.PreviewAdvert{
				ID:          advert.ID,
				SellerId:    advert.SellerId,
				CategoryId:  advert.CategoryId,
				Title:       advert.Title,
				Price:       advert.Price,
				ImageURL:    advert.ImageURL.UUID.String(),
				Status:      dto.AdvertStatus(advert.Status),
				HasDelivery: advert.HasDelivery,
				Location:    advert.Location,
			},
			IsSaved: advert.IsSaved,
			IsViewed: advert.IsViewed,
		}
		dtoAdverts = append(dtoAdverts, &advertDTO)
	}

	return dtoAdverts, nil	
}

func (s *AdvertService) GetById(advertId, userId uuid.UUID) (*dto.AdvertCard, error) {
	advert, err := s.advertRepo.GetById(advertId, userId)

	if err != nil {
		if errors.Is(err, repository.ErrAdvertNotFound) {
			return nil, entity.UsecaseWrap(ErrAdvertNotFound, ErrAdvertNotFound)
		}
		return nil, entity.UsecaseWrap(err, err)
	}

	advertDTO := dto.AdvertCard{
		Advert: dto.Advert{
			ID:          advert.ID,
			SellerId:    advert.SellerId,
			CategoryId:  advert.CategoryId,
			Title:       advert.Title,
			Price:       advert.Price,
			ImageURL:    advert.ImageURL.UUID.String(),
			Status:      dto.AdvertStatus(advert.Status),
			HasDelivery: advert.HasDelivery,
			Location:    advert.Location,
			CreatedAt:   advert.CreatedAt,
			UpdatedAt:   advert.UpdatedAt,
			ViewsNumber: advert.ViewsNumber,
			SavesNumber: advert.SavesNumber,
		},
		IsSaved: advert.IsSaved,
		IsViewed: advert.IsViewed,
	}
	return &advertDTO, nil
}

func (s *AdvertService) Add(advert *dto.AdvertRequest, userId uuid.UUID) (*dto.Advert, error) {
	if err := entity.ValidateAdvert(advert.Title,
		advert.Description,
		advert.Location,
		string(advert.Status),
		int(advert.Price)); err != nil {
		return nil, entity.UsecaseWrap(ErrAdvertBadRequest, ErrAdvertBadRequest)
	}

	seller, err := s.sellerRepo.GetByUserId(userId)
	if err != nil {
		return nil, entity.UsecaseWrap(err, repository.ErrSellerNotFound)
	}

	entityAdvert, err := s.advertRepo.Add(&entity.Advert{
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

	advertDTO := dto.Advert{
		ID: entityAdvert.ID,
		CategoryId: entityAdvert.CategoryId,
		SellerId: entityAdvert.SellerId,
		Title: entityAdvert.Title,
		Description: entityAdvert.Description,
		Price: entityAdvert.Price,
		Status: dto.AdvertStatus(entityAdvert.Status),
		HasDelivery: entityAdvert.HasDelivery,
		Location: entityAdvert.Location,
		ImageURL: entityAdvert.ImageURL.UUID.String(),
		CreatedAt: entityAdvert.CreatedAt,
		UpdatedAt: entityAdvert.UpdatedAt,
		ViewsNumber: entityAdvert.ViewsNumber,
		SavesNumber: entityAdvert.SavesNumber,
	}
	return &advertDTO, nil
}

func (s *AdvertService) Update(advert *dto.AdvertRequest, userId uuid.UUID, advertId uuid.UUID) error {
	if err := entity.ValidateAdvert(advert.Title,
		advert.Description,
		advert.Location,
		string(advert.Status),
		int(advert.Price)); err != nil {
		return entity.UsecaseWrap(ErrAdvertBadRequest, ErrAdvertBadRequest)
	}

	seller, err := s.sellerRepo.GetByUserId(userId)
	if err != nil {
		return entity.UsecaseWrap(err, repository.ErrSellerNotFound)
	}

	existingAdvert, err := s.advertRepo.GetById(advertId, userId)
	if err != nil {
		if errors.Is(err, repository.ErrAdvertNotFound) {
			return entity.UsecaseWrap(ErrAdvertNotFound, ErrAdvertNotFound)
		}
		return entity.UsecaseWrap(err, err)
	}
	if existingAdvert.SellerId != seller.ID {
		return entity.UsecaseWrap(ErrForbidden, ErrForbidden)
	}

	err = s.advertRepo.Update(&entity.Advert{
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

func (s *AdvertService) DeleteById(advertId uuid.UUID, userId uuid.UUID) error {
	seller, err := s.sellerRepo.GetByUserId(userId)
	if err != nil {
		return entity.UsecaseWrap(err, repository.ErrSellerNotFound)
	}

	existingAdvert, err := s.advertRepo.GetById(advertId, userId)
	if err != nil {
		return entity.UsecaseWrap(err, repository.ErrAdvertNotFound)
	}
	if existingAdvert.SellerId != seller.ID {
		return entity.UsecaseWrap(ErrForbidden, ErrForbidden)
	}

	if err := s.advertRepo.DeleteById(advertId); err != nil {
		if errors.Is(err, repository.ErrAdvertNotFound) {
			return entity.UsecaseWrap(ErrAdvertNotFound, ErrAdvertNotFound)
		}
		return entity.UsecaseWrap(err, err)
	}

	return nil
}

func (s *AdvertService) UpdateStatus(advertId, userId uuid.UUID, status dto.AdvertStatus) error {
	ctx := context.Background()
	tx, err := s.advertRepo.BeginTransaction()
	if err != nil {
		s.logger.Error("failed to begin transaction", zap.Error(err))
		return entity.UsecaseWrap(errors.New("failed to begin transaction"), err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	seller, err := s.sellerRepo.GetByUserId(userId)
	if err != nil {
		return entity.UsecaseWrap(err, repository.ErrSellerNotFound)
	}

	existingAdvert, err := s.advertRepo.GetById(advertId, userId)
	if err != nil {
		return entity.UsecaseWrap(err, repository.ErrAdvertNotFound)
	}
	if existingAdvert.SellerId != seller.ID {
		return entity.UsecaseWrap(ErrForbidden, ErrForbidden)
	}

	if err := s.advertRepo.UpdateStatus(tx, advertId, entity.AdvertStatus(status)); err != nil {
		if errors.Is(err, repository.ErrAdvertNotFound) {
			return entity.UsecaseWrap(ErrAdvertNotFound, ErrAdvertNotFound)
		}
		return entity.UsecaseWrap(err, err)
	}

	return nil
}

func (s *AdvertService) GetByCategoryId(categoryId, userId uuid.UUID) ([]*dto.PreviewAdvertCard, error) {
	adverts, err := s.advertRepo.GetByCategoryId(categoryId, userId)
	if err != nil {
		return nil, entity.UsecaseWrap(err, err)
	}

	dtoAdverts := make([]*dto.PreviewAdvertCard, 0, len(adverts))
	for _, advert := range adverts {
		advertDTO := dto.PreviewAdvertCard{
			Preview: dto.PreviewAdvert{
				ID:          advert.ID,
				SellerId:    advert.SellerId,
				CategoryId:  advert.CategoryId,
				Title:       advert.Title,
				Price:       advert.Price,
				ImageURL:    advert.ImageURL.UUID.String(),
				Status:      dto.AdvertStatus(advert.Status),
				HasDelivery: advert.HasDelivery,
				Location:    advert.Location,
			},
			IsSaved: advert.IsSaved,
			IsViewed: advert.IsViewed,
		}
		dtoAdverts = append(dtoAdverts, &advertDTO)
	}

	return dtoAdverts, nil
}

func (s *AdvertService) UploadImage(advertId uuid.UUID, imageId uuid.UUID, userId uuid.UUID) error {
	seller, err := s.sellerRepo.GetByUserId(userId)
	if err != nil {
		return entity.UsecaseWrap(err, repository.ErrSellerNotFound)
	}

	existingAdvert, err := s.advertRepo.GetById(advertId, userId)
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

func (s *AdvertService) GetSavedByUserId(userId uuid.UUID) ([]*dto.PreviewAdvertCard, error) {
	adverts, err := s.advertRepo.GetSavedByUserId(userId)
	if err != nil {
		return nil, entity.UsecaseWrap(err, err)
	}

	dtoAdverts := make([]*dto.PreviewAdvertCard, 0, len(adverts))
	for _, advert := range adverts {
		advertDTO := dto.PreviewAdvertCard{
			Preview: dto.PreviewAdvert{
				ID:          advert.ID,
				SellerId:    advert.SellerId,
				CategoryId:  advert.CategoryId,
				Title:       advert.Title,
				Price:       advert.Price,
				ImageURL:    advert.ImageURL.UUID.String(),
				Status:      dto.AdvertStatus(advert.Status),
				HasDelivery: advert.HasDelivery,
				Location:    advert.Location,
			},
			IsSaved: advert.IsSaved,
			IsViewed: advert.IsViewed,
		}
		dtoAdverts = append(dtoAdverts, &advertDTO)
	}

	return dtoAdverts, nil
}

func (s *AdvertService) AddToSaved(advertId, userId uuid.UUID) error {
	err := s.advertRepo.AddToSaved(advertId, userId)
	if err != nil {
		return entity.UsecaseWrap(err, err)
	}

	return nil
}

func (s *AdvertService) RemoveFromSaved(advertId, userId uuid.UUID) error {
	err := s.advertRepo.DeleteFromSaved(advertId, userId)
	if err != nil {
		return entity.UsecaseWrap(err, err)
	}
	
	return nil
}
