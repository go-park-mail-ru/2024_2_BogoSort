package entity

import (
	"github.com/google/uuid"
)

type Cart struct {
	ID     uuid.UUID  `json:"id"`
	UserID uuid.UUID  `json:"user_id"`
	Status CartStatus `json:"status"`
}

type CartStatus string

const (
	CartStatusActive   CartStatus = "active"
	CartStatusInactive CartStatus = "inactive"
)
