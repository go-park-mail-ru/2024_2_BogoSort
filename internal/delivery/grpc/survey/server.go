package survey

import (
	"context"
	proto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/survey/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"go.uber.org/zap"
)

type GrpcSurveyServer struct {
	proto.UnimplementedSurveyServiceServer
	answerUC usecase.AnswerUsecase
	questionRepo repository.QuestionRepository
}

func NewSurveyGrpcServer(answerUC usecase.AnswerUsecase, questionRepo repository.QuestionRepository) *GrpcSurveyServer {
	return &GrpcSurveyServer{
		answerUC: answerUC,
		questionRepo: questionRepo,
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
	pageType := ConvertEnumToDBPageType(req.Page)

	questions, err := s.questionRepo.GetByPageType(entity.PageType(pageType))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get questions: %v", err)
	}

	var protoQuestions []*proto.Question
	for _, question := range questions {
		zap.L().Info("question", zap.Any("question", question.Page))

		protoPageType := proto.PageType(proto.PageType_value[string(question.Page)])

		zap.L().Info("protoPageType", zap.Any("protoPageType", protoPageType))

		protoQuestions = append(protoQuestions, &proto.Question{
			Id: question.ID.String(),	
			Title: question.Title,
			Description: question.Description,
			Page: protoPageType,
			TriggerValue: int32(question.TriggerValue),
			LowerDescription: question.LowerDescription,
			UpperDescription: question.UpperDescription,
			ParentId: question.ParentID.UUID.String(),
		})
	}

	zap.L().Info("protoQuestions", zap.Any("protoQuestions", protoQuestions))

	var protoResp = &proto.GetQuestionsResponse{
		Questions: protoQuestions,
	}

	return protoResp, nil
}

func (s *GrpcSurveyServer) Ping(ctx context.Context, req *proto.NoContent) (*proto.NoContent, error) {
	return &proto.NoContent{}, nil
}
