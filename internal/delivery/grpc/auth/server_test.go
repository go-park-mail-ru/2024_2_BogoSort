package auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	proto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/auth/proto"
)

type MockAuth struct {
	mock.Mock
}

func (m *MockAuth) GetUserIdBySession(sessionID string) (uuid.UUID, error) {
	args := m.Called(sessionID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockAuth) CreateSession(userID uuid.UUID) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockAuth) Logout(sessionID string) error {
	args := m.Called(sessionID)
	return args.Error(0)
}

func TestServerGetUserIDBySession(t *testing.T) {
	mockAuth := new(MockAuth)
	server := NewGrpcServer(mockAuth)

	sessionID := "test-session-id"
	userID := uuid.New()
	mockAuth.On("GetUserIdBySession", sessionID).Return(userID, nil)

	session := &proto.Session{Id: sessionID}
	user, err := server.GetUserIDBySession(context.Background(), session)

	assert.NoError(t, err)
	assert.Equal(t, userID.String(), user.Id)
	mockAuth.AssertExpectations(t)
}

func TestServerGetUserIDBySession_Error(t *testing.T) {
	mockAuth := new(MockAuth)
	server := NewGrpcServer(mockAuth)

	sessionID := "test-session-id"
	mockAuth.On("GetUserIdBySession", sessionID).Return(uuid.Nil, assert.AnError)

	session := &proto.Session{Id: sessionID}
	user, err := server.GetUserIDBySession(context.Background(), session)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockAuth.AssertExpectations(t)
}

func TestServerCreateSession(t *testing.T) {
	mockAuth := new(MockAuth)
	server := NewGrpcServer(mockAuth)

	userID := uuid.New()
	sessionID := "test-session-id"
	mockAuth.On("CreateSession", userID).Return(sessionID, nil)

	user := &proto.User{Id: userID.String()}
	session, err := server.CreateSession(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, sessionID, session.Id)
	mockAuth.AssertExpectations(t)
}

func TestServerCreateSession_Error(t *testing.T) {
	mockAuth := new(MockAuth)
	server := NewGrpcServer(mockAuth)

	userID := uuid.New()
	mockAuth.On("CreateSession", userID).Return("", assert.AnError)

	user := &proto.User{Id: userID.String()}
	session, err := server.CreateSession(context.Background(), user)

	assert.Error(t, err)
	assert.Nil(t, session)
	mockAuth.AssertExpectations(t)
}

func TestServerDeleteSession(t *testing.T) {
	mockAuth := new(MockAuth)
	server := NewGrpcServer(mockAuth)

	sessionID := "test-session-id"
	mockAuth.On("Logout", sessionID).Return(nil)

	session := &proto.Session{Id: sessionID}
	_, err := server.DeleteSession(context.Background(), session)

	assert.NoError(t, err)
	mockAuth.AssertExpectations(t)
}

func TestServerDeleteSession_Error(t *testing.T) {
	mockAuth := new(MockAuth)
	server := NewGrpcServer(mockAuth)

	sessionID := "test-session-id"
	mockAuth.On("Logout", sessionID).Return(assert.AnError)

	session := &proto.Session{Id: sessionID}
	_, err := server.DeleteSession(context.Background(), session)

	assert.Error(t, err)
	mockAuth.AssertExpectations(t)
}

func TestServerPing(t *testing.T) {
	mockAuth := new(MockAuth)
	server := NewGrpcServer(mockAuth)

	_, err := server.Ping(context.Background(), &proto.NoContent{})

	assert.NoError(t, err)
}
