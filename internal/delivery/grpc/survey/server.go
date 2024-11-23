package survey

import (
	"context"

	proto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/survey/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcSurveyServer struct {
	proto.UnimplementedSurveyServiceServer
	answerUC     usecase.AnswerUsecase
	questionRepo repository.QuestionRepository
	statisticUC  usecase.StatisticUsecase
}

func NewSurveyGrpcServer(answerUC usecase.AnswerUsecase, questionRepo repository.QuestionRepository, statisticUC usecase.StatisticUsecase) *GrpcSurveyServer {
	return &GrpcSurveyServer{
		answerUC:     answerUC,
		questionRepo: questionRepo,
		statisticUC:  statisticUC,
	}
}

func (s *GrpcSurveyServer) AddAnswer(ctx context.Context, req *proto.AddAnswerRequest) (*proto.AddAnswerResponse, error) {
	answerID, err := s.answerUC.Add(&dto.AnswerRequest{
		UserID:     uuid.MustParse(req.UserId),
		QuestionID: uuid.MustParse(req.QuestionId),
		Value:      int(req.Value),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add answer: %v", err)
	}

	return &proto.AddAnswerResponse{
		Message: "Answer added successfully with ID: " + answerID.ID.String(),
	}, nil
}

func (s *GrpcSurveyServer) GetQuestions(ctx context.Context, req *proto.GetQuestionsRequest) (*proto.GetQuestionsResponse, error) {
	zap.L().Info("GetQuestions", zap.Any("req", req.Page))
	pageType := ConvertEnumToDBPageType(req.Page)

	questions, err := s.questionRepo.GetByPageType(entity.PageType(pageType))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get questions: %v", err)
	}

	zap.L().Info("questions", zap.Any("questions", questions))

	var protoQuestions []*proto.Question
	for _, question := range questions {
		protoQuestions = append(protoQuestions, &proto.Question{
			Id: question.ID.String(),
		})
	}

	return &proto.GetQuestionsResponse{
		Questions: protoQuestions,
	}, nil
}

func (s *GrpcSurveyServer) GetStats(ctx context.Context, req *proto.NoContent) (*proto.GetStatsResponse, error) {
	stats, err := s.statisticUC.GetStats()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get stats: %v", err)
	}

	protoStats := ConvertDBStatsToProto(stats)

	return &proto.GetStatsResponse{
		PageStats: protoStats,
	}, nil
}

func (s *GrpcSurveyServer) Ping(ctx context.Context, req *proto.NoContent) (*proto.NoContent, error) {
	return &proto.NoContent{}, nil
}
