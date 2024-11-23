package survey

import (
	"context"
	proto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/survey/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	proto.UnimplementedSurveyServiceServer
	surveyUC usecase.AnswerUsecase
}

func NewGrpcServer(surveyUC usecase.AnswerUsecase) *GrpcServer {
	return &GrpcServer{
		surveyUC: surveyUC,
	}
}

func (s *GrpcServer) AddAnswer(ctx context.Context, req *proto.AddAnswerRequest) (*proto.AddAnswerResponse, error) {
	answerID, err := s.surveyUC.Add(&dto.AnswerRequest{
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

func (s *GrpcServer) Ping(ctx context.Context, req *proto.NoContent) (*proto.NoContent, error) {
	return &proto.NoContent{}, nil
}
