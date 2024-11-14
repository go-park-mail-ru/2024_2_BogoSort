package auth

import (
	"context"

	"github.com/google/uuid"

	authProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/auth/proto"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
)

type GrpsServer struct {
	authProto.UnimplementedAuthServiceServer
	AuthUC usecase.Auth
}

func NewGrpsServer(authUC usecase.Auth) *GrpsServer {
	return &GrpsServer{AuthUC: authUC}
}

func (s *GrpsServer) GetUserIDBySession(_ context.Context, in *authProto.Session) (*authProto.User, error) {
	userID, err := s.AuthUC.GetUserIdBySession(in.Id)
	if err != nil {
		return nil, err
	}
	return &authProto.User{Id: userID.String()}, nil
}

func (s *GrpsServer) CreateSession(_ context.Context, in *authProto.User) (*authProto.Session, error) {
	userID, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, err
	}
	sessionID, err := s.AuthUC.CreateSession(userID)
	if err != nil {
		return nil, err
	}
	return &authProto.Session{Id: sessionID}, nil
}

func (s *GrpsServer) DeleteSession(_ context.Context, in *authProto.Session) (*authProto.NoContent, error) {
	err := s.AuthUC.Logout(in.Id)
	if err != nil {
		return nil, err
	}
	return &authProto.NoContent{}, nil
}

func (s *GrpsServer) Ping(_ context.Context, _ *authProto.NoContent) (*authProto.NoContent, error) {
	return &authProto.NoContent{}, nil
}
