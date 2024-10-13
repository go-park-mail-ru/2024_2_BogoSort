package storage

import (
	"sync"

	"errors"

	"github.com/google/uuid"
)

var ErrSessionNotExists = errors.New("session does not exist")

type SessionStorage struct {
	Sessions map[string]string
	mu       sync.Mutex
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		Sessions: make(map[string]string),
	}
}

func (s *SessionStorage) AddSession(email string) string {
	sessionID := uuid.New().String()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.Sessions[sessionID] = email

	return sessionID
}

func (s *SessionStorage) RemoveSession(email string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	if _, exists := s.Sessions[email]; !exists {
		return errors.New("session does not exist")
	}

	delete(s.Sessions, email)

	return nil
}

func (s *SessionStorage) SessionExists(sessionID string) bool {
	_, exists := s.Sessions[sessionID]

	return exists
}

func (s *SessionStorage) GetUserBySession(sessionID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	email, exists := s.Sessions[sessionID]

	if !exists {
		return "", ErrSessionNotExists
	}

	return email, nil
}
