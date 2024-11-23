package dto

import "github.com/google/uuid"

type AnswerRequest struct {
	Value int `json:"value"`
	QuestionID uuid.UUID `json:"question_id"`
	UserID uuid.UUID `json:"user_id"`
}

type AnswerResponse struct {
	ID uuid.UUID `json:"id"`
	Value int `json:"value"`
	QuestionID uuid.UUID `json:"question_id"`
	UserID uuid.UUID `json:"user_id"`
}
