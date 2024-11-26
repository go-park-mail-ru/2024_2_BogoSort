package http

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type AuthEndpoint struct {
	authUC         usecase.Auth
	sessionManager *utils.SessionManager
}

func NewAuthEndpoint(authUC usecase.Auth, sessionManager *utils.SessionManager) *AuthEndpoint {
	return &AuthEndpoint{
		authUC:         authUC,
		sessionManager: sessionManager,
	}
}

func (a *AuthEndpoint) Configure(router *mux.Router) {
	protected := router.PathPrefix("/api/v1").Subrouter()
	sessionMiddleware := middleware.NewAuthMiddleware(a.sessionManager)
	protected.Use(sessionMiddleware.SessionMiddleware)

	protected.HandleFunc("/logout", a.Logout).Methods(http.MethodPost)
}

func (a *AuthEndpoint) handleError(w http.ResponseWriter, err error, method string, data map[string]string) {
	logger := middleware.GetLogger(context.Background())
	logger.Error(method+" error", zap.Error(err), zap.Any("data", data))
	utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
}

// Logout
// @Summary User logout
// @Description Allows the user to log out of the system by deleting their session
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {string} string "You have successfully logged out"
// @Failure 400 {object} utils.ErrResponse "Invalid request or missing cookie"
// @Failure 401 {object} utils.ErrResponse "Unauthorized access"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/logout [post]
func (a *AuthEndpoint) Logout(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())
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
	err = a.sessionManager.DeleteSession(cookie.Value)
	if err != nil {
		a.handleError(w, err, "Logout", map[string]string{"userID": userID.String()})
		return
	}

	logger.Info("user logged out", zap.String("userID", userID.String()))
	w.Header().Set("X-authenticated", "false")
	utils.SendJSONResponse(w, http.StatusOK, "You have successfully logged out")
}
