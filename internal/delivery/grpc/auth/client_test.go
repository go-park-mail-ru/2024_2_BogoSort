package auth

import (
	"context"
	"errors"
	"testing"

	authProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/auth/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockAuthServiceClient struct {
	mock.Mock
}

func (m *MockAuthServiceClient) Ping(ctx context.Context, in *authProto.NoContent, opts ...grpc.CallOption) (*authProto.NoContent, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*authProto.NoContent), args.Error(1)
}

func (m *MockAuthServiceClient) GetUserIDBySession(ctx context.Context, in *authProto.Session, opts ...grpc.CallOption) (*authProto.User, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*authProto.User), args.Error(1)
}

func (m *MockAuthServiceClient) CreateSession(ctx context.Context, in *authProto.User, opts ...grpc.CallOption) (*authProto.Session, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*authProto.Session), args.Error(1)
}

func (m *MockAuthServiceClient) DeleteSession(ctx context.Context, in *authProto.Session, opts ...grpc.CallOption) (*authProto.NoContent, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*authProto.NoContent), args.Error(1)
}

func TestNewGrpcClient(t *testing.T) {
	mockConn := new(MockAuthServiceClient)
	mockConn.On("Ping", mock.Anything, mock.Anything).Return(&authProto.NoContent{}, nil)

	client := &GrpcClient{authManager: mockConn}

	_, err := client.authManager.Ping(context.Background(), &authProto.NoContent{})
	assert.NoError(t, err)

	mockConn.AssertExpectations(t)
}

func TestGetUserIDBySession(t *testing.T) {
	mockAuthClient := new(MockAuthServiceClient)
	client := &GrpcClient{authManager: mockAuthClient}

	sessionID := "test-session-id"
	userID := "test-user-id"

	mockAuthClient.On("GetUserIDBySession", mock.Anything, &authProto.Session{Id: sessionID}).Return(&authProto.User{Id: userID}, nil)

	result, err := client.GetUserIDBySession(sessionID)
	assert.NoError(t, err)
	assert.Equal(t, userID, result)

	mockAuthClient.AssertExpectations(t)
}

func TestGetUserIDBySession_Error(t *testing.T) {
	mockAuthClient := new(MockAuthServiceClient)
	client := &GrpcClient{authManager: mockAuthClient}

	sessionID := "test-session-id"

	mockAuthClient.On("GetUserIDBySession", mock.Anything, &authProto.Session{Id: sessionID}).Return(&authProto.User{}, errors.New("session not found"))

	result, err := client.GetUserIDBySession(sessionID)
	assert.Error(t, err)
	assert.Empty(t, result)

	mockAuthClient.AssertExpectations(t)
}

func TestCreateSession(t *testing.T) {
	mockAuthClient := new(MockAuthServiceClient)
	client := &GrpcClient{authManager: mockAuthClient}

	userID := uuid.New()
	sessionID := "test-session-id"

	mockAuthClient.On("CreateSession", mock.Anything, &authProto.User{Id: userID.String()}).Return(&authProto.Session{Id: sessionID}, nil)

	result, err := client.CreateSession(userID)
	assert.NoError(t, err)
	assert.Equal(t, sessionID, result)

	mockAuthClient.AssertExpectations(t)
}

func TestCreateSession_Error(t *testing.T) {
	mockAuthClient := new(MockAuthServiceClient)
	client := &GrpcClient{authManager: mockAuthClient}

	userID := uuid.New()

	mockAuthClient.On("CreateSession", mock.Anything, &authProto.User{Id: userID.String()}).Return(&authProto.Session{}, errors.New("failed to create session"))

	result, err := client.CreateSession(userID)
	assert.Error(t, err)
	assert.Empty(t, result)

	mockAuthClient.AssertExpectations(t)
}

func TestDeleteSession(t *testing.T) {
	mockAuthClient := new(MockAuthServiceClient)
	client := &GrpcClient{authManager: mockAuthClient}

	sessionID := "test-session-id"

	mockAuthClient.On("DeleteSession", mock.Anything, &authProto.Session{Id: sessionID}).Return(&authProto.NoContent{}, nil)

	err := client.DeleteSession(sessionID)
	assert.NoError(t, err)

	mockAuthClient.AssertExpectations(t)
}

func TestDeleteSession_Error(t *testing.T) {
	mockAuthClient := new(MockAuthServiceClient)
	client := &GrpcClient{authManager: mockAuthClient}

	sessionID := "test-session-id"

	mockAuthClient.On("DeleteSession", mock.Anything, &authProto.Session{Id: sessionID}).Return(&authProto.NoContent{}, errors.New("failed to delete session"))

	err := client.DeleteSession(sessionID)
	assert.Error(t, err)

	mockAuthClient.AssertExpectations(t)
}