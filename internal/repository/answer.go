package repository

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
)

type AnswerRepository interface {
	// Add добавляет ответ
	// Возможные ошибки:
	// ErrAnswerBadRequest - некорректные данные для создания ответа
	Add(answer *entity.Answer) (*entity.Answer, error)
	// GetAnswersByQuestionID возвращает все ответы для конкретного вопроса
	GetByQuestionID(questionID string) ([]entity.Answer, error)
}
