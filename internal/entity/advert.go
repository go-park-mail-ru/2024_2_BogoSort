package entity

import "github.com/google/uuid"

type Advert struct {
	ID          uuid.UUID    `db:"id"`
	SellerId    uuid.UUID    `db:"seller_id"`
	CategoryId  uuid.UUID    `db:"category_id"`
	Title       string       `db:"title"`
	Description string       `db:"description"`
	Price       uint         `db:"price"`
	ImageURL    string       `db:"image_url"`
	Status      AdvertStatus `db:"status"`
	HasDelivery bool         `db:"has_delivery"`
	Location    string       `db:"location"`
}

type AdvertStatus string

const (
	AdvertStatusActive   AdvertStatus = "active"
	AdvertStatusInactive              = "inactive"
)
