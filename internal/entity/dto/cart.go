package dto

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/google/uuid"
)

type AddAdvertToUserCartRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	AdvertID uuid.UUID `json:"advert_id"`
}

type DeleteAdvertFromUserCartRequest struct {
	CartID   uuid.UUID `json:"cart_id"`
	AdvertID uuid.UUID `json:"advert_id"`
}

type CartPurchase struct {
	SellerID uuid.UUID `json:"seller_id"`
	Adverts []PreviewAdvertCard `json:"adverts"`
}

type Cart struct {
	ID      uuid.UUID         `json:"id"`
	UserID  uuid.UUID         `json:"user_id"`
	CartPurchases []CartPurchase `json:"cart_purchases"`
	Status  entity.CartStatus `json:"status"`
}

type CartResponse struct {
	Cart Cart `json:"cart"`
}
