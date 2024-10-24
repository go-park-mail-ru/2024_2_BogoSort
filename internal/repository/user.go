package repository

import (
	"errors"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
)

type User interface {
	// AddUser добавляет пользователя в базу
	AddUser(email, password string) (*entity.User, error)
	// GetUserByEmail возвращает пользователя по его емейлу
	GetUserByEmail(email string) (*entity.User, error)
	// GetUserById возвращает пользователя по его id
	GetUserById(id int) (*entity.User, error)
	// UpdateUser обновляет данные пользователя
	UpdateUser(user *entity.User) error
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrHashPassword      = errors.New("failed to hash password")
	ErrUserAlreadyExists = errors.New("user already exists")
)
