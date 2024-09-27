package storage

import (
	"sync"
)

type SessionStorage struct {
	Sessions map[string] bool
	mu sync.Mutex
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		Sessions: make(map[string]bool),
	}
}

func (s *SessionStorage) AddSession(jwtToken string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Sessions[jwtToken] = true
}

func (s *SessionStorage) RemoveSession(jwtToken string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Sessions, jwtToken)
}

func (s *SessionStorage) SessionExists(jwtToken string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Sessions[jwtToken]
}
