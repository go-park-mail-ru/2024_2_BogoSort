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
// @Failure 500 {object} responses.ErrResponse
// @Router /api/v1/adverts [get]
func (authHandler *AdvertsHandler) GetAdvertsHandler(writer http.ResponseWriter, _ *http.Request) {
	adverts := authHandler.List.GetAdverts()

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
// @Success 200 {object} storage.Advert
// @Failure 400 {object} responses.ErrResponse
// @Failure 404 {object} responses.ErrResponse
// @Failure 500 {object} responses.ErrResponse
// @Router /api/v1/adverts/{id} [get]
func (authHandler *AdvertsHandler) GetAdvertByIDHandler(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]
	uintID, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		responses.SendErrorResponse(writer, http.StatusBadRequest, "Invalid ID")

		return
	}

	advert, err := authHandler.List.GetAdvertByID(uint(uintID))
	if err != nil {
		responses.SendErrorResponse(writer, http.StatusNotFound, "Advert not found")

		return
	}

	imageURL, err := authHandler.ImageService.GetImageURL(advert.ID)

	if err == nil {
		advert.ImageURL = imageURL
	} else {
		log.Println(err)
		responses.SendErrorResponse(writer, http.StatusInternalServerError, "Internal server error")

		return
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
// @Success 200 {object} storage.Advert
// @Failure 400 {object} responses.ErrResponse
// @Router /api/v1/adverts [post]
func (authHandler *AdvertsHandler) AddAdvertHandler(writer http.ResponseWriter, reader *http.Request) {
	var advert storage.Advert
	err := json.NewDecoder(reader.Body).Decode(&advert)

	if err != nil {
		responses.SendErrorResponse(writer, http.StatusBadRequest, err.Error())

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
// @Success 200 {object} storage.Advert
// @Failure 400 {object} responses.ErrResponse
// @Failure 404 {object} responses.ErrResponse
// @Failure 500 {object} responses.ErrResponse
// @Router /api/v1/adverts/{id} [put]
func (authHandler *AdvertsHandler) UpdateAdvertHandler(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]
	uintID, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		responses.SendErrorResponse(writer, http.StatusBadRequest, "Invalid ID")

		return
	}

	var advert storage.Advert
	err = json.NewDecoder(reader.Body).Decode(&advert)

	if err != nil {
		responses.SendErrorResponse(writer, http.StatusBadRequest, err.Error())

		return
	}

	if advert.ID != uint(uintID) {
		responses.SendErrorResponse(writer, http.StatusBadRequest, "Id in URL and JSON do not match")

		return
	}

	err = authHandler.List.Update(&advert)
	if err != nil {
		responses.SendErrorResponse(writer, http.StatusNotFound, err.Error())

		return
	}

	responses.SendJSONResponse(writer, http.StatusOK, advert)
}

// DeleteAdvertHandler godoc
// @Summary Delete an advert
// @Description Delete an advert by its ID
// @Tags adverts
// @Param id path int true "Advert ID"
// @Success 204
// @Failure 400 {object} responses.ErrResponse
// @Failure 500 {object} responses.ErrResponse
// @Router /api/v1/adverts/{id} [delete]
func (authHandler *AdvertsHandler) DeleteAdvertHandler(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]
	uintID, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		responses.SendErrorResponse(writer, http.StatusBadRequest, "Invalid ID")

		return
	}

	err = authHandler.List.DeleteAdvert(uint(uintID))
	if err != nil {
		responses.SendErrorResponse(writer, http.StatusInternalServerError, "Internal server error")

		return
	}

	responses.SendJSONResponse(writer, http.StatusNoContent, nil)
}