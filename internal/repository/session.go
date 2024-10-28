package repository

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
