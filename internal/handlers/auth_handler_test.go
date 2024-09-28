package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/storage"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/utils"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET_KEY", "your_very_long_and_secure_secret_key_here")

	config.Init()
	utils.InitJWT()

	code := m.Run()
	os.Exit(code)
}

func TestRegisterHandler(t *testing.T) {
	userStorage := storage.NewUserStorage()
	sessionStorage := storage.NewSessionStorage()
	authHandler := &AuthHandler{
		UserStorage:    userStorage,
		SessionStorage: sessionStorage,
	}

	reqBody, _ := json.Marshal(AuthData{Email: "newuser@example.com", Password: "password"})
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authHandler.RegisterHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response AuthResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("could not decode response: %v", err)
	}

	if response.User.Email != "newuser@example.com" {
		t.Errorf("expected email to be newuser@example.com, got %v", response.User.Email)
	}

	req, err = http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	req, err = http.NewRequest("GET", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}

	req, err = http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("invalid json")))
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
	sessionStorage := storage.NewSessionStorage()
	authHandler := &AuthHandler{
		UserStorage:    userStorage,
		SessionStorage: sessionStorage,
	}

	_, err := userStorage.CreateUser("newuser@example.com", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	reqBody, _ := json.Marshal(LoginCredentials{Email: "newuser@example.com", Password: "password"})
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authHandler.LoginHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response AuthResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("could not decode response: %v", err)
	}

	if response.User.Email != "newuser@example.com" {
		t.Errorf("expected email to be newuser@example.com, got %v", response.User.Email)
	}

	reqBody, _ = json.Marshal(LoginCredentials{Email: "newuser@example.com", Password: "wrongpassword"})
	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	req, err = http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}

	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	reqBody, _ = json.Marshal(LoginCredentials{Email: "newuser@example.com", Password: "password"})
	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestLogoutHandler(t *testing.T) {
    userStorage := storage.NewUserStorage()
    sessionStorage := storage.NewSessionStorage()
    authHandler := 	&AuthHandler{
        UserStorage:    userStorage,
        SessionStorage: sessionStorage,
    }
    req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
    handler := http.HandlerFunc(authHandler.LogoutHandler)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusUnauthorized {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
    }

    req, err = http.NewRequest("POST", "/logout", nil)
    if err != nil {
        t.Fatal(err)
    }

    invalidToken := "invalid.token.string"
    req.AddCookie(&http.Cookie{Name: "session_id", Value: invalidToken})

    rr = httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusUnauthorized {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
    }

    req, err = http.NewRequest("POST", "/logout", nil)
    if err != nil {
        t.Fatal(err)
    }

    newToken, err := utils.CreateToken("newuser@example.com")
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    req.AddCookie(&http.Cookie{Name: "session_id", Value: newToken})

    rr = httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusUnauthorized {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
    }

    req, err = http.NewRequest("GET", "/logout", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr = httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusMethodNotAllowed {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
    }
}
