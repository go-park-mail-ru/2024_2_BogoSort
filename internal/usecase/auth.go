package usecase

import (
	"errors"
	"github.com/google/uuid"
)

type Auth interface {
	// Logout удаляет сессию
	Logout(session string) error
	CreateSession(userId uuid.UUID) (string, error)
	GetUserIdBySession(session string) (uuid.UUID, error)
}

var (
	ErrSessionNotFound = errors.New("сессия не найдена")
)
