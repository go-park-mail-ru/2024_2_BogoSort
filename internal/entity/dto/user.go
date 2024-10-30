package dto

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Phone    string    `json:"phone"`
	AvatarId string    `json:"avatar_id"`
	Status   string    `json:"status" default:"active"`
}
