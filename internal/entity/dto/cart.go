//go:generate easyjson -all .
package dto

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/google/uuid"
)

//easyjson:json
type AddAdvertToUserCartRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	AdvertID uuid.UUID `json:"advert_id"`
}

//easyjson:json
type DeleteAdvertFromUserCartRequest struct {
	CartID   uuid.UUID `json:"cart_id"`
	AdvertID uuid.UUID `json:"advert_id"`
}

//easyjson:json
type CartPurchase struct {
	SellerID uuid.UUID           `json:"seller_id"`
	Adverts  []PreviewAdvertCard `json:"adverts"`
}

//easyjson:json
type Cart struct {
	ID            uuid.UUID         `json:"id"`
	UserID        uuid.UUID         `json:"user_id"`
	CartPurchases []CartPurchase    `json:"cart_purchases"`
	Status        entity.CartStatus `json:"status"`
}

//easyjson:json
type CartResponse struct {
	Cart Cart `json:"cart"`
}
