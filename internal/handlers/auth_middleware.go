package handlers

import (
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/responses"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		cookie, err := r.Cookie("session_id")
		if err == nil {
			tokenString = cookie.Value
		}

		if tokenString == "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			responses.SendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		_, err = utils.ValidateToken(tokenString)
		if err != nil {
			responses.SendErrorResponse(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		next.ServeHTTP(w, r)
	})
}