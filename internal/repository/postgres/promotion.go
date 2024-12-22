package postgres

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PromotionDB struct {
	db      DBExecutor
	ctx     context.Context
	timeout time.Duration
}

func NewPromotionRepository(db *pgxpool.Pool, ctx context.Context, timeout time.Duration) (repository.PromotionRepository, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return &PromotionDB{
		db:      db,
		ctx:     ctx,
		timeout: timeout,
	}, nil
}

const selectPromotionByIDQuery = `
	SELECT id, days, amount
	FROM promotion
`

func (r *PromotionDB) GetPromotionInfo() (*entity.Promotion, error) {
	var promotion entity.Promotion

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)

	err := r.db.QueryRow(ctx, selectPromotionByIDQuery).Scan(&promotion.ID, &promotion.Days, &promotion.Amount)
	if err != nil {
		logger.Error("failed to get promotion by id", zap.Error(err))
		return nil, entity.PSQLWrap(err, err)
	}

	return &promotion, nil
}
