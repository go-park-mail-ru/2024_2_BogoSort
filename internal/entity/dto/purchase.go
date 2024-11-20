package dto

import (
	"github.com/google/uuid"
)

type PurchaseRequest struct {
	CartID         uuid.UUID 	`json:"cart_id"`
	Address        string   	`json:"address"`
	PaymentMethod  PaymentMethod `json:"payment_method"`
	DeliveryMethod DeliveryMethod `json:"delivery_method"`
}

type PurchaseStatus string
type PaymentMethod string
type DeliveryMethod string

const (
	StatusPending PurchaseStatus = "pending"
	StatusCompleted PurchaseStatus = "completed"
	StatusFailed PurchaseStatus = "in_progress"
	StatusCanceled PurchaseStatus = "canceled"
)

const (
	PaymentMethodCard PaymentMethod = "card"
	PaymentMethodCash PaymentMethod = "cash"
)

const (
	DeliveryMethodPickup DeliveryMethod = "pickup"
	DeliveryMethodDelivery DeliveryMethod = "delivery"
)

type PurchaseResponse struct {
	ID uuid.UUID `json:"id"`
	CartID uuid.UUID `json:"cart_id"`
	Address string `json:"address"`
	Status PurchaseStatus `json:"status"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	DeliveryMethod DeliveryMethod `json:"delivery_method"`
}
