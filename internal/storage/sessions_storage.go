package storage

import (
	"sync"
)

type SessionStorage struct {
	Sessions map[string]string // email -> token
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

func (s *SessionStorage) RemoveSession(email string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Sessions, email)
}

func (s *SessionStorage) SessionExists(email string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.Sessions[email]
	return exists
}
