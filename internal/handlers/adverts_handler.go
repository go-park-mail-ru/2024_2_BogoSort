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

// GetAdvertsHandler godoc
// @Summary Get all adverts
// @Description Get a list of all adverts
// @Tags adverts
// @Produce json
// @Success 200 {array} storage.Advert
// @Router /api/v1/adverts [get]
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

	err := json.NewEncoder(w).Encode(adverts)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode adverts", http.StatusInternalServerError)
	}
}

// GetAdvertByIDHandler godoc
// @Summary Get an advert by ID
// @Description Get a single advert by its ID
// @Tags adverts
// @Produce json
// @Param id path int true "Advert ID"
// @Success 200 {object} storage.Advert
// @Failure 400 {object} responses.AuthErrResponse
// @Failure 404 {object} responses.AuthErrResponse
// @Router /api/v1/adverts/{id} [get]
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
	err = json.NewEncoder(w).Encode(advert)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode advert", http.StatusInternalServerError)
	}
}

// AddAdvertHandler godoc
// @Summary Add a new advert
// @Description Add a new advert to the list
// @Tags adverts
// @Accept json
// @Produce json
// @Param advert body storage.Advert true "Advert data"
// @Success 200 {object} storage.Advert
// @Failure 400 {object} responses.AuthErrResponse
// @Router /api/v1/adverts [post]
func (h *AdvertsHandler) AddAdvertHandler(w http.ResponseWriter, r *http.Request) {
	var advert storage.Advert
	err := json.NewDecoder(r.Body).Decode(&advert)
	if err != nil {
		responses.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	h.List.Add(&advert)
	err = json.NewEncoder(w).Encode(advert)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode advert", http.StatusInternalServerError)
	}
}

// UpdateAdvertHandler godoc
// @Summary Update an advert
// @Description Update an existing advert by its ID
// @Tags adverts
// @Accept json
// @Produce json
// @Param id path int true "Advert ID"
// @Param advert body storage.Advert true "Advert data"
// @Success 200 {object} storage.Advert
// @Failure 400 {object} responses.AuthErrResponse
// @Failure 404 {object} responses.AuthErrResponse
// @Router /api/v1/adverts/{id} [put]
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

	err = json.NewEncoder(w).Encode(advert)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode advert", http.StatusInternalServerError)
	}
}

// DeleteAdvertHandler godoc
// @Summary Delete an advert
// @Description Delete an advert by its ID
// @Tags adverts
// @Param id path int true "Advert ID"
// @Success 204
// @Failure 400 {object} responses.AuthErrResponse
// @Failure 500 {object} responses.AuthErrResponse
// @Router /api/v1/adverts/{id} [delete]
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