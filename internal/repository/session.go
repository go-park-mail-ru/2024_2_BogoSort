package repository

import (
	"errors"

	"github.com/google/uuid"
)

type Session interface {
	// Create создает сессию для пользователя
	Create(userID uuid.UUID) (string, error)
	// Get возвращает id пользователя по sessionID
	Get(sessionID string) (uuid.UUID, error)
	// Delete удаляет сессию
	Delete(sessionID string) error
}

var (
	ErrSessionNotFound       = errors.New("user with such session not found")
	ErrSessionCreationFailed = errors.New("failed to create session")
	ErrSessionCheckFailed    = errors.New("failed to check session")
	ErrSessionDeleteFailed   = errors.New("failed to delete session")
	ErrIncorrectID           = errors.New("incorrect id")
)
