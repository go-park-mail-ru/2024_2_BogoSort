package service

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo   repository.User
	sellerRepo repository.Seller
}

func NewUserService(userRepo repository.User, sellerRepo repository.Seller) *UserService {
	return &UserService{
		userRepo:   userRepo,
		sellerRepo: sellerRepo,
	}
}

func (u *UserService) handleRepoError(err error) error {
	switch {
	case errors.Is(err, repository.ErrUserAlreadyExists):
		return usecase.ErrUserAlreadyExists
	case errors.Is(err, repository.ErrUserNotFound):
		return usecase.ErrUserNotFound
	case errors.Is(err, repository.ErrSellerAlreadyExists):
		return repository.ErrSellerAlreadyExists
	case errors.Is(err, repository.ErrSellerNotFound):
		return repository.ErrSellerNotFound
	case err != nil:
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

	ctx := context.Background()
	tx, err := u.userRepo.BeginTransaction()
	if err != nil {
		logger := middleware.GetLogger(ctx)
		logger.Error("failed to begin transaction", zap.Error(err))
		return uuid.Nil, entity.UsecaseWrap(errors.New("failed to begin transaction"), err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	userID, err := u.userRepo.Add(tx, signupInfo.Email, hash, salt)
	if err != nil {
		err = u.handleRepoError(err)
		return uuid.Nil, err
	}

	_, err = u.sellerRepo.Add(tx, userID)
	if err != nil {
		err = u.handleRepoError(err)
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

	user, err := u.userRepo.GetByEmail(loginInfo.Email)
	if err != nil {
		return uuid.Nil, u.handleRepoError(err)
	}

	if !user.CheckPassword(loginInfo.Password) {
		return uuid.Nil, usecase.ErrInvalidCredentials
	}

	return user.ID, nil
}

func (u *UserService) UpdateInfo(user *dto.UserUpdate) error {
	entityUser := &entity.User{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Phone:    user.Phone,
	}

	err := u.userRepo.Update(entityUser)
	if err != nil {
		return u.handleRepoError(err)
	}
	return nil
}

func (u *UserService) Delete(userID uuid.UUID) error {
	err := u.userRepo.Delete(userID)
	if err != nil {
		return u.handleRepoError(err)
	}
	return nil
}

func (u *UserService) Get(userID uuid.UUID) (*dto.User, error) {
	entityUser, err := u.userRepo.GetById(userID)
	if err != nil {
		return nil, u.handleRepoError(err)
	}
	return &dto.User{
		ID:        entityUser.ID,
		Email:     entityUser.Email,
		Username:  entityUser.Username,
		Phone:     entityUser.Phone,
		AvatarId:  entityUser.AvatarId,
		Status:    entityUser.Status,
		CreatedAt: entityUser.CreatedAt,
		UpdatedAt: entityUser.UpdatedAt,
	}, nil
}

func (u *UserService) ChangePassword(userID uuid.UUID, password *dto.UpdatePassword) error {
	if err := entity.ValidatePassword(password.NewPassword); err != nil {
		return usecase.UserIncorrectDataError{Err: err}
	}
	user, err := u.userRepo.GetById(userID)
	switch {
	case err != nil:
		return u.handleRepoError(err)
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

	err = u.userRepo.Update(entityUser)
	if err != nil {
		return u.handleRepoError(err)
	}
	return nil
}

func (u *UserService) UploadImage(userID uuid.UUID, imageId uuid.UUID) error {
	if err := u.userRepo.UploadImage(userID, imageId); err != nil {
		return entity.UsecaseWrap(ErrAdvertBadRequest, ErrAdvertBadRequest)
	}

	return nil
}
