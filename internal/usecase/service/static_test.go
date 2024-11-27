package service

import (
	"bytes"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/chai2010/webp"
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

	mockRepo.EXPECT().GetMaxSize().Return(1024 * 1024) // 1MB

	id, err := service.UploadStatic(reader)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestStaticService_UploadStatic_InvalidImage(t *testing.T) {
	service, mockRepo, ctrl := setupStaticTest(t)
	defer ctrl.Finish()

	imageData := []byte("not a valid image")
	reader := bytes.NewReader(imageData)

	mockRepo.EXPECT().GetMaxSize().Return(1024 * 1024)

	id, err := service.UploadStatic(reader)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func generateValidWEBPImage() ([]byte, error) {
	width, height := 100, 100
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	blue := color.RGBA{0, 0, 255, 255}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, blue)
		}
	}

	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func TestStaticService_UploadStatic_Success(t *testing.T) {
	service, mockRepo, ctrl := setupStaticTest(t)
	defer ctrl.Finish()

	imageData, err := generateValidWEBPImage()
	if err != nil {
		t.Fatalf("Failed to generate valid WEBP image: %v", err)
	}
	reader := bytes.NewReader(imageData)

	mockRepo.EXPECT().GetMaxSize().Return(10 * 1024 * 1024) // 10MB
	mockRepo.EXPECT().Upload("images", gomock.Any(), gomock.Any()).Return(uuid.New(), nil)

	id, err := service.UploadStatic(reader)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func TestStaticService_GetAvatar_Success(t *testing.T) {
	service, mockRepo, ctrl := setupStaticTest(t)
	defer ctrl.Finish()

	expectedPath := "path/to/avatar.jpg"
	mockID := uuid.New()
	mockRepo.EXPECT().Get(mockID).Return(expectedPath, nil)

	path, err := service.GetAvatar(mockID)

	assert.NoError(t, err)
	assert.Equal(t, expectedPath, path)
}

func TestStaticService_GetStaticFile_FileNotFound(t *testing.T) {
	service, _, ctrl := setupStaticTest(t)
	defer ctrl.Finish()

	mockPath := "static/avatar.jpg"
	expectedPath, _ := filepath.Abs(mockPath)
	file, _ := os.Create(expectedPath)
	defer file.Close()

	reader, err := service.GetStaticFile(mockPath)

	assert.Error(t, err)
	assert.Nil(t, reader)
}
