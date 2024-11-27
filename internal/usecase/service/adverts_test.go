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

func setupAdvertService(t *testing.T) (*AdvertService, *mocks.MockAdvertRepository, *mocks.MockSeller, *mocks.MockUser, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	advertRepo := mocks.NewMockAdvertRepository(ctrl)
	sellerRepo := mocks.NewMockSeller(ctrl)
	userRepo := mocks.NewMockUser(ctrl)
	service := NewAdvertService(advertRepo, sellerRepo, userRepo)
	return service, advertRepo, sellerRepo, userRepo, ctrl
}

func TestAdvertService_GetById(t *testing.T) {
	service, advertRepo, _, _, ctrl := setupAdvertService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	advertID := uuid.New()
	expectedAdvert := &entity.Advert{
		ID:          advertID,
		SellerId:    uuid.New(),
		CategoryId:  uuid.New(),
		Title:       "Test Advert",
		Price:       100,
		Status:      entity.AdvertStatusActive,
		HasDelivery: true,
		Location:    "Test Location",
	}

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				advertRepo.EXPECT().GetById(advertID, userID).Return(expectedAdvert, nil)
			},
			expectedError: nil,
		},
		{
			name: "Advert Not Found",
			setupMocks: func() {
				advertRepo.EXPECT().GetById(advertID, userID).Return(nil, repository.ErrAdvertNotFound)
			},
			expectedError: ErrAdvertNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			advert, err := service.GetById(advertID, userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedAdvert.ID, advert.Advert.ID)
			}
		})
	}
}

func TestAdvertService_Add(t *testing.T) {
	service, advertRepo, sellerRepo, _, ctrl := setupAdvertService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	sellerID := uuid.New()
	advertRequest := &dto.AdvertRequest{
		Title:       "New Advert",
		Description: "Description",
		Price:       200,
		Status:      dto.AdvertStatusActive,
		Location:    "Location",
	}
	expectedAdvert := &entity.Advert{
		ID:          uuid.New(),
		SellerId:    sellerID,
		CategoryId:  uuid.New(),
		Title:       advertRequest.Title,
		Description: advertRequest.Description,
		Price:       advertRequest.Price,
		Status:      entity.AdvertStatus(advertRequest.Status),
		HasDelivery: advertRequest.HasDelivery,
		Location:    advertRequest.Location,
	}

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				sellerRepo.EXPECT().GetByUserId(userID).Return(&entity.Seller{ID: sellerID}, nil)
				advertRepo.EXPECT().Add(gomock.Any()).Return(expectedAdvert, nil)
			},
			expectedError: nil,
		},
		{
			name: "Seller Not Found",
			setupMocks: func() {
				sellerRepo.EXPECT().GetByUserId(userID).Return(nil, repository.ErrSellerNotFound)
			},
			expectedError: repository.ErrSellerNotFound,
		},
		{
			name: "Invalid Advert Data",
			setupMocks: func() {
				sellerRepo.EXPECT().GetByUserId(userID).Return(&entity.Seller{ID: sellerID}, nil)
				advertRepo.EXPECT().Add(gomock.Any()).Return(nil, ErrAdvertBadRequest)
			},
			expectedError: ErrAdvertBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			advert, err := service.Add(advertRequest, userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedAdvert.ID, advert.ID)
			}
		})
	}
}

func TestAdvertService_Update(t *testing.T) {
	service, advertRepo, sellerRepo, _, ctrl := setupAdvertService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	advertID := uuid.New()
	sellerID := uuid.New()
	advertRequest := &dto.AdvertRequest{
		Title:       "Updated Advert",
		Description: "Updated Description",
		Price:       300,
		Status:      dto.AdvertStatusActive,
		Location:    "Updated Location",
	}

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				sellerRepo.EXPECT().GetByUserId(userID).Return(&entity.Seller{ID: sellerID}, nil)
				advertRepo.EXPECT().GetById(advertID, userID).Return(&entity.Advert{SellerId: sellerID}, nil)
				advertRepo.EXPECT().Update(gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Advert Not Found",
			setupMocks: func() {
				sellerRepo.EXPECT().GetByUserId(userID).Return(&entity.Seller{ID: sellerID}, nil)
				advertRepo.EXPECT().GetById(advertID, userID).Return(nil, repository.ErrAdvertNotFound)
			},
			expectedError: ErrAdvertNotFound,
		},
		{
			name: "Forbidden",
			setupMocks: func() {
				sellerRepo.EXPECT().GetByUserId(userID).Return(&entity.Seller{ID: sellerID}, nil)
				advertRepo.EXPECT().GetById(advertID, userID).Return(&entity.Advert{SellerId: uuid.New()}, nil)
			},
			expectedError: ErrForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			err := service.Update(advertRequest, userID, advertID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAdvertService_DeleteById(t *testing.T) {
	service, advertRepo, sellerRepo, _, ctrl := setupAdvertService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	advertID := uuid.New()
	sellerID := uuid.New()

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				sellerRepo.EXPECT().GetByUserId(userID).Return(&entity.Seller{ID: sellerID}, nil)
				advertRepo.EXPECT().GetById(advertID, userID).Return(&entity.Advert{SellerId: sellerID}, nil)
				advertRepo.EXPECT().DeleteById(advertID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Advert Not Found",
			setupMocks: func() {
				sellerRepo.EXPECT().GetByUserId(userID).Return(&entity.Seller{ID: sellerID}, nil)
				advertRepo.EXPECT().GetById(advertID, userID).Return(nil, repository.ErrAdvertNotFound)
			},
			expectedError: repository.ErrAdvertNotFound,
		},
		{
			name: "Forbidden",
			setupMocks: func() {
				sellerRepo.EXPECT().GetByUserId(userID).Return(&entity.Seller{ID: sellerID}, nil)
				advertRepo.EXPECT().GetById(advertID, userID).Return(&entity.Advert{SellerId: uuid.New()}, nil)
			},
			expectedError: ErrForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			err := service.DeleteById(advertID, userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAdvertService_GetByUserId(t *testing.T) {
	service, advertRepo, sellerRepo, _, ctrl := setupAdvertService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	sellerID := uuid.New()
	expectedAdverts := []*entity.Advert{
		{
			ID:          uuid.New(),
			SellerId:    sellerID,
			CategoryId:  uuid.New(),
			Title:       "Advert 1",
			Price:       100,
			Status:      entity.AdvertStatusActive,
			HasDelivery: true,
			Location:    "Location 1",
		},
		{
			ID:          uuid.New(),
			SellerId:    sellerID,
			CategoryId:  uuid.New(),
			Title:       "Advert 2",
			Price:       200,
			Status:      entity.AdvertStatusInactive,
			HasDelivery: false,
			Location:    "Location 2",
		},
	}

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				sellerRepo.EXPECT().GetByUserId(userID).Return(&entity.Seller{ID: sellerID}, nil)
				advertRepo.EXPECT().GetByUserId(sellerID, userID).Return(expectedAdverts, nil)
			},
			expectedError: nil,
		},
		{
			name: "Seller Not Found",
			setupMocks: func() {
				sellerRepo.EXPECT().GetByUserId(userID).Return(nil, repository.ErrSellerNotFound)
			},
			expectedError: repository.ErrSellerNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			adverts, err := service.GetByUserId(userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(expectedAdverts), len(adverts))
			}
		})
	}
}

func TestAdvertService_Get(t *testing.T) {
	service, advertRepo, _, _, ctrl := setupAdvertService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	expectedAdverts := []*entity.Advert{
		{
			ID:          uuid.New(),
			SellerId:    uuid.New(),
			CategoryId:  uuid.New(),
			Title:       "Advert 1",
			Price:       100,
			Status:      entity.AdvertStatusActive,
			HasDelivery: true,
			Location:    "Location 1",
		},
		{
			ID:          uuid.New(),
			SellerId:    uuid.New(),
			CategoryId:  uuid.New(),
			Title:       "Advert 2",
			Price:       200,
			Status:      entity.AdvertStatusInactive,
			HasDelivery: false,
			Location:    "Location 2",
		},
	}

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				advertRepo.EXPECT().Get(10, 0, userID).Return(expectedAdverts, nil)
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			adverts, err := service.Get(10, 0, userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(expectedAdverts), len(adverts))
			}
		})
	}
}

func TestAdvertService_GetByCartId(t *testing.T) {
	service, advertRepo, _, _, ctrl := setupAdvertService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	cartID := uuid.New()
	expectedAdverts := []*entity.Advert{
		{
			ID:          uuid.New(),
			SellerId:    uuid.New(),
			CategoryId:  uuid.New(),
			Title:       "Advert 1",
			Price:       100,
			Status:      entity.AdvertStatusActive,
			HasDelivery: true,
			Location:    "Location 1",
		},
		{
			ID:          uuid.New(),
			SellerId:    uuid.New(),
			CategoryId:  uuid.New(),
			Title:       "Advert 2",
			Price:       200,
			Status:      entity.AdvertStatusInactive,
			HasDelivery: false,
			Location:    "Location 2",
		},
	}

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				advertRepo.EXPECT().GetByCartId(cartID, userID).Return(expectedAdverts, nil)
			},
			expectedError: nil,
		},
		{
			name: "Cart not found",
			setupMocks: func() {
				advertRepo.EXPECT().GetByCartId(cartID, userID).Return(nil, repository.ErrCartNotFound)
			},
			expectedError: repository.ErrCartNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			adverts, err := service.GetByCartId(cartID, userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(expectedAdverts), len(adverts))
			}
		})
	}
}

func TestAdvertService_UploadImage(t *testing.T) {
	service, advertRepo, sellerRepo, _, ctrl := setupAdvertService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	advertID := uuid.New()
	imageID := uuid.New()
	sellerID := uuid.New()

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				sellerRepo.EXPECT().GetByUserId(userID).Return(&entity.Seller{ID: sellerID}, nil)
				advertRepo.EXPECT().GetById(advertID, userID).Return(&entity.Advert{SellerId: sellerID}, nil)
				advertRepo.EXPECT().UploadImage(advertID, imageID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Advert Not Found",
			setupMocks: func() {
				sellerRepo.EXPECT().GetByUserId(userID).Return(&entity.Seller{ID: sellerID}, nil)
				advertRepo.EXPECT().GetById(advertID, userID).Return(nil, repository.ErrAdvertNotFound)
			},
			expectedError: repository.ErrAdvertNotFound,
		},
		{
			name: "Forbidden",
			setupMocks: func() {
				sellerRepo.EXPECT().GetByUserId(userID).Return(&entity.Seller{ID: sellerID}, nil)
				advertRepo.EXPECT().GetById(advertID, userID).Return(&entity.Advert{SellerId: uuid.New()}, nil)
			},
			expectedError: ErrForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			err := service.UploadImage(advertID, imageID, userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAdvertService_DeleteFromSaved(t *testing.T) {
	service, advertRepo, _, _, ctrl := setupAdvertService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	advertID := uuid.New()

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				advertRepo.EXPECT().CheckIfExists(advertID).Return(true, nil)
				advertRepo.EXPECT().DeleteFromSaved(advertID, userID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Advert Not Found",
			setupMocks: func() {
				advertRepo.EXPECT().CheckIfExists(advertID).Return(false, nil)
			},
			expectedError: ErrAdvertNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			err := service.RemoveFromSaved(advertID, userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAdvertService_GetSavedByUserId(t *testing.T) {
	service, advertRepo, _, _, ctrl := setupAdvertService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	expectedAdverts := []*entity.Advert{
		{
			ID:          uuid.New(),
			SellerId:    uuid.New(),
			CategoryId:  uuid.New(),
			Title:       "Advert 1",
			Price:       100,
			Status:      entity.AdvertStatusActive,
			HasDelivery: true,
			Location:    "Location 1",
		},
		{
			ID:          uuid.New(),
			SellerId:    uuid.New(),
			CategoryId:  uuid.New(),
			Title:       "Advert 2",
			Price:       200,
			Status:      entity.AdvertStatusInactive,
			HasDelivery: false,
			Location:    "Location 2",
		},
	}

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				advertRepo.EXPECT().GetSavedByUserId(userID).Return(expectedAdverts, nil)
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			adverts, err := service.GetSavedByUserId(userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(expectedAdverts), len(adverts))
			}
		})
	}
}

func TestAdvertService_Search(t *testing.T) {
	service, advertRepo, _, _, ctrl := setupAdvertService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	query := "test"
	expectedAdverts := []*entity.Advert{
		{
			ID:          uuid.New(),
			SellerId:    uuid.New(),
			CategoryId:  uuid.New(),
			Title:       "Advert 1",
			Price:       100,
			Status:      entity.AdvertStatusActive,
			HasDelivery: true,
			Location:    "Location 1",
		},
		{
			ID:          uuid.New(),
			SellerId:    uuid.New(),
			CategoryId:  uuid.New(),
			Title:       "Advert 2",
			Price:       200,
			Status:      entity.AdvertStatusInactive,
			HasDelivery: false,
			Location:    "Location 2",
		},
	}

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				advertRepo.EXPECT().Count().Return(2, nil)
				advertRepo.EXPECT().Search(query, gomock.Any(), gomock.Any(), userID).Return(expectedAdverts, nil).AnyTimes()
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			adverts, err := service.Search(query, 10, 10, 0, userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(expectedAdverts), len(adverts))
			}
		})
	}
}

func TestAdvertService_GetByCategoryId(t *testing.T) {
	service, advertRepo, _, _, ctrl := setupAdvertService(t)
	defer ctrl.Finish()

	categoryID := uuid.New()
	userID := uuid.New()
	expectedAdverts := []*entity.Advert{
		{
			ID:          uuid.New(),
			SellerId:    uuid.New(),
			CategoryId:  categoryID,
			Title:       "Advert 1",
			Price:       100,
			Status:      entity.AdvertStatusActive,
			HasDelivery: true,
			Location:    "Location 1",
		},
	}

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				advertRepo.EXPECT().GetByCategoryId(categoryID, userID).Return(expectedAdverts, nil)
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			adverts, err := service.GetByCategoryId(categoryID, userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(expectedAdverts), len(adverts))
			}
		})
	}
}
