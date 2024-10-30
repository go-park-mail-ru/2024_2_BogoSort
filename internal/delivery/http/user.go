package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	ErrInvalidRequestBody          = errors.New("invalid request body")
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
	logger         *zap.Logger
}

func NewUserEndpoints(userUC usecase.User, authUC usecase.Auth, sessionManager *utils.SessionManager, logger *zap.Logger) *UserEndpoints {
	return &UserEndpoints{
		userUC:         userUC,
		authUC:         authUC,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

func (u *UserEndpoints) Configure(router *mux.Router) {
	router.HandleFunc("/signup", u.Signup).Methods(http.MethodPost)
	router.HandleFunc("/login", u.Login).Methods(http.MethodPost)
	router.HandleFunc("/password", u.ChangePassword).Methods(http.MethodPost)
	router.HandleFunc("/profile/{user_id}", u.GetProfile).Methods(http.MethodGet)
	router.HandleFunc("/profile", u.UpdateProfile).Methods(http.MethodPut)
	router.HandleFunc("/me", u.GetMe).Methods(http.MethodGet)
}

func (u *UserEndpoints) Signup(w http.ResponseWriter, r *http.Request) {
	var credentials dto.Signup
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		u.logger.Error("error decoding signup request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return
	}

	userID, err := u.userUC.Signup(&credentials)
	var errUserIncorrectData usecase.UserIncorrectDataError

	switch {
	case errors.Is(err, usecase.ErrUserAlreadyExists):
		u.logger.Error("user already exists", zap.String("email", credentials.Email))
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrUserAlreadyExists.Error())
	case errors.Is(err, ErrUnauthorized):
		u.logger.Error("unauthorized request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusUnauthorized, ErrUnauthorized.Error())
	case errors.As(err, &errUserIncorrectData):
		u.logger.Error("user incorrect data", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, errUserIncorrectData.Error())
	case err != nil:
		u.logger.Error("error signing up", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	default:
		sessionID, err := u.sessionManager.CreateSession(userID)
		if err != nil {
			u.logger.Error("error creating session", zap.String("userID", userID.String()), zap.Error(err))
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		u.logger.Info("session created", zap.String("sessionID", sessionID), zap.String("userID", userID.String()))
		utils.SendJSONResponse(w, http.StatusOK, sessionID)
	}
}

func (u *UserEndpoints) Login(w http.ResponseWriter, r *http.Request) {
	var credentials dto.Login
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		u.logger.Error("error decoding login request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return
	}
	userID, err := u.userUC.Login(&credentials)
	var errUserIncorrectData usecase.UserIncorrectDataError
	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		u.logger.Error("user not found", zap.String("email", credentials.Email))
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrUserNotFound.Error())
	case errors.Is(err, usecase.ErrInvalidCredentials):
		u.logger.Error("invalid credentials", zap.String("email", credentials.Email))
		utils.SendErrorResponse(w, http.StatusUnauthorized, ErrInvalidCredentials.Error())
	case errors.As(err, &errUserIncorrectData):
		u.logger.Error("user incorrect data", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, errUserIncorrectData.Error())
	case err != nil:
		u.logger.Error("error logging in", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	default:
		sessionID, err := u.sessionManager.CreateSession(userID)
		if err != nil {
			u.logger.Error("error creating session", zap.String("userID", userID.String()), zap.Error(err))
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		u.logger.Info("session created", zap.String("sessionID", sessionID), zap.String("userID", userID.String()))
		utils.SendJSONResponse(w, http.StatusOK, sessionID)
	}
}

func (u *UserEndpoints) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var updatePassword dto.UpdatePassword
	if err := json.NewDecoder(r.Body).Decode(&updatePassword); err != nil {
		u.logger.Error("error decoding change password request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return
	}
	userID, err := u.sessionManager.GetUserID(r)
	if err != nil {
		u.logger.Error("unauthorized request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	err = u.userUC.ChangePassword(userID, &updatePassword)
	var errUserIncorrectData usecase.UserIncorrectDataError

	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		u.logger.Error("user not found", zap.String("userID", userID.String()))
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrUserNotFound.Error())
	case errors.Is(err, usecase.ErrInvalidCredentials):
		u.logger.Error("invalid credentials", zap.String("userID", userID.String()))
		utils.SendErrorResponse(w, http.StatusUnauthorized, ErrInvalidCredentials.Error())
	case errors.As(err, &errUserIncorrectData):
		u.logger.Error("user incorrect data", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, errUserIncorrectData.Error())
	case errors.Is(err, usecase.ErrOldAndNewPasswordAreTheSame):
		u.logger.Error("old and new password are the same", zap.String("userID", userID.String()))
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrOldAndNewPasswordAreTheSame.Error())
	case err != nil:
		u.logger.Error("error changing password", zap.String("userID", userID.String()), zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	default:
		u.logger.Info("password changed", zap.String("userID", userID.String()))
		utils.SendJSONResponse(w, http.StatusOK, "Пароль изменен успешно")
	}
}

func (u *UserEndpoints) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var user dto.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		u.logger.Error("error decoding update profile request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return
	}
	userID, err := u.sessionManager.GetUserID(r)
	if err != nil {
		u.logger.Error("unauthorized request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	err = u.userUC.UpdateInfo(&user)
	var errUserIncorrectData usecase.UserIncorrectDataError

	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		u.logger.Error("user not found", zap.String("userID", userID.String()))
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrUserNotFound.Error())
	case errors.Is(err, usecase.ErrInvalidCredentials):
		u.logger.Error("invalid credentials", zap.String("userID", userID.String()))
		utils.SendErrorResponse(w, http.StatusUnauthorized, ErrInvalidCredentials.Error())
	case errors.As(err, &errUserIncorrectData):
		u.logger.Error("user incorrect data", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, errUserIncorrectData.Error())
	case err != nil:
		u.logger.Error("error updating profile", zap.String("userID", userID.String()), zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	default:
		u.logger.Info("profile updated", zap.String("userID", userID.String()))
		utils.SendJSONResponse(w, http.StatusOK, "Профиль обновлен успешно")
	}
}

func (u *UserEndpoints) GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["user_id"])
	if err != nil {
		u.logger.Error("error parsing userID", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := u.userUC.GetUser(userID)
	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		u.logger.Error("user not found", zap.String("userID", userID.String()))
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrUserNotFound.Error())
	case err != nil:
		u.logger.Error("error getting user", zap.String("userID", userID.String()), zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, user)
}

func (u *UserEndpoints) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, err := u.sessionManager.GetUserID(r)
	if err != nil {
		u.logger.Error("unauthorized request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	user, err := u.userUC.GetUser(userID)
	if err != nil {
		u.logger.Error("error getting user", zap.String("userID", userID.String()), zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	u.logger.Info("user found", zap.String("userID", userID.String()))
	utils.SendJSONResponse(w, http.StatusOK, user)
}
