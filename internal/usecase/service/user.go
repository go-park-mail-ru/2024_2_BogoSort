package service

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
)

type UserService struct {
	userRepo repository.User
}

// ChangePassword implements usecase.User.
func (u *UserService) ChangePassword(userID string, password *dto.UpdatePassword) error {
	panic("unimplemented")
}

// GetUserByEmail implements usecase.User.
func (u *UserService) GetUserByEmail(email string) (*dto.User, error) {
	panic("unimplemented")
}

// GetUserById implements usecase.User.
func (u *UserService) GetUserById(userID string) (*dto.User, error) {
	panic("unimplemented")
}

func NewUserService(userRepo repository.User) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (u *UserService) Signup(user *dto.Signup) (string, error) {
	newUser, err := u.userRepo.AddUser(user.Email, user.Password)
	if err != nil {
		return "", err
	}

	return newUser, nil
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

func (u *UserService) DeleteUser(userID string) error {
	err := u.userRepo.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetUser(userID string) (*dto.User, error) {
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
