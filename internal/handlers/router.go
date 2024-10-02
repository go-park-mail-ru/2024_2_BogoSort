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

type AuthHandler struct {
	UserStorage *storage.UserStorage
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(recoveryMiddleware)

	userStorage := storage.NewUserStorage()
	advertsList := storage.NewAdvertsList()
	imageService := services.NewImageService()
	storage.FillAdverts(advertsList, imageService)

	advertsHandler := &AdvertsHandler{
		List:         advertsList,
		ImageService: imageService,
	}

	authHandler := &AuthHandler{
		UserStorage: userStorage,
	}

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
