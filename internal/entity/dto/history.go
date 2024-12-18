//go:generate easyjson -all .
package dto

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/google/uuid"
)

//easyjson:json
type PriceHistoryResponse struct {
	AdvertID      uuid.UUID             `json:"advert_id"`
	AdvertHistory []entity.PriceHistory `json:"advert_history"`
}
