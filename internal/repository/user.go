package repository

import "github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"

type User interface {
	CreateUser(email, password string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetAllUsers() ([]*User, error)
	ValidateUserByEmailAndPassword(email, password string) (*entity.User, error)
}
