package service

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
// 	"github.com/golang/mock/gomock"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// )

// func setupPurchaseService(t *testing.T) (*PurchaseService, *mocks.MockPurchaseRepository, *mocks.MockCart, *mocks.MockAdvertRepository, *gomock.Controller) {
// 	ctrl := gomock.NewController(t)
// 	purchaseRepo := mocks.NewMockPurchaseRepository(ctrl)
// 	cartRepo := mocks.NewMockCart(ctrl)
// 	advertRepo := mocks.NewMockAdvertRepository(ctrl)
// 	service := NewPurchaseService(purchaseRepo, advertRepo, cartRepo)
// 	return service, purchaseRepo, cartRepo, advertRepo, ctrl
// }

// func TestPurchaseService_Add(t *testing.T) {
// 	service, purchaseRepo, cartRepo, advertRepo, ctrl := setupPurchaseService(t)
// 	defer ctrl.Finish()

// 	userID := uuid.New()
// 	cartID := uuid.New()
// 	purchaseRequest := dto.PurchaseRequest{
// 		CartID:         cartID,
// 		Address:        "123 Test St",
// 		PaymentMethod:  dto.PaymentMethodCard,
// 		DeliveryMethod: dto.DeliveryMethodDelivery,
// 	}

// 	testCases := []struct {
// 		name          string
// 		setupMocks    func()
// 		expectedError error
// 	}{
// 		{
// 			name: "Success",
// 			setupMocks: func() {
// 				tx := mocks.NewMockTransaction(ctrl)
// 				purchaseRepo.EXPECT().BeginTransaction().Return(tx, nil)
// 				purchaseRepo.EXPECT().Add(tx, gomock.Any()).Return(&entity.Purchase{ID: uuid.New(), CartID: cartID}, nil)
// 				cartRepo.EXPECT().UpdateStatus(tx, cartID, entity.CartStatusInactive).Return(nil)
// 				advertRepo.EXPECT().GetByCartId(cartID, userID).Return([]*entity.Advert{{ID: uuid.New()}}, nil)
// 				advertRepo.EXPECT().UpdateStatus(tx, gomock.Any(), entity.AdvertStatusReserved).Return(nil)
// 				tx.EXPECT().Commit(context.Background()).Return(nil)
// 			},
// 			expectedError: nil,
// 		},
// 		{
// 			name: "Transaction Begin Error",
// 			setupMocks: func() {
// 				purchaseRepo.EXPECT().BeginTransaction().Return(nil, errors.New("transaction error"))
// 			},
// 			expectedError: errors.New("transaction error"),
// 		},
// 		{
// 			name: "Add Purchase Error",
// 			setupMocks: func() {
// 				tx := mocks.NewMockTransaction(ctrl)
// 				purchaseRepo.EXPECT().BeginTransaction().Return(tx, nil)
// 				purchaseRepo.EXPECT().Add(tx, gomock.Any()).Return(nil, errors.New("add purchase error"))
// 				tx.EXPECT().Rollback(context.Background()).Return(nil)
// 			},
// 			expectedError: errors.New("add purchase error"),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tc.setupMocks()

// 			purchase, err := service.Add(purchaseRequest, userID)

// 			if tc.expectedError != nil {
// 				assert.Error(t, err)
// 				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.NotNil(t, purchase)
// 			}
// 		})
// 	}
// }

// func TestPurchaseService_GetByUserId(t *testing.T) {
// 	service, purchaseRepo, _, _, ctrl := setupPurchaseService(t)
// 	defer ctrl.Finish()

// 	userID := uuid.New()
// 	expectedPurchases := []*entity.Purchase{
// 		{ID: uuid.New(), CartID: uuid.New(), Address: "123 Test St"},
// 	}

// 	testCases := []struct {
// 		name          string
// 		setupMocks    func()
// 		expectedError error
// 	}{
// 		{
// 			name: "Success",
// 			setupMocks: func() {
// 				purchaseRepo.EXPECT().GetByUserId(userID).Return(expectedPurchases, nil)
// 			},
// 			expectedError: nil,
// 		},
// 		{
// 			name: "Get Purchases Error",
// 			setupMocks: func() {
// 				purchaseRepo.EXPECT().GetByUserId(userID).Return(nil, errors.New("get purchases error"))
// 			},
// 			expectedError: errors.New("get purchases error"),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tc.setupMocks()

// 			purchases, err := service.GetByUserId(userID)

// 			if tc.expectedError != nil {
// 				assert.Error(t, err)
// 				assert.True(t, errors.Is(err, tc.expectedError), "expected error: %v, got: %v", tc.expectedError, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, len(expectedPurchases), len(purchases))
// 			}
// 		})
// 	}
// }
