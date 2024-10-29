package service

import (
	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/google/uuid"
)

type AuthService struct {
	sessionRepo repository.Session
}

func NewAuthService(authRepo repository.Session) *AuthService {
	return &AuthService{sessionRepo: authRepo}
}

func (a *AuthService) Logout(session string) error {
	err := a.sessionRepo.DeleteSession(session)
	switch {
	case errors.Is(err, repository.ErrSessionNotFound):
		return usecase.ErrUserNotFound
	case err != nil:
		return entity.UsecaseWrap(errors.New("не удалось удалить сессию"), err)
	}
	return nil
}

func (a *AuthService) CreateSession(userId uuid.UUID) (string, error) {
	session, err := a.sessionRepo.CreateSession(userId)
	if err != nil {
		return "", entity.UsecaseWrap(errors.New("не удалось создать сессию"), err)
	}
	return session, nil
}

func (a *AuthService) GetUserIdBySession(session string) (uuid.UUID, error) {
	userID, err := a.sessionRepo.GetSession(session)
	switch {
	case errors.Is(err, repository.ErrSessionNotFound):
		return uuid.Nil, usecase.ErrUserNotFound
	case err != nil:
		return uuid.Nil, entity.UsecaseWrap(errors.New("не удалось удалить сессию"), err)
	}

	return userID, nil
}
