package entity

import (
	"time"

	"github.com/google/uuid"
)

type PriceHistory struct {
	ID        uuid.UUID `db:"id"`
	AdvertID  uuid.UUID `db:"advert_id"`
	OldPrice  int       `db:"old_price"`
	NewPrice  int       `db:"new_price"`
	ChangedAt time.Time `db:"changed_at"`
}
