package repository

import (
	"errors"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
)

type User interface {
	// AddUser добавляет пользователя в базу
	AddUser(email string, hash, salt []byte) (uuid.UUID, error)
	// GetUserByEmail возвращает пользователя по его емейлу
	GetUserByEmail(email string) (*entity.User, error)
	// GetUserById возвращает пользователя по его id
	GetUserById(userId uuid.UUID) (*entity.User, error)
	// UpdateUser обновляет данные пользователя
	UpdateUser(user *entity.User) error
	// DeleteUser удаляет пользователя
	DeleteUser(userID uuid.UUID) error
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrHashPassword      = errors.New("error hashing password")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserIncorrectData = errors.New("incorrect data")
)
