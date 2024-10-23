package postgres

import (
	"errors"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"sync"

	"github.com/google/uuid"
)

var ErrSessionNotExists = errors.New("session does not exist")

type sessionRepository struct {
	sessions entity.Sessions
	mu       sync.Mutex
}

func NewSessionRepository() entity.SessionRepository {
	return &sessionRepository{
		sessions: entity.Sessions{
			Sessions: make(map[string]string),
		},
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
