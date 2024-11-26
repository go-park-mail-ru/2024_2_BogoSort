package service

import (
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupAdvertTestService(t *testing.T) (*AdvertService, *gomock.Controller, *mocks.MockAdvertRepository, *mocks.MockStaticRepository, *mocks.MockSeller) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockAdvertRepository(ctrl)
	mockStaticRepo := mocks.NewMockStaticRepository(ctrl)
	mockSellerRepo := mocks.NewMockSeller(ctrl)

	service := &AdvertService{
		advertRepo: mockRepo,
		sellerRepo: mockSellerRepo,
	}

	return service, ctrl, mockRepo, mockStaticRepo, mockSellerRepo
}

func TestAdvertService_AddAdvert_Success(t *testing.T) {
	service, ctrl, mockRepo, _, mockSellerRepo := setupAdvertTestService(t)
	defer ctrl.Finish()

	advertRequest := &dto.AdvertRequest{
		Title:       "Test Advert",
		Description: "Test Description",
		Price:       100,
		Location:    "Test Location",
		HasDelivery: true,
		CategoryId:  uuid.New(),
		Status:      "active",
	}

	sellerId := uuid.New()
	mockSellerRepo.EXPECT().GetByUserId(sellerId).Return(&entity.Seller{ID: sellerId}, nil)
	mockRepo.EXPECT().Add(gomock.Any()).Return(&entity.Advert{ID: uuid.New()}, nil)

	result, err := service.Add(advertRequest, sellerId)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestAdvertService_AddAdvert_SellerNotFound(t *testing.T) {
	service, ctrl, _, _, mockSellerRepo := setupAdvertTestService(t)
	defer ctrl.Finish()

	advertRequest := &dto.AdvertRequest{
		Title:       "Test Advert",
		Description: "Test Description",
		Price:       100,
		Location:    "Test Location",
		HasDelivery: true,
		CategoryId:  uuid.New(),
		Status:      "active",
	}

	sellerId := uuid.New()
	mockSellerRepo.EXPECT().GetByUserId(sellerId).Return(nil, repository.ErrSellerNotFound)

	result, err := service.Add(advertRequest, sellerId)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, errors.Unwrap(err))
}

func TestAdvertService_UpdateAdvert_Success(t *testing.T) {
	service, ctrl, mockRepo, _, mockSellerRepo := setupAdvertTestService(t)
	defer ctrl.Finish()

	advertId := uuid.New()
	userId := uuid.New()

	advertUpdate := &dto.AdvertRequest{
		Title:       "Updated Advert",
		Description: "Updated Description",
		Price:       150,
		Location:    "Updated Location",
		HasDelivery: false,
		CategoryId:  uuid.New(),
		Status:      "inactive",
	}

	sellerId := uuid.New()
	mockSellerRepo.EXPECT().GetByUserId(userId).Return(&entity.Seller{ID: sellerId}, nil)
	mockRepo.EXPECT().GetById(advertId, userId).Return(&entity.Advert{ID: advertId, SellerId: sellerId}, nil)
	mockRepo.EXPECT().Update(gomock.Any()).Return(nil)

	err := service.Update(advertUpdate, userId, advertId)
	assert.NoError(t, err)
}

func TestAdvertService_UpdateAdvert_NotFound(t *testing.T) {
	service, ctrl, mockRepo, _, mockSellerRepo := setupAdvertTestService(t)
	defer ctrl.Finish()

	advertId := uuid.New()
	userId := uuid.New()

	advertUpdate := &dto.AdvertRequest{
		Title:       "Updated Advert",
		Description: "Updated Description",
		Price:       150,
		Location:    "Updated Location",
		HasDelivery: false,
		CategoryId:  uuid.New(),
		Status:      "inactive",
	}

	sellerId := uuid.New()
	mockSellerRepo.EXPECT().GetByUserId(userId).Return(&entity.Seller{ID: sellerId}, nil)
	mockRepo.EXPECT().GetById(advertId, userId).Return(nil, repository.ErrAdvertNotFound)

	err := service.Update(advertUpdate, userId, advertId)
	assert.Error(t, err)
	assert.Nil(t, errors.Unwrap(err))
}

func TestAdvertService_UpdateAdvert_Forbidden(t *testing.T) {
	service, ctrl, mockRepo, _, mockSellerRepo := setupAdvertTestService(t)
	defer ctrl.Finish()

	advertId := uuid.New()
	userId := uuid.New()

	advertUpdate := &dto.AdvertRequest{
		Title:       "Updated Advert",
		Description: "Updated Description",
		Price:       150,
		Location:    "Updated Location",
		HasDelivery: false,
		CategoryId:  uuid.New(),
		Status:      "inactive",
	}

	sellerId := uuid.New()
	mockSellerRepo.EXPECT().GetByUserId(userId).Return(&entity.Seller{ID: sellerId}, nil)
	// Возвращаем объявление с другим SellerId
	mockRepo.EXPECT().GetById(advertId, userId).Return(&entity.Advert{ID: advertId, SellerId: uuid.New()}, nil)

	err := service.Update(advertUpdate, userId, advertId)
	assert.Error(t, err)
	assert.Nil(t, errors.Unwrap(err))
}

func TestAdvertService_DeleteById_Success(t *testing.T) {
	service, ctrl, mockRepo, _, mockSellerRepo := setupAdvertTestService(t)
	defer ctrl.Finish()

	advertId := uuid.New()
	userId := uuid.New()

	sellerId := uuid.New()
	mockSellerRepo.EXPECT().GetByUserId(userId).Return(&entity.Seller{ID: sellerId}, nil)
	mockRepo.EXPECT().GetById(advertId, userId).Return(&entity.Advert{ID: advertId, SellerId: sellerId}, nil)
	mockRepo.EXPECT().DeleteById(advertId).Return(nil)

	err := service.DeleteById(advertId, userId)
	assert.NoError(t, err)
}

func TestAdvertService_DeleteById_NotFound(t *testing.T) {
	service, ctrl, mockRepo, _, mockSellerRepo := setupAdvertTestService(t)
	defer ctrl.Finish()

	advertId := uuid.New()
	userId := uuid.New()

	sellerId := uuid.New()
	mockSellerRepo.EXPECT().GetByUserId(userId).Return(&entity.Seller{ID: sellerId}, nil)
	mockRepo.EXPECT().GetById(advertId, userId).Return(nil, repository.ErrAdvertNotFound)

	err := service.DeleteById(advertId, userId)
	assert.Error(t, err)
	assert.Nil(t, errors.Unwrap(err))
}

func TestAdvertService_DeleteById_Forbidden(t *testing.T) {
	service, ctrl, mockRepo, _, mockSellerRepo := setupAdvertTestService(t)
	defer ctrl.Finish()

	advertId := uuid.New()
	userId := uuid.New()

	sellerId := uuid.New()
	mockSellerRepo.EXPECT().GetByUserId(userId).Return(&entity.Seller{ID: sellerId}, nil)
	// Возвращаем объявление с другим SellerId
	mockRepo.EXPECT().GetById(advertId, userId).Return(&entity.Advert{ID: advertId, SellerId: uuid.New()}, nil)

	err := service.DeleteById(advertId, userId)
	assert.Error(t, err)
	assert.Nil(t, errors.Unwrap(err))
}

func TestAdvertService_UploadImage_Success(t *testing.T) {
	service, ctrl, mockRepo, _, mockSellerRepo := setupAdvertTestService(t)
	defer ctrl.Finish()

	advertId := uuid.New()
	imageId := uuid.New()
	userId := uuid.New()

	sellerId := uuid.New()
	mockSellerRepo.EXPECT().GetByUserId(userId).Return(&entity.Seller{ID: sellerId}, nil)
	mockRepo.EXPECT().GetById(advertId, userId).Return(&entity.Advert{ID: advertId, SellerId: sellerId}, nil)
	mockRepo.EXPECT().UploadImage(advertId, imageId).Return(nil)

	err := service.UploadImage(advertId, imageId, userId)
	assert.NoError(t, err)
}

func TestAdvertService_UploadImage_AdvertNotFound(t *testing.T) {
	service, ctrl, mockRepo, _, mockSellerRepo := setupAdvertTestService(t)
	defer ctrl.Finish()

	advertId := uuid.New()
	imageId := uuid.New()
	userId := uuid.New()

	sellerId := uuid.New()
	mockSellerRepo.EXPECT().GetByUserId(userId).Return(&entity.Seller{ID: sellerId}, nil)
	mockRepo.EXPECT().GetById(advertId, userId).Return(nil, repository.ErrAdvertNotFound)

	err := service.UploadImage(advertId, imageId, userId)
	assert.Error(t, err)
	assert.Nil(t, errors.Unwrap(err))
}

func TestAdvertService_UploadImage_Forbidden(t *testing.T) {
	service, ctrl, mockRepo, _, mockSellerRepo := setupAdvertTestService(t)
	defer ctrl.Finish()

	advertId := uuid.New()
	imageId := uuid.New()
	userId := uuid.New()

	sellerId := uuid.New()
	mockSellerRepo.EXPECT().GetByUserId(userId).Return(&entity.Seller{ID: sellerId}, nil)
	// Возвращаем объявление с другим SellerId
	mockRepo.EXPECT().GetById(advertId, userId).Return(&entity.Advert{ID: advertId, SellerId: uuid.New()}, nil)

	err := service.UploadImage(advertId, imageId, userId)
	assert.Error(t, err)
	assert.Nil(t, errors.Unwrap(err))
}

func TestAdvertService_UploadImage_UploadError(t *testing.T) {
	service, ctrl, mockRepo, _, mockSellerRepo := setupAdvertTestService(t)
	defer ctrl.Finish()

	advertId := uuid.New()
	imageId := uuid.New()
	userId := uuid.New()

	sellerId := uuid.New()
	mockSellerRepo.EXPECT().GetByUserId(userId).Return(&entity.Seller{ID: sellerId}, nil)
	mockRepo.EXPECT().GetById(advertId, userId).Return(&entity.Advert{ID: advertId, SellerId: sellerId}, nil)
	mockRepo.EXPECT().UploadImage(advertId, imageId).Return(repository.ErrAdvertBadRequest)

	err := service.UploadImage(advertId, imageId, userId)
	assert.Error(t, err)
	assert.Nil(t, errors.Unwrap(err))
}

func TestAdvertService_GetSavedByUserId_Success(t *testing.T) {
	service, ctrl, mockRepo, _, _ := setupAdvertTestService(t)
	defer ctrl.Finish()

	userId := uuid.New()
	mockRepo.EXPECT().GetSavedByUserId(userId).Return([]*entity.Advert{
		{ID: uuid.New()},
		{ID: uuid.New()},
	}, nil)

	adverts, err := service.GetSavedByUserId(userId)
	assert.NoError(t, err)
	assert.Len(t, adverts, 2)
}

func TestAdvertService_GetSavedByUserId_Error(t *testing.T) {
	service, ctrl, mockRepo, _, _ := setupAdvertTestService(t)
	defer ctrl.Finish()

	userId := uuid.New()
	mockRepo.EXPECT().GetSavedByUserId(userId).Return(nil, errors.New("repository error"))

	adverts, err := service.GetSavedByUserId(userId)
	assert.Error(t, err)
	assert.Nil(t, adverts)
}
