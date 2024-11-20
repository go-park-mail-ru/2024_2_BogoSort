package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func CSRFMiddleware(tk *utils.CryptToken, sm *utils.SessionManager) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
				token := r.Header.Get("X-CSRF-Token")
				if token == "" {
					http.Error(w, "CSRF token missing", http.StatusForbidden)
					return
				}

				sessionIDCookie, err := r.Cookie("session_id")
				if err != nil {
					http.Error(w, "Session ID missing", http.StatusForbidden)
					return
				}

				userID, err := sm.GetUserID(r)
				if err != nil {
					http.Error(w, "Invalid session", http.StatusForbidden)
					return
				}

				valid, err := tk.Check(uuid.MustParse(sessionIDCookie.Value), userID, token)
				if err != nil || !valid {
					http.Error(w, "Invalid CSRF token", http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
