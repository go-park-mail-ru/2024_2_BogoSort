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

func (u *UserService) handleRepoError(err error, context string) error {
	switch {
	case errors.Is(err, repository.ErrUserAlreadyExists):
		u.logger.Error("user already exists", zap.String("context", context))
		return usecase.ErrUserAlreadyExists
	case errors.Is(err, repository.ErrUserNotFound):
		u.logger.Error("user not found", zap.String("context", context))
		return usecase.ErrUserNotFound
	case err != nil:
		u.logger.Error("repository error", zap.String("context", context), zap.Error(err))
		return entity.UsecaseWrap(errors.New("repository error"), err)
	}
	return nil
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
		return uuid.Nil, entity.UsecaseWrap(errors.New("error hashing password"), err)
	}

	userID, err := u.userRepo.AddUser(signupInfo.Email, hash, salt)
	if err != nil {
		return uuid.Nil, u.handleRepoError(err, "Signup")
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
		return uuid.Nil, u.handleRepoError(err, "Login")
	}

	if !user.CheckPassword(loginInfo.Password) {
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
	if err != nil {
		return u.handleRepoError(err, "UpdateInfo")
	}
	return nil
}

func (u *UserService) DeleteUser(userID uuid.UUID) error {
	err := u.userRepo.DeleteUser(userID)
	if err != nil {
		return u.handleRepoError(err, "DeleteUser")
	}
	return nil
}

func (u *UserService) GetUser(userID uuid.UUID) (*dto.User, error) {
	entityUser, err := u.userRepo.GetUserById(userID)
	if err != nil {
		return nil, u.handleRepoError(err, "GetUser")
	}
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
		return usecase.UserIncorrectDataError{Err: err}
	}
	user, err := u.userRepo.GetUserById(userID)
	switch {
	case err != nil:
		return u.handleRepoError(err, "ChangePassword")
	case password.OldPassword == password.NewPassword:
		return usecase.ErrOldAndNewPasswordAreTheSame
	case !user.CheckPassword(password.OldPassword):
		return usecase.ErrInvalidCredentials
	}

	salt, hash, err := entity.HashPassword(password.NewPassword)
	if err != nil {
		return entity.UsecaseWrap(errors.New("error hashing password"), err)
	}

	entityUser := &entity.User{
		ID:           userID,
		PasswordHash: hash,
		PasswordSalt: salt,
	}

	err = u.userRepo.UpdateUser(entityUser)
	if err != nil {
		return u.handleRepoError(err, "ChangePassword")
	}
	return nil
}
