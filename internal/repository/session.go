package repository

import "errors"

type Session interface {
	// CreateSession создает сессию для пользователя
	CreateSession(userID string) error
	// GetSession возвращает id пользователя по sessionID
	GetSession(sessionID string) (string, error)
	// DeleteSession удаляет сессию
	DeleteSession(sessionID string) error
	// CheckSession проверяет сессию
	CheckSession(sessionID string) (string, error)
}

var (
	ErrSessionNotFound       = errors.New("пользователь с такой сессией не найден")
	ErrSessionCreationFailed = errors.New("не получилось создать сессию")
	ErrSessionCheckFailed    = errors.New("не получилось проверить сессию")
	ErrSessionDeleteFailed   = errors.New("не удалось удалить сессию")
)
