package handlers

import (
	"emporium/internal/storage"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(recoveryMiddleware)

	userStorage := storage.NewUserStorage()
	sessionStorage := storage.NewSessionStorage()

	authHandler := &AuthHandler{
		UserStorage:    userStorage,
		SessionStorage: sessionStorage,
	}

	log.Println("Server is running")

	router.HandleFunc("/register", authHandler.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
	router.HandleFunc("/logout", authHandler.LogoutHandler).Methods("POST")

	return router
}

// Обработка паник
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