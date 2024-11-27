package cart_purchase

import (
	"context"
	"testing"

	cartPurchaseProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/cart_purchase/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCartPurchaseServiceServer - мок для CartPurchaseServiceServer
type MockCartPurchaseServiceServer struct {
	mock.Mock
}

func (m *MockCartPurchaseServiceServer) AddPurchase(ctx context.Context, req *cartPurchaseProto.AddPurchaseRequest) (*cartPurchaseProto.AddPurchaseResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*cartPurchaseProto.AddPurchaseResponse), args.Error(1)
}

func (m *MockCartPurchaseServiceServer) GetPurchasesByUserID(ctx context.Context, req *cartPurchaseProto.GetPurchasesByUserIDRequest) (*cartPurchaseProto.GetPurchasesByUserIDResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*cartPurchaseProto.GetPurchasesByUserIDResponse), args.Error(1)
}

func (m *MockCartPurchaseServiceServer) GetCartByID(ctx context.Context, req *cartPurchaseProto.GetCartByIDRequest) (*cartPurchaseProto.GetCartByIDResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*cartPurchaseProto.GetCartByIDResponse), args.Error(1)
}

func (m *MockCartPurchaseServiceServer) GetCartByUserID(ctx context.Context, req *cartPurchaseProto.GetCartByUserIDRequest) (*cartPurchaseProto.GetCartByUserIDResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*cartPurchaseProto.GetCartByUserIDResponse), args.Error(1)
}

func (m *MockCartPurchaseServiceServer) AddAdvertToCart(ctx context.Context, req *cartPurchaseProto.AddAdvertToCartRequest) (*cartPurchaseProto.AddAdvertToCartResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*cartPurchaseProto.AddAdvertToCartResponse), args.Error(1)
}

func (m *MockCartPurchaseServiceServer) DeleteAdvertFromCart(ctx context.Context, req *cartPurchaseProto.DeleteAdvertFromCartRequest) (*cartPurchaseProto.DeleteAdvertFromCartResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*cartPurchaseProto.DeleteAdvertFromCartResponse), args.Error(1)
}

func (m *MockCartPurchaseServiceServer) CheckCartExists(ctx context.Context, req *cartPurchaseProto.CheckCartExistsRequest) (*cartPurchaseProto.CheckCartExistsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*cartPurchaseProto.CheckCartExistsResponse), args.Error(1)
}

func TestServerAddPurchase(t *testing.T) {
	mockServer := new(MockCartPurchaseServiceServer)

	req := &cartPurchaseProto.AddPurchaseRequest{
		CartId:        uuid.New().String(),
		Address:       "123 Test St",
		PaymentMethod: cartPurchaseProto.PaymentMethod(cartPurchaseProto.PaymentMethod_value["CREDIT_CARD"]),
		DeliveryMethod: cartPurchaseProto.DeliveryMethod(cartPurchaseProto.DeliveryMethod_value["STANDARD"]),
		UserId:       uuid.New().String(),
	}

	resp := &cartPurchaseProto.AddPurchaseResponse{
		Id:             uuid.New().String(),
		CartId:         req.CartId,
		Address:        req.Address,
		Status:         cartPurchaseProto.PurchaseStatus(cartPurchaseProto.PurchaseStatus_value["COMPLETED"]),
		PaymentMethod:  cartPurchaseProto.PaymentMethod(cartPurchaseProto.PaymentMethod_value["CREDIT_CARD"]),
		DeliveryMethod: cartPurchaseProto.DeliveryMethod(cartPurchaseProto.DeliveryMethod_value["STANDARD"]),
	}

	mockServer.On("AddPurchase", mock.Anything, req).Return(resp, nil)

	result, err := mockServer.AddPurchase(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, resp.Id, result.Id)
	mockServer.AssertExpectations(t)
}

func TestServerGetPurchasesByUserID(t *testing.T) {
	mockServer := new(MockCartPurchaseServiceServer)

	userID := uuid.New()
	req := &cartPurchaseProto.GetPurchasesByUserIDRequest{
		UserId: userID.String(),
	}

	resp := &cartPurchaseProto.GetPurchasesByUserIDResponse{
		Purchases: []*cartPurchaseProto.PurchaseResponse{
			{
				Id:             uuid.New().String(),
				CartId:         uuid.New().String(),
				Address:        "123 Test St",
				Status:         cartPurchaseProto.PurchaseStatus(cartPurchaseProto.PurchaseStatus_value["COMPLETED"]),
				PaymentMethod:  cartPurchaseProto.PaymentMethod(cartPurchaseProto.PaymentMethod_value["CREDIT_CARD"]),
				DeliveryMethod: cartPurchaseProto.DeliveryMethod(cartPurchaseProto.DeliveryMethod_value["STANDARD"]),
			},
		},
	}

	mockServer.On("GetPurchasesByUserID", mock.Anything, req).Return(resp, nil)

	result, err := mockServer.GetPurchasesByUserID(context.Background(), req)

	assert.NoError(t, err)
	assert.Len(t, result.Purchases, 1)
	mockServer.AssertExpectations(t)
}

func TestServerGetCartByID(t *testing.T) {
	mockServer := new(MockCartPurchaseServiceServer)

	cartID := uuid.New()
	req := &cartPurchaseProto.GetCartByIDRequest{
		CartId: cartID.String(),
	}

	resp := &cartPurchaseProto.GetCartByIDResponse{
		Cart: &cartPurchaseProto.Cart{
			Id:      cartID.String(),
			UserId:  uuid.New().String(),
			Status:  cartPurchaseProto.CartStatus(cartPurchaseProto.CartStatus_value["ACTIVE"]),
			Adverts: nil,
		},
	}

	mockServer.On("GetCartByID", mock.Anything, req).Return(resp, nil)

	result, err := mockServer.GetCartByID(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, cartID.String(), result.Cart.Id)
	mockServer.AssertExpectations(t)
}

func TestServerGetCartByUserID(t *testing.T) {
	mockServer := new(MockCartPurchaseServiceServer)

	userID := uuid.New()
	req := &cartPurchaseProto.GetCartByUserIDRequest{
		UserId: userID.String(),
	}

	resp := &cartPurchaseProto.GetCartByUserIDResponse{
		Cart: &cartPurchaseProto.Cart{
			Id:      uuid.New().String(),
			UserId:  userID.String(),
			Status:  cartPurchaseProto.CartStatus(cartPurchaseProto.CartStatus_value["ACTIVE"]),
			Adverts: nil,
		},
	}

	mockServer.On("GetCartByUserID", mock.Anything, req).Return(resp, nil)

	result, err := mockServer.GetCartByUserID(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, userID.String(), result.Cart.UserId)
	mockServer.AssertExpectations(t)
}

func TestServerAddAdvertToCart(t *testing.T) {
	mockServer := new(MockCartPurchaseServiceServer)

	req := &cartPurchaseProto.AddAdvertToCartRequest{
		UserId:   uuid.New().String(),
		AdvertId: uuid.New().String(),
	}

	resp := &cartPurchaseProto.AddAdvertToCartResponse{
		Message: "Advert added successfully",
	}

	mockServer.On("AddAdvertToCart", mock.Anything, req).Return(resp, nil)

	result, err := mockServer.AddAdvertToCart(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "Advert added successfully", result.Message)
	mockServer.AssertExpectations(t)
}

func TestServerDeleteAdvertFromCart(t *testing.T) {
	mockServer := new(MockCartPurchaseServiceServer)

	req := &cartPurchaseProto.DeleteAdvertFromCartRequest{
		CartId:   uuid.New().String(),
		AdvertId: uuid.New().String(),
	}

	resp := &cartPurchaseProto.DeleteAdvertFromCartResponse{
		Message: "Advert deleted successfully",
	}

	mockServer.On("DeleteAdvertFromCart", mock.Anything, req).Return(resp, nil)

	result, err := mockServer.DeleteAdvertFromCart(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "Advert deleted successfully", result.Message)
	mockServer.AssertExpectations(t)
}

func TestServerAddPurchase_Error(t *testing.T) {
	mockServer := new(MockCartPurchaseServiceServer)

	req := &cartPurchaseProto.AddPurchaseRequest{
		CartId:        uuid.New().String(),
		Address:       "123 Test St",
		PaymentMethod: cartPurchaseProto.PaymentMethod(cartPurchaseProto.PaymentMethod_value["CREDIT_CARD"]),
		DeliveryMethod: cartPurchaseProto.DeliveryMethod(cartPurchaseProto.DeliveryMethod_value["STANDARD"]),
		UserId:       uuid.New().String(),
	}

	mockServer.On("AddPurchase", mock.Anything, req).Return(&cartPurchaseProto.AddPurchaseResponse{}, assert.AnError)

	result, err := mockServer.AddPurchase(context.Background(), req)

	assert.Error(t, err)
	assert.NotNil(t, result)
	mockServer.AssertExpectations(t)
}

func TestServerGetPurchasesByUserID_Error(t *testing.T) {
	mockServer := new(MockCartPurchaseServiceServer)

	userID := uuid.New()
	req := &cartPurchaseProto.GetPurchasesByUserIDRequest{
		UserId: userID.String(),
	}

	mockServer.On("GetPurchasesByUserID", mock.Anything, req).Return(&cartPurchaseProto.GetPurchasesByUserIDResponse{}, assert.AnError)

	result, err := mockServer.GetPurchasesByUserID(context.Background(), req)

	assert.Error(t, err)
	assert.NotNil(t, result)
	mockServer.AssertExpectations(t)
}

func TestServerGetCartByID_Error(t *testing.T) {
	mockServer := new(MockCartPurchaseServiceServer)

	cartID := uuid.New()
	req := &cartPurchaseProto.GetCartByIDRequest{
		CartId: cartID.String(),
	}

	mockServer.On("GetCartByID", mock.Anything, req).Return(&cartPurchaseProto.GetCartByIDResponse{}, assert.AnError)

	result, err := mockServer.GetCartByID(context.Background(), req)

	assert.Error(t, err)
	assert.NotNil(t, result)
	mockServer.AssertExpectations(t)
}

func TestServerCheckCartExists(t *testing.T) {
	mockServer := new(MockCartPurchaseServiceServer)

	userID := uuid.New()
	req := &cartPurchaseProto.CheckCartExistsRequest{
		UserId: userID.String(),
	}

	resp := &cartPurchaseProto.CheckCartExistsResponse{
		CartId: uuid.New().String(),
	}

	mockServer.On("CheckCartExists", mock.Anything, req).Return(resp, nil)

	result, err := mockServer.CheckCartExists(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil.String(), result.CartId)
	mockServer.AssertExpectations(t)
}
