package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/responses"
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
		responses.SendErrorResponse(w, http.StatusBadRequest, "Invalid ID")
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
		responses.SendErrorResponse(w, http.StatusBadRequest, err.Error())
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
		responses.SendErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var advert storage.Advert
	err = json.NewDecoder(r.Body).Decode(&advert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if advert.ID != uint(uintID) {
		responses.SendErrorResponse(w, http.StatusBadRequest, "Id in URL and JSON do not match")
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
		responses.SendErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	err = h.List.DeleteAdvert(uint(uintID))
	if err != nil {
		responses.SendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
