package survey

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/survey/proto"
)

var (
	ErrSurveyNotFound = errors.New("survey not found")
)

type SurveyClient struct {
	client pb.SurveyServiceClient
	conn   *grpc.ClientConn
}

func NewSurveyClient(addr string) (*SurveyClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewSurveyServiceClient(conn)

	// Проверка доступности сервиса
	_, err = client.Ping(context.Background(), &pb.NoContent{})
	if err != nil {
		return nil, err
	}

	return &SurveyClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *SurveyClient) Close() error {
	return c.conn.Close()
}

func (c *SurveyClient) AddAnswer(ctx context.Context, userID uuid.UUID, questionID uuid.UUID, value int32) (string, error) {
	protoReq := &pb.AddAnswerRequest{
		UserId:     userID.String(),
		QuestionId: questionID.String(),
		Value:      value,
	}

	resp, err := c.client.AddAnswer(ctx, protoReq)
	if err != nil {
		return "", errors.Wrap(ErrSurveyNotFound, err.Error())
	}

	return resp.Message, nil
}

func (c *SurveyClient) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, &pb.NoContent{})
	return err
}
