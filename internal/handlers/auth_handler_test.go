package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"log"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/responses"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/storage"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/utils"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET_KEY", "your_very_long_and_secure_secret_key_here")

	err := config.Init()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}
	utils.InitJWT()

	code := m.Run()
	os.Exit(code)
}

func TestRegisterHandler(t *testing.T) {
	userStorage := storage.NewUserStorage()
	authHandler := &AuthHandler{
		UserStorage: userStorage,
	}

	reqBody, _ := json.Marshal(AuthData{Email: "newuser@example.com", Password: "Password1!"})
	req, err := http.NewRequest("POST", "/api/v1/signup", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authHandler.SignupHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response responses.AuthResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("could not decode response: %v", err)
	}

	if response.Email != "newuser@example.com" {
		t.Errorf("expected email to be newuser@example.com, got %v", response.Email)
	}

	req, err = http.NewRequest("POST", "/api/v1/signup", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	req, err = http.NewRequest("GET", "/api/v1/signup", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}

	req, err = http.NewRequest("POST", "/api/v1/signup", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestLoginHandler(t *testing.T) {
	userStorage := storage.NewUserStorage()
	authHandler := &AuthHandler{
		UserStorage: userStorage,
	}

	_, err := userStorage.CreateUser("newuser@example.com", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test login with valid body
	reqBody, _ := json.Marshal(LoginCredentials{Email: "newuser@example.com", Password: "password"})
	req, err := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authHandler.LoginHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response responses.AuthResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("could not decode response: %v", err)
	}

	if response.Email != "newuser@example.com" {
		t.Errorf("expected email to be newuser@example.com, got %v", response.Email)
	}

	// Test login with invalid password
	reqBody, _ = json.Marshal(LoginCredentials{Email: "newuser@example.com", Password: "wrongpassword"})
	req, err = http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	// Test login with invalid method
	req, err = http.NewRequest("GET", "/api/v1/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}

	// Test login with invalid JSON
	req, err = http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestSignupHandlerValidation(t *testing.T) {
	userStorage := storage.NewUserStorage()
	authHandler := &AuthHandler{
		UserStorage: userStorage,
	}

	invalidReqBodies := []AuthData{
		{Email: "", Password: "password"},
		{Email: "invalid-email", Password: "password"},
		{Email: "newuser@example.com", Password: ""},
	}

	for _, reqBody := range invalidReqBodies {
		body, _ := json.Marshal(reqBody)
		req, err := http.NewRequest("POST", "/api/v1/signup", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.SignupHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	}
}

func TestLoginHandlerValidation(t *testing.T) {
	userStorage := storage.NewUserStorage()
	authHandler := &AuthHandler{
		UserStorage: userStorage,
	}

	invalidReqBodies := []LoginCredentials{
		{Email: "", Password: "password"},
		{Email: "invalid-email", Password: "password"},
		{Email: "newuser@example.com", Password: ""},
	}

	for _, reqBody := range invalidReqBodies {
		body, _ := json.Marshal(reqBody)
		req, err := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.LoginHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	}
}

func TestLogoutHandler(t *testing.T) {
	userStorage := storage.NewUserStorage()
	authHandler := &AuthHandler{
		UserStorage: userStorage,
	}

	// Test logout with no session cookie
	req, err := http.NewRequest("POST", "/api/v1/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authHandler.LogoutHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	// Test logout with invalid session cookie
	req, err = http.NewRequest("POST", "/api/v1/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "invalid_token",
		Expires:  time.Now().Add(config.GetJWTExpirationTime()),
		HttpOnly: true,
	}
	req.AddCookie(cookie)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestAuthMiddleware(t *testing.T) {

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	ts := httptest.NewServer(AuthMiddleware(testHandler))
	defer ts.Close()

	tests := []struct {
		name           string
		setupRequest   func(req *http.Request)
		expectedStatus int
	}{
		{
			name: "No token",
			setupRequest: func(req *http.Request) {
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid token in cookie",
			setupRequest: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "session_id",
					Value: "invalid_token",
				})
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid token in header",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer invalid_token")
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Valid token in cookie",
			setupRequest: func(req *http.Request) {
				token, _ := utils.CreateToken("test@example.com")
				req.AddCookie(&http.Cookie{
					Name:    "session_id",
					Value:   token,
					Expires: time.Now().Add(config.GetJWTExpirationTime()),
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Valid token in header",
			setupRequest: func(req *http.Request) {
				token, _ := utils.CreateToken("test@example.com")
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", ts.URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			tt.setupRequest(req)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}