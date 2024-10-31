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

func (u *UserEndpoints) handleError(w http.ResponseWriter, err error, context string, additionalInfo map[string]string) {
	var errUserIncorrectData usecase.UserIncorrectDataError

	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		u.sendError(w, http.StatusNotFound, ErrUserNotFound, context, additionalInfo)
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
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя в системе
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param signup body dto.Signup true "Данные для регистрации"
// @Success 200 {string} string "SessionID"
// @Failure 400 {object} utils.ErrResponse "Некорректный запрос или пользователь уже существует"
// @Failure 401 {object} utils.ErrResponse "Несанкционированный запрос"
// @Failure 500 {object} utils.ErrResponse "Внутренняя ошибка сервера"
// @Router /signup [post]
func (u *UserEndpoints) Signup(w http.ResponseWriter, r *http.Request) {
	var credentials dto.Signup
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		u.logger.Error("error decoding signup request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return
	}

	userID, err := u.userUC.Signup(&credentials)
	if err != nil {
		u.handleError(w, err, "Signup", map[string]string{"email": credentials.Email})
		return
	}

	sessionID, err := u.sessionManager.CreateSession(userID)
	if err != nil {
		u.logger.Error("error creating session", zap.String("userID", userID.String()), zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	u.logger.Info("session created", zap.String("sessionID", sessionID), zap.String("userID", userID.String()))
	utils.SendJSONResponse(w, http.StatusOK, sessionID)

}

// Login
// @Summary Вход пользователя
// @Description Позволяет пользователю войти в систему
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param login body dto.Login true "Данные для входа"
// @Success 200 {string} string "SessionID"
// @Failure 400 {object} utils.ErrResponse "Некорректный запрос"
// @Failure 401 {object} utils.ErrResponse "Неверные учетные данные или несанкционированный доступ"
// @Failure 404 {object} utils.ErrResponse "Пользователь не найден"
// @Failure 500 {object} utils.ErrResponse "Внутренняя ошибка сервера"
// @Router /login [post]
func (u *UserEndpoints) Login(w http.ResponseWriter, r *http.Request) {
	var credentials dto.Login
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		u.logger.Error("error decoding login request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return
	}
	userID, err := u.userUC.Login(&credentials)
	if err != nil {
		u.handleError(w, err, "Login", map[string]string{"email": credentials.Email})
		return
	}

	sessionID, err := u.sessionManager.CreateSession(userID)
	if err != nil {
		u.logger.Error("error creating session", zap.String("userID", userID.String()), zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	u.logger.Info("session created", zap.String("sessionID", sessionID), zap.String("userID", userID.String()))
	utils.SendJSONResponse(w, http.StatusOK, sessionID)

}

// ChangePassword
// @Summary Изменение пароля пользователя
// @Description Позволяет пользователю изменить свой пароль
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param password body dto.UpdatePassword true "Данные для изменения пароля"
// @Success 200 {string} string "Пароль изменен успешно"
// @Failure 400 {object} utils.ErrResponse "Некорректные данные"
// @Failure 401 {object} utils.ErrResponse "Несанкционированный доступ"
// @Failure 404 {object} utils.ErrResponse "Пользователь не найден"
// @Failure 500 {object} utils.ErrResponse "Внутренняя ошибка сервера"
// @Router /password [post]
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
	if err != nil {
		u.handleError(w, err, "ChangePassword", map[string]string{"userID": userID.String()})
		return
	}

	u.logger.Info("password changed", zap.String("userID", userID.String()))
	utils.SendJSONResponse(w, http.StatusOK, "Пароль изменен успешно")
}

// UpdateProfile
// @Summary Обновление профиля пользователя
// @Description Позволяет пользователю обновить информацию своего профиля
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param profile body dto.User true "Данные профиля"
// @Success 200 {string} string "Профиль обновлен успешно"
// @Failure 400 {object} utils.ErrResponse "Некорректные данные"
// @Failure 401 {object} utils.ErrResponse "Несанкционированный доступ"
// @Failure 404 {object} utils.ErrResponse "Пользователь не найден"
// @Failure 500 {object} utils.ErrResponse "Внутренняя ошибка сервера"
// @Router /profile [put]
func (u *UserEndpoints) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var user dto.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		u.logger.Error("error decoding update profile request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return
	}
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
// @Summary Получение профиля пользователя
// @Description Возвращает информацию о пользователе по его ID
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Success 200 {object} dto.User "Профиль пользователя"
// @Failure 404 {object} utils.ErrResponse "Пользователь не найден"
// @Failure 500 {object} utils.ErrResponse "Внутренняя ошибка сервера"
// @Router /profile/{user_id} [get]
func (u *UserEndpoints) GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["user_id"])
	if err != nil {
		u.logger.Error("error parsing userID", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := u.userUC.GetUser(userID)
	if err != nil {
		u.handleError(w, err, "GetProfile", map[string]string{"userID": userID.String()})
		return
	}

	u.logger.Info("user found", zap.String("userID", userID.String()))
	utils.SendJSONResponse(w, http.StatusOK, user)
}

// GetMe
// @Summary Получение информации о текущем пользователе
// @Description Возвращает информацию о пользователе, текущий пользователь которого аутентифицирован
// @Tags Пользователи
// @Accept json
// @Produce json
// @Success 200 {object} dto.User "Информация о пользователе"
// @Failure 401 {object} utils.ErrResponse "Несанкционированный доступ"
// @Failure 500 {object} utils.ErrResponse "Внутренняя ошибка сервера"
// @Router /me [get]
func (u *UserEndpoints) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, err := u.sessionManager.GetUserID(r)
	if err != nil {
		u.logger.Error("unauthorized request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	user, err := u.userUC.GetUser(userID)
	if err != nil {
		u.handleError(w, err, "GetMe", map[string]string{"userID": userID.String()})
		return
	}

	u.logger.Info("user found", zap.String("userID", userID.String()))
	utils.SendJSONResponse(w, http.StatusOK, user)
}
