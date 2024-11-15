package http

// import (
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/auth"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/mocks"
// 	"github.com/golang/mock/gomock"
// 	"github.com/google/uuid"
// 	"go.uber.org/zap"
// )

// func setupAuthEndpoints(t *testing.T) (*AuthEndpoints, *mocks.MockAuth, *utils.SessionManager, *gomock.Controller) {
// 	ctrl := gomock.NewController(t)
// 	mockAuthUC := mocks.NewMockAuth(ctrl)
// 	grpcClient, err := auth.NewGrpcClient("localhost:50051")
// 	if err != nil {
// 		t.Fatalf("Failed to create gRPC client: %v", err)
// 	}
// 	logger := zap.NewNop()

// 	return NewAuthEndpoints(mockAuthUC, grpcClient, logger), mockAuthUC, nil, ctrl
// }

// func TestAuthEndpoints_Logout(t *testing.T) {
// 	endpoints, mockAuthUC, sessionManager, ctrl := setupAuthEndpoints(t)
// 	defer ctrl.Finish()

// 	t.Run("Success", func(t *testing.T) {
// 		userID := uuid.New()
// 		sessionID := uuid.New().String()

// 		mockAuthUC.
// 			EXPECT().
// 			CreateSession(userID).
// 			Return(sessionID, nil)

// 		mockAuthUC.
// 			EXPECT().
// 			GetUserIdBySession(sessionID).
// 			Return(userID, nil)

// 		mockAuthUC.
// 			EXPECT().
// 			Logout(sessionID).
// 			Return(nil)

// 		// Создаем сессию через вызов CreateSession
// 		sessionManager.CreateSession(userID)

// 		req := httptest.NewRequest("POST", "/api/v1/logout", nil)
// 		req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})
// 		rr := httptest.NewRecorder()

// 		endpoints.Logout(rr, req)

// 		if status := rr.Code; status != http.StatusOK {
// 			t.Errorf("Expected status %v, got %v", http.StatusOK, status)
// 		}

// 		expectedBody := `"You have successfully logged out"`
// 		actualBody := strings.TrimSpace(rr.Body.String())
// 		if actualBody != expectedBody {
// 			t.Errorf("Expected body %v, got %v", expectedBody, actualBody)
// 		}
// 	})

// 	t.Run("Missing session cookie", func(t *testing.T) {
// 		req := httptest.NewRequest("POST", "/api/v1/logout", nil)
// 		rr := httptest.NewRecorder()

// 		endpoints.Logout(rr, req)

// 		if status := rr.Code; status != http.StatusInternalServerError {
// 			t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
// 		}
// 	})

// 	t.Run("Logout error", func(t *testing.T) {
// 		userID := uuid.New()
// 		sessionID := uuid.New().String()

// 		mockAuthUC.
// 			EXPECT().
// 			CreateSession(userID).
// 			Return(sessionID, nil)

// 		mockAuthUC.
// 			EXPECT().
// 			GetUserIdBySession(sessionID).
// 			Return(userID, nil)

// 		mockAuthUC.
// 			EXPECT().
// 			Logout(sessionID).
// 			Return(errors.New("logout error"))

// 		// Создаем сессию через вызов CreateSession
// 		sessionManager.CreateSession(userID)

// 		req := httptest.NewRequest("POST", "/api/v1/logout", nil)
// 		req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})
// 		rr := httptest.NewRecorder()

// 		endpoints.Logout(rr, req)

// 		if status := rr.Code; status != http.StatusInternalServerError {
// 			t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
// 		}
// 	})
// }
