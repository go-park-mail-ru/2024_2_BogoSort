package storage

import (
	"errors"
	"sync"
	"emporium/pkg"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

type LoginCredentials struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type UserStorage struct {
	Users map[string]*User
	mu sync.Mutex
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		Users: map[string]*User{
			"test@test.com": {
				ID: 1,
				Name: "John",
				Surname: "Doe",
				Email: "test@test.com",
				PasswordHash: pkg.HashPassword("password"),
			},
		},
	}
}

func (s *UserStorage) CreateUser(email, name, surname, password string) (*User, error) {
	hash := pkg.HashPassword(password)

	if hash == "" {
		return nil, errors.New("failed to hash password")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.Users[email] = &User{
		ID: uint(len(s.Users) + 1),
		Name: name,
		Surname: surname,
		Email: email,
		PasswordHash: hash,
	}

	return s.Users[email], nil
}

func (s *UserStorage) GetUserByEmail(email string) (*User, error) {
	user, ok := s.Users[email]

	if !ok {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *UserStorage) ValidateUserByEmailAndPassword(email, password string) (*User, error) {
	user, err := s.GetUserByEmail(email)

	if err != nil {
		return nil, err
	}

	valid := pkg.ComparePassword(password, user.PasswordHash)

	if !valid {
		return nil, ErrInvalidPassword
	}

	return user, nil
}
