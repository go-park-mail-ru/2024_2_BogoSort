package http

import (
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/gorilla/mux"
)

var (
	ErrFailedToGetAdverts   = errors.New("failed to get adverts")
	ErrInvalidID            = errors.New("invalid ID")
	ErrAdvertNotFound       = errors.New("advert not found")
	ErrFailedToAddAdvert    = errors.New("failed to add advert")
	ErrFailedToUpdateAdvert = errors.New("failed to update advert")
	ErrFailedToDeleteAdvert = errors.New("failed to delete advert")
)

type AdvertEndpoints struct {
	AdvertsUseCase usecase.AdvertUseCase
	StaticUseCase  usecase.StaticUseCase
	logger         *zap.Logger
}

func NewAdvertEndpoints(advertsUseCase usecase.AdvertUseCase,
	staticUseCase usecase.StaticUseCase,
	logger *zap.Logger) *AdvertEndpoints {
	return &AdvertEndpoints{
		AdvertsUseCase: advertsUseCase,
		StaticUseCase:  staticUseCase,
		logger:         logger,
	}
}

func (h *AdvertEndpoints) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/adverts/{advertId}", h.GetAdvertById).Methods("GET")
	router.HandleFunc("/api/v1/adverts/seller/{sellerId}", h.GetAdvertsBySellerId).Methods("GET")
	router.HandleFunc("/api/v1/adverts/user/{userId}/saved", h.GetSavedAdvertsByUserId).Methods("GET")
	router.HandleFunc("/api/v1/adverts/cart/{cartId}", h.GetAdvertsByCartId).Methods("GET")
	router.HandleFunc("/api/v1/adverts", h.AddAdvert).Methods("POST")
	router.HandleFunc("/api/v1/adverts/{advertId}", h.UpdateAdvert).Methods("PUT")
	router.HandleFunc("/api/v1/adverts/{advertId}", h.DeleteAdvertById).Methods("DELETE")
	router.HandleFunc("/api/v1/adverts/{advertId}/status", h.UpdateAdvertStatus).Methods("PUT")
	router.HandleFunc("/api/v1/adverts/category/{categoryId}", h.GetAdvertsByCategoryId).Methods("GET")
	router.HandleFunc("/api/v1/adverts/{advertId}/image", h.UploadImage).Methods("PUT")
	router.HandleFunc("/api/v1/adverts", h.GetAdverts).Methods("GET")
}

// GetAdverts godoc
// @Summary Get all adverts
// @Description Get a list of all adverts
// @Tags adverts
// @Produce json
// @Param limit query int false "Limit the number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} dto.Advert "List of adverts"
// @Failure 400 {object} utils.ErrResponse "Invalid limit or offset"
// @Failure 500 {object} utils.ErrResponse "Failed to get adverts"
// @Router /api/v1/adverts [get]
func (h *AdvertEndpoints) GetAdverts(writer http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		h.logger.Error("invalid limit", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusBadRequest, "Invalid limit")
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		h.logger.Error("invalid offset", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusBadRequest, "Invalid offset")
		return
	}

	adverts, err := h.AdvertsUseCase.GetAdverts(limit, offset)
	if err != nil {
		h.logger.Error("failed to get adverts", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusInternalServerError, ErrFailedToGetAdverts.Error())
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, adverts)
}

// GetAdvertsBySellerId godoc
// @Summary Get adverts by seller ID
// @Description Get a list of adverts by seller ID
// @Tags adverts
// @Produce json
// @Param sellerId path string true "Seller ID"
// @Success 200 {array} dto.Advert "List of adverts"
// @Failure 400 {object} utils.ErrResponse "Invalid seller ID"
// @Failure 500 {object} utils.ErrResponse "Failed to get adverts by seller ID"
// @Router /api/v1/adverts/seller/{sellerId} [get]
func (h *AdvertEndpoints) GetAdvertsBySellerId(writer http.ResponseWriter, r *http.Request) {
	sellerIdStr := mux.Vars(r)["sellerId"]
	sellerId, err := uuid.Parse(sellerIdStr)
	if err != nil {
		h.logger.Error("invalid seller ID", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusBadRequest, "Invalid seller ID")
		return
	}

	adverts, err := h.AdvertsUseCase.GetAdvertsBySellerId(sellerId)
	if err != nil {
		h.logger.Error("failed to get adverts by seller ID", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusInternalServerError, "Failed to get adverts by seller ID")
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, adverts)
}

// GetSavedAdvertsByUserId godoc
// @Summary Get saved adverts by user ID
// @Description Get a list of saved adverts by user ID
// @Tags adverts
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {array} dto.Advert "List of saved adverts"
// @Failure 400 {object} utils.ErrResponse "Invalid user ID"
// @Failure 500 {object} utils.ErrResponse "Failed to get saved adverts by user ID"
// @Router /api/v1/adverts/user/{userId}/saved [get]
func (h *AdvertEndpoints) GetSavedAdvertsByUserId(writer http.ResponseWriter, r *http.Request) {
	userIdStr := mux.Vars(r)["userId"]
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		h.logger.Error("invalid user ID", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusBadRequest, "Invalid user ID")
		return
	}

	adverts, err := h.AdvertsUseCase.GetSavedAdvertsByUserId(userId)
	if err != nil {
		h.logger.Error("failed to get saved adverts by user ID", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusInternalServerError, "Failed to get saved adverts by user ID")
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, adverts)
}

// GetAdvertsByCartId godoc
// @Summary Get adverts by cart ID
// @Description Get a list of adverts in the specified cart
// @Tags adverts
// @Produce json
// @Param cartId path string true "Cart ID"
// @Success 200 {array} dto.Advert "List of adverts in cart"
// @Failure 400 {object} utils.ErrResponse "Invalid cart ID"
// @Failure 500 {object} utils.ErrResponse "Failed to get adverts by cart ID"
// @Router /api/v1/adverts/cart/{cartId} [get]
func (h *AdvertEndpoints) GetAdvertsByCartId(writer http.ResponseWriter, r *http.Request) {
	cartIdStr := mux.Vars(r)["cartId"]
	cartId, err := uuid.Parse(cartIdStr)
	if err != nil {
		h.logger.Error("invalid cart ID", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusBadRequest, "Invalid cart ID")
		return
	}

	adverts, err := h.AdvertsUseCase.GetAdvertsByCartId(cartId)
	if err != nil {
		h.logger.Error("failed to get adverts by cart ID", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusInternalServerError, "Failed to get adverts by cart ID")
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, adverts)
}

// GetAdvertById godoc
// @Summary Get an advert by ID
// @Description Get an advert by its ID
// @Tags adverts
// @Produce json
// @Param advertId path string true "Advert ID"
// @Success 200 {object} dto.Advert "Advert details"
// @Failure 400 {object} utils.ErrResponse "Invalid advert ID"
// @Failure 404 {object} utils.ErrResponse "Advert not found"
// @Failure 500 {object} utils.ErrResponse "Failed to get advert by ID"
// @Router /api/v1/adverts/{advertId} [get]
func (h *AdvertEndpoints) GetAdvertById(writer http.ResponseWriter, r *http.Request) {
	advertIdStr := mux.Vars(r)["advertId"]
	advertId, err := uuid.Parse(advertIdStr)
	if err != nil {
		h.logger.Error("http: invalid advert ID", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusBadRequest, "Invalid advert ID")
		return
	}

	advert, err := h.AdvertsUseCase.GetAdvertById(advertId)

	if err != nil {
		if errors.Is(err, ErrAdvertNotFound) {
			h.logger.Error("http: advert not found", zap.String("advert_id", advertId.String()))
			utils.SendErrorResponse(writer, http.StatusNotFound, "Advert not found")
		} else {
			h.logger.Error("http: failed to get advert by ID", zap.Error(err))
			utils.SendErrorResponse(writer, http.StatusInternalServerError, "Failed to get advert by ID")
		}
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, advert)
}

// AddAdvert godoc
// @Summary Add a new advert
// @Description Create a new advert
// @Tags adverts
// @Accept json
// @Produce json
// @Param advert body dto.Advert true "Advert data"
// @Success 201 {object} dto.Advert "Advert created"
// @Failure 400 {object} utils.ErrResponse "Invalid advert data"
// @Failure 500 {object} utils.ErrResponse "Failed to add advert"
// @Router /api/v1/adverts [post]
func (h *AdvertEndpoints) AddAdvert(writer http.ResponseWriter, r *http.Request) {
	var advert dto.Advert
	if err := json.NewDecoder(r.Body).Decode(&advert); err != nil {
		h.logger.Error("invalid advert data", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusBadRequest, "Invalid advert data")
		return
	}

	newAdvert, err := h.AdvertsUseCase.AddAdvert(&advert)
	if err != nil {
		h.logger.Error("failed to add advert", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusInternalServerError, "Failed to add advert")
		return
	}

	utils.SendJSONResponse(writer, http.StatusCreated, newAdvert)
}

// UpdateAdvert godoc
// @Summary Update an existing advert
// @Description Update an advert's information
// @Tags adverts
// @Accept json
// @Produce json
// @Param advertId path string true "Advert ID"
// @Param advert body dto.Advert true "Updated advert data"
// @Success 200 "Advert updated successfully"
// @Failure 400 {object} utils.ErrResponse "Invalid advert data"
// @Failure 404 {object} utils.ErrResponse "Advert not found"
// @Failure 500 {object} utils.ErrResponse "Failed to update advert"
// @Router /api/v1/adverts/{advertId} [put]
func (h *AdvertEndpoints) UpdateAdvert(writer http.ResponseWriter, r *http.Request) {
	var advert dto.Advert
	if err := json.NewDecoder(r.Body).Decode(&advert); err != nil {
		h.logger.Error("invalid advert data", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusBadRequest, "Invalid advert data")
		return
	}

	if err := h.AdvertsUseCase.UpdateAdvert(&advert); err != nil {
		if errors.Is(err, ErrAdvertNotFound) {
			utils.SendErrorResponse(writer, http.StatusNotFound, "Advert not found")
		} else {
			h.logger.Error("failed to update advert", zap.Error(err))
			utils.SendErrorResponse(writer, http.StatusInternalServerError, "Failed to update advert")
		}
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, "Advert updated successfully")
}

// DeleteAdvertById godoc
// @Summary Delete an advert by ID
// @Description Delete an advert by its ID
// @Tags adverts
// @Param advertId path string true "Advert ID"
// @Success 204 "Advert deleted"
// @Failure 400 {object} utils.ErrResponse "Invalid advert ID"
// @Failure 404 {object} utils.ErrResponse "Advert not found"
// @Failure 500 {object} utils.ErrResponse "Failed to delete advert"
// @Router /api/v1/adverts/{advertId} [delete]
func (h *AdvertEndpoints) DeleteAdvertById(writer http.ResponseWriter, r *http.Request) {
	advertIdStr := mux.Vars(r)["advertId"]
	advertId, err := uuid.Parse(advertIdStr)
	if err != nil {
		h.logger.Error("invalid advert ID", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusBadRequest, "Invalid advert ID")
		return
	}

	if err := h.AdvertsUseCase.DeleteAdvertById(advertId); err != nil {
		if errors.Is(err, ErrAdvertNotFound) {
			utils.SendErrorResponse(writer, http.StatusNotFound, "Advert not found")
		} else {
			h.logger.Error("failed to delete advert", zap.Error(err))
			utils.SendErrorResponse(writer, http.StatusInternalServerError, "Failed to delete advert")
		}
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, "Advert deleted")
}

// UpdateAdvertStatus godoc
// @Summary Update advert status
// @Description Update advert status by ID
// @Tags adverts
// @Param advertId path string true "Advert ID"
// @Param status body string true "New status"
// @Success 200 "Advert status updated"
// @Failure 400 {object} utils.ErrResponse "Invalid advert ID or status"
// @Failure 404 {object} utils.ErrResponse "Advert not found"
// @Failure 500 {object} utils.ErrResponse "Failed to update advert status"
// @Router /api/v1/adverts/{advertId}/status [put]
func (h *AdvertEndpoints) UpdateAdvertStatus(writer http.ResponseWriter, r *http.Request) {
	advertIdStr := mux.Vars(r)["advertId"]
	advertId, err := uuid.Parse(advertIdStr)
	if err != nil {
		h.logger.Error("invalid advert ID", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusBadRequest, "Invalid advert ID")
		return
	}

	status := r.FormValue("status")
	if status != string(dto.AdvertStatusActive) && status != string(dto.AdvertStatusInactive) {
		h.logger.Error("invalid advert status", zap.String("status", status))
		utils.SendErrorResponse(writer, http.StatusBadRequest, "Invalid advert status")
		return
	}

	if err := h.AdvertsUseCase.UpdateAdvertStatus(advertId, status); err != nil {
		if errors.Is(err, ErrAdvertNotFound) {
			utils.SendErrorResponse(writer, http.StatusNotFound, "Advert not found")
		} else {
			h.logger.Error("failed to update advert status", zap.Error(err))
			utils.SendErrorResponse(writer, http.StatusInternalServerError, "Failed to update advert status")
		}
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, "Advert status updated")
}

// GetAdvertsByCategoryId godoc
// @Summary Get adverts by category ID
// @Description Get a list of adverts by category ID
// @Tags adverts
// @Produce json
// @Param categoryId path string true "Category ID"
// @Success 200 {array} dto.Advert "List of adverts by category ID"
// @Failure 400 {object} utils.ErrResponse "Invalid category ID"
// @Failure 500 {object} utils.ErrResponse "Failed to get adverts by category ID"
// @Router /api/v1/adverts/category/{categoryId} [get]
func (h *AdvertEndpoints) GetAdvertsByCategoryId(writer http.ResponseWriter, r *http.Request) {
	categoryIdStr := mux.Vars(r)["categoryId"]
	categoryId, err := uuid.Parse(categoryIdStr)
	if err != nil {
		h.logger.Error("invalid category ID", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusBadRequest, "Invalid category ID")
		return
	}

	adverts, err := h.AdvertsUseCase.GetAdvertsByCategoryId(categoryId)
	if err != nil {
		h.logger.Error("failed to get adverts by category ID", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusInternalServerError, "Failed to get adverts by category ID")
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, adverts)
}

// UploadImage godoc
// @Summary Upload an image
// @Description Upload an image by ID
// @Tags adverts
// @Param advertId path string true "Advert ID"
// @Param image formData file true "Image file to upload"
// @Success 200 {string} string "Image uploaded"
// @Failure 400 {object} utils.ErrResponse "Invalid advert ID or file not attached"
// @Failure 500 {object} utils.ErrResponse "Failed to upload image"
// @Router /api/v1/adverts/{advertId}/image [put]
func (h *AdvertEndpoints) UploadImage(writer http.ResponseWriter, r *http.Request) {
	advertIdStr := mux.Vars(r)["advertId"]
	advertId, err := uuid.Parse(advertIdStr)
	if err != nil {
		h.logger.Error("invalid advert ID", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusBadRequest, "Invalid advert ID")
		return
	}

	fileHeader, _, err := r.FormFile("image")
	if err != nil {
		h.logger.Error("file not attached", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusBadRequest, "File not attached")
		return
	}

	data, err := io.ReadAll(fileHeader)
	if err != nil {
		h.logger.Error("failed to read file", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusInternalServerError, "Failed to read file")
		return
	}

	if err = fileHeader.Close(); err != nil {
		h.logger.Error("failed to close file", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusInternalServerError, "Failed to close file")
		return
	}

	imageId, err := h.StaticUseCase.UploadFile(data)
	if err != nil {
		h.logger.Error("failed to upload image", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusInternalServerError, "Failed to upload image")
		return
	}

	if err := h.AdvertsUseCase.UploadImage(advertId, imageId); err != nil {
		h.logger.Error("failed to upload image", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusInternalServerError, "Failed to upload image")
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, "Image uploaded")
}
