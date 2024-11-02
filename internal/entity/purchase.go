package entity

import "github.com/google/uuid"

type Purchase struct {
	CartID         uuid.UUID `db:"cart_id"`
	Address        string    `db:"address"`
	Status         string    `db:"status"`
	PaymentMethod  string    `db:"payment_method"`
	DeliveryMethod uuid.UUID `db:"delivery_method"`
}
