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
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
)

func setupCartService(t *testing.T) (*CartService, *mocks.MockCart, *mocks.MockAdvertRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	cartRepo := mocks.NewMockCart(ctrl)
	advertRepo := mocks.NewMockAdvertRepository(ctrl)
	logger := zap.NewNop()
	service := NewCartService(cartRepo, advertRepo, logger)
	return service, cartRepo, advertRepo, ctrl
}

func TestCartService_AddAdvertToUserCart_CartExists(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	advertID := uuid.New()
	cartID := uuid.New()

	cartRepo.EXPECT().GetCartByUserID(userID).Return(entity.Cart{ID: cartID}, nil)
	cartRepo.EXPECT().AddAdvertToCart(cartID, advertID).Return(nil)

	err := service.AddAdvertToUserCart(userID, advertID)

	assert.NoError(t, err)
}

func TestCartService_DeleteAdvertFromCart(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	cartID := uuid.New()
	advertID := uuid.New()

	cartRepo.EXPECT().DeleteAdvertFromCart(cartID, advertID).Return(nil)

	err := service.DeleteAdvertFromCart(cartID, advertID)

	assert.NoError(t, err)
}

func TestCartService_GetCartByID_Success(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	cartID := uuid.New()
	userID := uuid.New()
	adverts := []entity.Advert{
		{ID: uuid.New(), Title: "Advert 1", Price: 100},
		{ID: uuid.New(), Title: "Advert 2", Price: 200},
	}

	cartRepo.EXPECT().GetAdvertsByCartID(cartID).Return(adverts, nil)
	cartRepo.EXPECT().GetCartByID(cartID).Return(entity.Cart{ID: cartID, UserID: userID, Status: entity.CartStatusActive}, nil)

	cart, err := service.GetCartByID(cartID)

	assert.NoError(t, err)
	assert.Equal(t, cartID, cart.ID)
	assert.Equal(t, userID, cart.UserID)
	assert.Equal(t, len(adverts), len(cart.Adverts))
}

func TestCartService_GetCartByUserID_Success(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	cartID := uuid.New()

	cartRepo.EXPECT().GetCartByUserID(userID).Return(entity.Cart{ID: cartID}, nil)
	cartRepo.EXPECT().GetAdvertsByCartID(cartID).Return([]entity.Advert{
		{ID: uuid.New(), Title: "Advert 1", Price: 100},
	}, nil)
	cartRepo.EXPECT().GetCartByID(cartID).Return(entity.Cart{ID: cartID, UserID: userID, Status: entity.CartStatusActive}, nil)

	cart, err := service.GetCartByUserID(userID)

	assert.NoError(t, err)
	assert.Equal(t, cartID, cart.ID)
}

func TestCartService_CheckCartExists_Found(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	userID := uuid.New()

	cartRepo.EXPECT().GetCartByUserID(userID).Return(entity.Cart{ID: uuid.New()}, nil)

	exists, err := service.CheckCartExists(userID)

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestCartService_CheckCartExists_NotFound(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	userID := uuid.New()

	cartRepo.EXPECT().GetCartByUserID(userID).Return(entity.Cart{}, repository.ErrCartNotFound)

	exists, err := service.CheckCartExists(userID)

	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestCartService_AddAdvertToUserCart_CartNotExists(t *testing.T) {
	service, cartRepo, advertRepo, ctrl := setupCartService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	advertID := uuid.New()
	newCartID := uuid.New()

	// Настройка ожидаемых вызовов
	cartRepo.EXPECT().GetCartByUserID(userID).Return(entity.Cart{}, repository.ErrCartNotFound)
	cartRepo.EXPECT().CreateCart(userID).Return(newCartID, nil)
	advert := entity.Advert{
		ID:     advertID,
		Title:  "New Advert",
		Price:  150,
		Status: entity.AdvertStatusActive,
	}
	advertRepo.EXPECT().GetAdvertById(advertID).Return(&advert, nil)
	cartRepo.EXPECT().AddAdvertToCart(newCartID, advertID).Return(nil)

	// Выполнение действия
	err := service.AddAdvertToUserCart(userID, advertID)

	// Проверка результата
	assert.NoError(t, err)
}

func TestCartService_GetCartByID_NotFound(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	cartID := uuid.New()

	cartRepo.EXPECT().GetAdvertsByCartID(cartID).Return(nil, repository.ErrCartNotFound)

	cart, err := service.GetCartByID(cartID)

	assert.Error(t, err)
	assert.Equal(t, dto.Cart{}, cart)
}

func TestCartService_GetCartByUserID_ErrorGettingAdverts(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	cartID := uuid.New()

	cartRepo.EXPECT().GetCartByUserID(userID).Return(entity.Cart{ID: cartID}, nil)
	cartRepo.EXPECT().GetAdvertsByCartID(cartID).Return(nil, errors.New("database error"))

	cart, err := service.GetCartByUserID(userID)

	assert.Error(t, err)
	assert.Equal(t, dto.Cart{}, cart)
}

func TestCartService_CheckCartExists_Error(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	userID := uuid.New()

	cartRepo.EXPECT().GetCartByUserID(userID).Return(entity.Cart{}, errors.New("database error"))

	exists, err := service.CheckCartExists(userID)

	assert.Error(t, err)
	assert.False(t, exists)
}
