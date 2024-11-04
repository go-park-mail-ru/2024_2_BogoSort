package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"github.com/golang/mock/gomock"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
)

func TestGetCategories_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockCategoryUseCase(ctrl)
	logger, _ := zap.NewDevelopment()

	mockUseCase.EXPECT().GetCategories().Return([]entity.Category{{ID: uuid.New(), Title: "Category1"}}, nil)

	endpoints := NewCategoryEndpoints(mockUseCase, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	w := httptest.NewRecorder()

	endpoints.GetCategories(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var categories []entity.Category
	err := json.Unmarshal(w.Body.Bytes(), &categories)
	assert.NoError(t, err)
	assert.Len(t, categories, 1)
	assert.Equal(t, "Category1", categories[0].Title)
}

func TestGetCategories_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockCategoryUseCase(ctrl)
	logger, _ := zap.NewDevelopment()

	mockUseCase.EXPECT().GetCategories().Return(nil, errors.New("some error"))

	endpoints := NewCategoryEndpoints(mockUseCase, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	w := httptest.NewRecorder()

	endpoints.GetCategories(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}