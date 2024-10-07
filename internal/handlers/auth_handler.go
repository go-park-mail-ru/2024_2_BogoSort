package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/responses"
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

// SignupHandler godoc
// @Summary Signup a new user
// @Description Signup a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body AuthData true "User credentials"
// @Success 201 {object} responses.AuthResponse
// @Failure 400 {object} responses.ErrResponse "Invalid request body or data"
// @Failure 405 {object} responses.ErrResponse "Method not allowed"
// @Failure 500 {object} responses.ErrResponse "Failed to create user"
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

	if err := utils.ValidatePassword(credentials.Password); err != nil {
		responses.SendErrorResponse(w, http.StatusBadRequest, err.Error())
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

	sessionID := ah.SessionStorage.AddSession(user.Email)

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(config.GetJWTExpirationTime()),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	responses.SendJSONResponse(w, http.StatusCreated, responses.AuthResponse{
		Email:     user.Email,
		SessionID: sessionID,
		IsAuth:    true,
	})
}

// LoginHandler godoc
// @Summary Login a user
// @Description Login a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginCredentials true "User credentials"
// @Success 200 {object} responses.AuthResponse
// @Failure 400 {object} responses.ErrResponse "Invalid request body or data"
// @Failure 401 {object} responses.ErrResponse "Invalid credentials"
// @Failure 405 {object} responses.ErrResponse "Method not allowed"
// @Router /login [post]	
func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var credentials LoginCredentials
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
		responses.SendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	sessionID := ah.SessionStorage.AddSession(user.Email)

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(config.GetJWTExpirationTime()),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	responses.SendJSONResponse(w, http.StatusOK, responses.AuthResponse{
		Email:     user.Email,
		SessionID: sessionID,
		IsAuth:    true,
	})
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

	sessionID := cookie.Value

	if !ah.SessionStorage.SessionExists(sessionID) {
		responses.SendErrorResponse(w, http.StatusUnauthorized, "Session does not exist")
		return
	}

	err = ah.SessionStorage.RemoveSession(sessionID)
	if err != nil {
		responses.SendErrorResponse(w, http.StatusInternalServerError, "Failed to remove session")
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

// CheckAuth godoc
// @Summary Check if user is authenticated
// @Description Verify if the current session is valid and the user is authenticated
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} responses.AuthResponse
// @Failure 400 {object} responses.ErrResponse "Failed to retrieve cookie"
// @Failure 401 {object} responses.ErrResponse "No active session or session does not exist"
// @Router /check-auth [get]	
func (ah *AuthHandler) CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			responses.SendErrorResponse(w, http.StatusUnauthorized, "No active session")
			return
		}
		responses.SendErrorResponse(w, http.StatusBadRequest, "Failed to retrieve cookie")
		return
	}

	sessionID := cookie.Value

	if !ah.SessionStorage.SessionExists(sessionID) {
		responses.SendJSONResponse(w, http.StatusOK, responses.AuthResponse{
			Email:     "",
			SessionID: "",
			IsAuth:    false,
		})

		return
	}

	responses.SendJSONResponse(w, http.StatusOK, responses.AuthResponse{
		Email:     "",
		SessionID: sessionID,
		IsAuth:    true,
	})
}