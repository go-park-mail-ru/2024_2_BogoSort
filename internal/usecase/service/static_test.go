package service

import (
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
)

func setupStaticTest(t *testing.T) (*StaticService, *mocks.MockStaticRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockStaticRepository(ctrl)
	service := NewStaticService(mockRepo)
	return service, mockRepo, ctrl
}

func TestStaticService_UploadStatic_FileTooLarge(t *testing.T) {
	service, mockRepo, ctrl := setupStaticTest(t)
	defer ctrl.Finish()

	imageData := make([]byte, 2*1024*1024) // 2MB
	reader := bytes.NewReader(imageData)

	// Ожидание вызова GetMaxSize
	mockRepo.EXPECT().GetMaxSize().Return(1024 * 1024) // 1MB

	id, err := service.UploadStatic(reader)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
	// assert.Equal(t, usecase.ErrStaticTooBigFile, err)
}

func TestStaticService_UploadStatic_InvalidImage(t *testing.T) {
	service, mockRepo, ctrl := setupStaticTest(t)
	defer ctrl.Finish()

	imageData := []byte("not a valid image")
	reader := bytes.NewReader(imageData)

	// Ожидание вызова GetMaxSize
	mockRepo.EXPECT().GetMaxSize().Return(1024 * 1024) // 1MB

	id, err := service.UploadStatic(reader)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
	// assert.Equal(t, usecase.ErrStaticNotImage, err)
}

// func TestStaticService_UploadStatic_UploadError(t *testing.T) {
// 	service, mockRepo, ctrl := setupStaticTest(t)
// 	defer ctrl.Finish()

// 	imageData := []byte("test image data")
// 	reader := bytes.NewReader(imageData)

// 	// Ожидание вызова GetMaxSize
// 	mockRepo.EXPECT().GetMaxSize().Return(1024 * 1024) // 1MB

// 	// Ожидание вызова Upload с ошибкой
// 	mockRepo.EXPECT().Upload("images", gomock.Any(), imageData).Return(uuid.Nil, entity.UsecaseWrap(errors.New("upload error"), usecase.ErrStaticUploadFailed))

// 	id, err := service.UploadStatic(reader)

// 	assert.Error(t, err)
// 	assert.Equal(t, uuid.Nil, id)
// 	// assert.Equal(t, usecase.ErrStaticUploadFailed, err)
// }
