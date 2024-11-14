package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/auth"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type AuthEndpoints struct {
	authUC     usecase.Auth
	grpcClient *auth.GrpcClient
	logger     *zap.Logger
}

func NewAuthEndpoints(authUC usecase.Auth, grpcClient *auth.GrpcClient, logger *zap.Logger) *AuthEndpoints {
	return &AuthEndpoints{
		authUC:     authUC,
		grpcClient: grpcClient,
		logger:     logger,
	}
}

func (a *AuthEndpoints) Configure(router *mux.Router) {
	protected := router.PathPrefix("/api/v1").Subrouter()
	// sessionMiddleware := middleware.NewAuthMiddleware(a.sessionManager)
	// protected.Use(sessionMiddleware.SessionMiddleware)

	protected.HandleFunc("/logout", a.Logout).Methods(http.MethodPost)
}

func (a *AuthEndpoints) handleError(w http.ResponseWriter, err error, method string, data map[string]string) {
	a.logger.Error(method+" error", zap.Error(err), zap.Any("data", data))
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
func (a *AuthEndpoints) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		a.handleError(w, err, "Logout", nil)
		return
	}
	userID, err := a.grpcClient.GetUserIDBySession(cookie.Value)
	if err != nil {
		a.handleError(w, err, "Logout", nil)
		return
	}

	err = a.grpcClient.DeleteSession(cookie.Value)
	if err != nil {
		a.handleError(w, err, "Logout", map[string]string{"userID": userID})
		return
	}

	a.logger.Info("user logged out", zap.String("userID", userID))
	w.Header().Set("X-authenticated", "false")
	utils.SendJSONResponse(w, http.StatusOK, "You have successfully logged out")
}
