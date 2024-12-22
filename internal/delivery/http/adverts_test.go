package http

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
// 	"github.com/gorilla/mux"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/static"
// 	"github.com/microcosm-cc/bluemonday"
// 	"go.uber.org/zap"
// )

// type MockAdvertUseCase struct {
// 	mock.Mock
// }

// func (m *MockAdvertUseCase) Get(limit, offset int, userId uuid.UUID) ([]dto.PreviewAdvertCard, error) {
// 	args := m.Called(limit, offset, userId)
// 	return args.Get(0).([]dto.PreviewAdvertCard), args.Error(1)
// }

// func setupRouter() *mux.Router {
// 	router := mux.NewRouter()
// 	advertUC := new(MockAdvertUseCase)
// 	sessionManager := utils.NewSessionManager(nil, 0, false, zap.NewNop())
// 	policy := bluemonday.UGCPolicy()
// 	staticGrpcClient := static.NewStaticGrpcClient(nil)
// 	advertEndpoint := NewAdvertEndpoint(advertUC, staticGrpcClient, sessionManager, policy)
// 	advertEndpoint.ConfigureRoutes(router)
// 	advertEndpoint.ConfigureProtectedRoutes(router)

// 	return router
// }

// func TestGet(t *testing.T) {
// 	router := setupRouter()
// 	advertUC := new(MockAdvertUseCase)
// 	advertUC.On("Get", 10, 0, mock.Anything).Return([]dto.PreviewAdvertCard{}, nil)

// 	req, _ := http.NewRequest("GET", "/api/v1/adverts?limit=10&offset=0", nil)
// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)
// 	advertUC.AssertExpectations(t)
// }

// func TestGetBySellerId(t *testing.T) {
// 	router := setupRouter()
// 	advertUC := new(MockAdvertUseCase)
// 	sellerId := uuid.New()
// 	advertUC.On("GetBySellerId", mock.Anything, sellerId).Return([]dto.PreviewAdvertCard{}, nil)

// 	req, _ := http.NewRequest("GET", "/api/v1/adverts/seller/"+sellerId.String(), nil)
// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)
// 	advertUC.AssertExpectations(t)
// }

// func TestGetByCartId(t *testing.T) {
// 	router := setupRouter()
// 	advertUC := new(MockAdvertUseCase)
// 	cartId := uuid.New()
// 	advertUC.On("GetByCartId", cartId, mock.Anything).Return([]dto.PreviewAdvertCard{}, nil)

// 	req, _ := http.NewRequest("GET", "/api/v1/adverts/cart/"+cartId.String(), nil)
// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)
// 	advertUC.AssertExpectations(t)
// }

// // Добавьте аналогичные тесты для других методов, таких как GetSavedByUserId, GetById, Add, Update и т.д.

// func TestAdd(t *testing.T) {
// 	router := setupRouter()
// 	advertUC := new(MockAdvertUseCase)
// 	advert := dto.AdvertRequest{ /* заполните поля */ }
// 	advertUC.On("Add", &advert, mock.Anything).Return(&dto.Advert{}, nil)

// 	body, _ := json.Marshal(advert)
// 	req, _ := http.NewRequest("POST", "/api/v1/adverts", bytes.NewBuffer(body))
// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusCreated, rr.Code)
// 	advertUC.AssertExpectations(t)
// }

// ... другие тесты для методов Update, Delete, и т.д. ...
