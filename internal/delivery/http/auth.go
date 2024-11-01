package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type AuthEndpoints struct {
	authUC         usecase.Auth
	sessionManager *utils.SessionManager
	logger         *zap.Logger
}

func NewAuthEndpoints(authUC usecase.Auth, sessionManager *utils.SessionManager, logger *zap.Logger) *AuthEndpoints {
	return &AuthEndpoints{
		authUC:         authUC,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

func (a *AuthEndpoints) Configure(router *mux.Router) {
	router.HandleFunc("/logout", a.Logout).Methods(http.MethodPost)
}

func (a *AuthEndpoints) handleError(w http.ResponseWriter, err error, method string, data map[string]string) {
	a.logger.Error(method+" error", zap.Error(err), zap.Any("data", data))
	utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
}

// Logout
// @Summary Выход пользователя
// @Description Позволяет пользователю выйти из системы, удаляя его сессию
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Success 200 {string} string "Вы успешно вышли из системы"
// @Failure 400 {object} utils.ErrResponse "Некорректный запрос или отсутствие cookie"
// @Failure 401 {object} utils.ErrResponse "Несанкционированный доступ"
// @Failure 500 {object} utils.ErrResponse "Внутренняя ошибка сервера"
// @Router /logout [post]
func (a *AuthEndpoints) Logout(w http.ResponseWriter, r *http.Request) {
	userID, err := a.sessionManager.GetUserID(r)
	if err != nil {
		a.handleError(w, err, "Logout", nil)
		return
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		a.handleError(w, err, "Logout", nil)
		return
	}
	err = a.authUC.Logout(cookie.Value)
	if err != nil {
		a.handleError(w, err, "Logout", map[string]string{"userID": userID.String()})
		return
	}
	err = a.sessionManager.DeleteSession(cookie.Value)
	if err != nil {
		a.handleError(w, err, "Logout", map[string]string{"userID": userID.String()})
		return
	}
	a.logger.Info("user logged out", zap.String("userID", userID.String()))
	utils.SendJSONResponse(w, http.StatusOK, "Вы успешно вышли из системы")
}
