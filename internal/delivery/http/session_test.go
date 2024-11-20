package http

// import (
// 	"errors"
// 	// "net/http"
// 	// "net/http/httptest"
// 	"testing"

// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/mocks"
// 	"github.com/golang/mock/gomock"
// 	"github.com/google/uuid"
// 	"go.uber.org/zap"
// )

// func setupSessionManager(t *testing.T) (*utils.SessionManager, *mocks.MockAuth, *gomock.Controller) {
// 	ctrl := gomock.NewController(t)
// 	mockAuthUC := mocks.NewMockAuth(ctrl)
// 	logger := zap.NewNop()

// 	sessionManager := utils.NewSessionManager(mockAuthUC, 10, true, logger)
// 	return sessionManager, mockAuthUC, ctrl
// }

// func TestSessionManager_CreateSession(t *testing.T) {
// 	sessionManager, mockAuthUC, ctrl := setupSessionManager(t)
// 	defer ctrl.Finish()

// 	t.Run("Success", func(t *testing.T) {
// 		userID := uuid.New()
// 		sessionID := uuid.New().String()

// 		mockAuthUC.
// 			EXPECT().
// 			CreateSession(userID).
// 			Return(sessionID, nil)

// 		gotSessionID, err := sessionManager.CreateSession(userID)
// 		if err != nil {
// 			t.Errorf("Unexpected error: %v", err)
// 		}
// 		if gotSessionID != sessionID {
// 			t.Errorf("Expected session ID %v, got %v", sessionID, gotSessionID)
// 		}
// 	})

// 	t.Run("Error", func(t *testing.T) {
// 		userID := uuid.New()

// 		mockAuthUC.
// 			EXPECT().
// 			CreateSession(userID).
// 			Return("", errors.New("create session error"))

// 		_, err := sessionManager.CreateSession(userID)
// 		if err == nil {
// 			t.Error("Expected error, got nil")
// 		}
// 	})
// }

// // func TestSessionManager_GetUserID(t *testing.T) {
// // 	sessionManager, mockAuthUC, ctrl := setupSessionManager(t)
// // 	defer ctrl.Finish()

// // 	t.Run("Success", func(t *testing.T) {
// // 		userID := uuid.New()
// // 		sessionID := uuid.New().String()

// // 		mockAuthUC.
// // 			EXPECT().
// // 			GetUserIdBySession(sessionID).
// // 			Return(userID, nil)

// // 		req := httptest.NewRequest("GET", "/", nil)
// // 		req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})

// // 		gotUserID, err := sessionManager.GetUserID(req)
// // 		if err != nil {
// // 			t.Errorf("Unexpected error: %v", err)
// // 		}
// // 		if gotUserID != userID {
// // 			t.Errorf("Expected user ID %v, got %v", userID, gotUserID)
// // 		}
// // 	})

// // 	t.Run("Missing cookie", func(t *testing.T) {
// // 		req := httptest.NewRequest("GET", "/", nil)

// // 		_, err := sessionManager.GetUserID(req)
// // 		if err == nil {
// // 			t.Error("Expected error, got nil")
// // 		}
// // 	})

// // 	t.Run("Error", func(t *testing.T) {
// // 		sessionID := uuid.New().String()

// // 		mockAuthUC.
// // 			EXPECT().
// // 			GetUserIdBySession(sessionID).
// // 			Return(uuid.Nil, errors.New("get user ID error"))

// // 		req := httptest.NewRequest("GET", "/", nil)
// // 		req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})

// // 		_, err := sessionManager.GetUserID(req)
// // 		if err == nil {
// // 			t.Error("Expected error, got nil")
// // 		}
// // 	})
// // }

// func TestSessionManager_DeleteSession(t *testing.T) {
// 	sessionManager, mockAuthUC, ctrl := setupSessionManager(t)
// 	defer ctrl.Finish()

// 	t.Run("Success", func(t *testing.T) {
// 		sessionID := uuid.New().String()

// 		mockAuthUC.
// 			EXPECT().
// 			Logout(sessionID).
// 			Return(nil)

// 		err := sessionManager.DeleteSession(sessionID)
// 		if err != nil {
// 			t.Errorf("Unexpected error: %v", err)
// 		}
// 	})

// 	t.Run("Error", func(t *testing.T) {
// 		sessionID := uuid.New().String()

// 		mockAuthUC.
// 			EXPECT().
// 			Logout(sessionID).
// 			Return(errors.New("logout error"))

// 		err := sessionManager.DeleteSession(sessionID)
// 		if err == nil {
// 			t.Error("Expected error, got nil")
// 		}
// 	})
// }
