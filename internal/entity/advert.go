package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTitleLength       = errors.New("title length exceeds 255 characters")
	ErrDescriptionLength = errors.New("description length exceeds 3000 characters")
	ErrLocationLength    = errors.New("location length exceeds 150 characters")
	ErrStatusLength      = errors.New("status length exceeds 100 characters")
	ErrPriceNegative     = errors.New("price cannot be negative")
)

type Advert struct {
	ID          uuid.UUID     `db:"id"`
	SellerId    uuid.UUID     `db:"seller_id"`
	CategoryId  uuid.UUID     `db:"category_id"`
	Title       string        `db:"title"`
	Description string        `db:"description"`
	Price       uint          `db:"price"`
	ImageId     uuid.UUID `db:"image_id"`
	Status      AdvertStatus  `db:"status"`
	HasDelivery bool          `db:"has_delivery"`
	Location    string        `db:"location"`
	SavesNumber uint          
	ViewsNumber uint         
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
	IsSaved     bool          `db:"is_saved"`
	IsViewed    bool          `db:"is_viewed"`
}

type AdvertStatus string

const (
	AdvertStatusActive   AdvertStatus = "active"
	AdvertStatusInactive AdvertStatus = "inactive"
	AdvertStatusReserved AdvertStatus = "reserved"
)

func ValidateAdvert(title, description, location, status string, price int) error {
	if len(strings.TrimSpace(title)) > 255 {
		return ErrTitleLength
	}
	if len(strings.TrimSpace(description)) > 3000 {
		return ErrDescriptionLength
	}
	if len(strings.TrimSpace(location)) > 150 {
		return ErrLocationLength
	}
	if len(strings.TrimSpace(status)) > 100 {
		return ErrStatusLength
	}
	if price < 0 {
		return ErrPriceNegative
	}
	return nil
}
