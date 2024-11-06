package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
)

func setup(t *testing.T) (*PurchaseService, *mocks.MockPurchaseRepository, *mocks.MockCart, *mocks.MockAdvertRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	purchaseRepo := mocks.NewMockPurchaseRepository(ctrl)
	cartRepo := mocks.NewMockCart(ctrl)
	advertRepo := mocks.NewMockAdvertRepository(ctrl)
	logger := zap.NewNop()
	service := NewPurchaseService(purchaseRepo, advertRepo, cartRepo, logger)
	return service, purchaseRepo, cartRepo, advertRepo, ctrl
}

func TestPurchaseService_AddPurchase_FailureInBeginTransaction(t *testing.T) {
	service, purchaseRepo, _, _, ctrl := setup(t)
	defer ctrl.Finish()

	purchaseRepo.EXPECT().BeginTransaction().Return(nil, errors.New("begin transaction error"))

	purchaseRequest := dto.PurchaseRequest{
		CartID:         uuid.New(),
		Address:        "123 Street",
		PaymentMethod:  dto.PaymentMethodCard,
		DeliveryMethod: dto.DeliveryMethodPickup,
	}

	resp, err := service.AddPurchase(purchaseRequest)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to begin transaction")
}

func TestPurchaseService_GetPurchasesByUserID_Success(t *testing.T) {
	service, purchaseRepo, _, _, ctrl := setup(t)
	defer ctrl.Finish()

	userID := uuid.New()
	mockPurchases := []*entity.Purchase{
		{ID: uuid.New(), CartID: uuid.New(), Status: entity.StatusCompleted},
	}

	purchaseRepo.EXPECT().GetPurchasesByUserID(userID).Return(mockPurchases, nil)

	resp, err := service.GetPurchasesByUserID(userID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(mockPurchases), len(resp))
	assert.Equal(t, mockPurchases[0].ID, resp[0].ID)
}

func TestPurchaseService_GetPurchasesByUserID_Failure(t *testing.T) {
	service, purchaseRepo, _, _, ctrl := setup(t)
	defer ctrl.Finish()

	userID := uuid.New()

	purchaseRepo.EXPECT().GetPurchasesByUserID(userID).Return(nil, errors.New("database error"))

	resp, err := service.GetPurchasesByUserID(userID)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to get purchases")
}
