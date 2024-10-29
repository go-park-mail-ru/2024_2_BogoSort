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
	ErrSessionNotFound       = errors.New("пользователь с такой сессией не найден")
	ErrSessionCreationFailed = errors.New("не получилось создать сессию")
	ErrSessionCheckFailed    = errors.New("не получилось проверить сессию")
	ErrSessionDeleteFailed   = errors.New("не удалось удалить сессию")
	ErrIncorrectID           = errors.New("некорректный id")
)
