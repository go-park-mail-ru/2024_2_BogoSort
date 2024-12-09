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
		INSERT INTO purchase (seller_id, customer_id, adress, status, payment_method, delivery_method, cart_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id, seller_id, customer_id, adress, status, payment_method, delivery_method, cart_id`

	addPurchaseAdvertQuery = `
		INSERT INTO purchase_advert (purchase_id, advert_id)
		VALUES ($1, $2)`

	getPurchasesByUserIDQuery = `
		SELECT 
			p.id, 
			p.seller_id,
			p.customer_id,
			p.adress, 
			p.status, 
			p.payment_method, 
			p.delivery_method,
			p.cart_id,
			array_agg(a.id) as advert_ids,
			array_agg(a.title) as advert_titles,
			array_agg(a.description) as advert_descriptions,
			array_agg(a.price) as advert_prices,
			array_agg(a.seller_id) as advert_seller_ids,
			array_agg(a.image_id) as advert_image_ids,
			array_agg(a.category_id) as advert_category_ids,
			array_agg(a.status) as advert_statuses,
			array_agg(a.location) as advert_locations,
			array_agg(a.has_delivery) as advert_has_deliveries
		FROM purchase p
		LEFT JOIN purchase_advert pa ON p.id = pa.purchase_id
		LEFT JOIN advert a ON pa.advert_id = a.id
		WHERE p.customer_id = $1 
		GROUP BY p.id, p.seller_id, p.customer_id, p.adress, p.status, 
				 p.payment_method, p.delivery_method, p.cart_id, p.created_at
		ORDER BY MAX(p.created_at) DESC`
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
	logger.Info("adding purchase to db", 
		zap.String("seller_id", purchase.SellerID.String()),
		zap.String("customer_id", purchase.CustomerID.String()),
		zap.String("cart_id", purchase.CartID.String()))

	err := tx.QueryRow(r.ctx, addPurchaseQuery,
		purchase.SellerID,
		purchase.CustomerID,
		purchase.Address,
		purchase.Status,
		purchase.PaymentMethod,
		purchase.DeliveryMethod,
		purchase.CartID).
		Scan(&entityPurchase.ID,
			&entityPurchase.SellerID,
			&entityPurchase.CustomerID,
			&entityPurchase.Address,
			&entityPurchase.Status,
			&entityPurchase.PaymentMethod,
			&entityPurchase.DeliveryMethod,
			&entityPurchase.CartID)
	if err != nil {
		logger.Error("failed to create purchase", zap.Error(err))
		return nil, entity.PSQLWrap(err, err)
	}

	for _, advert := range purchase.Adverts {
		_, err = tx.Exec(r.ctx, addPurchaseAdvertQuery, entityPurchase.ID, advert.ID)
		if err != nil {
			logger.Error("failed to add advert to purchase", zap.Error(err))
			return nil, entity.PSQLWrap(err, err)
		}
	}

	entityPurchase.Adverts = purchase.Adverts
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
		var advertIDs []uuid.UUID
		var titles, descriptions, locations []string
		var prices []int
		var sellerIDs, imageIDs, categoryIDs []uuid.UUID
		var statuses []string
		var hasDeliveries []bool

		err := rows.Scan(
			&purchase.ID,
			&purchase.SellerID,
			&purchase.CustomerID,
			&purchase.Address,
			&purchase.Status,
			&purchase.PaymentMethod,
			&purchase.DeliveryMethod,
			&purchase.CartID,
			&advertIDs,
			&titles,
			&descriptions,
			&prices,
			&sellerIDs,
			&imageIDs,
			&categoryIDs,
			&statuses,
			&locations,
			&hasDeliveries,
		)
		if err != nil {
			logger.Error("failed to scan purchase row", zap.Error(err))
			return nil, entity.PSQLWrap(err, err)
		}

		purchase.Adverts = make([]entity.Advert, len(advertIDs))
		for i := range advertIDs {
			purchase.Adverts[i] = entity.Advert{
				ID:          advertIDs[i],
				Title:       titles[i],
				Description: descriptions[i],
				Price:       uint(prices[i]),
				SellerId:    sellerIDs[i],
				ImageId:     imageIDs[i],
				Status:      entity.AdvertStatus(statuses[i]),
				CategoryId:  categoryIDs[i],
				Location:    locations[i],
				HasDelivery: hasDeliveries[i],
			}
		}

		purchases = append(purchases, &purchase)
	}

	if err := rows.Err(); err != nil {
		logger.Error("rows iteration error", zap.Error(err))
		return nil, entity.PSQLWrap(err, err)
	}

	return purchases, nil
}
