package repository

import (
	"errors"

	"github.com/google/uuid"
)

type Session interface {
	// CreateSession создает сессию для пользователя
	CreateSession(userID uuid.UUID) (string, error)
	// GetSession возвращает id пользователя по sessionID
	GetSession(sessionID string) (uuid.UUID, error)
	// DeleteSession удаляет сессию
	DeleteSession(sessionID string) error
}

var (
	ErrSessionNotFound       = errors.New("user with such session not found")
	ErrSessionCreationFailed = errors.New("failed to create session")
	ErrSessionCheckFailed    = errors.New("failed to check session")
	ErrSessionDeleteFailed   = errors.New("failed to delete session")
	ErrIncorrectID           = errors.New("incorrect id")
)
