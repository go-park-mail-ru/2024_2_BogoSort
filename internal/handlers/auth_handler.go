package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/responses"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/utils"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

var (
	ErrEmptyRequestBody   = errors.New("empty request body")
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrInvalidRequestData = errors.New("invalid request data")
)

type AuthCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// SignupHandler godoc
// @Summary Signup a new user
// @Description Signup a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body AuthCredentials true "User credentials"
// @Success 201 {object} responses.AuthResponse
// @Failure 400 {object} responses.ErrResponse
// @Failure 405 {object} responses.ErrResponse
// @Failure 500 {object} responses.ErrResponse
// @Router /api/v1/signup [post]
func (ah *AuthHandler) SignupHandler(writer http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.SendErrorResponse(writer, http.StatusMethodNotAllowed, "Method not allowed")

		return
	}

	var credentials AuthCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		responses.SendErrorResponse(writer, http.StatusBadRequest, "Invalid request body")

		return
	}

	if err := validate.Struct(credentials); err != nil {
		responses.SendErrorResponse(writer, http.StatusBadRequest, "Invalid request data")

		return
	}

	if err := utils.ValidatePassword(credentials.Password); err != nil {
		responses.SendErrorResponse(writer, http.StatusBadRequest, err.Error())

		return
	}

	user, err := ah.UserStorage.CreateUser(credentials.Email, credentials.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			responses.SendErrorResponse(writer, http.StatusBadRequest, "User already exists")
		} else {
			responses.SendErrorResponse(writer, http.StatusInternalServerError, "Failed to create user")
		}

		return
	}

	tokenString, err := utils.CreateToken(user.Email)
	if err != nil {
		responses.SendErrorResponse(writer, http.StatusInternalServerError, "Failed to generate token")

		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    tokenString,
		Expires:  time.Now().Add(config.GetJWTExpirationTime()),
		HttpOnly: true,
	}

	http.SetCookie(writer, cookie)

	responses.SendJSONResponse(writer, http.StatusCreated, responses.AuthResponse{Token: tokenString, Email: user.Email})
}

// LoginHandler godoc
// @Summary Login a user
// @Description Login a user with email and password or with a valid session cookie or Authorization header
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body AuthCredentials false "User credentials"
// @Success 200 {object} responses.AuthResponse
// @Failure 400 {object} responses.ErrResponse
// @Failure 401 {object} responses.ErrResponse
// @Failure 405 {object} responses.ErrResponse
// @Failure 500 {object} responses.ErrResponse
// @Router /api/v1/login [post]
func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var (
		email string
		err   error
	)

	var credentials AuthCredentials
	if r.Body != nil && r.ContentLength != 0 {
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			responses.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := validate.Struct(credentials); err != nil {
			responses.SendErrorResponse(w, http.StatusBadRequest, "Invalid request data")
			return
		}

		user, err := ah.UserStorage.ValidateUserByEmailAndPassword(credentials.Email, credentials.Password)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Expires:  time.Now().Add(-time.Hour),
				HttpOnly: true,
			})
			responses.SendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials: user not found")
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
		responses.SendJSONResponse(w, http.StatusOK, responses.AuthResponse{Token: tokenString, Email: user.Email})
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil {
		email, err = utils.ValidateToken(cookie.Value)
		if err == nil {
			responses.SendJSONResponse(w, http.StatusOK, responses.AuthResponse{Token: cookie.Value, Email: email})
			return
		}
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		email, err = utils.ValidateToken(tokenString)
		if err == nil {
			responses.SendJSONResponse(w, http.StatusOK, responses.AuthResponse{Token: tokenString, Email: email})
			return
		}
	}

	responses.SendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
}

// LogoutHandler godoc
// @Summary Logout a user
// @Description Logout a user by invalidating the session cookie or Authorization header
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} responses.ErrResponse
// @Failure 401 {object} responses.ErrResponse
// @Failure 405 {object} responses.ErrResponse
// @Router /api/v1/logout [post]
func (ah *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")

		return
	}

	var (
		token string
		err   error
	)

	cookie, err := r.Cookie("session_id")
	if err == nil {
		token = cookie.Value
	} else {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			responses.SendErrorResponse(w, http.StatusUnauthorized, "No active session")

			return
		}
	}

	_, err = utils.ValidateToken(token)
	if err != nil {
		responses.SendErrorResponse(w, http.StatusUnauthorized, "Invalid token")

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
