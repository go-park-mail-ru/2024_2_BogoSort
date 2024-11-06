package utils

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SessionManager struct {
	SessionUC        usecase.Auth
	SessionAliveTime int
	SecureCookie     bool
	Logger           *zap.Logger
}

func NewSessionManager(authUC usecase.Auth, sessionAliveTime int, secureCookie bool, logger *zap.Logger) *SessionManager {
	return &SessionManager{
		SessionUC:        authUC,
		SessionAliveTime: sessionAliveTime,
		SecureCookie:     secureCookie,
		Logger:           logger,
	}
}

func (s *SessionManager) CreateSession(userID uuid.UUID) (string, error) {
	return s.SessionUC.CreateSession(userID)
}

func (s *SessionManager) SetSession(value string) (*http.Cookie, error) {
	expires := time.Now().Add(time.Duration(s.SessionAliveTime) * time.Second)
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    value,
		Expires:  expires,
		HttpOnly: true,
		Secure:   s.SecureCookie,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	return cookie, nil
}

func (s *SessionManager) GetUserID(r *http.Request) (uuid.UUID, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		s.Logger.Error("error getting cookie", zap.Error(err))
		return uuid.Nil, err
	}

	return s.SessionUC.GetUserIdBySession(cookie.Value)
}

func (s *SessionManager) DeleteSession(sessionID string) error {
	return s.SessionUC.Logout(sessionID)
}
