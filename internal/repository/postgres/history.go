// internal/repository/postgres/price_history.go
package postgres

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PriceHistoryDB struct {
	db      DBExecutor
	ctx     context.Context
	timeout time.Duration
}

func NewHistoryRepository(db *pgxpool.Pool, ctx context.Context, timeout time.Duration) (*PriceHistoryDB, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &PriceHistoryDB{
		db:      db,
		ctx:     ctx,
		timeout: timeout,
	}, nil
}

const (
	insertPriceChangeQuery = `
		INSERT INTO price_history (advert_id, old_price, new_price)
		VALUES ($1, $2, $3)`

	selectPriceHistoryQuery = `
		SELECT id, advert_id, old_price, new_price, changed_at
		FROM price_history
		WHERE advert_id = $1
		ORDER BY changed_at DESC`
)

func (p *PriceHistoryDB) AddAdvertPriceChange(advertID uuid.UUID, oldPrice, newPrice int) error {
	ctx, cancel := context.WithTimeout(p.ctx, p.timeout)
	defer cancel()

	_, err := p.db.Exec(ctx, insertPriceChangeQuery, advertID, oldPrice, newPrice)
	return err
}

func (p *PriceHistoryDB) GetAdvertPriceHistory(advertID uuid.UUID) ([]*entity.PriceHistory, error) {
	ctx, cancel := context.WithTimeout(p.ctx, p.timeout)
	defer cancel()

	rows, err := p.db.Query(ctx, selectPriceHistoryQuery, advertID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []*entity.PriceHistory
	for rows.Next() {
		var history entity.PriceHistory
		if err := rows.Scan(&history.ID, &history.AdvertID, &history.OldPrice, &history.NewPrice, &history.ChangedAt); err != nil {
			return nil, err
		}
		histories = append(histories, &history)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return histories, nil
}
