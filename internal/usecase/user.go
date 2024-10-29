package usecase

import (
	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/google/uuid"
)

type User interface {
	// Регистрация пользователя
	Signup(*dto.Signup) (uuid.UUID, error)
	// Авторизация пользователя
	Login(*dto.Login) (uuid.UUID, error)
	// Обновление данных пользователя
	UpdateInfo(*dto.User) error
	// Изменение пароля
	ChangePassword(userID uuid.UUID, password *dto.UpdatePassword) error
	// Получение данных пользователя
	GetUser(userID uuid.UUID) (*dto.User, error)
}

type UserIncorrectDataError struct {
	Err error
}

func (u UserIncorrectDataError) Error() string {
	return u.Err.Error()
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
