package service

import (
	"errors"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"go.uber.org/zap"
)

type AnswerService struct {
	answerRepo repository.AnswerRepository
	logger     *zap.Logger
}

var (
	ErrAnswerBadRequest = errors.New("bad request")
)

func NewAnswerService(answerRepo repository.AnswerRepository, logger *zap.Logger) *AnswerService {
	return &AnswerService{answerRepo: answerRepo, logger: logger}
}

func (u *AnswerService) Add(answer *dto.AnswerRequest) (*dto.AnswerResponse, error) {
	if err := entity.ValidateAnswer(answer.Value); err != nil {
		return nil, entity.UsecaseWrap(ErrAnswerBadRequest, err)
	}

	answerEntity := &entity.Answer{
		Value:      answer.Value,
		QuestionID: answer.QuestionID,
		UserID:     answer.UserID,
	}
	answerEntity, err := u.answerRepo.Add(answerEntity)
	if err != nil {
		return nil, entity.UsecaseWrap(err, err)
	}

	return &dto.AnswerResponse{
		ID:         answerEntity.ID,
		Value:      answerEntity.Value,
		QuestionID: answerEntity.QuestionID,
		UserID:     answerEntity.UserID,
	}, nil
}
