package service

import (
	"testing"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupCartTestService(t *testing.T) (*CartService, *gomock.Controller, *mocks.MockCart) {
	ctrl := gomock.NewController(t)
	mockCartRepo := mocks.NewMockCart(ctrl)
	logger := zap.NewNop()

	service := NewCartService(mockCartRepo, logger)

	return service, ctrl, mockCartRepo
}

func TestCartService_AddAdvertToUserCart(t *testing.T) {
	service, ctrl, mockCartRepo := setupCartTestService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	advertID := uuid.New()
	cartID := uuid.New()

	mockCartRepo.EXPECT().GetCartByUserID(userID).Return(entity.Cart{ID: cartID}, nil)
	mockCartRepo.EXPECT().AddAdvertToCart(cartID, advertID).Return(nil)

	err := service.AddAdvertToUserCart(userID, advertID)
	assert.NoError(t, err)
}

func TestCartService_AddAdvertToUserCart_CreateCart(t *testing.T) {
	service, ctrl, mockCartRepo := setupCartTestService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	advertID := uuid.New()
	cartID := uuid.New()

	mockCartRepo.EXPECT().GetCartByUserID(userID).Return(entity.Cart{}, repository.ErrCartNotFound)
	mockCartRepo.EXPECT().CreateCart(userID).Return(cartID, nil)
	mockCartRepo.EXPECT().AddAdvertToCart(cartID, advertID).Return(nil)

	err := service.AddAdvertToUserCart(userID, advertID)
	assert.NoError(t, err)
}

func TestCartService_GetCartByID(t *testing.T) {
	service, ctrl, mockCartRepo := setupCartTestService(t)
	defer ctrl.Finish()

	cartID := uuid.New()
	userID := uuid.New()
	adverts := []entity.Advert{{ID: uuid.New()}}

	mockCartRepo.EXPECT().GetAdvertsByCartID(cartID).Return(adverts, nil)
	mockCartRepo.EXPECT().GetCartByID(cartID).Return(entity.Cart{ID: cartID, UserID: userID}, nil)

	cart, err := service.GetCartByID(cartID)
	assert.NoError(t, err)
	assert.Equal(t, cartID, cart.ID)
	assert.Equal(t, userID, cart.UserID)
	assert.Len(t, cart.Adverts, 1)
}

func TestCartService_GetCartByUserID(t *testing.T) {
	service, ctrl, mockCartRepo := setupCartTestService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	cartID := uuid.New()

	mockCartRepo.EXPECT().GetCartByUserID(userID).Return(entity.Cart{ID: cartID}, nil)
	mockCartRepo.EXPECT().GetAdvertsByCartID(cartID).Return([]entity.Advert{{ID: uuid.New()}}, nil)
	mockCartRepo.EXPECT().GetCartByID(cartID).Return(entity.Cart{ID: cartID, UserID: userID}, nil)

	cart, err := service.GetCartByUserID(userID)
	assert.NoError(t, err)
	assert.Equal(t, cartID, cart.ID)
	assert.Equal(t, userID, cart.UserID)
	assert.Len(t, cart.Adverts, 1)
}
