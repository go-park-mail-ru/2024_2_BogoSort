package cart_purchase

// import (
// 	"context"
// 	"testing"
// 	"errors"

// 	cartPurchaseProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/cart_purchase/proto"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"google.golang.org/grpc"
// )

// type MockCartPurchaseServiceClient struct {
// 	mock.Mock
// }

// func (m *MockCartPurchaseServiceClient) Ping(ctx context.Context, in *cartPurchaseProto.NoContent, opts ...grpc.CallOption) (*cartPurchaseProto.NoContent, error) {
// 	args := m.Called(ctx, in)
// 	return args.Get(0).(*cartPurchaseProto.NoContent), args.Error(1)
// }

// func (m *MockCartPurchaseServiceClient) AddPurchase(ctx context.Context, in *cartPurchaseProto.AddPurchaseRequest, opts ...grpc.CallOption) (*cartPurchaseProto.AddPurchaseResponse, error) {
// 	args := m.Called(ctx, in)
// 	return args.Get(0).(*cartPurchaseProto.AddPurchaseResponse), args.Error(1)
// }

// func (m *MockCartPurchaseServiceClient) GetPurchasesByUserID(ctx context.Context, in *cartPurchaseProto.GetPurchasesByUserIDRequest, opts ...grpc.CallOption) (*cartPurchaseProto.GetPurchasesByUserIDResponse, error) {
// 	args := m.Called(ctx, in)
// 	return args.Get(0).(*cartPurchaseProto.GetPurchasesByUserIDResponse), args.Error(1)
// }

// func (m *MockCartPurchaseServiceClient) GetCartByID(ctx context.Context, in *cartPurchaseProto.GetCartByIDRequest, opts ...grpc.CallOption) (*cartPurchaseProto.GetCartByIDResponse, error) {
// 	args := m.Called(ctx, in)
// 	return args.Get(0).(*cartPurchaseProto.GetCartByIDResponse), args.Error(1)
// }

// func (m *MockCartPurchaseServiceClient) GetCartByUserID(ctx context.Context, in *cartPurchaseProto.GetCartByUserIDRequest, opts ...grpc.CallOption) (*cartPurchaseProto.GetCartByUserIDResponse, error) {
// 	args := m.Called(ctx, in)
// 	return args.Get(0).(*cartPurchaseProto.GetCartByUserIDResponse), args.Error(1)
// }

// func (m *MockCartPurchaseServiceClient) AddAdvertToCart(ctx context.Context, in *cartPurchaseProto.AddAdvertToCartRequest, opts ...grpc.CallOption) (*cartPurchaseProto.AddAdvertToCartResponse, error) {
// 	args := m.Called(ctx, in)
// 	return args.Get(0).(*cartPurchaseProto.AddAdvertToCartResponse), args.Error(1)
// }

// func (m *MockCartPurchaseServiceClient) DeleteAdvertFromCart(ctx context.Context, in *cartPurchaseProto.DeleteAdvertFromCartRequest, opts ...grpc.CallOption) (*cartPurchaseProto.DeleteAdvertFromCartResponse, error) {
// 	args := m.Called(ctx, in)
// 	return args.Get(0).(*cartPurchaseProto.DeleteAdvertFromCartResponse), args.Error(1)
// }

// func (m *MockCartPurchaseServiceClient) CheckCartExists(ctx context.Context, in *cartPurchaseProto.CheckCartExistsRequest, opts ...grpc.CallOption) (*cartPurchaseProto.CheckCartExistsResponse, error) {
// 	args := m.Called(ctx, in)
// 	return args.Get(0).(*cartPurchaseProto.CheckCartExistsResponse), args.Error(1)
// }

// func TestNewCartPurchaseClient(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	mockConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
// 	if err != nil {
// 		t.Fatalf("failed to dial: %v", err)
// 	}

// 	client := &CartPurchaseClient{
// 		client: mockClient,
// 		conn:   mockConn,
// 	}

// 	assert.NotNil(t, client)
// 	mockClient.AssertExpectations(t)
// }

// func TestAddPurchase(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	client := &CartPurchaseClient{client: mockClient}

// 	req := dto.PurchaseRequest{
// 		CartID:        uuid.New(),
// 		Address:       "123 Test St",
// 		PaymentMethod: dto.PaymentMethod("CREDIT_CARD"),
// 		DeliveryMethod: dto.DeliveryMethod("STANDARD"),
// 		UserID:       uuid.New(),
// 	}

// 	protoResp := &cartPurchaseProto.AddPurchaseResponse{
// 		Id:             uuid.New().String(),
// 		CartId:         req.CartID.String(),
// 		Address:        req.Address,
// 		Status:         cartPurchaseProto.PurchaseStatus(cartPurchaseProto.PurchaseStatus_value["COMPLETED"]),
// 		PaymentMethod:  cartPurchaseProto.PaymentMethod(cartPurchaseProto.PaymentMethod_value["CREDIT_CARD"]),
// 		DeliveryMethod: cartPurchaseProto.DeliveryMethod(cartPurchaseProto.DeliveryMethod_value["STANDARD"]),
// 	}

// 	mockClient.On("AddPurchase", mock.Anything, mock.Anything).Return(protoResp, nil)

// 	resp, err := client.AddPurchase(context.Background(), req)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.Equal(t, req.CartID, resp.CartID)
// 	mockClient.AssertExpectations(t)
// }

// func TestAddPurchase_Error(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	client := &CartPurchaseClient{client: mockClient}

// 	req := dto.PurchaseRequest{
// 		CartID:        uuid.New(),
// 		Address:       "123 Test St",
// 		PaymentMethod: dto.PaymentMethod("CREDIT_CARD"),
// 		DeliveryMethod: dto.DeliveryMethod("STANDARD"),
// 		UserID:       uuid.New(),
// 	}

// 	mockClient.On("AddPurchase", mock.Anything, mock.Anything).Return(&cartPurchaseProto.AddPurchaseResponse{}, errors.New("some error"))

// 	resp, err := client.AddPurchase(context.Background(), req)

// 	assert.Error(t, err)
// 	assert.Nil(t, resp)
// 	mockClient.AssertExpectations(t)
// }

// func TestGetPurchasesByUserID(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	client := &CartPurchaseClient{client: mockClient}

// 	userID := uuid.New()
// 	protoResp := &cartPurchaseProto.GetPurchasesByUserIDResponse{
// 		Purchases: []*cartPurchaseProto.PurchaseResponse{
// 			{
// 				Id:             uuid.New().String(),
// 				CartId:         uuid.New().String(),
// 				Address:        "123 Test St",
// 				Status:         cartPurchaseProto.PurchaseStatus(cartPurchaseProto.PurchaseStatus_value["COMPLETED"]),
// 				PaymentMethod:  cartPurchaseProto.PaymentMethod(cartPurchaseProto.PaymentMethod_value["CREDIT_CARD"]),
// 				DeliveryMethod: cartPurchaseProto.DeliveryMethod(cartPurchaseProto.DeliveryMethod_value["STANDARD"]),
// 			},
// 		},
// 	}

// 	mockClient.On("GetPurchasesByUserID", mock.Anything, mock.Anything).Return(protoResp, nil)

// 	resp, err := client.GetPurchasesByUserID(context.Background(), userID)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.Len(t, resp, 1)
// 	mockClient.AssertExpectations(t)
// }

// func TestGetPurchasesByUserID_Error(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	client := &CartPurchaseClient{client: mockClient}

// 	userID := uuid.New()
// 	mockClient.On("GetPurchasesByUserID", mock.Anything, mock.Anything).Return(&cartPurchaseProto.GetPurchasesByUserIDResponse{}, errors.New("some error"))

// 	resp, err := client.GetPurchasesByUserID(context.Background(), userID)

// 	assert.Error(t, err)
// 	assert.Nil(t, resp)
// 	mockClient.AssertExpectations(t)
// }

// func TestGetCartByID(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	client := &CartPurchaseClient{client: mockClient}

// 	cartID := uuid.New()
// 	protoResp := &cartPurchaseProto.GetCartByIDResponse{
// 		Cart: &cartPurchaseProto.Cart{
// 			Id:      cartID.String(),
// 			UserId:  uuid.New().String(),
// 			Status:  cartPurchaseProto.CartStatus(cartPurchaseProto.CartStatus_value["ACTIVE"]),
// 			Adverts: nil,
// 		},
// 	}

// 	mockClient.On("GetCartByID", mock.Anything, mock.Anything).Return(protoResp, nil)

// 	resp, err := client.GetCartByID(context.Background(), cartID)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.Equal(t, cartID, resp.ID)
// 	mockClient.AssertExpectations(t)
// }

// func TestGetCartByUserID(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	client := &CartPurchaseClient{client: mockClient}

// 	userID := uuid.New()
// 	protoResp := &cartPurchaseProto.GetCartByUserIDResponse{
// 		Cart: &cartPurchaseProto.Cart{
// 			Id:      uuid.New().String(),
// 			UserId:  userID.String(),
// 			Status:  cartPurchaseProto.CartStatus(cartPurchaseProto.CartStatus_value["ACTIVE"]),
// 			Adverts: nil,
// 		},
// 	}

// 	mockClient.On("GetCartByUserID", mock.Anything, mock.Anything).Return(protoResp, nil)

// 	resp, err := client.GetCartByUserID(context.Background(), userID)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	mockClient.AssertExpectations(t)
// }

// func TestGetCartByUserID_Error(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	client := &CartPurchaseClient{client: mockClient}

// 	userID := uuid.New()
// 	mockClient.On("GetCartByUserID", mock.Anything, mock.Anything).Return(&cartPurchaseProto.GetCartByUserIDResponse{}, errors.New("some error"))

// 	resp, err := client.GetCartByUserID(context.Background(), userID)

// 	assert.Error(t, err)
// 	assert.Nil(t, resp)
// 	mockClient.AssertExpectations(t)
// }

// func TestAddAdvertToCart(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	client := &CartPurchaseClient{client: mockClient}

// 	userID := uuid.New()
// 	advertID := uuid.New()
// 	protoResp := &cartPurchaseProto.AddAdvertToCartResponse{
// 		Message: "Advert added successfully",
// 	}

// 	mockClient.On("AddAdvertToCart", mock.Anything, mock.Anything).Return(protoResp, nil)

// 	message, err := client.AddAdvertToCart(context.Background(), userID, advertID)

// 	assert.NoError(t, err)
// 	assert.Equal(t, "Advert added successfully", message)
// 	mockClient.AssertExpectations(t)
// }

// func TestAddAdvertToCart_Error(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	client := &CartPurchaseClient{client: mockClient}

// 	userID := uuid.New()
// 	advertID := uuid.New()
// 	mockClient.On("AddAdvertToCart", mock.Anything, mock.Anything).Return(&cartPurchaseProto.AddAdvertToCartResponse{}, errors.New("some error"))

// 	message, err := client.AddAdvertToCart(context.Background(), userID, advertID)

// 	assert.Error(t, err)
// 	assert.Empty(t, message)
// 	mockClient.AssertExpectations(t)
// }

// func TestDeleteAdvertFromCart(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	client := &CartPurchaseClient{client: mockClient}

// 	cartID := uuid.New()
// 	advertID := uuid.New()
// 	protoResp := &cartPurchaseProto.DeleteAdvertFromCartResponse{
// 		Message: "Advert deleted successfully",
// 	}

// 	mockClient.On("DeleteAdvertFromCart", mock.Anything, mock.Anything).Return(protoResp, nil)

// 	message, err := client.DeleteAdvertFromCart(context.Background(), cartID, advertID)

// 	assert.NoError(t, err)
// 	assert.Equal(t, "Advert deleted successfully", message)
// 	mockClient.AssertExpectations(t)
// }

// func TestCheckCartExists(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	client := &CartPurchaseClient{client: mockClient}

// 	userID := uuid.New()
// 	protoResp := &cartPurchaseProto.CheckCartExistsResponse{
// 		CartId: uuid.New().String(),
// 	}

// 	mockClient.On("CheckCartExists", mock.Anything, mock.Anything).Return(protoResp, nil)

// 	cartID, err := client.CheckCartExists(context.Background(), userID)

// 	assert.NoError(t, err)
// 	assert.NotEqual(t, uuid.Nil, cartID)
// 	mockClient.AssertExpectations(t)
// }

// func TestPing(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	client := &CartPurchaseClient{client: mockClient}

// 	mockClient.On("Ping", mock.Anything, mock.Anything).Return(&cartPurchaseProto.NoContent{}, nil)

// 	err := client.Ping(context.Background())

// 	assert.NoError(t, err)
// 	mockClient.AssertExpectations(t)
// }

// func TestPing_Error(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	client := &CartPurchaseClient{client: mockClient}

// 	mockClient.On("Ping", mock.Anything, mock.Anything).Return(&cartPurchaseProto.NoContent{}, errors.New("some error"))

// 	err := client.Ping(context.Background())

// 	assert.Error(t, err)
// 	mockClient.AssertExpectations(t)
// }

// func TestAddAdvertToCart_EmptyResponse(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	client := &CartPurchaseClient{client: mockClient}

// 	userID := uuid.New()
// 	advertID := uuid.New()
// 	mockClient.On("AddAdvertToCart", mock.Anything, mock.Anything).Return(&cartPurchaseProto.AddAdvertToCartResponse{}, nil)

// 	message, err := client.AddAdvertToCart(context.Background(), userID, advertID)

// 	assert.NoError(t, err)
// 	assert.Empty(t, message)
// 	mockClient.AssertExpectations(t)
// }

// func TestDeleteAdvertFromCart_EmptyResponse(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	client := &CartPurchaseClient{client: mockClient}

// 	cartID := uuid.New()
// 	advertID := uuid.New()
// 	mockClient.On("DeleteAdvertFromCart", mock.Anything, mock.Anything).Return(&cartPurchaseProto.DeleteAdvertFromCartResponse{}, nil)

// 	message, err := client.DeleteAdvertFromCart(context.Background(), cartID, advertID)

// 	assert.NoError(t, err)
// 	assert.Empty(t, message)
// 	mockClient.AssertExpectations(t)
// }

// func TestCheckCartExists_NotFound(t *testing.T) {
// 	mockClient := new(MockCartPurchaseServiceClient)
// 	client := &CartPurchaseClient{client: mockClient}

// 	userID := uuid.New()
// 	mockClient.On("CheckCartExists", mock.Anything, mock.Anything).Return(&cartPurchaseProto.CheckCartExistsResponse{}, errors.New("cart not found"))

// 	cartID, err := client.CheckCartExists(context.Background(), userID)

// 	assert.Error(t, err)
// 	assert.Equal(t, uuid.Nil, cartID)
// 	mockClient.AssertExpectations(t)
// }
