package entity

import "time"

type User struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	PasswordHash []byte
	Username     string    `json:"username"`
	Phone        string    `json:"phone"`
	AvatarId     string    `json:"avatar_id"`
	Status       uint      `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
