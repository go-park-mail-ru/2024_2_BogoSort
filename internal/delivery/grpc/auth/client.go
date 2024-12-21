package auth

import (
	"context"

	authProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/auth/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	authManager authProto.AuthServiceClient
}

func NewGrpcClient(addr string) (*GrpcClient, error) {
	//nolint:staticcheck // Suppressing deprecation warning for grpc.Dial
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	authManager := authProto.NewAuthServiceClient(conn)

	_, err = authManager.Ping(context.Background(), &authProto.NoContent{})
	if err != nil {
		return nil, err
	}

	return &GrpcClient{authManager: authManager}, nil
}

func (c *GrpcClient) GetUserIDBySession(sessionID string) (string, error) {
	user, err := c.authManager.GetUserIDBySession(context.Background(), &authProto.Session{Id: sessionID})
	if err != nil {
		return "", err
	}
	return user.Id, nil
}

func (c *GrpcClient) CreateSession(userID uuid.UUID) (string, error) {
	session, err := c.authManager.CreateSession(context.Background(), &authProto.User{Id: userID.String()})
	if err != nil {
		return "", err
	}
	return session.Id, nil
}

func (c *GrpcClient) DeleteSession(sessionID string) error {
	_, err := c.authManager.DeleteSession(context.Background(), &authProto.Session{Id: sessionID})
	return err
}
