package service

import (
	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthService struct {
	sessionRepo repository.Session
	logger      *zap.Logger
}

func NewAuthService(authRepo repository.Session, logger *zap.Logger) *AuthService {
	return &AuthService{sessionRepo: authRepo, logger: logger}
}

func (a *AuthService) Logout(session string) error {
	err := a.sessionRepo.DeleteSession(session)
	switch {
	case errors.Is(err, repository.ErrSessionNotFound):
		return usecase.ErrUserNotFound
	case err != nil:
		a.logger.Error("error deleting session", zap.String("sessionID", session), zap.Error(err))
		return entity.UsecaseWrap(errors.New("error deleting session"), err)
	}
	a.logger.Info("session deleted", zap.String("sessionID", session))
	return nil
}

func (a *AuthService) CreateSession(userId uuid.UUID) (string, error) {
	session, err := a.sessionRepo.CreateSession(userId)
	if err != nil {
		a.logger.Error("error creating session", zap.String("userID", userId.String()), zap.Error(err))
		return "", entity.UsecaseWrap(errors.New("error creating session"), err)
	}
	a.logger.Info("session created", zap.String("sessionID", session), zap.String("userID", userId.String()))
	return session, nil
}

func (a *AuthService) GetUserIdBySession(session string) (uuid.UUID, error) {
	userID, err := a.sessionRepo.GetSession(session)
	switch {
	case errors.Is(err, repository.ErrSessionNotFound):
		a.logger.Error("session not found", zap.String("sessionID", session))
		return uuid.Nil, usecase.ErrUserNotFound
	case err != nil:
		a.logger.Error("error getting userID by session", zap.String("sessionID", session), zap.Error(err))
		return uuid.Nil, entity.UsecaseWrap(errors.New("error getting userID by session"), err)
	}
	a.logger.Info("userID found by session", zap.String("sessionID", session), zap.String("userID", userID.String()))
	return userID, nil
}
