package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/storage"

	"github.com/gorilla/mux"
)

func (h *AdvertsHandler) GetAdvertsHandler(w http.ResponseWriter, r *http.Request) {
	adverts := h.List.GetAdverts()

	for i := range adverts {
		imageURL, err := h.ImageService.GetImageURL(adverts[i].ID)
		if err != nil {
			log.Println(err)
			continue
		}
		adverts[i].ImageURL = imageURL
	}

	json.NewEncoder(w).Encode(adverts)
}

func (h *AdvertsHandler) GetAdvertByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}
	advert, err := h.List.GetAdvertByID(uint(uintID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(advert)
}

func (h *AdvertsHandler) AddAdvertHandler(w http.ResponseWriter, r *http.Request) {
	var advert storage.Advert
	err := json.NewDecoder(r.Body).Decode(&advert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.List.Add(&advert)
	json.NewEncoder(w).Encode(advert)
}

func (h *AdvertsHandler) UpdateAdvertHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	var advert storage.Advert
	err = json.NewDecoder(r.Body).Decode(&advert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if advert.ID != uint(uintID) {
		http.Error(w, "ID в URL и JSON не совпадают", http.StatusBadRequest)
		return
	}

	err = h.List.Update(&advert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(advert)
}

func (h *AdvertsHandler) DeleteAdvertHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	err = h.List.DeleteAdvert(uint(uintID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
