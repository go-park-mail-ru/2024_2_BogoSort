package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"github.com/golang/mock/gomock"
)

func TestGetStaticById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStaticUC := mocks.NewMockStaticUseCase(ctrl)
	logger, _ := zap.NewDevelopment()

	endpoints := NewStaticEndpoint(mockStaticUC, logger)

	staticID := uuid.New()
	mockStaticUC.EXPECT().GetStaticURL(staticID).Return("http://example.com/staticfile", nil)

	req, err := http.NewRequest("GET", "/api/v1/files/"+staticID.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	endpoints.ConfigureRoutes(router)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "http://example.com/staticfile")
}

func TestGetStaticById_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStaticUC := mocks.NewMockStaticUseCase(ctrl)
	logger, _ := zap.NewDevelopment()
	endpoints := NewStaticEndpoint(mockStaticUC, logger)

	req, err := http.NewRequest("GET", "/api/v1/files/invalid-id", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	endpoints.ConfigureRoutes(router)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetStaticById_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStaticUC := mocks.NewMockStaticUseCase(ctrl)
	logger, _ := zap.NewDevelopment()
	endpoints := NewStaticEndpoint(mockStaticUC, logger)

	staticID := uuid.New()
	mockStaticUC.EXPECT().GetStaticURL(staticID).Return("", ErrStaticFileNotFound)

	req, err := http.NewRequest("GET", "/api/v1/files/"+staticID.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	endpoints.ConfigureRoutes(router)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetStaticById_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStaticUC := mocks.NewMockStaticUseCase(ctrl)
	logger, _ := zap.NewDevelopment()
	endpoints := NewStaticEndpoint(mockStaticUC, logger)

	staticID := uuid.New()
	mockStaticUC.EXPECT().GetStaticURL(staticID).Return("", errors.New("internal error"))

	req, err := http.NewRequest("GET", "/api/v1/files/"+staticID.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	endpoints.ConfigureRoutes(router)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}