package delivery

import (
	"context"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	delivery "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/service"
	"github.com/jackc/pgx/v5"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(recoveryMiddleware)

	var cfg config.Config
	conn, err := pgx.Connect(context.Background(), cfg.Postgres.GetConnectURL())
	if err != nil {
		return nil
	}

	userRepo := postgres.NewUserRepository(conn)
	userService := service.NewUserService(userRepo)
	userHandler := delivery.NewUserHandler(userService)

	router.HandleFunc("/api/v1/user", userHandler.GetUser).Methods("GET")
	router.HandleFunc("/api/v1/user", userHandler.DeleteUser).Methods("DELETE")
	router.HandleFunc("/api/v1/signup", userHandler.Signup).Methods("POST")

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

// func isAuthenticated(r *http.Request, authHandler *http3.AuthHandler) bool {
// 	cookie, err := r.Cookie("session_id")
// 	if err != nil || cookie == nil {
// 		log.Println("No session cookie found")

// 		return false
// 	}

// 	exists := authHandler.SessionRepo.SessionExists(cookie.Value)
// 	log.Printf("Session exists: %v for session_id: %s", exists, cookie.Value)

// 	return exists
// }

// func authMiddleware(authHandler *http3.AuthHandler) mux.MiddlewareFunc {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			if isAuthenticated(r, authHandler) {
// 				w.Header().Set("X-Authenticated", "true")
// 			} else {
// 				w.Header().Set("X-Authenticated", "false")
// 			}

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }
