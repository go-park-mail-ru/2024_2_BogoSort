package postgres

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PurchaseDB struct {
	db      DBExecutor
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
		WHERE c.user_id = $1 
		ORDER BY p.created_at DESC`
)

func NewPurchaseRepository(db *pgxpool.Pool, ctx context.Context, timeout time.Duration) (repository.PurchaseRepository, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &PurchaseDB{
		db:      db,
		ctx:     ctx,
		timeout: timeout,
	}, nil
}

func (c *PurchaseDB) BeginTransaction() (pgx.Tx, error) {
	logger := middleware.GetLogger(c.ctx)
	logger.Info("beginning transaction")

	tx, err := c.db.Begin(c.ctx)
	if err != nil {
		logger.Error("failed to begin transaction", zap.Error(err))
		return nil, err
	}
	return tx, nil
}

func (r *PurchaseDB) Add(tx pgx.Tx, purchase *entity.Purchase) (*entity.Purchase, error) {
	var entityPurchase entity.Purchase
	logger := middleware.GetLogger(r.ctx)
	logger.Info("adding purchase to db", zap.String("cart_id", purchase.CartID.String()))

	err := tx.QueryRow(r.ctx, addPurchaseQuery, purchase.CartID, purchase.Address, purchase.Status, purchase.PaymentMethod, purchase.DeliveryMethod).
		Scan(&entityPurchase.ID, &entityPurchase.CartID, &entityPurchase.Address, &entityPurchase.Status, &entityPurchase.PaymentMethod, &entityPurchase.DeliveryMethod)
	if err != nil {
		logger.Error("failed to create purchase", zap.Error(err))
		return nil, entity.PSQLWrap(err, err)
	}

	return &entityPurchase, nil
}

func (r *PurchaseDB) GetByUserId(userID uuid.UUID) ([]*entity.Purchase, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	logger := middleware.GetLogger(r.ctx)
	logger.Info("getting purchases by user id from db", zap.String("user_id", userID.String()))

	rows, err := r.db.Query(ctx, getPurchasesByUserIDQuery, userID)
	if err != nil {
		logger.Error("failed to execute getPurchasesByUserIDQuery", zap.Error(err))
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
			logger.Error("failed to scan purchase row", zap.Error(err))
			return nil, entity.PSQLWrap(err, err)
		}

		purchases = append(purchases, &purchase)
	}

	if err := rows.Err(); err != nil {
		logger.Error("rows iteration error", zap.Error(err))
		return nil, entity.PSQLWrap(err, err)
	}

	return purchases, nil
}
