package usecase

import (
	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/google/uuid"
)

type User interface {
	// Signup регистрация пользователя
	Signup(*dto.Signup) (uuid.UUID, error)
	// Login авторизация пользователя
	Login(*dto.Login) (uuid.UUID, error)
	// UpdateInfo обновление данных пользователя
	UpdateInfo(*dto.UserUpdate) error
	// ChangePassword изменение пароля
	ChangePassword(userID uuid.UUID, password *dto.UpdatePassword) error
	// Get получение данных пользователя
	Get(userID uuid.UUID) (*dto.User, error)
	// UploadImage обновление аватара пользователя
	UploadImage(userID uuid.UUID, imageID uuid.UUID) error
}

type UserIncorrectDataError struct {
	Err error
}

func (u UserIncorrectDataError) Error() string {
	return u.Err.Error()
}

var (
	ErrUserNotFound                = errors.New("user not found")
	ErrUserAlreadyExists           = errors.New("user already exists")
	ErrInvalidCredentials          = errors.New("invalid credentials")
	ErrOldAndNewPasswordAreTheSame = errors.New("old and new password are the same")
)
