package usecase

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
)

type AnswerUsecase interface {
	// Add добавляет ответ на вопрос
	Add(answer *dto.AnswerRequest) (*dto.AnswerResponse, error)
}
