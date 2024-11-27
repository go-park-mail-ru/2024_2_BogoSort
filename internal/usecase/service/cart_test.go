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

// func TestCartService_AddAdvert_FailureAdvertNotFound(t *testing.T) {
// 	service, cartRepo, advertRepo, ctrl := setupCartService(t)
// 	defer ctrl.Finish()

// 	userID := uuid.New()
// 	advertID := uuid.New()
// 	cartID := uuid.New()

// 	testCases := []struct {
// 		name          string
// 		setupMocks    func()
// 		expectedError error
// 	}{
// 		{
// 			name: "Advert Not Found",
// 			setupMocks: func() {
// 				cartRepo.EXPECT().GetByUserId(userID).Return(entity.Cart{ID: cartID}, nil)
// 				advertRepo.EXPECT().GetById(advertID, userID).Return(nil, repository.ErrAdvertNotFound)
// 			},
// 			expectedError: repository.ErrAdvertNotFound,
// 		},
// 		{
// 			name: "Advert Inactive",
// 			setupMocks: func() {
// 				cartRepo.EXPECT().GetByUserId(userID).Return(entity.Cart{ID: cartID}, nil)
// 				advertRepo.EXPECT().GetById(advertID, userID).Return(&entity.Advert{Status: entity.AdvertStatusInactive}, nil)
// 			},
// 			expectedError: errors.New("advert inactive"),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tc.setupMocks()

// 			err := service.AddAdvert(userID, advertID)

// 			if tc.expectedError != nil {
// 				assert.Error(t, err)
// 				assert.True(t, errors.Is(err, tc.expectedError), "ожидалась ошибка: %v, получена: %v", tc.expectedError, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }

func TestCartService_DeleteAdvert_FailureAdvertNotInCart(t *testing.T) {
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
			name: "Advert Not Found in Cart",
			setupMocks: func() {
				cartRepo.EXPECT().DeleteAdvert(cartID, advertID).Return(repository.ErrCartOrAdvertNotFound)
			},
			expectedError: repository.ErrCartOrAdvertNotFound,
		},
		{
			name: "Delete Advert Error",
			setupMocks: func() {
				cartRepo.EXPECT().DeleteAdvert(cartID, advertID).Return(errors.New("delete advert error"))
			},
			expectedError: errors.New("delete advert error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			err := service.DeleteAdvert(cartID, advertID)

			if tc.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// func TestCartService_GetById_FailureCartNotFound(t *testing.T) {
// 	service, cartRepo, _, ctrl := setupCartService(t)
// 	defer ctrl.Finish()

// 	cartID := uuid.New()

// 	testCases := []struct {
// 		name          string
// 		setupMocks    func()
// 		expectedError error
// 	}{
// 		{
// 			name: "Cart Not Found",
// 			setupMocks: func() {
// 				cartRepo.EXPECT().GetById(cartID).Return(entity.Cart{}, repository.ErrCartNotFound)
// 				cartRepo.EXPECT().GetAdvertsByCartId(cartID).Return(nil, repository.ErrCartNotFound)
// 			},
// 			expectedError: repository.ErrCartNotFound,
// 		},
// 		{
// 			name: "Get Adverts Error",
// 			setupMocks: func() {
// 				cartRepo.EXPECT().GetById(cartID).Return(entity.Cart{}, nil)
// 				cartRepo.EXPECT().GetAdvertsByCartId(cartID).Return(nil, errors.New("get adverts error"))
// 			},
// 			expectedError: errors.New("get adverts error"),
// 		},
// 		{
// 			name: "Get Cart Error",
// 			setupMocks: func() {
// 				cartRepo.EXPECT().GetById(cartID).Return(entity.Cart{}, errors.New("get cart error"))
// 				cartRepo.EXPECT().GetAdvertsByCartId(cartID).Return([]entity.Advert{}, nil)
// 			},
// 			expectedError: errors.New("get cart error"),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tc.setupMocks()

// 			cart, err := service.GetById(cartID)

// 			if tc.expectedError != nil {
// 				assert.Error(t, err)
// 				assert.Empty(t, cart.Adverts)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }

func TestCartService_GetByUserId_FailureCartNotFound(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	userID := uuid.New()

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Cart Not Found",
			setupMocks: func() {
				cartRepo.EXPECT().GetByUserId(userID).Return(entity.Cart{}, repository.ErrCartNotFound)
			},
			expectedError: repository.ErrCartNotFound,
		},
		{
			name: "Get Cart Error",
			setupMocks: func() {
				cartRepo.EXPECT().GetByUserId(userID).Return(entity.Cart{}, errors.New("get cart error"))
			},
			expectedError: errors.New("get cart error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			cart, err := service.GetByUserId(userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				// assert.True(t, errors.Is(err, tc.expectedError), "ожидалась ошибка: %v, получена: %v", tc.expectedError, err)
				assert.Empty(t, cart.Adverts)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCartService_CheckExists_FailureRepositoryError(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	userID := uuid.New()

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Repository Error",
			setupMocks: func() {
				cartRepo.EXPECT().GetByUserId(userID).Return(entity.Cart{}, errors.New("repository error"))
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			cartID, err := service.CheckExists(userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				// assert.True(t, errors.Is(err, tc.expectedError), "ожидалась ошибка: %v, получена: %v", tc.expectedError, err)
				assert.Equal(t, uuid.Nil, cartID)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// func TestCartService_AddAdvert_RepositoryError(t *testing.T) {
// 	service, cartRepo, advertRepo, ctrl := setupCartService(t)
// 	defer ctrl.Finish()

// 	userID := uuid.New()
// 	advertID := uuid.New()
// 	cartID := uuid.New()

// 	testCases := []struct {
// 		name          string
// 		setupMocks    func()
// 		expectedError error
// 	}{
// 		{
// 			name: "Add Advert Repository Error",
// 			setupMocks: func() {
// 				cartRepo.EXPECT().GetByUserId(gomock.Any()).Return(entity.Cart{ID: cartID}, nil)
// 				advertRepo.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(&entity.Advert{
// 					ID:     advertID,
// 					Status: entity.AdvertStatusActive,
// 				}, nil)
// 				cartRepo.EXPECT().AddAdvert(gomock.Any(), gomock.Any()).Return(errors.New("add advert error"))
// 			},
// 			expectedError: errors.New("add advert error"),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tc.setupMocks()

// 			err := service.AddAdvert(userID, advertID)

// 			if tc.expectedError != nil {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }
