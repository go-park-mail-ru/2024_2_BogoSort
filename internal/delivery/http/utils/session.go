package utils

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/auth"
	proto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/auth/proto"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrSessionExpired = errors.New("session has expired")
)

type SessionManager struct {
	GrpcServer       *auth.GrpcServer
	SessionAliveTime int
	SecureCookie     bool
	Logger           *zap.Logger
}

func NewSessionManager(grpcServer *auth.GrpcServer, sessionAliveTime int, secureCookie bool, logger *zap.Logger) *SessionManager {
	return &SessionManager{
		GrpcServer:       grpcServer,
		SessionAliveTime: sessionAliveTime,
		SecureCookie:     secureCookie,
		Logger:           logger,
	}
}

func (s *SessionManager) CreateSession(userID uuid.UUID) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session, err := s.GrpcServer.CreateSession(ctx, &proto.User{Id: userID.String()})
	if err != nil {
		return "", err
	}
	return session.GetId(), nil
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cookie, err := r.Cookie("session_id")
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := s.GrpcServer.GetUserIDBySession(ctx, &proto.Session{Id: cookie.Value})
	if err != nil {
		s.Logger.Error("session expired or not found", zap.String("sessionID", cookie.Value))
		s.DeleteSession(ctx, cookie.Value)
		return uuid.Nil, ErrSessionExpired
	}

	return uuid.MustParse(userID.GetId()), nil
}

func (s *SessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := s.GrpcServer.DeleteSession(ctx, &proto.Session{Id: sessionID})
	if err != nil {
		return err
	}

	return nil
}
