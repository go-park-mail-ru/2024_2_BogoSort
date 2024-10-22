package user

import (
	"database/sql"
	"errors"
	"log"
	"sync"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/domain"

	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/utils"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrHashPassword      = errors.New("failed to hash password")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type userRepository struct {
	db       *sql.DB
	sessions map[string]string
	mu       sync.Mutex
}

func NewUserRepository() domain.UserRepository {
	return &userRepository{
		db:       db,
		sessions: make(map[string]string),
		mu:       sync.Mutex{},
	}
}

func (u *userRepository) CreateUser(email, password string) (*domain.User, error) {
	hash := utils.HashPassword(password)

	if hash == "" {
		return nil, ErrHashPassword
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	var exist bool
	err := u.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exist)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, ErrUserAlreadyExists
	}

	var id int
	err = u.db.QueryRow(
		"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id",
		email, hash,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		ID:           uint(id),
		Email:        email,
		PasswordHash: hash,
	}

	log.Printf("User created: %v", email)

	return newUser, nil
}

func (s *userRepository) GetUserByEmail(email string) (*domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[email]

	if !ok {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *userRepository) ValidateUserByEmailAndPassword(email, password string) (*domain.User, error) {
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

func (s *userRepository) GetAllUsers() ([]*domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users := make([]*domain.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	return users, nil
}

func (s *userRepository) GetUserBySession(sessionID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	email, exists := s.sessions[sessionID]

	if !exists {
		// return "", domain.ErrSessionDoesNotExist
	}

	return email, nil
}
