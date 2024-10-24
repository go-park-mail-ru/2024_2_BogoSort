package http

// import (
// 	"encoding/json"
// 	"errors"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery"
// 	domain2 "github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
// 	postgres2 "github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/postgres"
// 	"net/http"
// 	"time"

// 	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/utils"
// 	"github.com/go-playground/validator/v10"
// )

// var validate = validator.New()

// var (
// 	ErrInvalidRequestBody     = errors.New("invalid request body: unable to parse JSON")
// 	ErrInvalidRequestData     = errors.New("invalid request data: validation failed")
// 	ErrInvalidCredentials     = errors.New("invalid credentials: incorrect email or password")
// 	ErrMethodNotAllowed       = errors.New("method not allowed: only POST method is supported")
// 	ErrNoActiveSession        = errors.New("no active session: session cookie is missing")
// 	ErrSessionDoesNotExist    = errors.New("session does not exist: invalid session ID")
// 	ErrFailedToRetrieveCookie = errors.New("failed to retrieve cookie: unable to read session cookie")
// 	ErrFailedToCreateUser     = errors.New("failed to create user: internal server error")
// 	ErrUserAlreadyExists      = errors.New("user already exists: email is already registered")
// 	ErrFailedToRemoveSession  = errors.New("failed to remove session: internal server error")
// 	ErrUserNotFound           = errors.New("user not found: no user with provided credentials")
// )

// type AuthHandler struct {
// 	UserRepo    domain2.UserRepository
// 	SessionRepo domain2.SessionRepository
// }

// func NewAuthHandler() *AuthHandler {
// 	userRepo := postgres2.NewUserRepository()
// 	sessionRepo := postgres2.NewSessionRepository()

// 	return &AuthHandler{
// 		UserRepo:    userRepo,
// 		SessionRepo: sessionRepo,
// 	}
// }

// // SignupHandler godoc
// // @Summary Signup a new user
// // @Description Signup a new user with email and password
// // @Tags session
// // @Accept json
// // @Produce json
// // @Param credentials body AuthData true "User credentials"
// // @Success 201 {object} AuthResponse
// // @Header 200 {string} X-Authenticated "true"
// // @Failure 400 {object} ErrResponse "Invalid request body or data"
// // @Failure 400 {object} ErrResponse "User already exists"
// // @Failure 405 {object} ErrResponse "Method not allowed"
// // @Failure 500 {object} ErrResponse "Failed to create user"
// // @Router /signup [post]
// func (ah *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		delivery.SendErrorResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed.Error())
// 		return
// 	}

// 	var credentials domain2.AuthData
// 	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
// 		delivery.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody.Error())
// 		return
// 	}

// 	if err := validate.Struct(credentials); err != nil {
// 		delivery.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestData.Error())
// 		return
// 	}

// 	if err := utils.ValidatePassword(credentials.Password); err != nil {
// 		delivery.SendErrorResponse(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	existingUser, err := ah.UserRepo.GetUserByEmail(credentials.Email)
// 	if err == nil && existingUser != nil {
// 		delivery.SendErrorResponse(w, http.StatusBadRequest, ErrUserAlreadyExists.Error())
// 		return
// 	}

// 	user, err := ah.UserRepo.CreateUser(credentials.Email, credentials.Password)
// 	if err != nil {
// 		if errors.Is(err, ErrUserAlreadyExists) {
// 			delivery.SendErrorResponse(w, http.StatusBadRequest, ErrUserAlreadyExists.Error())
// 		} else {
// 			delivery.SendErrorResponse(w, http.StatusInternalServerError, ErrFailedToCreateUser.Error())
// 		}
// 		return
// 	}

// 	sessionID := ah.SessionRepo.AddSession(user.Email)

// 	cookie := &http.Cookie{
// 		Name:     "session_id",
// 		Value:    sessionID,
// 		Expires:  time.Now().Add(config.GetSessionExpirationTime()),
// 		HttpOnly: true,
// 	}

// 	http.SetCookie(w, cookie)

// 	w.Header().Set("X-Authenticated", "true")

// 	delivery.SendJSONResponse(w, http.StatusCreated, delivery.AuthResponse{
// 		Email:     user.Email,
// 		SessionID: sessionID,
// 	})
// }

// // LoginHandler godoc
// // @Summary Login a user
// // @Description Login a user with email and password
// // @Tags session
// // @Accept json
// // @Produce json
// // @Param credentials body LoginCredentials true "User credentials"
// // @Success 200 {object} AuthResponse
// // @Header 200 {string} X-Authenticated "true"
// // @Failure 400 {object} ErrResponse "Invalid request body or data"
// // @Failure 405 {object} ErrResponse "Method not allowed"
// // @Router /login [post]
// func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		delivery.SendErrorResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed.Error())
// 		return
// 	}

// 	var credentials domain2.LoginCredentials
// 	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
// 		delivery.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody.Error())
// 		return
// 	}

// 	if err := validate.Struct(credentials); err != nil {
// 		delivery.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestData.Error())
// 		return
// 	}

// 	user, err := ah.UserRepo.ValidateUserByEmailAndPassword(credentials.Email, credentials.Password)
// 	if err != nil {
// 		if errors.Is(err, ErrUserNotFound) {
// 			delivery.SendErrorResponse(w, http.StatusBadRequest, ErrUserNotFound.Error())
// 		} else {
// 			delivery.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidCredentials.Error())
// 		}
// 		return
// 	}

// 	sessionID := ah.SessionRepo.AddSession(user.Email)

// 	cookie := &http.Cookie{
// 		Name:     "session_id",
// 		Value:    sessionID,
// 		Expires:  time.Now().Add(config.GetSessionExpirationTime()),
// 		HttpOnly: true,
// 	}

// 	http.SetCookie(w, cookie)

// 	delivery.SendJSONResponse(w, http.StatusOK, delivery.AuthResponse{
// 		Email:     user.Email,
// 		SessionID: sessionID,
// 	})
// }

// // LogoutHandler godoc
// // @Summary Logout a user
// // @Description Logout a user by invalidating the session
// // @Tags session
// // @Accept json
// // @Produce json
// // @Success 200 {object} map[string]string "Logged out successfully"
// // @Header 200 {string} X-Authenticated "false"
// // @Failure 401 {object} ErrResponse "No active session or session does not exist"
// // @Failure 405 {object} ErrResponse "Method not allowed"
// // @Router /logout [post]
// func (ah *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		delivery.SendErrorResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed.Error())
// 		return
// 	}

// 	cookie, err := r.Cookie("session_id")
// 	if err != nil {
// 		if errors.Is(err, http.ErrNoCookie) {
// 			delivery.SendErrorResponse(w, http.StatusUnauthorized, ErrNoActiveSession.Error())
// 			return
// 		}
// 		delivery.SendErrorResponse(w, http.StatusBadRequest, ErrFailedToRetrieveCookie.Error())
// 		return
// 	}

// 	sessionID := cookie.Value

// 	if !ah.SessionRepo.SessionExists(sessionID) {
// 		delivery.SendErrorResponse(w, http.StatusUnauthorized, ErrSessionDoesNotExist.Error())
// 		return
// 	}

// 	err = ah.SessionRepo.RemoveSession(sessionID)
// 	if err != nil {
// 		delivery.SendErrorResponse(w, http.StatusInternalServerError, ErrFailedToRemoveSession.Error())
// 		return
// 	}

// 	cookie = &http.Cookie{
// 		Name:     "session_id",
// 		Value:    "",
// 		Expires:  time.Now().Add(-1 * time.Hour),
// 		HttpOnly: true,
// 	}
// 	http.SetCookie(w, cookie)

// 	w.Header().Set("X-Authenticated", "false")

// 	delivery.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
// }
