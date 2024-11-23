package service

import (
	"math"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
)

type StatisticService struct {
	answerRepo   repository.AnswerRepository
	questionRepo repository.QuestionRepository
}

func NewStatisticService(answerRepo repository.AnswerRepository, questionRepo repository.QuestionRepository) *StatisticService {
	return &StatisticService{answerRepo: answerRepo, questionRepo: questionRepo}
}

func (s *StatisticService) GetStats() (*dto.GetStatsResponse, error) {
	stats := dto.GetStatsResponse{}
	for _, page := range entity.PageTypeValues {
		questions, err := s.questionRepo.GetByPageType(entity.PageType(page))
		if err != nil {
			return nil, err
		}
		pageStats := dto.PageStats{
			Page:          string(page),
			QuestionStats: []dto.QuestionStats{},
		}
		questionStats := []dto.QuestionStats{}
		for _, question := range questions {
			answers, err := s.answerRepo.GetByQuestionID(question.ID.String())
			if err != nil {
				return nil, err
			}
			answerStats := []dto.AnswerStats{}
			for _, answer := range answers {
				answerStats = append(answerStats, dto.AnswerStats{
					Value: answer.Value,
					Count: len(answers),
				})
			}
			avgValue := 0
			sumValues := 0
			for _, answer := range answers {
				sumValues += answer.Value
			}
			if len(answers) > 0 {
				avgValue = int(math.Round(float64(sumValues) / float64(len(answers))))
			}
			questionStats = append(questionStats, dto.QuestionStats{
				ID:          question.ID,
				Title:       question.Title,
				AvgValue:    avgValue,
				AnswerStats: answerStats,
			})
		}
		pageStats.QuestionStats = questionStats
		stats.PageStats = append(stats.PageStats, pageStats)
	}

	return &stats, nil

}
