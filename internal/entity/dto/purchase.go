package dto

import (
	"github.com/google/uuid"
)

type PurchaseRequest struct {
	CartID         uuid.UUID `json:"cart_id"`
	Address        string    `json:"address"`
	PaymentMethod  string    `json:"payment_method"`
	DeliveryMethod uuid.UUID `json:"delivery_method"`
}

type PurchaseResponse struct {
	Success bool `json:"success"`
}
