package entity

import "github.com/google/uuid"

type Purchase struct {
	ID             uuid.UUID      `db:"id"`
	SellerID       uuid.UUID      `db:"seller_id"`
	CustomerID     uuid.UUID      `db:"customer_id"`
	CartID         uuid.UUID      `db:"cart_id"`
	Address        string         `db:"address"`
	Adverts        []Advert       `db:"adverts"`
	Status         PurchaseStatus `db:"status"`
	PaymentMethod  PaymentMethod  `db:"payment_method"`
	DeliveryMethod DeliveryMethod `db:"delivery_method"`
}

type (
	PurchaseStatus string
	PaymentMethod  string
	DeliveryMethod string
)

const (
	StatusPending   PurchaseStatus = "pending"
	StatusCompleted PurchaseStatus = "completed"
	StatusFailed    PurchaseStatus = "in_progress"
	StatusCanceled  PurchaseStatus = "canceled"
)

const (
	PaymentMethodCard PaymentMethod = "card"
	PaymentMethodCash PaymentMethod = "cash"
)

const (
	DeliveryMethodPickup   DeliveryMethod = "pickup"
	DeliveryMethodDelivery DeliveryMethod = "delivery"
)
