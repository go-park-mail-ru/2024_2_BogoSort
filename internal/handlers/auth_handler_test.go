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
	authHandler := &AuthHandler{
		UserStorage: userStorage,
	}

	reqBody, _ := json.Marshal(AuthData{Email: "newuser@example.com", Password: "password"})
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authHandler.SignupHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response AuthResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("could not decode response: %v", err)
	}

	if response.Email != "newuser@example.com" {
		t.Errorf("expected email to be newuser@example.com, got %v", response.Email)
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
	authHandler := &AuthHandler{
		UserStorage: userStorage,
	}

	_, err := userStorage.CreateUser("newuser@example.com", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test login with valid body
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

	if response.Email != "newuser@example.com" {
		t.Errorf("expected email to be newuser@example.com, got %v", response.Email)
	}

	// Test login with invalid password
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

	// Test login with invalid method
	req, err = http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}

	// Test login with invalid JSON
	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
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
		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
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
		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
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