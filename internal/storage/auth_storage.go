package storage

import (
	"errors"
	"log"
	"sync"

	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/utils"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrHashPassword      = errors.New("failed to hash password")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type User struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type UserStorage struct {
	Users map[string]*User
	mu    sync.Mutex
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		Users: map[string]*User{
			"test@test.com": {
				ID:           1,
				Email:        "test@test.com",
				PasswordHash: utils.HashPassword("password"),
			},
		},
		mu: sync.Mutex{},
	}
}

func (s *UserStorage) CreateUser(email, password string) (*User, error) {
	hash := utils.HashPassword(password)

	if hash == "" {
		return nil, ErrHashPassword
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Users[email]; exists {
		return nil, ErrUserAlreadyExists
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
	s.mu.Lock()
	defer s.mu.Unlock()

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
	s.mu.Lock()
	defer s.mu.Unlock()

	users := make([]*User, 0, len(s.Users))
	for _, user := range s.Users {
		users = append(users, user)
	}

	return users, nil
}
