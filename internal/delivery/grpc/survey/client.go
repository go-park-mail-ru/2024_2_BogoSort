package survey

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/survey/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"go.uber.org/zap"
)

var (
	ErrSurveyNotFound = errors.New("survey not found")
)

type SurveyClient struct {
	client pb.SurveyServiceClient
	conn   *grpc.ClientConn
}

func NewSurveyGrpcClient(addr string) (*SurveyClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewSurveyServiceClient(conn)

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

func (c *SurveyClient) AddAnswer(ctx context.Context, request *dto.PostAnswersRequest) (string, error) {
	protoReq := &pb.AddAnswerRequest{
		UserId:     request.Answer.UserID.String(),
		QuestionId: request.Answer.QuestionID.String(),
		Value:      int32(request.Answer.Value),
	}

	resp, err := c.client.AddAnswer(ctx, protoReq)
	if err != nil {
		return "", errors.Wrap(ErrSurveyNotFound, err.Error())
	}

	return resp.Message, nil
}

func (c *SurveyClient) GetQuestions(ctx context.Context, request *dto.GetQuestionsRequest) ([]*entity.Question, error) {
	page, err := ConvertDBPageTypeToEnum(request.Page)
	if err != nil {
		return nil, errors.Wrap(ErrSurveyNotFound, err.Error())
	}

	protoReq := &pb.GetQuestionsRequest{
		Page: page,
	}

	zap.L().Info("protoReq", zap.Any("protoReq", protoReq))

	resp, err := c.client.GetQuestions(ctx, protoReq)
	if err != nil {
		return nil, errors.Wrap(ErrSurveyNotFound, err.Error())
	}

	zap.L().Info("resp", zap.Any("resp", resp))

	var questions []*entity.Question
	for _, pq := range resp.Questions {
		questionPageType := ConvertEnumToDBPageType(pq.Page)

		questions = append(questions, &entity.Question{
			ID:               uuid.MustParse(pq.Id),
			Page:             entity.PageType(questionPageType),
			Title:            pq.Title,
			Description:      pq.Description,
			TriggerValue:     int(pq.TriggerValue),
			LowerDescription: pq.LowerDescription,
			UpperDescription: pq.UpperDescription,
			ParentID:         uuid.NullUUID{UUID: uuid.MustParse(pq.ParentId), Valid: true},
		})
	}

	zap.L().Info("questions", zap.Any("questions", questions))

	return questions, nil
}

func (c *SurveyClient) GetStats(ctx context.Context) (*dto.GetStatsResponse, error) {
	resp, err := c.client.GetStats(ctx, &pb.NoContent{})
	zap.L().Info("resp", zap.Any("resp", resp))
	if err != nil {
		return nil, errors.Wrap(ErrSurveyNotFound, err.Error())
	}
	protoStats := ConvertProtoStatsToDB(resp)
	return protoStats, nil
}

func (c *SurveyClient) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, &pb.NoContent{})
	return err
}
