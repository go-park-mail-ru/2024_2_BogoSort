package utils

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/auth"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrSessionExpired = errors.New("session has expired")
)

type SessionManager struct {
	GrpcClient       *auth.GrpcClient
	SessionAliveTime int
	SecureCookie     bool
	Logger           *zap.Logger
}

func NewSessionManager(grpcClient *auth.GrpcClient, sessionAliveTime int, secureCookie bool, logger *zap.Logger) *SessionManager {
	return &SessionManager{
		GrpcClient:       grpcClient,
		SessionAliveTime: sessionAliveTime,
		SecureCookie:     secureCookie,
		Logger:           logger,
	}
}

func (s *SessionManager) CreateSession(userID uuid.UUID) (string, error) {
	return s.GrpcClient.CreateSession(userID)
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
		return uuid.Nil, err
	}

	userID, err := s.GrpcClient.GetUserIDBySession(cookie.Value)
	if err != nil {
		s.Logger.Error("session expired or not found", zap.String("sessionID", cookie.Value))
		s.DeleteSession(cookie.Value)
		return uuid.Nil, ErrSessionExpired
	}

	return uuid.MustParse(userID), nil
}

func (s *SessionManager) DeleteSession(sessionID string) error {
	return s.GrpcClient.DeleteSession(sessionID)
}
