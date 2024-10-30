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

func (a *AuthEndpoints) Logout(w http.ResponseWriter, r *http.Request) {
	userID, err := a.sessionManager.GetUserID(r)
	if err != nil {
		a.logger.Error("unauthorized request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		a.logger.Error("error getting cookie", zap.String("userID", userID.String()), zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = a.authUC.Logout(cookie.Value)
	if err != nil {
		a.logger.Error("error logging out", zap.String("userID", userID.String()), zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = a.sessionManager.DeleteSession(cookie.Value)
	if err != nil {
		a.logger.Error("error deleting session", zap.String("userID", userID.String()), zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	a.logger.Info("user logged out", zap.String("userID", userID.String()))
	utils.SendJSONResponse(w, http.StatusOK, "Вы успешно вышли из системы")
}
