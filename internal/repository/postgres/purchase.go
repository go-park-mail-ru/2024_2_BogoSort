package postgres

import (
	"context"
	"time"
	"go.uber.org/zap"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
)

type PurchaseDB struct {
	db      *pgxpool.Pool
	logger  *zap.Logger
	ctx     context.Context
	timeout time.Duration
}

const (
	addPurchaseQuery = `
		INSERT INTO purchase (cart_id, adress, status, payment_method, delivery_method) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, cart_id, adress, status, payment_method, delivery_method`

	getPurchasesByUserIDQuery = `
		SELECT 
			p.id, 
			p.cart_id, 
			p.adress, 
			p.status, 
			p.payment_method, 
			p.delivery_method
		FROM purchase p
		INNER JOIN cart c ON p.cart_id = c.id
		WHERE c.user_id = $1`
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

func (r *PurchaseDB) GetPurchasesByUserID(userID uuid.UUID) ([]*entity.Purchase, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	rows, err := r.db.Query(ctx, getPurchasesByUserIDQuery, userID)
	if err != nil {
		r.logger.Error("failed to execute getPurchasesByUserIDQuery", zap.Error(err))
		return nil, entity.PSQLWrap(err, err)
	}
	defer rows.Close()

	var purchases []*entity.Purchase
	for rows.Next() {
		var purchase entity.Purchase

		err := rows.Scan(
			&purchase.ID,
			&purchase.CartID,
			&purchase.Address,
			&purchase.Status,
			&purchase.PaymentMethod,
			&purchase.DeliveryMethod,
		)
		if err != nil {
			r.logger.Error("failed to scan purchase row", zap.Error(err))
			return nil, entity.PSQLWrap(err, err)
		}

		purchases = append(purchases, &purchase)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("rows iteration error", zap.Error(err))
		return nil, entity.PSQLWrap(err, err)
	}

	return purchases, nil
}