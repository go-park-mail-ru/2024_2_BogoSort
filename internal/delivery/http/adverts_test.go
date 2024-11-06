package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
)

func NewTestSessionManager() *utils.SessionManager {
	return utils.NewSessionManager(nil, 0, false, nil)
}

func setupAdvertEndpoints(t *testing.T) (*AdvertEndpoints, *mocks.MockAdvertUseCase, *mocks.MockStaticUseCase, *utils.SessionManager, *zap.Logger) {
	ctrl := gomock.NewController(t)
	mockAdvertUseCase := mocks.NewMockAdvertUseCase(ctrl)
	mockStaticUseCase := mocks.NewMockStaticUseCase(ctrl)
	sessionManager := NewTestSessionManager()
	logger, err := zap.NewDevelopment()
	policy := bluemonday.UGCPolicy()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	endpoints := NewAdvertEndpoints(mockAdvertUseCase, mockStaticUseCase, sessionManager, logger, policy)
	return endpoints, mockAdvertUseCase, mockStaticUseCase, sessionManager, logger
}

func TestAdvertEndpoints(t *testing.T) {
	endpoints, mockAdvertUseCase, _, _, _ := setupAdvertEndpoints(t)
	defer endpoints.logger.Sync()

	t.Run("GetAdverts", func(t *testing.T) {
		t.Run("Success with limit and offset", func(t *testing.T) {
			limit := 10
			offset := 5
			adverts := []*dto.AdvertResponse{
				{
					ID:          uuid.New(),
					Title:       "Advert 1",
					Description: "Description 1",
					SellerId:    uuid.New(),
				},
			}

			mockAdvertUseCase.
				EXPECT().
				GetAdverts(limit, offset).
				Return(adverts, nil)

			req := httptest.NewRequest("GET", "/api/v1/adverts?limit=10&offset=5", nil)
			rr := httptest.NewRecorder()

			endpoints.GetAdverts(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("Expected status %v, got %v", http.StatusOK, status)
			}

			var gotAdverts []dto.AdvertResponse
			if err := parseJSONResponse(rr.Body, &gotAdverts); err != nil {
				t.Errorf("Failed to parse response: %v", err)
			}

			if len(gotAdverts) != len(adverts) {
				t.Errorf("Expected %v adverts, got %v", len(adverts), len(gotAdverts))
			}
		})

		t.Run("Invalid limit", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/adverts?limit=-1&offset=5", nil)
			rr := httptest.NewRecorder()

			endpoints.GetAdverts(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
			}

			var errResp utils.ErrResponse
			if err := parseJSONResponse(rr.Body, &errResp); err != nil {
				t.Errorf("Failed to parse error response: %v", err)
			}
		})

		t.Run("Invalid offset", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/adverts?limit=10&offset=-5", nil)
			rr := httptest.NewRecorder()

			endpoints.GetAdverts(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
			}

			var errResp utils.ErrResponse
			if err := parseJSONResponse(rr.Body, &errResp); err != nil {
				t.Errorf("Failed to parse error response: %v", err)
			}
		})

		t.Run("Internal server error", func(t *testing.T) {
			limit := 10
			offset := 5

			mockAdvertUseCase.
				EXPECT().
				GetAdverts(limit, offset).
				Return(nil, errors.New("database error"))

			req := httptest.NewRequest("GET", "/api/v1/adverts?limit=10&offset=5", nil)
			rr := httptest.NewRecorder()

			endpoints.GetAdverts(rr, req)

			if status := rr.Code; status != http.StatusInternalServerError {
				t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
			}

			var errResp utils.ErrResponse
			if err := parseJSONResponse(rr.Body, &errResp); err != nil {
				t.Errorf("Failed to parse error response: %v", err)
			}
		})
	})

	t.Run("GetAdvertsBySellerId", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			sellerID := uuid.New()
			adverts := []*dto.AdvertResponse{
				{
					ID:          uuid.New(),
					Title:       "Advert 1",
					Description: "Description 1",
					SellerId:    sellerID,
				},
			}

			mockAdvertUseCase.
				EXPECT().
				GetAdvertsByUserId(sellerID).
				Return(adverts, nil)

			req := httptest.NewRequest("GET", "/api/v1/adverts/seller/"+sellerID.String(), nil)
			req = mux.SetURLVars(req, map[string]string{
				"sellerId": sellerID.String(),
			})
			rr := httptest.NewRecorder()

			endpoints.GetAdvertsBySellerId(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("Expected status %v, got %v", http.StatusOK, status)
			}

			var gotAdverts []dto.AdvertResponse
			if err := parseJSONResponse(rr.Body, &gotAdverts); err != nil {
				t.Errorf("Failed to parse response: %v", err)
			}

			if len(gotAdverts) != len(adverts) {
				t.Errorf("Expected %v adverts, got %v", len(adverts), len(gotAdverts))
			}
		})

		t.Run("Invalid seller ID", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/adverts/seller/invalid-uuid", nil)
			req = mux.SetURLVars(req, map[string]string{
				"sellerId": "invalid-uuid",
			})
			rr := httptest.NewRecorder()

			endpoints.GetAdvertsBySellerId(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
			}

			var errResp utils.ErrResponse
			if err := parseJSONResponse(rr.Body, &errResp); err != nil {
				t.Errorf("Failed to parse error response: %v", err)
			}
		})

		t.Run("Internal server error", func(t *testing.T) {
			sellerID := uuid.New()

			mockAdvertUseCase.
				EXPECT().
				GetAdvertsByUserId(sellerID).
				Return(nil, errors.New("database error"))

			req := httptest.NewRequest("GET", "/api/v1/adverts/seller/"+sellerID.String(), nil)
			req = mux.SetURLVars(req, map[string]string{
				"sellerId": sellerID.String(),
			})
			rr := httptest.NewRecorder()

			endpoints.GetAdvertsBySellerId(rr, req)

			if status := rr.Code; status != http.StatusInternalServerError {
				t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
			}

			var errResp utils.ErrResponse
			if err := parseJSONResponse(rr.Body, &errResp); err != nil {
				t.Errorf("Failed to parse error response: %v", err)
			}
		})
	})

	t.Run("GetAdvertsByCartId", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			cartID := uuid.New()
			adverts := []*dto.AdvertResponse{
				{
					ID:          uuid.New(),
					Title:       "Advert 1",
					Description: "Description 1",
					SellerId:    uuid.New(),
				},
			}

			mockAdvertUseCase.
				EXPECT().
				GetAdvertsByCartId(cartID).
				Return(adverts, nil)

			req := httptest.NewRequest("GET", "/api/v1/adverts/cart/"+cartID.String(), nil)
			req = mux.SetURLVars(req, map[string]string{
				"cartId": cartID.String(),
			})
			rr := httptest.NewRecorder()

			endpoints.GetAdvertsByCartId(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("Expected status %v, got %v", http.StatusOK, status)
			}

			var gotAdverts []dto.AdvertResponse
			if err := parseJSONResponse(rr.Body, &gotAdverts); err != nil {
				t.Errorf("Failed to parse response: %v", err)
			}

			if len(gotAdverts) != len(adverts) {
				t.Errorf("Expected %v adverts, got %v", len(adverts), len(gotAdverts))
			}
		})

		t.Run("Invalid cart ID", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/adverts/cart/invalid-uuid", nil)
			req = mux.SetURLVars(req, map[string]string{
				"cartId": "invalid-uuid",
			})
			rr := httptest.NewRecorder()

			endpoints.GetAdvertsByCartId(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
			}

			var errResp utils.ErrResponse
			if err := parseJSONResponse(rr.Body, &errResp); err != nil {
				t.Errorf("Failed to parse error response: %v", err)
			}
		})

		t.Run("Internal server error", func(t *testing.T) {
			cartID := uuid.New()

			mockAdvertUseCase.
				EXPECT().
				GetAdvertsByCartId(cartID).
				Return(nil, errors.New("database error"))

			req := httptest.NewRequest("GET", "/api/v1/adverts/cart/"+cartID.String(), nil)
			req = mux.SetURLVars(req, map[string]string{
				"cartId": cartID.String(),
			})
			rr := httptest.NewRecorder()

			endpoints.GetAdvertsByCartId(rr, req)

			if status := rr.Code; status != http.StatusInternalServerError {
				t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
			}

			var errResp utils.ErrResponse
			if err := parseJSONResponse(rr.Body, &errResp); err != nil {
				t.Errorf("Failed to parse error response: %v", err)
			}
		})
	})

	t.Run("GetAdvertById", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			advertID := uuid.New()
			advert := dto.AdvertResponse{
				ID:          advertID,
				Title:       "Test Advert",
				Description: "This is a test advert.",
				SellerId:    uuid.New(),
			}

			mockAdvertUseCase.
				EXPECT().
				GetAdvertById(advertID).
				Return(&advert, nil)

			req := httptest.NewRequest("GET", "/api/v1/adverts/"+advertID.String(), nil)
			req = mux.SetURLVars(req, map[string]string{
				"advertId": advertID.String(),
			})
			rr := httptest.NewRecorder()

			endpoints.GetAdvertById(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("Expected status %v, got %v", http.StatusOK, status)
			}

			var gotAdvert dto.AdvertResponse
			if err := parseJSONResponse(rr.Body, &gotAdvert); err != nil {
				t.Errorf("Failed to parse response: %v", err)
			}

			if gotAdvert != advert {
				t.Errorf("Expected advert %v, got %v", advert, gotAdvert)
			}
		})

		t.Run("Invalid advert ID", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/adverts/invalid-uuid", nil)
			req = mux.SetURLVars(req, map[string]string{
				"advertId": "invalid-uuid",
			})
			rr := httptest.NewRecorder()

			endpoints.GetAdvertById(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
			}

			var errResp utils.ErrResponse
			if err := parseJSONResponse(rr.Body, &errResp); err != nil {
				t.Errorf("Failed to parse error response: %v", err)
			}
		})

		t.Run("Advert not found", func(t *testing.T) {
			advertID := uuid.New()

			mockAdvertUseCase.
				EXPECT().
				GetAdvertById(advertID).
				Return(nil, ErrAdvertNotFound)

			req := httptest.NewRequest("GET", "/api/v1/adverts/"+advertID.String(), nil)
			req = mux.SetURLVars(req, map[string]string{
				"advertId": advertID.String(),
			})
			rr := httptest.NewRecorder()

			endpoints.GetAdvertById(rr, req)

			if status := rr.Code; status != http.StatusNotFound {
				t.Errorf("Expected status %v, got %v", http.StatusNotFound, status)
			}

			var errResp utils.ErrResponse
			if err := parseJSONResponse(rr.Body, &errResp); err != nil {
				t.Errorf("Failed to parse error response: %v", err)
			}
		})

		t.Run("Internal server error", func(t *testing.T) {
			advertID := uuid.New()

			mockAdvertUseCase.
				EXPECT().
				GetAdvertById(advertID).
				Return(nil, errors.New("database error"))

			req := httptest.NewRequest("GET", "/api/v1/adverts/"+advertID.String(), nil)
			req = mux.SetURLVars(req, map[string]string{
				"advertId": advertID.String(),
			})
			rr := httptest.NewRecorder()

			endpoints.GetAdvertById(rr, req)

			if status := rr.Code; status != http.StatusInternalServerError {
				t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
			}

			var errResp utils.ErrResponse
			if err := parseJSONResponse(rr.Body, &errResp); err != nil {
				t.Errorf("Failed to parse error response: %v", err)
			}
		})
	})

	t.Run("UpdateAdvert", func(t *testing.T) {
		t.Run("Invalid advert data", func(t *testing.T) {
			advertID := uuid.New()
			body := []byte(`{"title": "Incomplete Advert"`)
			req := httptest.NewRequest("PUT", "/api/v1/adverts/"+advertID.String(), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{
				"advertId": advertID.String(),
			})
			rr := httptest.NewRecorder()

			endpoints.UpdateAdvert(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
			}

			var errResp utils.ErrResponse
			if err := parseJSONResponse(rr.Body, &errResp); err != nil {
				t.Errorf("Failed to parse error response: %v", err)
			}
		})
	})

	t.Run("DeleteAdvertById", func(t *testing.T) {
		t.Run("Invalid advert ID", func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/api/v1/adverts/invalid-uuid", nil)
			req = mux.SetURLVars(req, map[string]string{
				"advertId": "invalid-uuid",
			})
			rr := httptest.NewRecorder()

			endpoints.DeleteAdvertById(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
			}

			var errResp utils.ErrResponse
			if err := parseJSONResponse(rr.Body, &errResp); err != nil {
				t.Errorf("Failed to parse error response: %v", err)
			}
		})
	})

	t.Run("UpdateAdvertStatus", func(t *testing.T) {
		testCases := []struct {
			name         string
			advertID     uuid.UUID
			statusInput  string
			mockBehavior func()
			expectedCode int
		}{
			{
				name:         "Invalid status value",
				advertID:     uuid.New(),
				statusInput:  "unknown_status",
				mockBehavior: func() {},
				expectedCode: http.StatusBadRequest,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.mockBehavior()

				reqBody := bytes.NewBufferString(`{"status":"` + tc.statusInput + `"}`)
				req := httptest.NewRequest("PUT", "/api/v1/adverts/"+tc.advertID.String()+"/status", reqBody)
				req.Header.Set("Content-Type", "application/json")
				req = mux.SetURLVars(req, map[string]string{
					"advertId": tc.advertID.String(),
				})
				rr := httptest.NewRecorder()

				endpoints.UpdateAdvertStatus(rr, req)

				if status := rr.Code; status != tc.expectedCode {
					t.Errorf("Expected status %v, got %v", tc.expectedCode, status)
				}

				if tc.expectedCode != http.StatusOK {
					var errResp utils.ErrResponse
					if err := parseJSONResponse(rr.Body, &errResp); err != nil {
						t.Errorf("Failed to parse error response: %v", err)
					}
				}
			})
		}
	})

	t.Run("GetAdvertsByCategoryId", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			categoryID := uuid.New()
			adverts := []*dto.AdvertResponse{
				{
					ID:          uuid.New(),
					Title:       "Advert 1",
					Description: "Description 1",
					SellerId:    uuid.New(),
				},
			}

			mockAdvertUseCase.
				EXPECT().
				GetAdvertsByCategoryId(categoryID).
				Return(adverts, nil)

			req := httptest.NewRequest("GET", "/api/v1/adverts/category/"+categoryID.String(), nil)
			req = mux.SetURLVars(req, map[string]string{
				"categoryId": categoryID.String(),
			})
			rr := httptest.NewRecorder()

			endpoints.GetAdvertsByCategoryId(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("Expected status %v, got %v", http.StatusOK, status)
			}

			var gotAdverts []dto.AdvertResponse
			if err := parseJSONResponse(rr.Body, &gotAdverts); err != nil {
				t.Errorf("Failed to parse response: %v", err)
			}

			if len(gotAdverts) != len(adverts) {
				t.Errorf("Expected %v adverts, got %v", len(adverts), len(gotAdverts))
			}
		})

		t.Run("Invalid category ID", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/adverts/category/invalid-uuid", nil)
			req = mux.SetURLVars(req, map[string]string{
				"categoryId": "invalid-uuid",
			})
			rr := httptest.NewRecorder()

			endpoints.GetAdvertsByCategoryId(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
			}

			var errResp utils.ErrResponse
			if err := parseJSONResponse(rr.Body, &errResp); err != nil {
				t.Errorf("Failed to parse error response: %v", err)
			}
		})

		t.Run("Internal server error", func(t *testing.T) {
			categoryID := uuid.New()

			mockAdvertUseCase.
				EXPECT().
				GetAdvertsByCategoryId(categoryID).
				Return(nil, errors.New("database error"))

			req := httptest.NewRequest("GET", "/api/v1/adverts/category/"+categoryID.String(), nil)
			req = mux.SetURLVars(req, map[string]string{
				"categoryId": categoryID.String(),
			})
			rr := httptest.NewRecorder()

			endpoints.GetAdvertsByCategoryId(rr, req)

			if status := rr.Code; status != http.StatusInternalServerError {
				t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
			}

			var errResp utils.ErrResponse
			if err := parseJSONResponse(rr.Body, &errResp); err != nil {
				t.Errorf("Failed to parse error response: %v", err)
			}
		})
	})

	t.Run("UploadImage", func(t *testing.T) {
		t.Run("Invalid advert ID", func(t *testing.T) {
			req := httptest.NewRequest("PUT", "/api/v1/adverts/invalid-uuid/image", nil)
			req = mux.SetURLVars(req, map[string]string{
				"advertId": "invalid-uuid",
			})
			rr := httptest.NewRecorder()

			endpoints.UploadImage(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
			}

			var errResp utils.ErrResponse
			if err := parseJSONResponse(rr.Body, &errResp); err != nil {
				t.Errorf("Failed to parse error response: %v", err)
			}
		})
	})
}