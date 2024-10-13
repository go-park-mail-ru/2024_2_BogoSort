package delivery

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("SESSION_EXPIRATION_TIME", "1h")
	config.InitFromEnv()

	code := m.Run()
	os.Exit(code)
}

func TestSignupHandler(t *testing.T) {
	userStorage := storage.NewUserStorage()
	sessionStorage := storage.NewSessionStorage()
	authHandler := &AuthHandler{
		UserStorage:    userStorage,
		SessionStorage: sessionStorage,
	}

	t.Run("Successful signup", func(t *testing.T) {
		testSuccessfulSignup(t, authHandler)
	})

	t.Run("Duplicate signup", func(t *testing.T) {
		testDuplicateSignup(t, authHandler)
	})

	t.Run("Method not allowed", func(t *testing.T) {
		testMethodNotAllowed(t, authHandler)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		testInvalidJSON(t, authHandler)
	})

	t.Run("Invalid password", func(t *testing.T) {
		testInvalidPassword(t, authHandler)
	})
}

func testSuccessfulSignup(t *testing.T, authHandler *AuthHandler) {
	t.Helper()

	reqBody, err := json.Marshal(AuthData{Email: "newuser@example.com", Password: "Password1!"})
	assert.NoError(t, err, "failed to marshal request body")

	req, err := http.NewRequest(http.MethodPost, "/api/v1/signup", bytes.NewBuffer(reqBody))
	assert.NoError(t, err, "failed to create request")

	rr := httptest.NewRecorder()
	http.HandlerFunc(authHandler.SignupHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "handler returned wrong status code")

	var response responses.AuthResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err, "failed to decode response")

	assert.Equal(t, "newuser@example.com", response.Email, "expected email to be newuser@example.com")
	assert.NotEmpty(t, response.SessionID, "expected session ID to be set")
}

func testDuplicateSignup(t *testing.T, authHandler *AuthHandler) {
	t.Helper()

	reqBody, err := json.Marshal(AuthData{Email: "newuser@example.com", Password: "Password1!"})
	assert.NoError(t, err, "failed to marshal request body")

	req, err := http.NewRequest(http.MethodPost, "/api/v1/signup", bytes.NewBuffer(reqBody))
	assert.NoError(t, err, "failed to create request")

	rr := httptest.NewRecorder()
	http.HandlerFunc(authHandler.SignupHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
	assert.Contains(t, rr.Body.String(), ErrUserAlreadyExists.Error(), "expected user already exists error")
}

func testMethodNotAllowed(t *testing.T, authHandler *AuthHandler) {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, "/api/v1/signup", nil)
	assert.NoError(t, err, "failed to create request")

	rr := httptest.NewRecorder()
	http.HandlerFunc(authHandler.SignupHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "handler returned wrong status code")
	assert.Contains(t, rr.Body.String(), ErrMethodNotAllowed.Error(), "expected method not allowed error")
}

func testInvalidJSON(t *testing.T, authHandler *AuthHandler) {
	t.Helper()

	req, err := http.NewRequest(http.MethodPost, "/api/v1/signup", bytes.NewBufferString("invalid json"))
	assert.NoError(t, err, "failed to create request")

	rr := httptest.NewRecorder()
	http.HandlerFunc(authHandler.SignupHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
	assert.Contains(t, rr.Body.String(), ErrInvalidRequestBody.Error(), "expected invalid request body error")
}

func testInvalidPassword(t *testing.T, authHandler *AuthHandler) {
	t.Helper()

	reqBody, err := json.Marshal(AuthData{Email: "newuser@example.com", Password: "short"})
	assert.NoError(t, err, "failed to marshal request body")

	req, err := http.NewRequest(http.MethodPost, "/api/v1/signup", bytes.NewBuffer(reqBody))
	assert.NoError(t, err, "failed to create request")

	rr := httptest.NewRecorder()
	http.HandlerFunc(authHandler.SignupHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
	assert.Contains(t, rr.Body.String(), "password is too short", "expected password validation error")
}

func TestLoginHandler(t *testing.T) {
	userStorage := storage.NewUserStorage()
	sessionStorage := storage.NewSessionStorage()
	authHandler := &AuthHandler{
		UserStorage:    userStorage,
		SessionStorage: sessionStorage,
	}

	_, err := userStorage.CreateUser("newuser@example.com", "Password1!")
	assert.NoError(t, err, "expected no error")

	t.Run("Successful login", func(t *testing.T) {
		testSuccessfulLogin(t, authHandler)
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		testInvalidCredentials(t, authHandler)
	})

	t.Run("Method not allowed", func(t *testing.T) {
		testLoginMethodNotAllowed(t, authHandler)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		testLoginInvalidJSON(t, authHandler)
	})
}

func testSuccessfulLogin(t *testing.T, authHandler *AuthHandler) {
	t.Helper()

	reqBody, err := json.Marshal(LoginCredentials{Email: "newuser@example.com", Password: "Password1!"})
	assert.NoError(t, err, "failed to marshal request body")

	req, err := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBody))
	assert.NoError(t, err, "failed to create request")

	rr := httptest.NewRecorder()
	http.HandlerFunc(authHandler.LoginHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	var response responses.AuthResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err, "could not decode response")

	assert.Equal(t, "newuser@example.com", response.Email, "expected email to be newuser@example.com")
	assert.NotEmpty(t, response.SessionID, "expected session ID to be set")
}

func testInvalidCredentials(t *testing.T, authHandler *AuthHandler) {
	t.Helper()

	reqBody, err := json.Marshal(LoginCredentials{Email: "newuser@example.com", Password: "wrongpassword"})
	assert.NoError(t, err, "failed to marshal request body")

	req, err := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBody))
	assert.NoError(t, err, "failed to create request")

	rr := httptest.NewRecorder()
	http.HandlerFunc(authHandler.LoginHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
	assert.Contains(t, rr.Body.String(), ErrInvalidCredentials.Error(), "expected invalid credentials error")
}

func testLoginMethodNotAllowed(t *testing.T, authHandler *AuthHandler) {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, "/api/v1/login", nil)
	assert.NoError(t, err, "failed to create request")

	rr := httptest.NewRecorder()
	http.HandlerFunc(authHandler.LoginHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "handler returned wrong status code")
	assert.Contains(t, rr.Body.String(), ErrMethodNotAllowed.Error(), "expected method not allowed error")
}

func testLoginInvalidJSON(t *testing.T, authHandler *AuthHandler) {
	t.Helper()

	req, err := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString("invalid json"))
	assert.NoError(t, err, "failed to create request")

	rr := httptest.NewRecorder()
	http.HandlerFunc(authHandler.LoginHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
	assert.Contains(t, rr.Body.String(), ErrInvalidRequestBody.Error(), "expected invalid request body error")
}

func TestLogoutHandler(t *testing.T) {
	userStorage := storage.NewUserStorage()
	sessionStorage := storage.NewSessionStorage()
	authHandler := &AuthHandler{
		UserStorage:    userStorage,
		SessionStorage: sessionStorage,
	}

	t.Run("Logout without session", func(t *testing.T) {
		testLogoutWithoutSession(t, authHandler)
	})

	t.Run("Logout with invalid session", func(t *testing.T) {
		testLogoutWithInvalidSession(t, authHandler)
	})

	t.Run("Successful logout", func(t *testing.T) {
		testSuccessfulLogout(t, authHandler)
	})
}

func testLogoutWithoutSession(t *testing.T, authHandler *AuthHandler) {
	t.Helper()

	req, err := http.NewRequest(http.MethodPost, "/api/v1/logout", nil)
	assert.NoError(t, err, "failed to create request")

	rr := httptest.NewRecorder()
	http.HandlerFunc(authHandler.LogoutHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code, "handler returned wrong status code")
	assert.Contains(t, rr.Body.String(), ErrNoActiveSession.Error(), "expected no active session error")
}

func testLogoutWithInvalidSession(t *testing.T, authHandler *AuthHandler) {
	t.Helper()

	req, err := http.NewRequest(http.MethodPost, "/api/v1/logout", nil)
	assert.NoError(t, err, "failed to create request")

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "invalid_token",
		Expires:  time.Now().Add(config.GetSessionExpirationTime()),
		HttpOnly: true,
	}
	req.AddCookie(cookie)

	rr := httptest.NewRecorder()
	http.HandlerFunc(authHandler.LogoutHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code, "handler returned wrong status code")
	assert.Contains(t, rr.Body.String(), ErrSessionDoesNotExist.Error(), "expected session does not exist error")
}

func testSuccessfulLogout(t *testing.T, authHandler *AuthHandler) {
	t.Helper()

	// Create a user and login to get a valid session
	user, err := authHandler.UserStorage.CreateUser("logoutuser@example.com", "Password1!")
	assert.NoError(t, err, "failed to create user")

	sessionID := authHandler.SessionStorage.AddSession(user.Email)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/logout", nil)
	assert.NoError(t, err, "failed to create request")

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(config.GetSessionExpirationTime()),
		HttpOnly: true,
	}
	req.AddCookie(cookie)

	rr := httptest.NewRecorder()
	http.HandlerFunc(authHandler.LogoutHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
}
