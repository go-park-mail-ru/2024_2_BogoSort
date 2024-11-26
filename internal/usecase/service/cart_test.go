package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
)

func setupCartService(t *testing.T) (*CartService, *mocks.MockCart, *mocks.MockAdvertRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	cartRepo := mocks.NewMockCart(ctrl)
	advertRepo := mocks.NewMockAdvertRepository(ctrl)
	service := NewCartService(cartRepo, advertRepo)
	return service, cartRepo, advertRepo, ctrl
}

func TestCartService_AddAdvertToUserCart_CartExists(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	advertID := uuid.New()
	cartID := uuid.New()

	cartRepo.EXPECT().GetByUserId(userID).Return(entity.Cart{ID: cartID}, nil)
	cartRepo.EXPECT().AddAdvert(cartID, advertID).Return(nil)

	err := service.AddAdvert(userID, advertID)

	assert.NoError(t, err)
}

func TestCartService_DeleteAdvertFromCart(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	cartID := uuid.New()
	advertID := uuid.New()

	cartRepo.EXPECT().DeleteAdvert(cartID, advertID).Return(nil)

	err := service.DeleteAdvert(cartID, advertID)

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

	cartRepo.EXPECT().GetAdvertsByCartId(cartID).Return(adverts, nil)
	cartRepo.EXPECT().GetById(cartID).Return(entity.Cart{ID: cartID, UserID: userID, Status: entity.CartStatusActive}, nil)

	cart, err := service.GetById(cartID)

	assert.NoError(t, err)
	assert.Equal(t, cartID, cart.ID)
	assert.Equal(t, userID, cart.UserID)
	assert.Equal(t, len(adverts), len(cart.Adverts))
}

func TestCartService_CheckCartExists_Found(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	userID := uuid.New()

	cartRepo.EXPECT().GetByUserId(userID).Return(entity.Cart{ID: uuid.New()}, nil)

	cartID, err := service.CheckExists(userID)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, cartID)
}

func TestCartService_CheckCartExists_NotFound(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	userID := uuid.New()

	cartRepo.EXPECT().GetByUserId(userID).Return(entity.Cart{}, repository.ErrCartNotFound)

	cartID, err := service.CheckExists(userID)

	assert.NoError(t, err)
	assert.Equal(t, uuid.Nil, cartID)
}

func TestCartService_GetCartByID_NotFound(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	cartID := uuid.New()

	cartRepo.EXPECT().GetAdvertsByCartId(cartID).Return(nil, repository.ErrCartNotFound)

	cart, err := service.GetById(cartID)

	assert.Error(t, err)
	assert.Equal(t, dto.Cart{}, cart)
}

// func TestCartService_GetCartByUserID_ErrorGettingAdverts(t *testing.T) {
// 	service, cartRepo, _, ctrl := setupCartService(t)
// 	defer ctrl.Finish()

// 	userID := uuid.New()
// 	cartID := uuid.New()

// 	cartRepo.EXPECT().GetByUserId(userID).Return(entity.Cart{ID: cartID}, nil)
// 	cartRepo.EXPECT().GetAdvertsByCartId(cartID).Return(nil, errors.New("database error"))

// 	cart, err := service.cartRepo.GetByUserId(userID)

// 	assert.Error(t, err)
// 	assert.Equal(t, dto.Cart{}, cart)
// }

func TestCartService_CheckCartExists_Error(t *testing.T) {
	service, cartRepo, _, ctrl := setupCartService(t)
	defer ctrl.Finish()

	userID := uuid.New()

	cartRepo.EXPECT().GetByUserId(userID).Return(entity.Cart{}, errors.New("database error"))

	cartID, err := service.CheckExists(userID)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, cartID)
}
