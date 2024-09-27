package handlers

import (
	"emporium/internal/storage"
	"log"

	"github.com/gorilla/mux"
)

type AdvertsHandler struct {
	List *storage.AdvertsList
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	advertsList := storage.NewAdvertsList()
	storage.FillAdverts(advertsList)
	advertsHandler := &AdvertsHandler{
		List: advertsList,
	}

	log.Println("Server is running")
	router.HandleFunc("/adverts", advertsHandler.GetAdvertsHandler).Methods("GET")
	router.HandleFunc("/adverts/{id}", advertsHandler.GetAdvertByIDHandler).Methods("GET")
	router.HandleFunc("/adverts", advertsHandler.AddAdvertHandler).Methods("POST")
	router.HandleFunc("/adverts/{id}", advertsHandler.UpdateAdvertHandler).Methods("PUT")
	router.HandleFunc("/adverts/{id}", advertsHandler.DeleteAdvertHandler).Methods("DELETE")

	return router
}
