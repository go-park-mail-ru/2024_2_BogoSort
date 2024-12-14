package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentDB struct {
	DB      DBExecutor
	ctx     context.Context
	timeout time.Duration
}

func NewPaymentRepository(db *pgxpool.Pool, ctx context.Context, timeout time.Duration) (repository.PaymentRepository, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &PaymentDB{
		DB:      db,
		ctx:     ctx,
		timeout: timeout,
	}, nil
}

const (
	OrderStatusInProcess = "in_process"
	OrderStatusCanceled  = "canceled"
	OrderStatusCompleted = "completed"
)

const insertOrderQuery = `
	INSERT INTO orders (order_id, amount, payment_id, status) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id, order_id, amount, payment_id, status, created_at, updated_at
`

func (r *PaymentDB) InsertOrder(orderID, amount, paymentID, status string) (*entity.Order, error) {
	var order entity.Order

	err := r.DB.QueryRow(context.Background(), insertOrderQuery, orderID, amount, paymentID, status).Scan(
		&order.ID,
		&order.OrderID,
		&order.Amount,
		&order.PaymentID,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order: %w", err)
	}

	return &order, nil
}

const updateOrderStatusQuery = `
	UPDATE orders
	SET status = $1, updated_at = CURRENT_TIMESTAMP
	WHERE order_id = $2
	RETURNING id, order_id, amount, payment_id, status, created_at, updated_at
`

func (r *PaymentDB) UpdateOrderStatus(orderID string, status string) (*entity.Order, error) {
	var order entity.Order

	err := r.DB.QueryRow(context.Background(), updateOrderStatusQuery, status, orderID).Scan(
		&order.ID,
		&order.OrderID,
		&order.Amount,
		&order.PaymentID,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return &order, nil
}

const selectOrderByIDQuery = `
	SELECT id, order_id, amount, payment_id, status, created_at, updated_at
	FROM orders
	WHERE order_id = $1
`

func (r *PaymentDB) GetOrderByID(orderID string) (*entity.Order, error) {
	var order entity.Order

	err := r.DB.QueryRow(context.Background(), selectOrderByIDQuery, orderID).Scan(
		&order.ID,
		&order.OrderID,
		&order.Amount,
		&order.PaymentID,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by id: %w", err)
	}

	return &order, nil
}

const selectOrdersInProcessQuery = `
	SELECT id, order_id, amount, payment_id, status, created_at, updated_at
	FROM orders
	WHERE status = 'in_process'
`

func (r *PaymentDB) GetOrdersInProcess() ([]entity.Order, error) {
	var orders []entity.Order

	rows, err := r.DB.Query(context.Background(), selectOrdersInProcessQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.ID, &order.OrderID, &order.Amount, &order.PaymentID, &order.Status, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return orders, nil
}
