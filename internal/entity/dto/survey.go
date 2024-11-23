package dto

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
)

type GetQuestionsRequest struct {
	Page string `json:"page"`
}

type GetQuestionsResponse struct {
	Questions []entity.Question `json:"questions"`
}

type PostAnswersRequest struct {
	Answer entity.Answer `json:"answer"`
}

type PostAnswersResponse struct {
	Message string `json:"message"`
}

