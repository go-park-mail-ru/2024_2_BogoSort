package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/responses"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/storage"

	"github.com/gorilla/mux"
)

var (
	ErrFailedToGetAdverts     = errors.New("failed to get adverts")
	ErrInvalidID              = errors.New("invalid ID")
	ErrAdvertNotFound         = errors.New("advert not found")
	ErrFailedToAddAdvert      = errors.New("failed to add advert")
	ErrFailedToUpdateAdvert   = errors.New("failed to update advert")
	ErrFailedToDeleteAdvert   = errors.New("failed to delete advert")
)

// GetAdvertsHandler godoc
// @Summary Get all adverts
// @Description Get a list of all adverts
// @Tags adverts
// @Produce json
// @Success 200 {array} storage.Advert "List of adverts"
// @Failure 500 {object} responses.ErrResponse "Failed to get adverts"
// @Router /api/v1/adverts [get]
func (authHandler *AdvertsHandler) GetAdvertsHandler(writer http.ResponseWriter, _ *http.Request) {
	adverts, err := authHandler.List.GetAdverts()
	if err != nil {
		responses.SendErrorResponse(writer, http.StatusInternalServerError, ErrFailedToGetAdverts.Error())
		return
	}

	for index := range adverts {
		imageURL, err := authHandler.ImageService.GetImageURL(adverts[index].ID)
		if err != nil {
			log.Println(err)

			continue
		}

		adverts[index].ImageURL = imageURL
	}

	responses.SendJSONResponse(writer, http.StatusOK, adverts)
}

// GetAdvertByIDHandler godoc
// @Summary Get an advert by ID
// @Description Get a single advert by its ID
// @Tags adverts
// @Produce json
// @Param id path int true "Advert ID"
// @Success 200 {object} storage.Advert "Advert details"
// @Failure 400 {object} responses.ErrResponse "Invalid ID"
// @Failure 404 {object} responses.ErrResponse "Advert not found"
// @Failure 500 {object} responses.ErrResponse "Failed to get advert"
// @Router /api/v1/adverts/{id} [get]
func (authHandler *AdvertsHandler) GetAdvertByIDHandler(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]
	uintID, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		responses.SendErrorResponse(writer, http.StatusBadRequest, ErrInvalidID.Error())
		return
	}

	advert, err := authHandler.List.GetAdvertByID(uint(uintID))
	if err != nil {
		responses.SendErrorResponse(writer, http.StatusNotFound, ErrAdvertNotFound.Error())
		return
	}

	imageURL, err := authHandler.ImageService.GetImageURL(advert.ID)
	if err != nil {
		log.Println("изображение не найдено:", err)
	} else {
		advert.ImageURL = imageURL
	}

	responses.SendJSONResponse(writer, http.StatusOK, advert)
}

// AddAdvertHandler godoc
// @Summary Add a new advert
// @Description Add a new advert to the list
// @Tags adverts
// @Accept json
// @Produce json
// @Param advert body storage.Advert true "Advert data"
// @Success 200 {object} storage.Advert "Advert added successfully"
// @Failure 400 {object} responses.ErrResponse "Failed to add advert"
// @Router /api/v1/adverts [post]
func (authHandler *AdvertsHandler) AddAdvertHandler(writer http.ResponseWriter, reader *http.Request) {
	var advert storage.Advert
	err := json.NewDecoder(reader.Body).Decode(&advert)

	if err != nil {
		responses.SendErrorResponse(writer, http.StatusBadRequest, ErrFailedToAddAdvert.Error())
		return
	}

	authHandler.List.Add(&advert)

	responses.SendJSONResponse(writer, http.StatusOK, advert)
}

// UpdateAdvertHandler godoc
// @Summary Update an advert
// @Description Update an existing advert by its ID
// @Tags adverts
// @Accept json
// @Produce json
// @Param id path int true "Advert ID"
// @Param advert body storage.Advert true "Advert data"
// @Success 200 {object} storage.Advert "Advert updated successfully"
// @Failure 400 {object} responses.ErrResponse "Invalid ID or data"
// @Failure 404 {object} responses.ErrResponse "Advert not found"
// @Failure 500 {object} responses.ErrResponse "Failed to update advert"
// @Router /api/v1/adverts/{id} [put]
func (authHandler *AdvertsHandler) UpdateAdvertHandler(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]
	uintID, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		responses.SendErrorResponse(writer, http.StatusBadRequest, ErrInvalidID.Error())
		return
	}

	var advert storage.Advert
	err = json.NewDecoder(reader.Body).Decode(&advert)

	if err != nil {
		responses.SendErrorResponse(writer, http.StatusBadRequest, ErrFailedToUpdateAdvert.Error())
		return
	}

	if advert.ID != uint(uintID) {
		responses.SendErrorResponse(writer, http.StatusBadRequest, "Id in URL and JSON do not match")

		return
	}

	err = authHandler.List.Update(&advert)
	if err != nil {
		responses.SendErrorResponse(writer, http.StatusNotFound, ErrFailedToUpdateAdvert.Error())

		return
	}

	responses.SendJSONResponse(writer, http.StatusOK, advert)
}

// DeleteAdvertHandler godoc
// @Summary Delete an advert
// @Description Delete an advert by its ID
// @Tags adverts
// @Param id path int true "Advert ID"
// @Success 204 "Advert deleted successfully"
// @Failure 400 {object} responses.ErrResponse "Invalid ID"
// @Failure 500 {object} responses.ErrResponse "Failed to delete advert"
// @Router /api/v1/adverts/{id} [delete]
func (authHandler *AdvertsHandler) DeleteAdvertHandler(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]
	uintID, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		responses.SendErrorResponse(writer, http.StatusBadRequest, ErrInvalidID.Error())

		return
	}

	err = authHandler.List.DeleteAdvert(uint(uintID))
	if err != nil {
		responses.SendErrorResponse(writer, http.StatusInternalServerError, ErrFailedToDeleteAdvert.Error())

		return
	}

	responses.SendJSONResponse(writer, http.StatusNoContent, nil)
}
