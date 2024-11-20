package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func parseJSONResponse(body *bytes.Buffer, v interface{}) error {
	return json.NewDecoder(body).Decode(v)
}

func setupSellerEndpoints(t *testing.T) (*SellerEndpoint, *mocks.MockSeller, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockSellerRepo := mocks.NewMockSeller(ctrl)
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	endpoints := NewSellerEndpoint(mockSellerRepo, logger)
	return endpoints, mockSellerRepo, ctrl
}

func TestSellerEndpoints_GetSellerByID(t *testing.T) {
	endpoints, mockSellerRepo, ctrl := setupSellerEndpoints(t)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		sellerID := uuid.New()
		seller := entity.Seller{
			ID:          sellerID,
			UserID:      uuid.New(),
			Description: "Test Seller Description",
		}

		mockSellerRepo.
			EXPECT().
			GetById(sellerID).
			Return(&seller, nil)

		req := httptest.NewRequest("GET", "/api/v1/seller/"+sellerID.String(), nil)
		req = mux.SetURLVars(req, map[string]string{
			"seller_id": sellerID.String(),
		})
		rr := httptest.NewRecorder()

		endpoints.GetByID(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status %v, got %v", http.StatusOK, status)
		}

		var gotSeller entity.Seller
		if err := parseJSONResponse(rr.Body, &gotSeller); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		if gotSeller != seller {
			t.Errorf("Expected seller %v, got %v", seller, gotSeller)
		}
	})

	t.Run("Invalid seller ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/seller/invalid-uuid", nil)
		req = mux.SetURLVars(req, map[string]string{
			"seller_id": "invalid-uuid",
		})
		rr := httptest.NewRecorder()

		endpoints.GetByID(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
		}

		var errResp utils.ErrResponse
		if err := parseJSONResponse(rr.Body, &errResp); err != nil {
			t.Errorf("Failed to parse error response: %v", err)
		}
	})

	t.Run("Seller not found", func(t *testing.T) {
		sellerID := uuid.New()

		mockSellerRepo.
			EXPECT().
			GetById(sellerID).
			Return(nil, repository.ErrSellerNotFound)

		req := httptest.NewRequest("GET", "/api/v1/seller/"+sellerID.String(), nil)
		req = mux.SetURLVars(req, map[string]string{
			"seller_id": sellerID.String(),
		})
		rr := httptest.NewRecorder()

		endpoints.GetByID(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Expected status %v, got %v", http.StatusNotFound, status)
		}

		var errResp utils.ErrResponse
		if err := parseJSONResponse(rr.Body, &errResp); err != nil {
			t.Errorf("Failed to parse error response: %v", err)
		}
	})

	t.Run("Internal server error", func(t *testing.T) {
		sellerID := uuid.New()

		mockSellerRepo.
			EXPECT().
			GetById(sellerID).
			Return(nil, errors.New("database error"))

		req := httptest.NewRequest("GET", "/api/v1/seller/"+sellerID.String(), nil)
		req = mux.SetURLVars(req, map[string]string{
			"seller_id": sellerID.String(),
		})
		rr := httptest.NewRecorder()

		endpoints.GetByID(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
		}

		var errResp utils.ErrResponse
		if err := parseJSONResponse(rr.Body, &errResp); err != nil {
			t.Errorf("Failed to parse error response: %v", err)
		}
	})
}

func TestSellerEndpoints_GetSellerByUserID(t *testing.T) {
	endpoints, mockSellerRepo, ctrl := setupSellerEndpoints(t)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		seller := entity.Seller{
			ID:          uuid.New(),
			UserID:      userID,
			Description: "Test Seller Description",
		}

		mockSellerRepo.
			EXPECT().
			GetByUserId(userID).
			Return(&seller, nil)

		req := httptest.NewRequest("GET", "/api/v1/seller/user/"+userID.String(), nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id": userID.String(),
		})
		rr := httptest.NewRecorder()

		endpoints.GetByUserID(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status %v, got %v", http.StatusOK, status)
		}

		var gotSeller entity.Seller
		if err := parseJSONResponse(rr.Body, &gotSeller); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		if gotSeller != seller {
			t.Errorf("Expected seller %v, got %v", seller, gotSeller)
		}
	})

	t.Run("Invalid user ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/seller/user/invalid-uuid", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id": "invalid-uuid",
		})
		rr := httptest.NewRecorder()

		endpoints.GetByUserID(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
		}

		var errResp utils.ErrResponse
		if err := parseJSONResponse(rr.Body, &errResp); err != nil {
			t.Errorf("Failed to parse error response: %v", err)
		}
	})

	t.Run("Seller not found", func(t *testing.T) {
		userID := uuid.New()

		mockSellerRepo.
			EXPECT().
			GetByUserId(userID).
			Return(nil, repository.ErrSellerNotFound)

		req := httptest.NewRequest("GET", "/api/v1/seller/user/"+userID.String(), nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id": userID.String(),
		})
		rr := httptest.NewRecorder()

		endpoints.GetByUserID(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Expected status %v, got %v", http.StatusNotFound, status)
		}

		var errResp utils.ErrResponse
		if err := parseJSONResponse(rr.Body, &errResp); err != nil {
			t.Errorf("Failed to parse error response: %v", err)
		}
	})

	t.Run("Internal server error", func(t *testing.T) {
		userID := uuid.New()

		mockSellerRepo.
			EXPECT().
			GetByUserId(userID).
			Return(nil, errors.New("database error"))

		req := httptest.NewRequest("GET", "/api/v1/seller/user/"+userID.String(), nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id": userID.String(),
		})
		rr := httptest.NewRecorder()

		endpoints.GetByUserID(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
		}

		var errResp utils.ErrResponse
		if err := parseJSONResponse(rr.Body, &errResp); err != nil {
			t.Errorf("Failed to parse error response: %v", err)
		}
	})
}
