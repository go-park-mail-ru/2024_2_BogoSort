package service

import (
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupCartService(t *testing.T) (*CartService, *mocks.MockCart, *mocks.MockAdvertRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	cartRepo := mocks.NewMockCart(ctrl)
	advertRepo := mocks.NewMockAdvertRepository(ctrl)
	service := NewCartService(cartRepo, advertRepo)
	return service, cartRepo, advertRepo, ctrl
}

func TestCartService_AddAdvert(t *testing.T) {
	service, cartRepo, advertRepo, ctrl := setupCartService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	advertID := uuid.New()
	cartID := uuid.New()

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				cartRepo.EXPECT().GetByUserId(userID).Return(entity.Cart{ID: cartID}, nil)
				advertRepo.EXPECT().GetById(advertID, userID).Return(&entity.Advert{Status: entity.AdvertStatusActive}, nil)
				cartRepo.EXPECT().AddAdvert(cartID, advertID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Cart Not Found",
			setupMocks: func() {
				cartRepo.EXPECT().GetByUserId(userID).Return(entity.Cart{}, repository.ErrCartNotFound)
				cartRepo.EXPECT().Create(userID).Return(cartID, nil)
				cartRepo.EXPECT().AddAdvert(cartID, advertID).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			err := service.AddAdvert(userID, advertID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCartService_DeleteAdvert(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	cartID := uuid.New()
	advertID := uuid.New()

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				cartRepo.EXPECT().DeleteAdvert(cartID, advertID).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			err := service.DeleteAdvert(cartID, advertID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCartService_GetById(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	cartID := uuid.New()
	expectedAdverts := []entity.Advert{
		{
			ID:          uuid.New(),
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
				cartRepo.EXPECT().GetAdvertsByCartId(cartID).Return(expectedAdverts, nil)
				cartRepo.EXPECT().GetById(cartID).Return(entity.Cart{ID: cartID, UserID: uuid.New(), Status: "active"}, nil)
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			cart, err := service.GetById(cartID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(expectedAdverts), len(cart.Adverts))
			}
		})
	}
}

func TestCartService_GetByUserId(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	userID := uuid.New()

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func() {
				cartRepo.EXPECT().GetByUserId(userID).Return(entity.Cart{ID: uuid.New(), UserID: userID, Status: "active"}, nil)
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			cart, err := service.GetByUserId(userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, userID, cart.UserID)
			}
		})
	}
}

func TestCartService_CheckExists(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	cartID := uuid.New()

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Exists",
			setupMocks: func() {
				cartRepo.EXPECT().GetByUserId(userID).Return(entity.Cart{ID: cartID}, nil)
			},
			expectedError: nil,
		},
		{
			name: "Does Not Exist",
			setupMocks: func() {
				cartRepo.EXPECT().GetByUserId(userID).Return(entity.Cart{}, repository.ErrCartNotFound)
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			cartID, err := service.CheckExists(userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				if tc.name == "Exists" {
					assert.Equal(t, cartID, cartID)
				} else {
					assert.Equal(t, uuid.Nil, cartID)
				}
			}
		})
	}
}
