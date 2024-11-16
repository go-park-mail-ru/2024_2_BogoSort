package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"bytes"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/static"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserInvalidRequestBody      = errors.New("invalid request body")
	ErrUserAlreadyExists           = errors.New("user already exists")
	ErrUserNotFound                = errors.New("user not found")
	ErrUserIncorrectData           = errors.New("user incorrect data")
	ErrInvalidCredentials          = errors.New("invalid credentials")
	ErrUnauthorized                = errors.New("unauthorized request")
	ErrOldAndNewPasswordAreTheSame = errors.New("old and new password are the same")
)

type UserEndpoints struct {
	userUC         usecase.User
	authUC         usecase.Auth
	sessionManager *utils.SessionManager
	staticGrpcClient static.StaticGrpcClient
	logger         *zap.Logger
	policy         *bluemonday.Policy
}

func NewUserEndpoints(userUC usecase.User, authUC usecase.Auth, sessionManager *utils.SessionManager, staticGrpcClient static.StaticGrpcClient, logger *zap.Logger, policy *bluemonday.Policy) *UserEndpoints {
	return &UserEndpoints{
		userUC:         userUC,
		authUC:         authUC,
		sessionManager: sessionManager,
		staticGrpcClient: staticGrpcClient,
		logger:         logger,
		policy:         policy,
	}
}

func (u *UserEndpoints) ConfigureProtectedRoutes(router *mux.Router) {
	protected := router.PathPrefix("/api/v1").Subrouter()
	sessionMiddleware := middleware.NewAuthMiddleware(u.sessionManager)
	protected.Use(sessionMiddleware.SessionMiddleware)

	protected.HandleFunc("/password", u.ChangePassword).Methods(http.MethodPost)
	protected.HandleFunc("/profile", u.UpdateProfile).Methods(http.MethodPut)
	protected.HandleFunc("/me", u.GetMe).Methods(http.MethodGet)
	protected.HandleFunc("/user/{user_id}/image", u.UploadImage).Methods(http.MethodPut)
}

func (u *UserEndpoints) ConfigureUnprotectedRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/signup", u.Signup).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/login", u.Login).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/profile/{user_id}", u.GetProfile).Methods(http.MethodGet)
}

func (u *UserEndpoints) handleError(w http.ResponseWriter, err error, context string, additionalInfo map[string]string) {
	var errUserIncorrectData usecase.UserIncorrectDataError

	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		u.sendError(w, http.StatusNotFound, ErrUserNotFound, context, additionalInfo)
	case errors.Is(err, usecase.ErrUserAlreadyExists):
		u.sendError(w, http.StatusBadRequest, ErrUserAlreadyExists, context, additionalInfo)
	case errors.Is(err, usecase.ErrInvalidCredentials):
		u.sendError(w, http.StatusUnauthorized, ErrInvalidCredentials, context, additionalInfo)
	case errors.As(err, &errUserIncorrectData):
		u.sendError(w, http.StatusBadRequest, errUserIncorrectData, context, additionalInfo)
	case errors.Is(err, usecase.ErrOldAndNewPasswordAreTheSame):
		u.sendError(w, http.StatusBadRequest, ErrOldAndNewPasswordAreTheSame, context, additionalInfo)
	case err != nil:
		u.sendError(w, http.StatusInternalServerError, err, context, additionalInfo)
	}
}

func (u *UserEndpoints) sendError(w http.ResponseWriter, statusCode int, err error, context string, additionalInfo map[string]string) {
	u.logger.Error(err.Error(), zap.String("context", context), zap.Any("info", additionalInfo))
	utils.SendErrorResponse(w, statusCode, err.Error())
}

// Signup
// @Summary User registration
// @Description Creates a new user in the system
// @Tags Users
// @Accept json
// @Produce json
// @Param signup body dto.Signup true "Registration data"
// @Success 200 {string} string "SessionID"
// @Failure 400 {object} utils.ErrResponse "Invalid request or user already exists"
// @Failure 401 {object} utils.ErrResponse "Unauthorized request"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/signup [post]
func (u *UserEndpoints) Signup(w http.ResponseWriter, r *http.Request) {
	var credentials dto.Signup
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		u.sendError(w, http.StatusBadRequest, err, "error decoding signup request", nil)
		return
	}

	utils.SanitizeRequestSignup(&credentials, u.policy)

	userID, err := u.userUC.Signup(&credentials)
	if err != nil {
		u.handleError(w, err, "Signup", map[string]string{"email": credentials.Email})
		return
	}

	sessionID, err := u.sessionManager.CreateSession(userID)
	if err != nil {
		u.sendError(w, http.StatusInternalServerError, err, "error creating session", map[string]string{"userID": userID.String()})
		return
	}
	u.logger.Info("session created", zap.String("sessionID", sessionID), zap.String("userID", userID.String()))

	cookie, err := u.sessionManager.SetSession(sessionID)
	if err != nil {
		u.logger.Error("error setting session cookie", zap.Error(err))
		u.sendError(w, http.StatusInternalServerError, err, "error setting session cookie", nil)
		return
	}
	http.SetCookie(w, cookie)

	w.Header().Set("X-authenticated", "true")
	utils.SendJSONResponse(w, http.StatusOK, "Signup successful")
}

// Login
// @Summary User login
// @Description Allows a user to log into the system
// @Tags Users
// @Accept json
// @Produce json
// @Param login body dto.Login true "Login data"
// @Success 200 {string} string "SessionID"
// @Failure 400 {object} utils.ErrResponse "Invalid request"
// @Failure 401 {object} utils.ErrResponse "Invalid credentials or unauthorized access"
// @Failure 404 {object} utils.ErrResponse "User not found"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/login [post]
func (u *UserEndpoints) Login(w http.ResponseWriter, r *http.Request) {
	var credentials dto.Login
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		u.sendError(w, http.StatusBadRequest, err, "error decoding login request", nil)
		return
	}
	utils.SanitizeRequestLogin(&credentials, u.policy)

	userID, err := u.userUC.Login(&credentials)
	if err != nil {
		u.handleError(w, err, "Login", map[string]string{"email": credentials.Email})
		return
	}

	sessionID, err := u.sessionManager.CreateSession(userID)
	if err != nil {
		u.sendError(w, http.StatusInternalServerError, err, "error creating session", map[string]string{"userID": userID.String()})
		return
	}
	u.logger.Info("session created", zap.String("sessionID", sessionID), zap.String("userID", userID.String()))

	cookie, err := u.sessionManager.SetSession(sessionID)
	if err != nil {
		u.sendError(w, http.StatusInternalServerError, err, "error setting session cookie", nil)
		return
	}
	http.SetCookie(w, cookie)
	w.Header().Set("X-authenticated", "true")
	utils.SendJSONResponse(w, http.StatusOK, sessionID)
}

// ChangePassword
// @Summary Change user password
// @Description Allows a user to change their password
// @Tags Users
// @Accept json
// @Produce json
// @Param password body dto.UpdatePassword true "Password change data"
// @Success 200 {string} string "Password changed successfully"
// @Failure 400 {object} utils.ErrResponse "Invalid data"
// @Failure 401 {object} utils.ErrResponse "Unauthorized access"
// @Failure 404 {object} utils.ErrResponse "User not found"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/password [post]
func (u *UserEndpoints) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var updatePassword dto.UpdatePassword
	if err := json.NewDecoder(r.Body).Decode(&updatePassword); err != nil {
		u.sendError(w, http.StatusBadRequest, err, "error decoding change password request", nil)
		return
	}
	utils.SanitizeRequestChangePassword(&updatePassword, u.policy)

	userID, err := u.sessionManager.GetUserID(r)
	if err != nil {
		u.sendError(w, http.StatusUnauthorized, err, "unauthorized request", nil)
		return
	}
	err = u.userUC.ChangePassword(userID, &updatePassword)
	if err != nil {
		u.handleError(w, err, "ChangePassword", map[string]string{"userID": userID.String()})
		return
	}

	u.logger.Info("password changed", zap.String("userID", userID.String()))
	utils.SendJSONResponse(w, http.StatusOK, "Пароль изменен успешно")
}

// UpdateProfile
// @Summary Update user profile
// @Description Allows a user to update their profile information
// @Tags Users
// @Accept json
// @Produce json
// @Param profile body dto.UserUpdate true "Profile data"
// @Success 200 {string} string "Profile updated successfully"
// @Failure 400 {object} utils.ErrResponse "Invalid data"
// @Failure 401 {object} utils.ErrResponse "Unauthorized access"
// @Failure 404 {object} utils.ErrResponse "User not found"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/profile [put]
func (u *UserEndpoints) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var user dto.UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		u.sendError(w, http.StatusBadRequest, err, "error decoding update profile request", nil)
		return
	}
	utils.SanitizeRequestUserUpdate(&user, u.policy)
	userID, err := u.sessionManager.GetUserID(r)
	if err != nil {
		u.handleError(w, err, "UpdateProfile", nil)
		return
	}

	err = u.userUC.UpdateInfo(&user)
	if err != nil {
		u.handleError(w, err, "UpdateProfile", map[string]string{"userID": userID.String()})
		return
	}

	u.logger.Info("profile updated", zap.String("userID", userID.String()))
	utils.SendJSONResponse(w, http.StatusOK, "Профиль обновлен успешно")
}

// GetProfile
// @Summary Get user profile
// @Description Returns user information by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} dto.User "User profile"
// @Failure 404 {object} utils.ErrResponse "User not found"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/profile/{user_id} [get]
func (u *UserEndpoints) GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["user_id"])
	if err != nil {
		u.sendError(w, http.StatusBadRequest, err, "error parsing userID", nil)
		return
	}
	user, err := u.userUC.GetUser(userID)
	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		u.sendError(w, http.StatusNotFound, err, "user not found", nil)
	case err != nil:
		u.handleError(w, err, "GetProfile", map[string]string{"userID": userID.String()})
	}
	utils.SanitizeResponseUser(user, u.policy)
	utils.SendJSONResponse(w, http.StatusOK, user)
}

// GetMe
// @Summary Get current user information
// @Description Returns information about the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} dto.User "User information"
// @Failure 401 {object} utils.ErrResponse "Unauthorized access"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/me [get]
func (u *UserEndpoints) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, err := u.sessionManager.GetUserID(r)
	if err != nil {
		u.sendError(w, http.StatusUnauthorized, err, "unauthorized request", nil)
		return
	}
	user, err := u.userUC.GetUser(userID)
	if err != nil {
		u.handleError(w, err, "GetMe", map[string]string{"userID": userID.String()})
		return
	}
	utils.SanitizeResponseUser(user, u.policy)
	utils.SendJSONResponse(w, http.StatusOK, user)
}

// UploadImage godoc
// @Summary Upload an image for an advert
// @Description Upload an image associated with an advert by its ID
// @Tags adverts
// @Param user_id path string true "User ID"
// @Param image formData file true "Image file to upload"
// @Success 200 {string} string "Image uploaded"
// @Failure 400 {object} utils.ErrResponse "Invalid user ID or file not attached"
// @Failure 500 {object} utils.ErrResponse "Failed to upload image"
// @Router /api/v1/user/{user_id}/image [put]
func (h *UserEndpoints) UploadImage(writer http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["user_id"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, err, "invalid advert ID", nil)
		return
	}

	fileHeader, _, err := r.FormFile("image")
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, err, "file not attached", nil)
		return
	}

	data, err := io.ReadAll(fileHeader)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to read file", nil)
		return
	}

	if err = fileHeader.Close(); err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to close file", nil)
		return
	}

	imageId, err := h.staticGrpcClient.UploadStatic(bytes.NewReader(data))
	if err != nil {
		if status, ok := status.FromError(err); ok {
			switch status.Code() {
			case codes.DeadlineExceeded:
				h.sendError(writer, http.StatusGatewayTimeout, ErrTimeout, "upload image timeout deadline exceeded", nil)
			case codes.ResourceExhausted:
				h.sendError(writer, http.StatusRequestEntityTooLarge, ErrTooLargeFile, "file size exceeds limit", nil)
			default:
				h.sendError(writer, http.StatusInternalServerError, ErrFailedToUploadFile, "failed to upload image", nil)
			}
		} else {
			h.sendError(writer, http.StatusInternalServerError, ErrFailedToUploadFile, "failed to upload image", nil)
		}
		return
	}

	if err := h.userUC.UploadImage(userID, imageId); err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to upload image", nil)
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, "Image uploaded")
}
