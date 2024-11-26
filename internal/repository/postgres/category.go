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

type CategoryDB struct {
	DB      DBExecutor
	ctx     context.Context
	timeout time.Duration
}

const (
	getCategoryQuery = `
		SELECT id, title FROM category`
)

func NewCategoryRepository(db *pgxpool.Pool, logger *zap.Logger, ctx context.Context, timeout time.Duration) (repository.CategoryRepository, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &CategoryDB{
		DB:      db,
		ctx:     ctx,
		timeout: timeout,
	}, nil
}

func (c *CategoryDB) Get() ([]*entity.Category, error) {
	var categories []*entity.Category

	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()
	logger := middleware.GetLogger(c.ctx)
	logger.Info("getting categories from db")

	rows, err := c.DB.Query(ctx, getCategoryQuery)
	if err != nil {
		logger.Error("failed to execute query", zap.Error(err))
		return nil, entity.PSQLWrap(err)
	}
	defer rows.Close()

	for rows.Next() {
		var dbCategory entity.Category
		if err := rows.Scan(&dbCategory.ID, &dbCategory.Title); err != nil {
			logger.Error("failed to scan row", zap.Error(err))
			return nil, entity.PSQLWrap(err)
		}
		categories = append(categories, &dbCategory)
	}

	if err := rows.Err(); err != nil {
		logger.Error("error iterating over rows", zap.Error(err))
		return nil, entity.PSQLWrap(err)
	}

	return categories, nil
}
