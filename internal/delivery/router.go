package delivery

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	http3 "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/connector"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(recoveryMiddleware)

	var cfg config.Config
	psql, err := connector.GetPostgresConnector(cfg.Postgres.GetConnectURL())

	if err != nil {
		return err
	}

	userRepo := postgres.NewUserRepository(psql)

	advertsHandler := http3.NewAdvertsHandler()
	authHandler := advertsHandler.NewAuthHandler()

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

func isAuthenticated(r *http.Request, authHandler *http3.AuthHandler) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie == nil {
		log.Println("No session cookie found")

		return false
	}

	exists := authHandler.SessionRepo.SessionExists(cookie.Value)
	log.Printf("Session exists: %v for session_id: %s", exists, cookie.Value)

	return exists
}

func authMiddleware(authHandler *http3.AuthHandler) mux.MiddlewareFunc {
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
