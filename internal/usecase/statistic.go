package usecase

import "github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"

type StatisticUsecase interface {
	GetStats() (*dto.GetStatsResponse, error)
}
