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
)

var (
	ErrInvalidRequestBody = errors.New("некорректные данные запроса")
	ErrUserAlreadyExists  = errors.New("пользователь уже существует")
	ErrUserNotFound       = errors.New("пользователь не найден")
)

type UserEndpoints struct {
	userUC         usecase.User
	authUC         usecase.Auth
	sessionManager *utils.SessionManager
}

func NewUserEndpoints(userUC usecase.User, authUC usecase.Auth, sessionManager *utils.SessionManager) *UserEndpoints {
	return &UserEndpoints{
		userUC:         userUC,
		authUC:         authUC,
		sessionManager: sessionManager,
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
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return
	}

	userID, err := u.userUC.Signup(&credentials)
	var errUserIncorrectData usecase.UserIncorrectDataError

	switch {
	case errors.Is(err, usecase.ErrUserAlreadyExists):
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrUserAlreadyExists.Error())
	case errors.As(err, &errUserIncorrectData):
		utils.SendErrorResponse(w, http.StatusBadRequest, errUserIncorrectData.Error())
	case err != nil:
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	default:
		sessionID, err := u.sessionManager.CreateSession(userID)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.SendJSONResponse(w, http.StatusOK, sessionID)
	}
}

func (u *UserEndpoints) Login(w http.ResponseWriter, r *http.Request) {
	var credentials dto.Login
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return
	}
	userID, err := u.userUC.Login(&credentials)
	var errUserIncorrectData usecase.UserIncorrectDataError

	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrUserNotFound.Error())
	case errors.As(err, &errUserIncorrectData):
		utils.SendErrorResponse(w, http.StatusBadRequest, errUserIncorrectData.Error())
	case err != nil:
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	default:
		sessionID, err := u.sessionManager.CreateSession(userID)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.SendJSONResponse(w, http.StatusOK, sessionID)
	}
}

func (u *UserEndpoints) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var updatePassword dto.UpdatePassword
	if err := json.NewDecoder(r.Body).Decode(&updatePassword); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return
	}
	userID, err := u.sessionManager.GetUserID(r)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	err = u.userUC.ChangePassword(userID, &updatePassword)
	var errUserIncorrectData usecase.UserIncorrectDataError

	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrUserNotFound.Error())
	case errors.As(err, &errUserIncorrectData):
		utils.SendErrorResponse(w, http.StatusBadRequest, errUserIncorrectData.Error())
	case err != nil:
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	default:
		utils.SendJSONResponse(w, http.StatusOK, "Password changed successfully")
	}
}

func (u *UserEndpoints) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var user dto.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return
	}
	_, err := u.sessionManager.GetUserID(r)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	err = u.userUC.UpdateInfo(&user)
	var errUserIncorrectData usecase.UserIncorrectDataError

	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrUserNotFound.Error())
	case errors.As(err, &errUserIncorrectData):
		utils.SendErrorResponse(w, http.StatusBadRequest, errUserIncorrectData.Error())
	case err != nil:
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	default:
		utils.SendJSONResponse(w, http.StatusOK, "Profile updated successfully")
	}
}

func (u *UserEndpoints) GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["user_id"])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := u.userUC.GetUser(userID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, user)
}

func (u *UserEndpoints) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, err := u.sessionManager.GetUserID(r)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	user, err := u.userUC.GetUser(userID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, user)
}
