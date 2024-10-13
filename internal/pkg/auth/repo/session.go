package repository

import (
	"errors"
	"sync"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/pkg/domain"
	"github.com/google/uuid"
)

var ErrSessionNotExists = errors.New("session does not exist")

type sessionRepository struct {
	sessions domain.Sessions
	mu       sync.Mutex
}

type SessionRepository interface {
	AddSession(email string) string
	RemoveSession(email string) error
	GetSession(email string) (string, error)
}

func NewSessionStorage() *domain.Sessions {
	return &domain.Sessions{
		Sessions: make(map[string]string),
	}
}

func (s *sessionRepository) AddSession(email string) string {
	sessionID := uuid.New().String()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions.Sessions[sessionID] = email

	return sessionID
}

func (s *sessionRepository) RemoveSession(email string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	if _, exists := s.sessions.Sessions[email]; !exists {
		return errors.New("session does not exist")
	}

	delete(s.sessions.Sessions, email)

	return nil
}

func (s *sessionRepository) SessionExists(sessionID string) bool {
	_, exists := s.sessions.Sessions[sessionID]

	return exists
}
