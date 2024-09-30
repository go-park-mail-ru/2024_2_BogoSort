package handlers

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/services"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/storage"
	"github.com/gorilla/mux"
)

type AdvertsHandler struct {
	List         *storage.AdvertsList
	ImageService *services.ImageService
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	userStorage := storage.NewUserStorage()
	authHandler := &AuthHandler{UserStorage: userStorage}
	advertsList := storage.NewAdvertsList()
	imageService := services.NewImageService()
	advertsHandler := &AdvertsHandler{List: advertsList, ImageService: imageService}

	router.HandleFunc("/api/v1/signup", authHandler.SignupHandler).Methods("POST")
	router.HandleFunc("/api/v1/login", authHandler.LoginHandler).Methods("POST")
	router.HandleFunc("/api/v1/adverts", advertsHandler.GetAdvertsHandler).Methods("GET")

	protected := router.PathPrefix("/api/v1").Subrouter()
	protected.Use(AuthMiddleware)
	protected.HandleFunc("/logout", authHandler.LogoutHandler).Methods("POST")
	protected.HandleFunc("/adverts/{id}", advertsHandler.GetAdvertByIDHandler).Methods("GET")
	protected.HandleFunc("/adverts", advertsHandler.AddAdvertHandler).Methods("POST")
	protected.HandleFunc("/adverts/{id}", advertsHandler.UpdateAdvertHandler).Methods("PUT")
	protected.HandleFunc("/adverts/{id}", advertsHandler.DeleteAdvertHandler).Methods("DELETE")

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
