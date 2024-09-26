package handlers

import (
	"emporium/internal/storage"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetAdvertsHandler(w http.ResponseWriter, r *http.Request) {
	adverts := storage.GetAdverts()
	json.NewEncoder(w).Encode(adverts)
}

func GetAdvertByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}
	advert, err := storage.GetAdvertByID(uint(uintID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(advert)
}

func AddAdvertHandler(w http.ResponseWriter, r *http.Request) {
	var advert storage.Advert
	err := json.NewDecoder(r.Body).Decode(&advert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	storage.AddAdvert(advert)
	json.NewEncoder(w).Encode(advert)
}
