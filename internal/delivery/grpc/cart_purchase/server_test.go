package cart_purchase

// import (
// 	"context"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
// 	cartPurchaseProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/cart_purchase/proto"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
// )

// type MockCartService struct {
// 	mock.Mock
// }

// func (m *MockCartService) DeleteAdvert(cartID uuid.UUID, advertID uuid.UUID) error {
// 	args := m.Called(cartID, advertID)
// 	return args.Error(0)
// }

// func (m *MockCartService) CheckExists(userID uuid.UUID) (uuid.UUID, error) {
// 	args := m.Called(userID)
// 	return args.Get(0).(uuid.UUID), args.Error(1)
// }

// func (m *MockCartService) GetById(cartID uuid.UUID) (dto.Cart, error) {
// 	args := m.Called(cartID)
// 	return args.Get(0).(dto.Cart), args.Error(1)
// }

// func (m *MockCartService) GetByUserId(userID uuid.UUID) (dto.Cart, error) {
// 	args := m.Called(userID)
// 	return args.Get(0).(dto.Cart), args.Error(1)
// }

// func (m *MockCartService) AddAdvert(userID uuid.UUID, advertID uuid.UUID) error {
// 	args := m.Called(userID, advertID)
// 	return args.Error(0)
// }

// type MockPurchaseService struct {
// 	mock.Mock
// }

// func (m *MockPurchaseService) Add(req dto.PurchaseRequest, userID uuid.UUID) (*dto.PurchaseResponse, error) {
// 	args := m.Called(req, userID)
// 	return args.Get(0).(*dto.PurchaseResponse), args.Error(1)
// }

// func (m *MockPurchaseService) GetPurchasesByUserID(userID uuid.UUID) ([]dto.PurchaseResponse, error) {
// 	args := m.Called(userID)
// 	return args.Get(0).([]dto.PurchaseResponse), args.Error(1)
// }

// func (m *MockPurchaseService) GetByUserId(userID uuid.UUID) ([]*dto.PurchaseResponse, error) {
// 	args := m.Called(userID)
// 	return args.Get(0).([]*dto.PurchaseResponse), args.Error(1)
// }

// func TestServerDeleteAdvertFromCart(t *testing.T) {
// 	mockCartUC := new(MockCartService)
// 	mockService := new(MockPurchaseService)
// 	server := NewGrpcServer(mockCartUC, mockService)

// 	req := &cartPurchaseProto.DeleteAdvertFromCartRequest{
// 		CartId:   uuid.New().String(),
// 		AdvertId: uuid.New().String(),
// 	}

// 	mockCartUC.On("DeleteAdvert", mock.Anything, mock.Anything).Return(nil)

// 	result, err := server.DeleteAdvertFromCart(context.Background(), req)

// 	assert.NoError(t, err)
// 	assert.Equal(t, "advert deleted from user cart", result.Message)
// 	mockCartUC.AssertExpectations(t)
// }

// func TestServerCheckCartExists(t *testing.T) {
// 	mockCartUC := new(MockCartService)
// 	mockService := new(MockPurchaseService)
// 	server := NewGrpcServer(mockCartUC, mockService)

// 	userID := uuid.New()
// 	req := &cartPurchaseProto.CheckCartExistsRequest{
// 		UserId: userID.String(),
// 	}

// 	resp := &cartPurchaseProto.CheckCartExistsResponse{
// 		CartId: uuid.New().String(),
// 	}

// 	mockCartUC.On("CheckExists", mock.Anything).Return(uuid.MustParse(resp.CartId), nil)

// 	result, err := server.CheckCartExists(context.Background(), req)

// 	assert.NoError(t, err)
// 	assert.NotEqual(t, uuid.Nil.String(), result.CartId)
// 	mockCartUC.AssertExpectations(t)
// }

// func TestServerGetCartByID(t *testing.T) {
// 	mockCartUC := new(MockCartService)
// 	mockService := new(MockPurchaseService)
// 	server := NewGrpcServer(mockCartUC, mockService)

// 	cartID := uuid.New()
// 	req := &cartPurchaseProto.GetCartByIDRequest{
// 		CartId: cartID.String(),
// 	}

// 	mockCart := dto.Cart{
// 		ID: cartID,
// 		UserID: uuid.New(),
// 		Status: entity.CartStatusActive,
// 		Adverts: nil,
// 	}

// 	mockCartUC.On("GetById", mock.Anything).Return(mockCart, nil)

// 	result, err := server.GetCartByID(context.Background(), req)

// 	assert.NoError(t, err)
// 	assert.Equal(t, cartID.String(), result.Cart.Id)
// 	mockCartUC.AssertExpectations(t)
// }
