package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/storage"

	"github.com/gorilla/mux"
)

// GetAdvertsHandler godoc
// @Summary Get all adverts
// @Description Retrieves a list of all adverts
// @Tags adverts
// @Produce json
// @Success 200 {array} storage.Advert
// @Router /adverts [get]
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

// GetAdvertByIDHandler godoc
// @Summary Get an advert by ID
// @Description Retrieves a single advert by its ID
// @Tags adverts
// @Produce json
// @Param id path int true "Advert ID"
// @Success 200 {object} storage.Advert
// @Failure 400 {object} string "Bad Request"
// @Failure 404 {object} string "Not Found"
// @Router /adverts/{id} [get]
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

// AddAdvertHandler godoc
// @Summary Add a new advert
// @Description Adds a new advert to the list
// @Tags adverts
// @Accept json
// @Produce json
// @Param advert body storage.Advert true "Advert object"
// @Success 200 {object} storage.Advert
// @Failure 400 {object} string "Bad Request"
// @Router /adverts [post]
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

// UpdateAdvertHandler godoc
// @Summary Update an existing advert
// @Description Updates an existing advert by its ID
// @Tags adverts
// @Accept json
// @Produce json
// @Param id path int true "Advert ID"
// @Param advert body storage.Advert true "Updated Advert object"
// @Success 200 {object} storage.Advert
// @Failure 400 {object} string "Bad Request"
// @Failure 404 {object} string "Not Found"
// @Router /adverts/{id} [put]
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

// DeleteAdvertHandler godoc
// @Summary Delete an advert
// @Description Deletes an advert by its ID
// @Tags adverts
// @Param id path int true "Advert ID"
// @Success 204 "No Content"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /adverts/{id} [delete]
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
