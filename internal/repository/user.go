package repository

import (
	"errors"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/jackc/pgx/v5"
)

type User interface {
	BeginTransaction() (pgx.Tx, error)
	// Add добавляет пользователя в базу в рамках транзакции
	Add(tx pgx.Tx, email string, hash, salt []byte) (uuid.UUID, error)
	// GetByEmail возвращает пользователя по его емейлу
	GetByEmail(email string) (*entity.User, error)
	// GetById возвращает пользователя по его id
	GetById(userId uuid.UUID) (*entity.User, error)
	// Update обновляет данные пользователя в рамках транзакции
	Update(user *entity.User) error
	// Delete удаляет пользователя
	Delete(userID uuid.UUID) error
	// CheckIfExists проверяет, существует ли пользователь
	CheckIfExists(userId uuid.UUID) (bool, error)
	// UploadImage обновляет аватар пользователя
	UploadImage(userID uuid.UUID, imageId uuid.UUID) error
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrHashPassword      = errors.New("error hashing password")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserIncorrectData = errors.New("incorrect data")
)
