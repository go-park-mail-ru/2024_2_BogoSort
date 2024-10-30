package service

import (
	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo repository.User
	logger   *zap.Logger
}

func NewUserService(userRepo repository.User, logger *zap.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (u *UserService) Signup(signupInfo *dto.Signup) (uuid.UUID, error) {
	if err := entity.ValidateEmail(signupInfo.Email); err != nil {
		u.logger.Error("invalid email", zap.String("email", signupInfo.Email), zap.Error(err))
		return uuid.Nil, usecase.UserIncorrectDataError{Err: err}
	}
	if err := entity.ValidatePassword(signupInfo.Password); err != nil {
		u.logger.Error("invalid password", zap.String("password", signupInfo.Password), zap.Error(err))
		return uuid.Nil, usecase.UserIncorrectDataError{Err: err}
	}

	salt, hash, err := entity.HashPassword(signupInfo.Password)
	if err != nil {
		u.logger.Error("error hashing password", zap.String("password", signupInfo.Password), zap.Error(err))
		return uuid.Nil, entity.UsecaseWrap(errors.New("error hashing password"), err)
	}

	userID, err := u.userRepo.AddUser(signupInfo.Email, hash, salt)
	switch {
	case errors.Is(err, repository.ErrUserAlreadyExists):
		u.logger.Error("user already exists", zap.String("email", signupInfo.Email))
		return uuid.Nil, usecase.ErrUserAlreadyExists
	case err != nil:
		u.logger.Error("error adding user", zap.String("email", signupInfo.Email), zap.Error(err))
		return uuid.Nil, entity.UsecaseWrap(errors.New("error adding user"), err)
	}

	return userID, nil
}

func (u *UserService) Login(loginInfo *dto.Login) (uuid.UUID, error) {
	if err := entity.ValidateEmail(loginInfo.Email); err != nil {
		u.logger.Error("invalid email", zap.String("email", loginInfo.Email), zap.Error(err))
		return uuid.Nil, usecase.UserIncorrectDataError{Err: err}
	}
	if err := entity.ValidatePassword(loginInfo.Password); err != nil {
		u.logger.Error("invalid password", zap.String("password", loginInfo.Password), zap.Error(err))
		return uuid.Nil, usecase.UserIncorrectDataError{Err: err}
	}

	user, err := u.userRepo.GetUserByEmail(loginInfo.Email)
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		u.logger.Error("user not found", zap.String("email", loginInfo.Email))
		return uuid.Nil, usecase.ErrUserNotFound
	case err != nil:
		u.logger.Error("error getting user", zap.String("email", loginInfo.Email), zap.Error(err))
		return uuid.Nil, entity.UsecaseWrap(errors.New("error getting user"), err)
	}

	if !user.CheckPassword(loginInfo.Password) {
		u.logger.Error("invalid credentials", zap.String("email", loginInfo.Email))
		return uuid.Nil, usecase.ErrInvalidCredentials
	}

	u.logger.Info("user logged in", zap.String("email", loginInfo.Email), zap.String("userID", user.ID.String()))
	return user.ID, nil
}

func (u *UserService) UpdateInfo(user *dto.User) error {
	entityUser := &entity.User{
		Email:    user.Email,
		Username: user.Username,

		Phone:    user.Phone,
		AvatarId: user.AvatarId,
		Status:   user.Status,
	}

	err := u.userRepo.UpdateUser(entityUser)
	switch {
	case err != nil:
		u.logger.Error("error updating user", zap.String("userID", entityUser.ID.String()), zap.Error(err))
		return entity.UsecaseWrap(errors.New("error updating user"), err)
	}
	u.logger.Info("user updated", zap.String("userID", entityUser.ID.String()))
	return nil
}

func (u *UserService) DeleteUser(userID uuid.UUID) error {
	err := u.userRepo.DeleteUser(userID)
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		u.logger.Error("user not found", zap.String("userID", userID.String()))
		return usecase.ErrUserNotFound
	case err != nil:
		u.logger.Error("error deleting user", zap.String("userID", userID.String()), zap.Error(err))
		return entity.UsecaseWrap(errors.New("error deleting user"), err)
	}
	u.logger.Info("user deleted", zap.String("userID", userID.String()))
	return nil
}

func (u *UserService) GetUser(userID uuid.UUID) (*dto.User, error) {
	entityUser, err := u.userRepo.GetUserById(userID)
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		u.logger.Error("user not found", zap.String("userID", userID.String()))
		return nil, usecase.ErrUserNotFound
	case err != nil:
		u.logger.Error("error getting user", zap.String("userID", userID.String()), zap.Error(err))
		return nil, entity.UsecaseWrap(errors.New("error getting user"), err)
	}
	u.logger.Info("user found", zap.String("userID", userID.String()))
	return &dto.User{
		Email:    entityUser.Email,
		Username: entityUser.Username,
		Phone:    entityUser.Phone,
		AvatarId: entityUser.AvatarId,
		Status:   entityUser.Status,
	}, nil
}

func (u *UserService) ChangePassword(userID uuid.UUID, password *dto.UpdatePassword) error {
	if err := entity.ValidatePassword(password.NewPassword); err != nil {
		u.logger.Error("invalid password", zap.String("password", password.NewPassword), zap.Error(err))
		return usecase.UserIncorrectDataError{Err: err}
	}
	user, err := u.userRepo.GetUserById(userID)
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		u.logger.Error("user not found", zap.String("userID", userID.String()))
		return usecase.ErrUserNotFound
	case password.OldPassword == password.NewPassword:
		u.logger.Error("old and new password are the same", zap.String("userID", userID.String()))
		return usecase.ErrOldAndNewPasswordAreTheSame
	case err != nil:
		u.logger.Error("error getting user", zap.String("userID", userID.String()), zap.Error(err))
		return entity.UsecaseWrap(errors.New("error getting user"), err)
	case !user.CheckPassword(password.OldPassword):
		u.logger.Error("invalid credentials", zap.String("userID", userID.String()))
		return usecase.ErrInvalidCredentials
	}

	salt, hash, err := entity.HashPassword(password.NewPassword)
	if err != nil {
		u.logger.Error("error hashing password", zap.String("userID", userID.String()), zap.Error(err))
		return entity.UsecaseWrap(errors.New("error hashing password"), err)
	}

	entityUser := &entity.User{
		ID:           userID,
		PasswordHash: hash,
		PasswordSalt: salt,
	}

	err = u.userRepo.UpdateUser(entityUser)
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		u.logger.Error("user not found", zap.String("userID", userID.String()))
		return usecase.ErrUserNotFound
	case err != nil:
		u.logger.Error("error updating user", zap.String("userID", userID.String()), zap.Error(err))
		return entity.UsecaseWrap(errors.New("error updating user"), err)
	}
	u.logger.Info("password changed", zap.String("userID", userID.String()))
	return nil
}
