package service

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"testing"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupStaticTestService(t *testing.T) (*StaticService, *gomock.Controller, *mocks.MockStaticRepository) {
	ctrl := gomock.NewController(t)
	mockStaticRepo := mocks.NewMockStaticRepository(ctrl)
	logger := zap.NewNop()

	service := NewStaticService(mockStaticRepo, logger)

	return service, ctrl, mockStaticRepo
}

func generateJPEGImage(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, nil)
	return buf.Bytes()
}

func generatePNGImage(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func TestStaticService_GetAvatar_Success(t *testing.T) {
	service, ctrl, mockStaticRepo := setupStaticTestService(t)
	defer ctrl.Finish()

	staticID := uuid.New()
	expectedPath := "/path/to/avatar.jpg"

	mockStaticRepo.EXPECT().Get(staticID).Return(expectedPath, nil)

	path, err := service.GetAvatar(staticID)

	assert.NoError(t, err)
	assert.Equal(t, expectedPath, path)
}

func TestStaticService_GetAvatar_Error(t *testing.T) {
	service, ctrl, mockStaticRepo := setupStaticTestService(t)
	defer ctrl.Finish()

	staticID := uuid.New()
	expectedError := errors.New("static not found")

	mockStaticRepo.EXPECT().Get(staticID).Return("", expectedError)

	path, err := service.GetAvatar(staticID)

	assert.Error(t, err)
	assert.Equal(t, "", path)
}

func TestStaticService_UploadFile_SuccessJPEG(t *testing.T) {
	service, ctrl, mockStaticRepo := setupStaticTestService(t)
	defer ctrl.Finish()

	data := generateJPEGImage(200, 200)

	expectedID := uuid.New()

	mockStaticRepo.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedID, nil)

	id, err := service.UploadFile(data)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
}

func TestStaticService_UploadFile_SuccessPNG(t *testing.T) {
	service, ctrl, mockStaticRepo := setupStaticTestService(t)
	defer ctrl.Finish()

	data := generatePNGImage(200, 200)

	expectedID := uuid.New()

		mockStaticRepo.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedID, nil)

	id, err := service.UploadFile(data)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
}

func TestStaticService_UploadFile_InvalidContentType(t *testing.T) {
	service, ctrl, _ := setupStaticTestService(t)
	defer ctrl.Finish()

	data := []byte("this is not an image")

	id, err := service.UploadFile(data)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestStaticService_UploadFile_DecodeError(t *testing.T) {
	service, ctrl, _ := setupStaticTestService(t)
	defer ctrl.Finish()

	data := []byte("\xff\xd8\xff")

	id, err := service.UploadFile(data)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestStaticService_UploadFile_SmallImage(t *testing.T) {
	service, ctrl, _ := setupStaticTestService(t)
	defer ctrl.Finish()

	data := generateJPEGImage(50, 50)

	id, err := service.UploadFile(data)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestStaticService_UploadFile_EncodeError(t *testing.T) {
	_, ctrl, _ := setupStaticTestService(t)
	defer ctrl.Finish()

	t.Skip("Skipping UploadFile_EncodeError test as it requires modifying the service to inject encoder dependencies")
}

func TestStaticService_UploadFile_UploadStaticError(t *testing.T) {
	service, ctrl, mockStaticRepo := setupStaticTestService(t)
	defer ctrl.Finish()

	data := generateJPEGImage(200, 200)

	mockStaticRepo.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, errors.New("upload failed"))

	id, err := service.UploadFile(data)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestStaticService_GetStaticURL_Success(t *testing.T) {
	service, ctrl, mockStaticRepo := setupStaticTestService(t)
	defer ctrl.Finish()

	staticID := uuid.New()
	expectedURL := "https://example.com/images/avatar.jpg"

	mockStaticRepo.EXPECT().Get(staticID).Return(expectedURL, nil)

	url, err := service.GetStaticURL(staticID)

	assert.NoError(t, err)
	assert.Equal(t, expectedURL, url)
}

func TestStaticService_GetStaticURL_Error(t *testing.T) {
	service, ctrl, mockStaticRepo := setupStaticTestService(t)
	defer ctrl.Finish()

	staticID := uuid.New()
	expectedError := errors.New("static not found")

	mockStaticRepo.EXPECT().Get(staticID).Return("", expectedError)

	url, err := service.GetStaticURL(staticID)

	assert.Error(t, err)
	assert.Equal(t, "", url)
}