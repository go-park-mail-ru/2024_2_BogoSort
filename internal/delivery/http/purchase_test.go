package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAddPurchase_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPurchaseUC := mocks.NewMockPurchase(ctrl)
	logger, _ := zap.NewDevelopment()

	endpoints := NewPurchaseEndpoint(mockPurchaseUC, logger)

	purchaseRequest := dto.PurchaseRequest{}
	purchaseResponse := dto.PurchaseResponse{}

	mockPurchaseUC.EXPECT().AddPurchase(purchaseRequest).Return(&purchaseResponse, nil)

	body, _ := json.Marshal(purchaseRequest)
	req, err := http.NewRequest("POST", "/api/v1/purchase", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	endpoints.ConfigureRoutes(router)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response dto.PurchaseResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.Equal(t, &purchaseResponse, &response)
}

func TestAddPurchase_DecodeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPurchaseUC := mocks.NewMockPurchase(ctrl)
	logger, _ := zap.NewDevelopment()
	endpoints := NewPurchaseEndpoint(mockPurchaseUC, logger)

	req, err := http.NewRequest("POST", "/api/v1/purchase", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	endpoints.ConfigureRoutes(router)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAddPurchase_AddPurchaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPurchaseUC := mocks.NewMockPurchase(ctrl)
	logger, _ := zap.NewDevelopment()
	endpoints := NewPurchaseEndpoint(mockPurchaseUC, logger)

	purchaseRequest := dto.PurchaseRequest{}

	mockPurchaseUC.EXPECT().AddPurchase(purchaseRequest).Return(&dto.PurchaseResponse{}, errors.New("some error"))

	body, _ := json.Marshal(purchaseRequest)
	req, err := http.NewRequest("POST", "/api/v1/purchase", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	endpoints.ConfigureRoutes(router)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
