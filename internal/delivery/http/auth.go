package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/gorilla/mux"
)

type AuthEndpoints struct {
	authUC         usecase.Auth
	sessionManager *utils.SessionManager
}

func NewAuthEndpoints(authUC usecase.Auth, sessionManager *utils.SessionManager) *AuthEndpoints {
	return &AuthEndpoints{
		authUC:         authUC,
		sessionManager: sessionManager,
	}
}

func (a *AuthEndpoints) Configure(router *mux.Router) {
	router.HandleFunc("/logout", a.Logout).Methods(http.MethodPost)
}

func (a *AuthEndpoints) Logout(w http.ResponseWriter, r *http.Request) {
	_, err := a.sessionManager.GetUserID(r)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = a.authUC.Logout(cookie.Value)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = a.sessionManager.DeleteSession(cookie.Value)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, "Вы успешно вышли из системы")
}
