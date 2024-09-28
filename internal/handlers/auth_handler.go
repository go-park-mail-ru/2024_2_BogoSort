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

func (ah *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{Token: tokenString, User: *user})
}

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

	json.NewEncoder(w).Encode(AuthResponse{Token: tokenString, User: *user})
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

	ah.SessionStorage.RemoveSession(email)

	cookie.Expires = time.Now().Add(-1 * time.Hour)
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		email, err := utils.ValidateToken(cookie.Value)
		if err != nil {
			sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		r.Header.Set("User", email)
		next.ServeHTTP(w, r)
	}
}

func sendErrorResponse(w http.ResponseWriter, code int, status string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(AuthErrResponse{Code: code, Status: status})
}

