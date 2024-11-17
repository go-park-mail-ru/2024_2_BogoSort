package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func setupCartEndpoints(t *testing.T) (*CartEndpoint, *mocks.MockCart, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockCartUC := mocks.NewMockCart(ctrl)
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	endpoints := NewCartEndpoint(mockCartUC, logger)
	return endpoints, mockCartUC, ctrl
}

func TestCartEndpoints_GetCartByID(t *testing.T) {
	endpoints, mockCartUC, ctrl := setupCartEndpoints(t)
	defer ctrl.Finish()

	t.Run("Invalid cart ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/cart/invalid-uuid", nil)
		req = mux.SetURLVars(req, map[string]string{
			"cart_id": "invalid-uuid",
		})
		rr := httptest.NewRecorder()

		endpoints.GetByID(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
		}

		var errResp utils.ErrResponse
		if err := parseJSONResponse(rr.Body, &errResp); err != nil {
			t.Errorf("Failed to parse error response: %v", err)
		}
	})

	t.Run("Cart not found", func(t *testing.T) {
		cartID := uuid.New()

		mockCartUC.
			EXPECT().
			GetById(cartID).
			Return(dto.Cart{}, repository.ErrCartNotFound) 

		req := httptest.NewRequest("GET", "/api/v1/cart/"+cartID.String(), nil)
		req = mux.SetURLVars(req, map[string]string{
			"cart_id": cartID.String(),
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
		cartID := uuid.New()

		mockCartUC.
			EXPECT().
			GetById(cartID).
			Return(dto.Cart{}, errors.New("database error")) 

		req := httptest.NewRequest("GET", "/api/v1/cart/"+cartID.String(), nil)
		req = mux.SetURLVars(req, map[string]string{
			"cart_id": cartID.String(),
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

func TestCartEndpoints_GetCartByUserID(t *testing.T) {
	endpoints, mockCartUC, ctrl := setupCartEndpoints(t)
	defer ctrl.Finish()

	t.Run("Invalid user ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/cart/user/invalid-uuid", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id": "invalid-uuid",
		})
		rr := httptest.NewRecorder()

		endpoints.GetByUserID(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
		}

		var errResp utils.ErrResponse
		if err := parseJSONResponse(rr.Body, &errResp); err != nil {
			t.Errorf("Failed to parse error response: %v", err)
		}
	})

	t.Run("Cart not found", func(t *testing.T) {
		userID := uuid.New()

		mockCartUC.
			EXPECT().
			GetByUserId(userID).
			Return(dto.Cart{}, repository.ErrCartNotFound) 

		req := httptest.NewRequest("GET", "/api/v1/cart/user/"+userID.String(), nil)
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

		mockCartUC.
			EXPECT().
			GetByUserId(userID).
			Return(dto.Cart{}, errors.New("database error"))

		req := httptest.NewRequest("GET", "/api/v1/cart/user/"+userID.String(), nil)
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

func TestCartEndpoints_AddAdvertToCart(t *testing.T) {
	endpoints, mockCartUC, ctrl := setupCartEndpoints(t)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		reqBody := dto.AddAdvertToUserCartRequest{
			UserID:   uuid.New(),
			AdvertID: uuid.New(),
		}

		mockCartUC.
			EXPECT().
			AddAdvert(reqBody.UserID, reqBody.AdvertID).
			Return(nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/cart/add", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		endpoints.AddToCart(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status %v, got %v", http.StatusOK, status)
		}

		var resp map[string]string
		if err := parseJSONResponse(rr.Body, &resp); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		if resp["message"] != "advert added to user cart" {
			t.Errorf("Expected message 'advert added to user cart', got %v", resp["message"])
		}
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/cart/add", bytes.NewBuffer([]byte("invalid body")))
		rr := httptest.NewRecorder()

		endpoints.AddToCart(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
		}

		var errResp utils.ErrResponse
		if err := parseJSONResponse(rr.Body, &errResp); err != nil {
			t.Errorf("Failed to parse error response: %v", err)
		}
	})

	t.Run("Cart not found", func(t *testing.T) {
		reqBody := dto.AddAdvertToUserCartRequest{
			UserID:   uuid.New(),
			AdvertID: uuid.New(),
		}

		mockCartUC.
			EXPECT().
			AddAdvert(reqBody.UserID, reqBody.AdvertID).
			Return(repository.ErrCartNotFound)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/cart/add", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		endpoints.AddToCart(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Expected status %v, got %v", http.StatusNotFound, status)
		}

		var errResp utils.ErrResponse
		if err := parseJSONResponse(rr.Body, &errResp); err != nil {
			t.Errorf("Failed to parse error response: %v", err)
		}
	})

	t.Run("Internal server error", func(t *testing.T) {
		reqBody := dto.AddAdvertToUserCartRequest{
			UserID:   uuid.New(),
			AdvertID: uuid.New(),
		}

		mockCartUC.
			EXPECT().
			AddAdvert(reqBody.UserID, reqBody.AdvertID).
			Return(errors.New("database error"))

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/cart/add", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		endpoints.AddToCart(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
		}

		var errResp utils.ErrResponse
		if err := parseJSONResponse(rr.Body, &errResp); err != nil {
			t.Errorf("Failed to parse error response: %v", err)
		}
	})
}

func TestCartEndpoints_DeleteAdvertFromCart(t *testing.T) {
	endpoints, mockCartUC, ctrl := setupCartEndpoints(t)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		reqBody := dto.DeleteAdvertFromUserCartRequest{
			CartID:   uuid.New(),
			AdvertID: uuid.New(),
		}

		mockCartUC.
			EXPECT().
			DeleteAdvert(reqBody.CartID, reqBody.AdvertID).
			Return(nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("DELETE", "/api/v1/cart/delete", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		endpoints.DeleteFromCart(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status %v, got %v", http.StatusOK, status)
		}

		var resp map[string]string
		if err := parseJSONResponse(rr.Body, &resp); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		if resp["message"] != "advert deleted from user cart" {
			t.Errorf("Expected message 'advert deleted from user cart', got %v", resp["message"])
		}
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/cart/delete", bytes.NewBuffer([]byte("invalid body")))
		rr := httptest.NewRecorder()

		endpoints.DeleteFromCart(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
		}

		var errResp utils.ErrResponse
		if err := parseJSONResponse(rr.Body, &errResp); err != nil {
			t.Errorf("Failed to parse error response: %v", err)
		}
	})

	t.Run("Cart or advert not found", func(t *testing.T) {
		reqBody := dto.DeleteAdvertFromUserCartRequest{
			CartID:   uuid.New(),
			AdvertID: uuid.New(),
		}

		mockCartUC.
			EXPECT().
			DeleteAdvert(reqBody.CartID, reqBody.AdvertID).
			Return(repository.ErrCartOrAdvertNotFound)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("DELETE", "/api/v1/cart/delete", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		endpoints.DeleteFromCart(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Expected status %v, got %v", http.StatusNotFound, status)
		}

		var errResp utils.ErrResponse
		if err := parseJSONResponse(rr.Body, &errResp); err != nil {
			t.Errorf("Failed to parse error response: %v", err)
		}
	})

	t.Run("Internal server error", func(t *testing.T) {
		reqBody := dto.DeleteAdvertFromUserCartRequest{
			CartID:   uuid.New(),
			AdvertID: uuid.New(),
		}

		mockCartUC.
			EXPECT().
			DeleteAdvert(reqBody.CartID, reqBody.AdvertID).
			Return(errors.New("database error"))

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("DELETE", "/api/v1/cart/delete", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		endpoints.DeleteFromCart(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
		}

		var errResp utils.ErrResponse
		if err := parseJSONResponse(rr.Body, &errResp); err != nil {
			t.Errorf("Failed to parse error response: %v", err)
		}
	})
}
