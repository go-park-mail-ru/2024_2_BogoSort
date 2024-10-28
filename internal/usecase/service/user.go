package service

import (
	"errors"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo repository.User
}

// ChangePassword implements usecase.User.
func (u *UserService) ChangePassword(userID uuid.UUID, password *dto.UpdatePassword) error {
	panic("unimplemented")
}

// GetUserByEmail implements usecase.User.
func (u *UserService) GetUserByEmail(email string) (*dto.User, error) {
	panic("unimplemented")
}

// GetUserById implements usecase.User.
func (u *UserService) GetUserById(userID uuid.UUID) (*dto.User, error) {
	panic("unimplemented")
}

func NewUserService(userRepo repository.User) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (u *UserService) Signup(signupInfo *dto.Signup) (uuid.UUID, error) {
	if err := entity.ValidateEmail(signupInfo.Email); err != nil {
		return uuid.Nil, usecase.UserIncorrectDataError{Err: err}
	}
	if err := entity.ValidatePassword(signupInfo.Password); err != nil {
		return uuid.Nil, usecase.UserIncorrectDataError{Err: err}
	}

	salt, hash, err := entity.HashPassword(signupInfo.Password)
	if err != nil {
		return uuid.Nil, entity.UsecaseWrap(errors.New("ошибка при хешировании пароля"), err)
	}

	user, err := repository.User.AddUser(signupInfo.Email, hash, salt)

}

func (u *UserService) Login(user *dto.Login) error {
	panic("not implemented")
}

func (u *UserService) UpdateInfo(user *dto.User) error {
	entityUser := &entity.User{
		Email:    user.Email,
		Username: user.Username,
		Phone:    user.Phone,
		AvatarId: user.AvatarId,
		Status:   uint(user.Status),
	}

	err := u.userRepo.UpdateUser(entityUser)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) DeleteUser(userID uuid.UUID) error {
	err := u.userRepo.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetUser(userID uuid.UUID) (*dto.User, error) {
	entityUser, err := u.userRepo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	return &dto.User{
		Email:    entityUser.Email,
		Username: entityUser.Username,
		Phone:    entityUser.Phone,
		AvatarId: entityUser.AvatarId,
		Status:   int(entityUser.Status),
	}, nil
}
