package domain

type User struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
}
