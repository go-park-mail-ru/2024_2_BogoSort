package http

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/mocks"
// 	"github.com/golang/mock/gomock"
// 	"github.com/google/uuid"
// 	// "github.com/gorilla/mux"
// 	// "github.com/microcosm-cc/bluemonday"
// 	// "go.uber.org/zap"
// )

// func setupUserEndpoint(t *testing.T) (*UserEndpoint, *mocks.MockUser, *mocks.MockAuth, *gomock.Controller) {
// 	// ctrl := gomock.NewController(t)
// 	// mockUserUC := mocks.NewMockUser(ctrl)
// 	// mockAuthUC := mocks.NewMockAuth(ctrl)
// 	// logger, _ := zap.NewDevelopment()
// 	// policy := bluemonday.UGCPolicy()

// 	// sessionManager := &utils.SessionManager{
// 	// 	SessionUC:        mockAuthUC,
// 	// 	SessionAliveTime: 1,
// 	// 	SecureCookie:     false,
// 	// 	Logger:           logger,
// 	// }

// 	// endpoints := NewUserEndpoints(mockUserUC, mockAuthUC, sessionManager, nil, logger, policy)
// 	// return endpoints, mockUserUC, mockAuthUC, ctrl
// 	return nil, nil, nil, nil
// }

// func TestUserEndpoint_Signup(t *testing.T) {
// 	endpoints, mockUserUC, mockAuthUC, ctrl := setupUserEndpoint(t)
// 	defer ctrl.Finish()

// 	t.Run("Success", func(t *testing.T) {
// 		credentials := dto.Signup{
// 			Email:    "test@example.com",
// 			Password: "password123",
// 		}
// 		userID := uuid.New()

// 		mockUserUC.EXPECT().Signup(&credentials).Return(userID, nil)
// 		mockAuthUC.EXPECT().CreateSession(userID).Return("session-id", nil)

// 		body, _ := json.Marshal(credentials)
// 		req := httptest.NewRequest("POST", "/api/v1/signup", bytes.NewBuffer(body))
// 		rr := httptest.NewRecorder()

// 		endpoints.Signup(rr, req)

// 		if status := rr.Code; status != http.StatusOK {
// 			t.Errorf("Expected status %v, got %v", http.StatusOK, status)
// 		}
// 	})

// 	t.Run("Invalid request body", func(t *testing.T) {
// 		req := httptest.NewRequest("POST", "/api/v1/signup", bytes.NewBuffer([]byte("invalid json")))
// 		rr := httptest.NewRecorder()

// 		endpoints.Signup(rr, req)

// 		if status := rr.Code; status != http.StatusBadRequest {
// 			t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
// 		}
// 	})
// }

// func TestUserEndpoint_Login(t *testing.T) {
// 	endpoints, mockUserUC, mockAuthUC, ctrl := setupUserEndpoint(t)
// 	defer ctrl.Finish()

// 	t.Run("Success", func(t *testing.T) {
// 		credentials := dto.Login{
// 			Email:    "test@example.com",
// 			Password: "password123",
// 		}
// 		userID := uuid.New()

// 		mockUserUC.EXPECT().Login(&credentials).Return(userID, nil)
// 		mockAuthUC.EXPECT().CreateSession(userID).Return("session-id", nil)

// 		body, _ := json.Marshal(credentials)
// 		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
// 		rr := httptest.NewRecorder()

// 		endpoints.Login(rr, req)

// 		if status := rr.Code; status != http.StatusOK {
// 			t.Errorf("Expected status %v, got %v", http.StatusOK, status)
// 		}
// 	})

// 	t.Run("Invalid credentials", func(t *testing.T) {
// 		credentials := dto.Login{
// 			Email:    "test@example.com",
// 			Password: "wrongpassword",
// 		}

// 		mockUserUC.EXPECT().Login(&credentials).Return(uuid.Nil, ErrInvalidCredentials)

// 		body, _ := json.Marshal(credentials)
// 		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
// 		rr := httptest.NewRecorder()

// 		endpoints.Login(rr, req)

// 		if status := rr.Code; status != http.StatusInternalServerError {
// 			t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
// 		}
// 	})
// }

// func TestUserEndpoints_GetProfile(t *testing.T) {
// 	endpoints, mockUserUC, _, ctrl := setupUserEndpoints(t)
// 	defer ctrl.Finish()

// 	t.Run("Success", func(t *testing.T) {
// 		userID := uuid.New()
// 		user := dto.User{
// 			ID:       userID,
// 			Email:    "test@example.com",
// 			Username: "testuser",
// 		}

// 		mockUserUC.EXPECT().GetUser(userID).Return(&user, nil)

// 		req := httptest.NewRequest("GET", "/api/v1/profile/"+userID.String(), nil)
// 		req = mux.SetURLVars(req, map[string]string{
// 			"user_id": userID.String(),
// 		})
// 		rr := httptest.NewRecorder()

// 		endpoints.GetProfile(rr, req)

// 		if status := rr.Code; status != http.StatusOK {
// 			t.Errorf("Expected status %v, got %v", http.StatusOK, status)
// 		}
// 	})

// 	t.Run("User not found", func(t *testing.T) {
// 		userID := uuid.New()

// 		mockUserUC.EXPECT().GetUser(userID).Return(nil, ErrUserNotFound)

// 		req := httptest.NewRequest("GET", "/api/v1/profile/"+userID.String(), nil)
// 		req = mux.SetURLVars(req, map[string]string{
// 			"user_id": userID.String(),
// 		})
// 		rr := httptest.NewRecorder()

// 		endpoints.GetProfile(rr, req)

// 		if status := rr.Code; status != http.StatusInternalServerError {
// 			t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
// 		}
// 	})
// }

// func TestUserEndpoints_UpdateProfile(t *testing.T) {
// 	endpoints, mockUserUC, mockAuthUC, ctrl := setupUserEndpoints(t)
// 	defer ctrl.Finish()

// 	t.Run("Success", func(t *testing.T) {
// 		userID := uuid.New()
// 		user := dto.User{
// 			ID:       userID,
// 			Email:    "test@example.com",
// 			Username: "testuser",
// 		}

// 		mockAuthUC.EXPECT().GetUserIdBySession("session-id").Return(userID, nil)

// 		mockUserUC.EXPECT().UpdateInfo(&user).Return(nil)

// 		body, _ := json.Marshal(user)
// 		req := httptest.NewRequest("PUT", "/api/v1/profile", bytes.NewBuffer(body))
// 		req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-id"})
// 		rr := httptest.NewRecorder()

// 		endpoints.UpdateProfile(rr, req)

// 		if status := rr.Code; status != http.StatusOK {
// 			t.Errorf("Expected status %v, got %v", http.StatusOK, status)
// 		}
// 	})

// 	t.Run("Internal server error", func(t *testing.T) {
// 		user := dto.User{
// 			Email:    "test@example.com",
// 			Username: "testuser",
// 		}

// 		body, _ := json.Marshal(user)
// 		req := httptest.NewRequest("PUT", "/api/v1/profile", bytes.NewBuffer(body))
// 		rr := httptest.NewRecorder()

// 		endpoints.UpdateProfile(rr, req)

// 		if status := rr.Code; status != http.StatusInternalServerError {
// 			t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
// 		}
// 	})
// }

// func TestUserEndpoints_ChangePassword(t *testing.T) {
// 	endpoints, mockUserUC, mockAuthUC, ctrl := setupUserEndpoints(t)
// 	defer ctrl.Finish()

// 	t.Run("Success", func(t *testing.T) {
// 		userID := uuid.New()
// 		updatePassword := dto.UpdatePassword{
// 			OldPassword: "oldpassword",
// 			NewPassword: "newpassword",
// 		}

// 		mockAuthUC.EXPECT().GetUserIdBySession("session-id").Return(userID, nil)

// 		mockUserUC.EXPECT().ChangePassword(userID, &updatePassword).Return(nil)

// 		body, _ := json.Marshal(updatePassword)
// 		req := httptest.NewRequest("POST", "/api/v1/password", bytes.NewBuffer(body))
// 		req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-id"})
// 		rr := httptest.NewRecorder()

// 		endpoints.ChangePassword(rr, req)

// 		if status := rr.Code; status != http.StatusOK {
// 			t.Errorf("Expected status %v, got %v", http.StatusOK, status)
// 		}
// 	})

// 	t.Run("Unauthorized", func(t *testing.T) {
// 		updatePassword := dto.UpdatePassword{
// 			OldPassword: "oldpassword",
// 			NewPassword: "newpassword",
// 		}

// 		body, _ := json.Marshal(updatePassword)
// 		req := httptest.NewRequest("POST", "/api/v1/password", bytes.NewBuffer(body))
// 		rr := httptest.NewRecorder()

// 		endpoints.ChangePassword(rr, req)

// 		if status := rr.Code; status != http.StatusUnauthorized {
// 			t.Errorf("Expected status %v, got %v", http.StatusUnauthorized, status)
// 		}
// 	})
// }
