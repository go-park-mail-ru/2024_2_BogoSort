package handlers

import (
	"emporium/internal/services"
	"emporium/internal/storage"
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

	advertsList := storage.NewAdvertsList()
	imageService := services.NewImageService()
	storage.FillAdverts(advertsList, imageService)

	advertsHandler := &AdvertsHandler{
		List:         advertsList,
		ImageService: imageService,
	}

	log.Println("Server is running")
	router.HandleFunc("/adverts", advertsHandler.GetAdvertsHandler).Methods("GET")
	router.HandleFunc("/adverts/{id}", advertsHandler.GetAdvertByIDHandler).Methods("GET")
	router.HandleFunc("/adverts", advertsHandler.AddAdvertHandler).Methods("POST")
	router.HandleFunc("/adverts/{id}", advertsHandler.UpdateAdvertHandler).Methods("PUT")
	router.HandleFunc("/adverts/{id}", advertsHandler.DeleteAdvertHandler).Methods("DELETE")

	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	return router
}
