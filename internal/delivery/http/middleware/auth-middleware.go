package middleware

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
)

type AuthMiddleware struct {
	sessionManager *utils.SessionManager
}

func NewAuthMiddleware(sm *utils.SessionManager) *AuthMiddleware {
	return &AuthMiddleware{
		sessionManager: sm,
	}
}

func (m *AuthMiddleware) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := m.sessionManager.GetUserID(r)
		if err != nil {
			if errors.Is(err, utils.ErrSessionExpired) {
				utils.SendErrorResponse(w, http.StatusUnauthorized, "Session has expired")
				return
			}
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := m.sessionManager.GetUserID(r)
		if err != nil {
			w.Header().Set("X-authenticated", "false")
		} else {
			w.Header().Set("X-authenticated", "true")
		}
		next.ServeHTTP(w, r)
	})
}
