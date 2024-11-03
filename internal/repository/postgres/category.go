package postgres

import (
	"context"
	"time"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"go.uber.org/zap"
)

type CategoryDB struct {
	DB      DBExecutor
	logger  *zap.Logger
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
		logger:  logger,
		ctx:     ctx,
		timeout: timeout,
	}, nil
}

func (c *CategoryDB) GetCategories() ([]*entity.Category, error) {
	var categories []*entity.Category

	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	rows, err := c.DB.Query(ctx, getCategoryQuery) 
	if err != nil {
		c.logger.Error("failed to execute query", zap.Error(err))
		return nil, entity.PSQLWrap(err)
	}
	defer rows.Close()

	for rows.Next() {
		var dbCategory entity.Category 
		if err := rows.Scan(&dbCategory.ID, &dbCategory.Title); err != nil { 
			c.logger.Error("failed to scan row", zap.Error(err))
			return nil, entity.PSQLWrap(err)
		}
		categories = append(categories, &dbCategory)
	}

	if err := rows.Err(); err != nil {
		c.logger.Error("error iterating over rows", zap.Error(err))
		return nil, entity.PSQLWrap(err)
	}

	return categories, nil
}
