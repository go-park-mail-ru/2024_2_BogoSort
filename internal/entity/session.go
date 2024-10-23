package entity

import "sync"

type Sessions struct {
	Sessions map[string]string
	Mu       sync.Mutex
}

type SessionRepository interface {
	AddSession(email string) string
	RemoveSession(email string) error
	SessionExists(sessionID string) bool
}
