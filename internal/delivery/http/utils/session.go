package utils

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SessionManager struct {
	sessionUC        usecase.Auth
	sessionAliveTime int
	secureCookie     bool
	logger           *zap.Logger
}

func NewSessionManager(authUC usecase.Auth, sessionAliveTime int, secureCookie bool, logger *zap.Logger) *SessionManager {
	return &SessionManager{
		sessionUC:        authUC,
		sessionAliveTime: sessionAliveTime,
		secureCookie:     secureCookie,
		logger:           logger,
	}
}

func (s *SessionManager) CreateSession(userID uuid.UUID) (string, error) {
	session, err := s.sessionUC.CreateSession(userID)
	if err != nil {
		s.logger.Error("error creating session", zap.String("userID", userID.String()), zap.Error(err))
		return "", err
	}

	return session, nil
}

func (s *SessionManager) SetSession(value string, expires time.Time) (*http.Cookie, error) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    value,
		Expires:  expires,
		HttpOnly: true,
		Secure:   s.secureCookie,
	}

	return cookie, nil
}

func (s *SessionManager) GetUserID(r *http.Request) (uuid.UUID, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		s.logger.Error("error getting cookie", zap.Error(err))
		return uuid.Nil, err
	}

	return s.sessionUC.GetUserIdBySession(cookie.Value)
}

func (s *SessionManager) DeleteSession(sessionID string) error {
	err := s.sessionUC.Logout(sessionID)
	if err != nil {
		s.logger.Error("error logging out", zap.String("sessionID", sessionID), zap.Error(err))
		return err
	}
	return nil
}
