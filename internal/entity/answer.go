package entity

import (
	"errors"

	"github.com/google/uuid"
)

type Answer struct {
	ID uuid.UUID `json:"id"`
	Value int `json:"value"`
	QuestionID uuid.UUID `json:"question_id"`
	UserID uuid.UUID `json:"user_id"`
}

var ErrValueNegative = errors.New("value is negative")

func ValidateAnswer(value int) error {
	if value < 0 {
		return ErrValueNegative
	}

	return nil
}
