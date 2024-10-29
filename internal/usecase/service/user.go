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

	userID, err := u.userRepo.AddUser(signupInfo.Email, hash, salt)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (u *UserService) Login(loginInfo *dto.Login) (uuid.UUID, error) {
	if err := entity.ValidateEmail(loginInfo.Email); err != nil {
		return uuid.Nil, usecase.UserIncorrectDataError{Err: err}
	}
	if err := entity.ValidatePassword(loginInfo.Password); err != nil {
		return uuid.Nil, usecase.UserIncorrectDataError{Err: err}
	}

	user, err := u.userRepo.GetUserByEmail(loginInfo.Email)
	if err != nil {
		return uuid.Nil, usecase.ErrUserNotFound
	}

	return user.ID, nil
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

func (u *UserService) ChangePassword(userID uuid.UUID, password *dto.UpdatePassword) error {
	// Start Generation Here
	if len(password.NewPassword) < 6 {
		return errors.New("пароль слишком короткий")
	}

	salt, hash, err := entity.HashPassword(password.NewPassword)
	if err != nil {
		return err
	}

	entityUser := &entity.User{
		ID:           userID,
		PasswordHash: hash,
		PasswordSalt: salt,
	}

	err = u.userRepo.UpdateUser(entityUser)
	if err != nil {
		return err
	}

	return nil
}
