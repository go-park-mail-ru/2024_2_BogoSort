package handlers

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"
// 	"time"

// 	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/responses"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/storage"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/utils"
// )

// func TestMain(m *testing.M) {
// 	os.Setenv("JWT_SECRET_KEY", "your_very_long_and_secure_secret_key_here")
// 	os.Setenv("JWT_EXPIRATION_TIME", "1h")
// 	os.Setenv("JWT_ISSUER", "test_issuer")

// 	config.InitFromEnv()
// 	utils.InitJWT()

// 	code := m.Run()
// 	os.Exit(code)
// }

// func TestSignupHandler(t *testing.T) {
// 	userStorage := storage.NewUserStorage()
// 	authHandler := &AuthHandler{UserStorage: userStorage}

// 	t.Run("Successful signup", func(t *testing.T) {
// 		testSuccessfulSignup(t, authHandler)
// 	})

// 	t.Run("Duplicate signup", func(t *testing.T) {
// 		testDuplicateSignup(t, authHandler)
// 	})

// 	t.Run("Method not allowed", func(t *testing.T) {
// 		testMethodNotAllowed(t, authHandler)
// 	})

// 	t.Run("Invalid JSON", func(t *testing.T) {
// 		testInvalidJSON(t, authHandler)
// 	})
// }

// func testSuccessfulSignup(t *testing.T, authHandler *AuthHandler) {
// 	t.Helper()

// 	reqBody, err := json.Marshal(AuthCredentials{Email: "newuser@example.com", Password: "Password1!"})

// 	if err != nil {
// 		t.Fatalf("failed to marshal request body: %v", err)
// 	}

// 	req, err := http.NewRequest(http.MethodPost, "/api/v1/signup", bytes.NewBuffer(reqBody))
// 	if err != nil {
// 		t.Fatalf("failed to create request: %v", err)
// 	}

// 	rr := httptest.NewRecorder()
// 	http.HandlerFunc(authHandler.SignupHandler).ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusCreated {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
// 	}

// 	var response responses.AuthResponse
// 	err = json.NewDecoder(rr.Body).Decode(&response)

// 	if err != nil {
// 		t.Errorf("failed to decode response: %v", err)
// 	}

// 	if response.Email != "newuser@example.com" {
// 		t.Errorf("expected email to be newuser@example.com, got %v", response.Email)
// 	}
// }

// func testDuplicateSignup(t *testing.T, authHandler *AuthHandler) {
// 	t.Helper()

// 	reqBody, err := json.Marshal(AuthCredentials{Email: "newuser@example.com", Password: "Password1!"})
// 	if err != nil {
// 		t.Fatalf("failed to marshal request body: %v", err)
// 	}

// 	req, err := http.NewRequest(http.MethodPost, "/api/v1/signup", bytes.NewBuffer(reqBody))

// 	if err != nil {
// 		t.Fatalf("failed to create request: %v", err)
// 	}

// 	rr := httptest.NewRecorder()

// 	http.HandlerFunc(authHandler.SignupHandler).ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusBadRequest {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
// 	}
// }

// func testMethodNotAllowed(t *testing.T, authHandler *AuthHandler) {
// 	t.Helper()

// 	req, err := http.NewRequest(http.MethodGet, "/api/v1/signup", nil)

// 	if err != nil {
// 		t.Fatalf("failed to create request: %v", err)
// 	}

// 	rr := httptest.NewRecorder()

// 	http.HandlerFunc(authHandler.SignupHandler).ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusMethodNotAllowed {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
// 	}
// }

// func testInvalidJSON(t *testing.T, authHandler *AuthHandler) {
// 	t.Helper()

// 	req, err := http.NewRequest(http.MethodPost, "/api/v1/signup", bytes.NewBufferString("invalid json"))

// 	if err != nil {
// 		t.Fatalf("failed to create request: %v", err)
// 	}

// 	rr := httptest.NewRecorder()

// 	http.HandlerFunc(authHandler.SignupHandler).ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusBadRequest {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
// 	}
// }

// func TestLoginHandler(t *testing.T) {
// 	userStorage := storage.NewUserStorage()
// 	authHandler := &AuthHandler{
// 		UserStorage: userStorage,
// 	}

// 	_, err := userStorage.CreateUser("newuser@example.com", "password")
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	reqBody, err := json.Marshal(AuthCredentials{Email: "newuser@example.com", Password: "password"})
// 	if err != nil {
// 		t.Fatalf("failed to marshal request body: %v", err)
// 	}

// 	req, err := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBody))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(authHandler.LoginHandler)
// 	handler.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	var response responses.AuthResponse
// 	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
// 		t.Errorf("could not decode response: %v", err)
// 	}

// 	if response.Email != "newuser@example.com" {
// 		t.Errorf("expected email to be newuser@example.com, got %v", response.Email)
// 	}

// 	reqBody, err = json.Marshal(AuthCredentials{Email: "newuser@example.com", Password: "wrongpassword"})
// 	if err != nil {
// 		t.Fatalf("failed to marshal request body: %v", err)
// 	}

// 	req, err = http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBody))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr = httptest.NewRecorder()
// 	handler.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusBadRequest {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
// 	}

// 	// Test login with invalid method
// 	req, err = http.NewRequest(http.MethodGet, "/api/v1/login", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr = httptest.NewRecorder()
// 	handler.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusMethodNotAllowed {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
// 	}

// 	// Test login with invalid JSON
// 	req, err = http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString("invalid json"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr = httptest.NewRecorder()

// 	handler.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusBadRequest {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
// 	}
// }

// func TestSignupHandlerValidation(t *testing.T) {
// 	userStorage := storage.NewUserStorage()
// 	authHandler := &AuthHandler{
// 		UserStorage: userStorage,
// 	}

// 	invalidReqBodies := []AuthCredentials{
// 		{Email: "", Password: "password"},
// 		{Email: "invalid-email", Password: "password"},
// 		{Email: "newuser@example.com", Password: ""},
// 	}

// 	for _, reqBody := range invalidReqBodies {
// 		body, err := json.Marshal(reqBody)
// 		if err != nil {
// 			t.Fatalf("failed to marshal request body: %v", err)
// 		}

// 		req, err := http.NewRequest(http.MethodPost, "/api/v1/signup", bytes.NewBuffer(body))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(authHandler.SignupHandler)
// 		handler.ServeHTTP(rr, req)

// 		if status := rr.Code; status != http.StatusBadRequest {
// 			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
// 		}
// 	}
// }

// func TestLoginHandlerValidation(t *testing.T) {
// 	userStorage := storage.NewUserStorage()
// 	authHandler := &AuthHandler{
// 		UserStorage: userStorage,
// 	}

// 	invalidReqBodies := []AuthCredentials{
// 		{Email: "", Password: "password"},
// 		{Email: "invalid-email", Password: "password"},
// 		{Email: "newuser@example.com", Password: ""},
// 	}

// 	for _, reqBody := range invalidReqBodies {
// 		body, err := json.Marshal(reqBody)
// 		if err != nil {
// 			t.Fatalf("failed to marshal request body: %v", err)
// 		}

// 		req, err := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(body))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(authHandler.LoginHandler)
// 		handler.ServeHTTP(rr, req)

// 		if status := rr.Code; status != http.StatusBadRequest {
// 			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
// 		}
// 	}
// }

// func TestLogoutHandler(t *testing.T) {
// 	userStorage := storage.NewUserStorage()
// 	authHandler := &AuthHandler{
// 		UserStorage: userStorage,
// 	}

// 	// Test logout with no session cookie
// 	req, err := http.NewRequest(http.MethodPost, "/api/v1/logout", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(authHandler.LogoutHandler)
// 	handler.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusUnauthorized {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
// 	}

// 	// Test logout with invalid session cookie
// 	req, err = http.NewRequest(http.MethodPost, "/api/v1/logout", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	cookie := &http.Cookie{
// 		Name:     "session_id",
// 		Value:    "invalid_token",
// 		Expires:  time.Now().Add(config.GetJWTExpirationTime()),
// 		HttpOnly: true,
// 	}
// 	req.AddCookie(cookie)

// 	rr = httptest.NewRecorder()
// 	handler.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusUnauthorized {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
// 	}
// }
