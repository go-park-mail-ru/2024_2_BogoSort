package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/responses"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/utils"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

var (
	ErrInvalidRequestBody       = errors.New("invalid request body: unable to parse JSON")
	ErrInvalidRequestData       = errors.New("invalid request data: validation failed")
	ErrInvalidCredentials       = errors.New("invalid credentials: incorrect email or password")
	ErrMethodNotAllowed         = errors.New("method not allowed: only POST method is supported")
	ErrNoActiveSession          = errors.New("no active session: session cookie is missing")
	ErrSessionDoesNotExist      = errors.New("session does not exist: invalid session ID")
	ErrFailedToRetrieveCookie   = errors.New("failed to retrieve cookie: unable to read session cookie")
	ErrFailedToCreateUser       = errors.New("failed to create user: internal server error")
	ErrUserAlreadyExists        = errors.New("user already exists: email is already registered")
	ErrFailedToRemoveSession    = errors.New("failed to remove session: internal server error")
	ErrUserNotFound             = errors.New("user not found: no user with provided credentials")
)

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
// @Failure 400 {object} responses.ErrResponse "User already exists"
// @Failure 405 {object} responses.ErrResponse "Method not allowed"
// @Failure 500 {object} responses.ErrResponse "Failed to create user"
// @Router /signup [post]
func (ah *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.SendErrorResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed.Error())
		return
	}

	var credentials AuthData
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		responses.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return
	}

	if err := validate.Struct(credentials); err != nil {
		responses.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestData.Error())
		return
	}

	if err := utils.ValidatePassword(credentials.Password); err != nil {
		responses.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	existingUser, err := ah.UserStorage.GetUserByEmail(credentials.Email)
	if err == nil && existingUser != nil {
		responses.SendErrorResponse(w, http.StatusBadRequest, ErrUserAlreadyExists.Error())
		return
	}

	user, err := ah.UserStorage.CreateUser(credentials.Email, credentials.Password)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			responses.SendErrorResponse(w, http.StatusBadRequest, ErrUserAlreadyExists.Error())
		} else {
			responses.SendErrorResponse(w, http.StatusInternalServerError, ErrFailedToCreateUser.Error())
		}
		return
	}

	sessionID := ah.SessionStorage.AddSession(user.Email)

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(config.GetSessionExpirationTime()),
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
// @Failure 405 {object} responses.ErrResponse "Method not allowed"
// @Router /login [post]	
func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.SendErrorResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed.Error())
		return
	}

	var credentials LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		responses.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return
	}

	if err := validate.Struct(credentials); err != nil {
		responses.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestData.Error())
		return
	}

	user, err := ah.UserStorage.ValidateUserByEmailAndPassword(credentials.Email, credentials.Password)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			responses.SendErrorResponse(w, http.StatusBadRequest, ErrUserNotFound.Error())
		} else {
			responses.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidCredentials.Error())
		}
		return
	}

	sessionID := ah.SessionStorage.AddSession(user.Email)

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(config.GetSessionExpirationTime()),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	responses.SendJSONResponse(w, http.StatusOK, responses.AuthResponse{
		Email:     user.Email,
		SessionID: sessionID,
		IsAuth:    true,
	})
}

// LogoutHandler godoc
// @Summary Logout a user
// @Description Logout a user by invalidating the session
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Logged out successfully"
// @Failure 401 {object} responses.ErrResponse "No active session or session does not exist"
// @Failure 405 {object} responses.ErrResponse "Method not allowed"
// @Router /logout [post]
func (ah *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.SendErrorResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed.Error())
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			responses.SendErrorResponse(w, http.StatusUnauthorized, ErrNoActiveSession.Error())
			return
		}
		responses.SendErrorResponse(w, http.StatusBadRequest, ErrFailedToRetrieveCookie.Error())
		return
	}

	sessionID := cookie.Value

	if !ah.SessionStorage.SessionExists(sessionID) {
		responses.SendErrorResponse(w, http.StatusUnauthorized, ErrSessionDoesNotExist.Error())
		return
	}

	err = ah.SessionStorage.RemoveSession(sessionID)
	if err != nil {
		responses.SendErrorResponse(w, http.StatusInternalServerError, ErrFailedToRemoveSession.Error())
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

// CheckAuthHandler godoc
// @Summary Check if user is authenticated
// @Description Verify if the current session is valid and the user is authenticated
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} responses.AuthResponse "User is authenticated"
// @Failure 400 {object} responses.ErrResponse "Failed to retrieve cookie"
// @Failure 401 {object} responses.ErrResponse "No active session or session does not exist"
// @Router /check-auth [get]	
func (ah *AuthHandler) CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			responses.SendErrorResponse(w, http.StatusUnauthorized, ErrNoActiveSession.Error())
			return
		}
		responses.SendErrorResponse(w, http.StatusBadRequest, ErrFailedToRetrieveCookie.Error())
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