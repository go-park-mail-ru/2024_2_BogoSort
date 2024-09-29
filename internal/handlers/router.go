package handlers

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/storage"
  "github.com/go-park-mail-ru/2024_2_BogoSort/internal/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type AdvertsHandler struct {
	List         *storage.AdvertsList
	ImageService *services.ImageService
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(recoveryMiddleware)

	userStorage := storage.NewUserStorage()
	sessionStorage := storage.NewSessionStorage()
  	advertsList := storage.NewAdvertsList()
	imageService := services.NewImageService()
	storage.FillAdverts(advertsList, imageService)

	advertsHandler := &AdvertsHandler{
		List:         advertsList,
		ImageService: imageService,
	}

	authHandler := &AuthHandler{
		UserStorage:    userStorage,
		SessionStorage: sessionStorage,
	}

	log.Println("Server is running")

	router.HandleFunc("/register", authHandler.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
	router.HandleFunc("/logout", authHandler.LogoutHandler).Methods("POST")
  	router.HandleFunc("/adverts", advertsHandler.GetAdvertsHandler).Methods("GET")
	router.HandleFunc("/adverts/{id}", advertsHandler.GetAdvertByIDHandler).Methods("GET")
	router.HandleFunc("/adverts", advertsHandler.AddAdvertHandler).Methods("POST")
	router.HandleFunc("/adverts/{id}", advertsHandler.UpdateAdvertHandler).Methods("PUT")
	router.HandleFunc("/adverts/{id}", advertsHandler.DeleteAdvertHandler).Methods("DELETE")
  
  fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

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
