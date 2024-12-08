package entity

import (
	"github.com/google/uuid"
)

type CartPurchase struct {
	SellerID uuid.UUID `json:"seller_id"`
	Adverts []Advert `json:"adverts"`
}

type Cart struct {
	ID     uuid.UUID  `json:"id"`
	UserID uuid.UUID  `json:"user_id"`
	CartPurchases []CartPurchase `json:"cart_purchases"`
	Status CartStatus `json:"status"`
}

type CartStatus string

const (
	CartStatusActive   CartStatus = "active"
	CartStatusInactive CartStatus = "inactive"
)
