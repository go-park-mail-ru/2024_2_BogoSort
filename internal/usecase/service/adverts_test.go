package service

import (
	"testing"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupAdvertTestService(t *testing.T) (*AdvertService, *gomock.Controller, *mocks.MockAdvertRepository, *mocks.MockStaticRepository, *mocks.MockSeller) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockAdvertRepository(ctrl)
	mockStaticRepo := mocks.NewMockStaticRepository(ctrl)
	mockSellerRepo := mocks.NewMockSeller(ctrl)
	logger := zap.NewNop()

	service := &AdvertService{
		advertRepo: mockRepo,
		staticRepo: mockStaticRepo,
		sellerRepo: mockSellerRepo,
		logger:     logger,
	}

	return service, ctrl, mockRepo, mockStaticRepo, mockSellerRepo
}

func TestAdvertService_AddAdvert(t *testing.T) {
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
	mockSellerRepo.EXPECT().GetSellerByUserID(sellerId).Return(&entity.Seller{ID: sellerId}, nil)
	mockRepo.EXPECT().AddAdvert(gomock.Any()).Return(&entity.Advert{ID: uuid.New()}, nil)

	result, err := service.AddAdvert(advertRequest, sellerId)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestAdvertService_Cases(t *testing.T) {
	service, ctrl, mockRepo, _, mockSellerRepo := setupAdvertTestService(t)
	defer ctrl.Finish()

	testCases := []struct {
		name     string
		testFunc func()
	}{
		{
			name: "GetAdverts",
			testFunc: func() {
				mockRepo.EXPECT().GetAdverts(10, 0).Return([]*entity.Advert{{}}, nil)
				adverts, err := service.GetAdverts(10, 0)
				assert.NoError(t, err)
				assert.Len(t, adverts, 1)
			},
		},
		{
			name: "GetAdvertById",
			testFunc: func() {
				advertId := uuid.New()
				mockRepo.EXPECT().GetAdvertById(advertId).Return(&entity.Advert{ID: advertId}, nil)
				advert, err := service.GetAdvertById(advertId)
				assert.NoError(t, err)
				assert.Equal(t, advertId, advert.ID)
			},
		},
		{
			name: "UpdateAdvert",
			testFunc: func() {
				advertId := uuid.New()
				advertRequest := &dto.AdvertRequest{
					Title:       "Updated Advert",
					Description: "Updated Description",
					Price:       150,
					Location:    "Updated Location",
					HasDelivery: false,
					CategoryId:  uuid.New(),
					Status:      "inactive",
				}
				sellerId := uuid.New()
				mockSellerRepo.EXPECT().GetSellerByUserID(sellerId).Return(&entity.Seller{ID: sellerId}, nil)
				mockRepo.EXPECT().GetAdvertById(advertId).Return(&entity.Advert{ID: advertId, SellerId: sellerId}, nil)
				mockRepo.EXPECT().UpdateAdvert(gomock.Any()).Return(nil)
				err := service.UpdateAdvert(advertRequest, sellerId, advertId)
				assert.NoError(t, err)
			},
		},
		{
			name: "DeleteAdvertById",
			testFunc: func() {
				advertId := uuid.New()
				sellerId := uuid.New()
				mockSellerRepo.EXPECT().GetSellerByUserID(sellerId).Return(&entity.Seller{ID: sellerId}, nil)
				mockRepo.EXPECT().GetAdvertById(advertId).Return(&entity.Advert{ID: advertId, SellerId: sellerId}, nil)
				mockRepo.EXPECT().DeleteAdvertById(advertId).Return(nil)
				err := service.DeleteAdvertById(advertId, sellerId)
				assert.NoError(t, err)
			},
		},
		// {
		// 	name: "UpdateAdvertStatus",
		// 	testFunc: func() {
		// 		advertId := uuid.New()
		// 		sellerId := uuid.New()
		// 		mockSellerRepo.EXPECT().GetSellerByUserID(sellerId).Return(&entity.Seller{ID: sellerId}, nil)
		// 		mockRepo.EXPECT().GetAdvertById(advertId).Return(&entity.Advert{ID: advertId, SellerId: sellerId}, nil)
		// 		mockRepo.EXPECT().UpdateAdvertStatus(advertId, "inactive").Return(nil)
		// 		err := service.UpdateAdvertStatus(advertId, "inactive", sellerId)
		// 		assert.NoError(t, err)
		// 	},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc()
		})
	}
}

func TestAdvertService_GetAdvertsByUserId(t *testing.T) {
	service, ctrl, mockRepo, _, mockSellerRepo := setupAdvertTestService(t)
	defer ctrl.Finish()

	sellerId := uuid.New()
	mockSellerRepo.EXPECT().GetSellerByUserID(sellerId).Return(&entity.Seller{ID: sellerId}, nil)
	mockRepo.EXPECT().GetAdvertsBySellerId(sellerId).Return([]*entity.Advert{{ID: uuid.New()}}, nil)

	adverts, err := service.GetAdvertsByUserId(sellerId)
	assert.NoError(t, err)
	assert.Len(t, adverts, 1)
}

func TestAdvertService_GetAdvertsByCartId(t *testing.T) {
	service, ctrl, mockRepo, _, _ := setupAdvertTestService(t)
	defer ctrl.Finish()

	cartId := uuid.New()
	mockRepo.EXPECT().GetAdvertsByCartId(cartId).Return([]*entity.Advert{{ID: uuid.New()}}, nil)

	adverts, err := service.GetAdvertsByCartId(cartId)
	assert.NoError(t, err)
	assert.Len(t, adverts, 1)
}

func TestAdvertService_GetAdvertsByCategoryId(t *testing.T) {
	service, ctrl, mockRepo, _, _ := setupAdvertTestService(t)
	defer ctrl.Finish()

	categoryId := uuid.New()
	mockRepo.EXPECT().GetAdvertsByCategoryId(categoryId).Return([]*entity.Advert{{ID: uuid.New()}}, nil)

	adverts, err := service.GetAdvertsByCategoryId(categoryId)
	assert.NoError(t, err)
	assert.Len(t, adverts, 1)
}

func TestAdvertService_UploadImage(t *testing.T) {
	service, ctrl, mockRepo, _, mockSellerRepo := setupAdvertTestService(t)
	defer ctrl.Finish()

	advertId := uuid.New()
	imageId := uuid.New()
	sellerId := uuid.New()

	mockSellerRepo.EXPECT().GetSellerByUserID(sellerId).Return(&entity.Seller{ID: sellerId}, nil)
	mockRepo.EXPECT().GetAdvertById(advertId).Return(&entity.Advert{ID: advertId, SellerId: sellerId}, nil)
	mockRepo.EXPECT().UploadImage(advertId, imageId).Return(nil)

	err := service.UploadImage(advertId, imageId, sellerId)
	assert.NoError(t, err)
}
