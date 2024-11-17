package dto

import (
	"time"

	"github.com/google/uuid"
)

type AdvertRequest struct {
	CategoryId  uuid.UUID    `json:"category_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Price       uint         `json:"price"`
	Status      AdvertStatus `json:"status"`
	HasDelivery bool         `json:"has_delivery"`
	Location    string       `json:"location"`
}

type AdvertResponse struct {
	ID          uuid.UUID    `json:"id"`
	SellerId    uuid.UUID    `json:"seller_id"`
	CategoryId  uuid.UUID    `json:"category_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Price       uint         `json:"price"`
	ImageId     uuid.UUID    `json:"image_id" default:"11111111-1111-1111-1111-111111111111"`
	Status      AdvertStatus `json:"status"`
	HasDelivery bool         `json:"has_delivery"`
	Location    string       `json:"location"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type AdvertStatus string

const (
	AdvertStatusActive   AdvertStatus = "active"
	AdvertStatusInactive AdvertStatus = "inactive"
	AdvertStatusReserved AdvertStatus = "reserved"
)
