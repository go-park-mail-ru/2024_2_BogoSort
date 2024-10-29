package utils

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/google/uuid"
)

type SessionManager struct {
	sessionUC        usecase.Auth
	sessionAliveTime int
	secureCookie     bool
}

func NewSessionManager(authUC usecase.Auth, sessionAliveTime int, secureCookie bool) *SessionManager {
	return &SessionManager{
		sessionUC:        authUC,
		sessionAliveTime: sessionAliveTime,
		secureCookie:     secureCookie,
	}
}

func (s *SessionManager) CreateSession(userID uuid.UUID) (string, error) {
	session, err := s.sessionUC.CreateSession(userID)
	if err != nil {
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
		return uuid.Nil, err
	}

	return s.sessionUC.GetUserIdBySession(cookie.Value)
}

func (s *SessionManager) DeleteSession(sessionID string) error {
	err := s.sessionUC.Logout(sessionID)
	if err != nil {
		return err
	}
	return nil
}
