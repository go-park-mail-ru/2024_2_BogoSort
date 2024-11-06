package service

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	pgxmock "github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupPurchaseTestService(t *testing.T) (*PurchaseService, *gomock.Controller, *mocks.MockPurchaseRepository, *mocks.MockCart, pgxmock.PgxConnIface) {
	ctrl := gomock.NewController(t)
	mockPurchaseRepo := mocks.NewMockPurchaseRepository(ctrl)
	mockCartRepo := mocks.NewMockCart(ctrl)
	logger := zap.NewNop()

	service := &PurchaseService{
		purchaseRepo: mockPurchaseRepo,
		cartRepo:     mockCartRepo,
		logger:       logger,
	}

	mockConn, err := pgxmock.NewConn()
	assert.NoError(t, err)

	return service, ctrl, mockPurchaseRepo, mockCartRepo, mockConn
}

func TestPurchaseService_AddPurchase_Success(t *testing.T) {
	service, ctrl, mockPurchaseRepo, mockCartRepo, mockConn := setupPurchaseTestService(t)
	defer ctrl.Finish()
	defer mockConn.Close(context.Background())

	purchaseRequest := dto.PurchaseRequest{
		CartID:         uuid.New(),
		Address:        "Test Address",
		PaymentMethod:  "credit_card",
		DeliveryMethod: "standard",
	}

	mockPurchaseRepo.EXPECT().BeginTransaction().Return(mockConn, nil)

	mockPurchaseRepo.EXPECT().AddPurchase(mockConn, gomock.Any()).Return(&entity.Purchase{
		ID:      uuid.New(),
		CartID:  purchaseRequest.CartID,
		Address: purchaseRequest.Address,
		Status:  entity.StatusPending,
	}, nil)

	mockCartRepo.EXPECT().UpdateCartStatus(mockConn, purchaseRequest.CartID, entity.CartStatusInactive).Return(nil)

	mockConn.ExpectCommit()

	result, err := service.AddPurchase(purchaseRequest)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, purchaseRequest.CartID, result.CartID)

	err = mockConn.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPurchaseService_AddPurchase_BeginTransactionError(t *testing.T) {
	service, ctrl, mockPurchaseRepo, _, _ := setupPurchaseTestService(t)
	defer ctrl.Finish()

	purchaseRequest := dto.PurchaseRequest{
		CartID:         uuid.New(),
		Address:        "Test Address",
		PaymentMethod:  "credit_card",
		DeliveryMethod: "standard",
	}

	mockPurchaseRepo.EXPECT().BeginTransaction().Return(nil, errors.New("transaction error"))

	result, err := service.AddPurchase(purchaseRequest)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPurchaseService_AddPurchase_AddPurchaseError(t *testing.T) {
	service, ctrl, mockPurchaseRepo, _, mockConn := setupPurchaseTestService(t)
	defer ctrl.Finish()
	defer mockConn.Close(context.Background())

	purchaseRequest := dto.PurchaseRequest{
		CartID:         uuid.New(),
		Address:        "Test Address",
		PaymentMethod:  "credit_card",
		DeliveryMethod: "standard",
	}

	mockPurchaseRepo.EXPECT().BeginTransaction().Return(mockConn, nil)

	mockPurchaseRepo.EXPECT().AddPurchase(mockConn, gomock.Any()).Return(nil, errors.New("add purchase error"))

	mockConn.ExpectRollback()

	result, err := service.AddPurchase(purchaseRequest)
	assert.Error(t, err)
	assert.Nil(t, result)

	err = mockConn.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPurchaseService_AddPurchase_UpdateCartStatusError(t *testing.T) {
	service, ctrl, mockPurchaseRepo, mockCartRepo, mockConn := setupPurchaseTestService(t)
	defer ctrl.Finish()
	defer mockConn.Close(context.Background())

	purchaseRequest := dto.PurchaseRequest{
		CartID:         uuid.New(),
		Address:        "Test Address",
		PaymentMethod:  "credit_card",
		DeliveryMethod: "standard",
	}

	mockPurchaseRepo.EXPECT().BeginTransaction().Return(mockConn, nil)

	mockPurchaseRepo.EXPECT().AddPurchase(mockConn, gomock.Any()).Return(&entity.Purchase{
		ID:      uuid.New(),
		CartID:  purchaseRequest.CartID,
		Address: purchaseRequest.Address,
		Status:  entity.StatusPending,
	}, nil)

	mockCartRepo.EXPECT().UpdateCartStatus(mockConn, purchaseRequest.CartID, entity.CartStatusInactive).Return(errors.New("update cart status error"))

	mockConn.ExpectRollback()

	result, err := service.AddPurchase(purchaseRequest)
	assert.Error(t, err)
	assert.Nil(t, result)

	err = mockConn.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPurchaseService_purchaseEntityToDTO(t *testing.T) {
	service, _, _, _, _ := setupPurchaseTestService(t)

	purchase := &entity.Purchase{
		ID:             uuid.New(),
		CartID:         uuid.New(),
		Address:        "Test Address",
		Status:         entity.StatusPending,
		PaymentMethod:  entity.PaymentMethod("credit_card"),
		DeliveryMethod: entity.DeliveryMethod("standard"),
	}

	dtoResponse, err := service.purchaseEntityToDTO(purchase)
	assert.NoError(t, err)
	assert.Equal(t, purchase.ID, dtoResponse.ID)
	assert.Equal(t, purchase.CartID, dtoResponse.CartID)
	assert.Equal(t, purchase.Address, dtoResponse.Address)
	assert.Equal(t, dto.PurchaseStatus(purchase.Status), dtoResponse.Status)
	assert.Equal(t, dto.PaymentMethod(purchase.PaymentMethod), dtoResponse.PaymentMethod)
	assert.Equal(t, dto.DeliveryMethod(purchase.DeliveryMethod), dtoResponse.DeliveryMethod)
}
