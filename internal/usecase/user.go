package usecase

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
)

type User interface {
	// Регистрация пользователя
	Signup(*dto.Signup) (string, error)
	// Авторизация пользователя
	Login(*dto.Login) error
	// Обновление данных пользователя
	UpdateInfo(*dto.User) error
	// Удаление пользователя
	DeleteUser(userID string) error
	// Изменение пароля
	ChangePassword(userID string, password *dto.UpdatePassword) error
	// Получение данных пользователя
	GetUserById(userID string) (*dto.User, error)
	// Получение данных пользователя по email
	GetUserByEmail(email string) (*dto.User, error)
}
