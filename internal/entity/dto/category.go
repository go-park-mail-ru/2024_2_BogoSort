//go:generate easyjson -all .
package dto

import "github.com/google/uuid"

//easyjson:json
type Category struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}
