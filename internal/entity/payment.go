package entity

import (
	"time"
)

type Order struct {
	ID        int       `db:"id"`
	OrderID   string    `db:"order_id"`
	Amount    string    `db:"amount"`
	PaymentID string    `db:"payment_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"update_at"`
}