package entity

import "github.com/google/uuid"

type Category struct {
	ID   uuid.UUID  `db:"id"`
	Title string    `db:"title"`
}
