package service

import (
	"testing"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupCategoryTestService(t *testing.T) (*CategoryService, *gomock.Controller, *mocks.MockCategoryRepository) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCategoryRepository(ctrl)
	logger := zap.NewNop()

	service := NewCategoryService(mockRepo, logger)

	return service, ctrl, mockRepo
}

func TestCategoryService_GetCategories_Success(t *testing.T) {
	service, ctrl, mockRepo := setupCategoryTestService(t)
	defer ctrl.Finish()

	expectedCategories := []*entity.Category{{ID: uuid.New(), Title: "Category1"}}
	mockRepo.EXPECT().Get().Return(expectedCategories, nil)

	categories, err := service.Get()

	assert.NoError(t, err)
	assert.Equal(t, expectedCategories, categories)
}

func TestCategoryService_GetCategories_Error(t *testing.T) {
	service, ctrl, mockRepo := setupCategoryTestService(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().Get().Return(nil, assert.AnError)

	categories, err := service.Get()

	assert.Error(t, err)
	assert.Nil(t, categories)
}
