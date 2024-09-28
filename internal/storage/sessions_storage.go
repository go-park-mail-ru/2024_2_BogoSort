package storage

import (
	"sync"

	"errors"
)

type SessionStorage struct {
	Sessions map[string]string 
	mu       sync.Mutex
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		Sessions: make(map[string]string),
	}
}

func (s *SessionStorage) AddSession(email, token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Sessions[email] = token
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

func (s *SessionStorage) SessionExists(email string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.Sessions[email]
	return exists
}
