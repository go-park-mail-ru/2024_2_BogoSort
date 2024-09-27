package storage

import (
	"errors"
	"sync"
	"emporium/internal/utils"
	"log"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

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
				PasswordHash: utils.HashPassword("password"),
			},
		},
	}
}

func (s *UserStorage) CreateUser(email, password string) (*User, error) {
	hash := utils.HashPassword(password)

	if hash == "" {
		return nil, errors.New("failed to hash password")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if user already exists
	if _, exists := s.Users[email]; exists {
		return nil, errors.New("user already exists")
	}

	newUser := &User{
		ID:           uint(len(s.Users) + 1),
		Email:        email,
		PasswordHash: hash,
	}

	s.Users[email] = newUser

	log.Printf("User created: %v", email)

	return newUser, nil
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
		log.Printf("User not found: %v", email)
		return nil, err
	}

	valid := utils.ComparePassword(password, user.PasswordHash)

	if !valid {
		log.Printf("Invalid password for user: %v", email)
		return nil, ErrInvalidPassword
	}

	log.Printf("User validated: %v", email)

	return user, nil
}

func (s *UserStorage) GetAllUsers() ([]*User, error) {
	users := make([]*User, 0, len(s.Users))

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, user := range s.Users {
		users = append(users, user)
	}

	return users, nil
}