package router

import (
	"log"
	"net/http"

	advertsHandler "github.com/go-park-mail-ru/2024_2_BogoSort/internal/pkg/delivery/adverts"
	authHandler "github.com/go-park-mail-ru/2024_2_BogoSort/internal/pkg/delivery/auth"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(recoveryMiddleware)

	advertsHandler := advertsHandler.NewAdvertsHandler()
	authHandler := authHandler.NewAuthHandler()

	router.Use(authMiddleware(authHandler))

	router.HandleFunc("/api/v1/signup", authHandler.SignupHandler).Methods("POST")
	router.HandleFunc("/api/v1/login", authHandler.LoginHandler).Methods("POST")
	router.HandleFunc("/api/v1/logout", authHandler.LogoutHandler).Methods("POST")
	router.HandleFunc("/api/v1/adverts", advertsHandler.GetAdvertsHandler).Methods("GET")
	router.HandleFunc("/api/v1/adverts/{id}", advertsHandler.GetAdvertByIDHandler).Methods("GET")
	router.HandleFunc("/api/v1/adverts", advertsHandler.AddAdvertHandler).Methods("POST")
	router.HandleFunc("/api/v1/adverts/{id}", advertsHandler.UpdateAdvertHandler).Methods("PUT")
	router.HandleFunc("/api/v1/adverts/{id}", advertsHandler.DeleteAdvertHandler).Methods("DELETE")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return router
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Panic occurred:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func isAuthenticated(r *http.Request, authHandler *authHandler.AuthHandler) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie == nil {
		log.Println("No session cookie found")

		return false
	}

	exists := authHandler.SessionRepo.SessionExists(cookie.Value)
	log.Printf("Session exists: %v for session_id: %s", exists, cookie.Value)

	return exists
}

func authMiddleware(authHandler *authHandler.AuthHandler) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isAuthenticated(r, authHandler) {
				w.Header().Set("X-Authenticated", "true")
			} else {
				w.Header().Set("X-Authenticated", "false")
			}

			next.ServeHTTP(w, r)
		})
	}
}
