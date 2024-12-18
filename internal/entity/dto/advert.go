//go:generate easyjson -all .
package dto

import (
	"time"

	"github.com/google/uuid"
)

//easyjson:json
type AdvertRequest struct {
	CategoryId  uuid.UUID    `json:"category_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Price       uint         `json:"price"`
	Status      AdvertStatus `json:"status"`
	HasDelivery bool         `json:"has_delivery"`
	Location    string       `json:"location"`
}

//easyjson:json
type PreviewAdvert struct {
	ID            uuid.UUID    `json:"id"`
	SellerId      uuid.UUID    `json:"seller_id"`
	CategoryId    uuid.UUID    `json:"category_id"`
	Title         string       `json:"title"`
	Price         uint         `json:"price"`
	ImageId       uuid.UUID    `json:"image_id"`
	Status        AdvertStatus `json:"status"`
	Location      string       `json:"location"`
	HasDelivery   bool         `json:"has_delivery"`
	PromotedUntil time.Time    `json:"promoted_until"`
}

//easyjson:json
type PreviewAdvertCard struct {
	Preview  PreviewAdvert `json:"preview"`
	IsSaved  bool          `json:"is_saved"`
	IsViewed bool          `json:"is_viewed"`
}

//easyjson:json
type MyPreviewAdvertCard struct {
	Preview     PreviewAdvert `json:"preview"`
	ViewsNumber uint          `json:"views_number"`
	SavesNumber uint          `json:"saves_number"`
}

//easyjson:json
type Advert struct {
	ID            uuid.UUID    `json:"id"`
	SellerId      uuid.UUID    `json:"seller_id"`
	CategoryId    uuid.UUID    `json:"category_id"`
	Title         string       `json:"title"`
	Description   string       `json:"description"`
	Price         uint         `json:"price"`
	ImageId       uuid.UUID    `json:"image_id"`
	Status        AdvertStatus `json:"status"`
	HasDelivery   bool         `json:"has_delivery"`
	Location      string       `json:"location"`
	SavesNumber   uint         `json:"saves_number"`
	ViewsNumber   uint         `json:"views_number"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	PromotedUntil time.Time    `json:"promoted_until"`
}

//easyjson:json
type AdvertCard struct {
	Advert   Advert `json:"advert"`
	IsSaved  bool   `json:"is_saved"`
	IsViewed bool   `json:"is_viewed"`
}

type AdvertStatus string

const (
	AdvertStatusActive   AdvertStatus = "active"
	AdvertStatusInactive AdvertStatus = "inactive"
	AdvertStatusReserved AdvertStatus = "reserved"
)
