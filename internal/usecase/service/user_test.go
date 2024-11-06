package service

import (
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupUserTestService(t *testing.T) (*UserService, *gomock.Controller, *mocks.MockUser, *mocks.MockSeller) {
	ctrl := gomock.NewController(t)
	mockUserRepo := mocks.NewMockUser(ctrl)
	mockSellerRepo := mocks.NewMockSeller(ctrl)
	logger := zap.NewNop()

	service := NewUserService(mockUserRepo, mockSellerRepo, logger)

	return service, ctrl, mockUserRepo, mockSellerRepo
}

func createTestUser(password string) (*entity.User, []byte, []byte, error) {
	salt, hash, err := entity.HashPassword(password)
	if err != nil {
		return nil, nil, nil, err
	}
	user := &entity.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: hash,
		PasswordSalt: salt,
		Username:     "testuser",
		Phone:        "1234567890",
		AvatarId:     uuid.Nil,
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return user, salt, hash, nil
}

func TestUserService_Signup_EmailValidationError(t *testing.T) {
	service, ctrl, _, _ := setupUserTestService(t)
	defer ctrl.Finish()

	signupInfo := &dto.Signup{
		Email:    "invalid-email",
		Password: "SecureP@ssw0rd",
	}

	result, err := service.Signup(signupInfo)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email")
	assert.Equal(t, uuid.Nil, result)
}

func TestUserService_Signup_PasswordValidationError(t *testing.T) {
	service, ctrl, _, _ := setupUserTestService(t)
	defer ctrl.Finish()

	signupInfo := &dto.Signup{
		Email:    "test@example.com",
		Password: "short",
	}

	result, err := service.Signup(signupInfo)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password must contain at least 8 characters")
	assert.Equal(t, uuid.Nil, result)
}

func TestUserService_Login_Success(t *testing.T) {
	service, ctrl, mockUserRepo, _ := setupUserTestService(t)
	defer ctrl.Finish()

	loginInfo := &dto.Login{
		Email:    "test@example.com",
		Password: "SecureP@ssw0rd",
	}

	user, _, _, err := createTestUser(loginInfo.Password)
	assert.NoError(t, err)

	mockUserRepo.EXPECT().
		GetUserByEmail(loginInfo.Email).
		Return(user, nil).
		Times(1)

	result, err := service.Login(loginInfo)

	assert.NoError(t, err)
	assert.Equal(t, user.ID, result)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	service, ctrl, mockUserRepo, _ := setupUserTestService(t)
	defer ctrl.Finish()

	loginInfo := &dto.Login{
		Email:    "test@example.com",
		Password: "WrongPassword",
	}

	user, _, _, err := createTestUser("SecureP@ssw0rd")
	assert.NoError(t, err)

	mockUserRepo.EXPECT().
		GetUserByEmail(loginInfo.Email).
		Return(user, nil).
		Times(1)

	result, err := service.Login(loginInfo)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, result)
	assert.True(t, errors.Is(err, usecase.ErrInvalidCredentials))
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	service, ctrl, mockUserRepo, _ := setupUserTestService(t)
	defer ctrl.Finish()

	loginInfo := &dto.Login{
		Email:    "nonexistent@example.com",
		Password: "SecureP@ssw0rd",
	}

	mockUserRepo.EXPECT().
		GetUserByEmail(loginInfo.Email).
		Return(nil, repository.ErrUserNotFound).
		Times(1)

	result, err := service.Login(loginInfo)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, result)
	assert.True(t, errors.Is(err, usecase.ErrUserNotFound))
}

func TestUserService_GetUser_Success(t *testing.T) {
	service, ctrl, mockUserRepo, _ := setupUserTestService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	user := &entity.User{
		ID:        userID,
		Email:     "test@example.com",
		Username:  "testuser",
		Phone:     "1234567890",
		AvatarId:  uuid.Nil,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUserRepo.EXPECT().
		GetUserById(userID).
		Return(user, nil).
		Times(1)

	result, err := service.GetUser(userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.Username, result.Username)
	assert.Equal(t, user.Phone, result.Phone)
	assert.Equal(t, user.AvatarId, result.AvatarId)
	assert.Equal(t, user.Status, result.Status)
}

func TestUserService_GetUser_UserNotFound(t *testing.T) {
	service, ctrl, mockUserRepo, _ := setupUserTestService(t)
	defer ctrl.Finish()

	userID := uuid.New()

	mockUserRepo.EXPECT().
		GetUserById(userID).
		Return(nil, repository.ErrUserNotFound).
		Times(1)

	result, err := service.GetUser(userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, usecase.ErrUserNotFound))
}

func TestUserService_UpdateInfo_Success(t *testing.T) {
	service, ctrl, mockUserRepo, _ := setupUserTestService(t)
	defer ctrl.Finish()

	userDTO := &dto.User{
		ID:       uuid.MustParse("c7d1f50c-2ce7-4a07-b777-8f21bd9912cb"),
		Email:    "updated@example.com",
		Username: "updateduser",
		Phone:    "0987654321",
		AvatarId: uuid.MustParse("6b7badc9-d8f9-438f-983f-a2de27c2dc08"),
		Status:   "active",
	}

	mockUserRepo.EXPECT().
		UpdateUser(gomock.Any()).
		DoAndReturn(func(user *entity.User) error {
			assert.Equal(t, userDTO.ID, user.ID)
			assert.Equal(t, userDTO.Email, user.Email)
			assert.Equal(t, userDTO.Username, user.Username)
			assert.Equal(t, userDTO.Phone, user.Phone)
			assert.Equal(t, userDTO.AvatarId, user.AvatarId)
			assert.Equal(t, userDTO.Status, user.Status)
			return nil
		}).
		Times(1)

	err := service.UpdateInfo(userDTO)

	assert.NoError(t, err)
}

func TestUserService_DeleteUser_Success(t *testing.T) {
	service, ctrl, mockUserRepo, _ := setupUserTestService(t)
	defer ctrl.Finish()

	userID := uuid.New()

	mockUserRepo.EXPECT().
		DeleteUser(userID).
		Return(nil).
		Times(1)

	err := service.DeleteUser(userID)

	assert.NoError(t, err)
}

func TestUserService_DeleteUser_Error(t *testing.T) {
	service, ctrl, mockUserRepo, _ := setupUserTestService(t)
	defer ctrl.Finish()

	userID := uuid.New()

	mockUserRepo.EXPECT().
		DeleteUser(userID).
		Return(errors.New("delete error")).
		Times(1)

	err := service.DeleteUser(userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository error")
}

func TestUserService_ChangePassword_Success(t *testing.T) {
	service, ctrl, mockUserRepo, _ := setupUserTestService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	oldPassword := "OldP@ssw0rd"
	newPassword := "NewSecur3P@ss"

	user, _, _, err := createTestUser(oldPassword)
	assert.NoError(t, err)
	user.ID = userID

	updatePassword := &dto.UpdatePassword{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	mockUserRepo.EXPECT().
		GetUserById(userID).
		Return(user, nil).
		Times(1)

	mockUserRepo.EXPECT().
		UpdateUser(gomock.Any()).
		Return(nil).
		Times(1)

	err = service.ChangePassword(userID, updatePassword)

	assert.NoError(t, err)
}

func TestUserService_ChangePassword_InvalidOldPassword(t *testing.T) {
	service, ctrl, mockUserRepo, _ := setupUserTestService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	actualPassword := "ActualP@ssw0rd"
	wrongOldPassword := "WrongP@ssw0rd"
	newPassword := "NewSecur3P@ss"

	user, _, _, err := createTestUser(actualPassword)
	assert.NoError(t, err)
	user.ID = userID

	updatePassword := &dto.UpdatePassword{
		OldPassword: wrongOldPassword,
		NewPassword: newPassword,
	}

	mockUserRepo.EXPECT().
		GetUserById(userID).
		Return(user, nil).
		Times(1)

	err = service.ChangePassword(userID, updatePassword)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, usecase.ErrInvalidCredentials))
}

func TestUserService_ChangePassword_GetUserError(t *testing.T) {
	service, ctrl, mockUserRepo, _ := setupUserTestService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	oldPassword := "OldP@ssw0rd"
	newPassword := "NewSecur3P@ss"

	updatePassword := &dto.UpdatePassword{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	mockUserRepo.EXPECT().
		GetUserById(userID).
		Return(nil, repository.ErrUserNotFound).
		Times(1)

	err := service.ChangePassword(userID, updatePassword)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, usecase.ErrUserNotFound))
}

func TestUserService_ChangePassword_UpdateUserError(t *testing.T) {
	service, ctrl, mockUserRepo, _ := setupUserTestService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	oldPassword := "OldP@ssw0rd"
	newPassword := "NewSecur3P@ss"

	user, _, _, err := createTestUser(oldPassword)
	assert.NoError(t, err)
	user.ID = userID

	updatePassword := &dto.UpdatePassword{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	mockUserRepo.EXPECT().
		GetUserById(userID).
		Return(user, nil).
		Times(1)

	mockUserRepo.EXPECT().
		UpdateUser(gomock.Any()).
		Return(errors.New("update error")).
		Times(1)

	err = service.ChangePassword(userID, updatePassword)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository error")
}

func TestUserService_UploadImage_Success(t *testing.T) {
	service, ctrl, mockUserRepo, _ := setupUserTestService(t)
	defer ctrl.Finish()

	userID := uuid.New()
	imageID := uuid.New()

	mockUserRepo.EXPECT().
		UploadImage(userID, imageID).
		Return(nil).
		Times(1)

	err := service.UploadImage(userID, imageID)

	assert.NoError(t, err)
}

func TestUserService_UploadImage_Error(t *testing.T) {
	service, ctrl, mockUserRepo, _ := setupUserTestService(t)
	defer ctrl.Finish()

	mockUserRepo.EXPECT().
		UploadImage(gomock.Any(), gomock.Any()).
		Return(errors.New("upload error")).
		Times(1)

	err := service.UploadImage(uuid.New(), uuid.New())

	assert.Error(t, err)
}
