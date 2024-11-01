package dto

import "github.com/google/uuid"

type Advert struct {
	ID          uuid.UUID    `json:"id"`
	SellerId    uuid.UUID    `json:"seller_id"`
	CategoryId  uuid.UUID    `json:"category_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Price       uint         `json:"price"`
	ImageURL    string `json:"image_url"`
	Status      AdvertStatus `json:"status"`
	HasDelivery bool         `json:"has_delivery"`
	Location    string       `json:"location"`
}

type AdvertStatus string

const (
	AdvertStatusActive   AdvertStatus = "active"
	AdvertStatusInactive              = "inactive"
)