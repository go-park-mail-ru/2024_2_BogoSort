package service

import (
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupAuthService(t *testing.T) (*AuthService, *mocks.MockSession, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	sessionRepo := mocks.NewMockSession(ctrl)
	service := NewAuthService(sessionRepo)
	return service, sessionRepo, ctrl
}

func TestAuthService_Logout(t *testing.T) {
	service, sessionRepo, ctrl := setupAuthService(t)
	defer ctrl.Finish()

	session := "valid-session-id"

	gomock.InOrder(
		sessionRepo.EXPECT().Delete(session).Return(nil),
	)

	err := service.Logout(session)

	assert.NoError(t, err)
}

func TestAuthService_Logout_SessionNotFound(t *testing.T) {
	service, sessionRepo, ctrl := setupAuthService(t)
	defer ctrl.Finish()

	session := "invalid-session-id"

	sessionRepo.EXPECT().Delete(session).Return(repository.ErrSessionNotFound)

	err := service.Logout(session)
	assert.ErrorIs(t, err, usecase.ErrUserNotFound)
}

func TestAuthService_CreateSession(t *testing.T) {
	service, sessionRepo, ctrl := setupAuthService(t)
	defer ctrl.Finish()

	userId := uuid.New()
	expectedSession := "new-session-id"

	sessionRepo.EXPECT().Create(userId).Return(expectedSession, nil)

	session, err := service.CreateSession(userId)
	assert.NoError(t, err)
	assert.Equal(t, expectedSession, session)
}

func TestAuthService_GetUserIdBySession(t *testing.T) {
	service, sessionRepo, ctrl := setupAuthService(t)
	defer ctrl.Finish()

	session := "valid-session-id"
	expectedUserId := uuid.New()

	sessionRepo.EXPECT().Get(session).Return(expectedUserId, nil)

	userId, err := service.GetUserIdBySession(session)
	assert.NoError(t, err)
	assert.Equal(t, expectedUserId, userId)
}

func TestAuthService_GetUserIdBySession_SessionNotFound(t *testing.T) {
	service, sessionRepo, ctrl := setupAuthService(t)
	defer ctrl.Finish()

	session := "invalid-session-id"

	sessionRepo.EXPECT().Get(session).Return(uuid.Nil, repository.ErrSessionNotFound)

	userId, err := service.GetUserIdBySession(session)
	assert.ErrorIs(t, err, usecase.ErrUserNotFound)
	assert.Equal(t, uuid.Nil, userId)
}

func TestAuthService_CreateSession_FailureInRepository(t *testing.T) {
	service, sessionRepo, ctrl := setupAuthService(t)
	defer ctrl.Finish()

	userId := uuid.New()
	expectedError := errors.New("repository error")
	sessionRepo.EXPECT().Create(userId).Return("", expectedError)

	session, err := service.CreateSession(userId)

	assert.Error(t, err)
	assert.Equal(t, "", session)
}

func TestAuthService_GetUserIdBySession_RepositoryError(t *testing.T) {
	service, sessionRepo, ctrl := setupAuthService(t)
	defer ctrl.Finish()

	session := "some-session-id"
	expectedError := errors.New("database error")

	sessionRepo.EXPECT().Get(session).Return(uuid.Nil, expectedError)

	userId, err := service.GetUserIdBySession(session)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, userId)
}

func TestAuthService_Logout_FailureInRepository(t *testing.T) {
	service, sessionRepo, ctrl := setupAuthService(t)
	defer ctrl.Finish()

	session := "some-session-id"
	expectedError := errors.New("delete session error")

	sessionRepo.EXPECT().Delete(session).Return(expectedError)

	err := service.Logout(session)

	assert.Error(t, err)
}
