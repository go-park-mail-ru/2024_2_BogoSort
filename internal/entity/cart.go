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

type Advert struct {
	ID          uuid.UUID     `db:"id"`
	SellerId    uuid.UUID     `db:"seller_id"`
	CategoryId  uuid.UUID     `db:"category_id"`
	Title       string        `db:"title"`
	Description string        `db:"description"`
	Price       uint          `db:"price"`
	ImageURL    uuid.NullUUID `db:"image_url"`
	Status      AdvertStatus  `db:"status"`
	HasDelivery bool          `db:"has_delivery"`
	Location    string        `db:"location"`
}

type AdvertStatus string

const (
	AdvertStatusActive   AdvertStatus = "active"
	AdvertStatusInactive AdvertStatus = "inactive"
)
