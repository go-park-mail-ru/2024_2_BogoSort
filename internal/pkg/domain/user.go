package domain

type User struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type UserRepository interface {
	CreateUser(email, password string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetAllUsers() ([]*User, error)
	ValidateUserByEmailAndPassword(email, password string) (*User, error)
}
