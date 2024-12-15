//go:generate easyjson -all .
package dto

import (
	"time"

	"github.com/google/uuid"
)

//easyjson:json
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Phone     string    `json:"phone"`
	AvatarId  uuid.UUID `json:"avatar_id" default:"00000000-0000-0000-0000-000000000000"`
	Status    string    `json:"status" default:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//easyjson:json
type UserUpdate struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Phone    string    `json:"phone"`
}
