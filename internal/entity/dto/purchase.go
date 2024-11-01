package dto

import "github.com/google/uuid"

type PurchaseRequest struct {
	CartID uuid.UUID `json:"cart_id"`
}

type PurchaseResponse struct {
	Success bool `json:"success"`
}
