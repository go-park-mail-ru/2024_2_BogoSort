package domain

import "sync"

type Sessions struct {
	Sessions map[string]string
	Mu       sync.Mutex
}

type SessionRepository interface {
	AddSession(email string) string
	RemoveSession(email string) error
	GetSession(email string) (string, error)
	SessionExists(sessionID string) bool
}

type AuthData struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
