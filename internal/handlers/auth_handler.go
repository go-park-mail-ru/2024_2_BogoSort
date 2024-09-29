package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/responses"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/storage"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/utils"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type AuthData struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthHandler struct {
	UserStorage *storage.UserStorage
}

// SignupHandler godoc
// @Summary Signup a new user
// @Description Signup a new user with email and password
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
		responses.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var credentials AuthData
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		responses.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.Struct(credentials); err != nil {
		responses.SendErrorResponse(w, http.StatusBadRequest, "Invalid request data")
		return
	}

	user, err := ah.UserStorage.CreateUser(credentials.Email, credentials.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			responses.SendErrorResponse(w, http.StatusBadRequest, "User already exists")
		} else {
			responses.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		}
		return
	}

	tokenString, err := utils.CreateToken(user.Email)
	if err != nil {
		responses.SendErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    tokenString,
		Expires:  time.Now().Add(config.GetJWTExpirationTime()),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	responses.SendJSONResponse(w, http.StatusCreated, responses.AuthResponse{Token: tokenString, Email: user.Email})
}

// LoginHandler godoc
// @Summary Login a user
// @Description Login a user with email and password or with a valid session cookie or Authorization header
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginCredentials false "User credentials"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} AuthErrResponse
// @Failure 401 {object} AuthErrResponse
// @Failure 405 {object} AuthErrResponse
// @Failure 500 {object} AuthErrResponse
// @Router /login [post]
func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var email string
	var err error

	cookie, err := r.Cookie("session_id")
	if err == nil {
		email, err = utils.ValidateToken(cookie.Value)
		if err == nil {
			sendJSONResponse(w, http.StatusOK, AuthResponse{Token: cookie.Value, Email: email})
			return
		}
	}

	authHeader := r.Header.Get("Authorization")
    if authHeader != "" {
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        email, err = utils.ValidateToken(tokenString)
        if err == nil {
            sendJSONResponse(w, http.StatusOK, AuthResponse{Token: tokenString, Email: email})
            return
        }
    }

	var credentials LoginCredentials
	if r.Body != nil && r.ContentLength != 0 {
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			sendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := validate.Struct(credentials); err != nil {
			sendErrorResponse(w, http.StatusBadRequest, "Invalid request data")
			return
		}

		user, err := ah.UserStorage.ValidateUserByEmailAndPassword(credentials.Email, credentials.Password)
		if err != nil {
			sendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		tokenString, err := utils.CreateToken(user.Email)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
			return
		}

		cookie = &http.Cookie{
			Name:     "session_id",
			Value:    tokenString,
			Expires:  time.Now().Add(config.GetJWTExpirationTime()),
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		sendJSONResponse(w, http.StatusOK, AuthResponse{Token: tokenString, Email: user.Email})
		return
	}

	sendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
}

func (ah *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			responses.SendErrorResponse(w, http.StatusUnauthorized, "No active session")
			return
		}
		responses.SendErrorResponse(w, http.StatusBadRequest, "Failed to retrieve cookie")
		return
	}

	_, err = utils.ValidateToken(cookie.Value)
	if err != nil {
		responses.SendErrorResponse(w, http.StatusUnauthorized, "Invalid token")
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

	responses.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}
