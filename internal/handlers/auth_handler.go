package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/storage"
	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/utils"
)

type AuthData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandler struct {
	UserStorage    *storage.UserStorage
	SessionStorage *storage.SessionStorage
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  storage.User `json:"user"`
}

type AuthErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

// RegisterHandler godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body AuthData true "User credentials"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} AuthErrResponse
// @Failure 405 {object} AuthErrResponse
// @Failure 500 {object} AuthErrResponse
// @Router /signup [post]
func (ah *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var credentials AuthData
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := ah.UserStorage.CreateUser(credentials.Email, credentials.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			sendErrorResponse(w, http.StatusBadRequest, "User already exists")
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		}
		return
	}

	tokenString, err := utils.CreateToken(user.Email)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    tokenString,
		Expires:  time.Now().Add(config.GetJWTExpirationTime()),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	sendJSONResponse(w, http.StatusCreated, AuthResponse{Token: tokenString, User: *user})
}

// LoginHandler godoc
// @Summary Login a user
// @Description Login a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginCredentials true "User credentials"
// @Success

func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var credentials LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := ah.UserStorage.ValidateUserByEmailAndPassword(credentials.Email, credentials.Password)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if ah.SessionStorage.SessionExists(user.Email) {
		sendErrorResponse(w, http.StatusBadRequest, "User already authenticated")
		return
	}

	tokenString, err := utils.CreateToken(user.Email)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    tokenString,
		Expires:  time.Now().Add(config.GetJWTExpirationTime()),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	ah.SessionStorage.AddSession(user.Email, tokenString)

	sendJSONResponse(w, http.StatusOK, AuthResponse{Token: tokenString, User: *user})
}

func (ah *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			sendErrorResponse(w, http.StatusUnauthorized, "No active session")
			return
		}
		sendErrorResponse(w, http.StatusBadRequest, "Failed to retrieve cookie")
		return
	}

	email, err := utils.ValidateToken(cookie.Value)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	if !ah.SessionStorage.SessionExists(email) {
		sendErrorResponse(w, http.StatusUnauthorized, "Session does not exist")
		return
	}

	err = ah.SessionStorage.RemoveSession(email)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to remove session")
		return
	}

	cookie = &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	sendJSONResponse(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func sendErrorResponse(w http.ResponseWriter, code int, status string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(AuthErrResponse{Code: code, Status: status})
}

func sendJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
