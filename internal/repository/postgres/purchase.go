package postgres

import (
	"context"
	"time"
	"go.uber.org/zap"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
)

type PurchaseDB struct {
	db *pgxpool.Pool
	logger *zap.Logger
	ctx context.Context
	timeout time.Duration
}

const (
	addPurchaseQuery = `
		INSERT INTO purchase (cart_id, adress, status, payment_method, delivery_method) VALUES ($1, $2, $3, $4, $5) RETURNING id, cart_id, adress, status, payment_method, delivery_method`
)

func NewPurchaseRepository(db *pgxpool.Pool, logger *zap.Logger, ctx context.Context, timeout time.Duration) (repository.PurchaseRepository, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &PurchaseDB{
		db:      db,
		logger:  logger,
		ctx:     ctx,
		timeout: timeout,
	}, nil
}

func (c *PurchaseDB) BeginTransaction() (pgx.Tx, error) {
	tx, err := c.db.Begin(c.ctx)
	if err != nil {
		c.logger.Error("failed to begin transaction", zap.Error(err))
		return nil, err
	}
	return tx, nil
}

func (r *PurchaseDB) AddPurchase(tx pgx.Tx, purchase *entity.Purchase) (*entity.Purchase, error) {
	var entityPurchase entity.Purchase

	err := tx.QueryRow(r.ctx, addPurchaseQuery, purchase.CartID, purchase.Address, purchase.Status, purchase.PaymentMethod, purchase.DeliveryMethod).
		Scan(&entityPurchase.ID, &entityPurchase.CartID, &entityPurchase.Address, &entityPurchase.Status, &entityPurchase.PaymentMethod, &entityPurchase.DeliveryMethod)
	if err != nil {
		r.logger.Error("failed to create purchase", zap.Error(err))
		return nil, entity.PSQLWrap(err, err)
	}

	return &entityPurchase, nil
}