package entity

import "time"

type Promotion struct {
	ID        int
	Days      int
	Amount    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
