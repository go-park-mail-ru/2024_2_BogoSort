package entity

import "time"

type User struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	PasswordHash []byte    `db:"password_hash"`
	PasswordSalt []byte    `db:"password_salt"`
	Username     string    `db:"username"`
	Phone        string    `db:"phone"`
	AvatarId     string    `db:"avatar_id"`
	Status       uint      `db:"status"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
