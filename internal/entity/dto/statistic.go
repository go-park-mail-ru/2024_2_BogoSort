package dto

import "github.com/google/uuid"

type GetStatsResponse struct {
	PageStats []PageStats `json:"page_stats"`
}

type PageStats struct {
	Page          string          `json:"page"`
	QuestionStats []QuestionStats `json:"question_stats"`
}

type QuestionStats struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	AvgValue    int           `json:"avg_value"`
	AnswerStats []AnswerStats `json:"answer_stats"`
}

type AnswerStats struct {
	Value int `json:"value"`
	Count int `json:"count"`
}
