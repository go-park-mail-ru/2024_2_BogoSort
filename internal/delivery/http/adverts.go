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
	ErrForbidden            = errors.New("forbidden")
)

type AdvertEndpoints struct {
	advertUseCase  usecase.AdvertUseCase
	staticUseCase  usecase.StaticUseCase
	sessionManager *utils.SessionManager
	logger         *zap.Logger
}

func NewAdvertEndpoints(advertUseCase usecase.AdvertUseCase,
	staticUseCase usecase.StaticUseCase,
	sessionManager *utils.SessionManager,
	logger *zap.Logger) *AdvertEndpoints {
	return &AdvertEndpoints{
		advertUseCase: advertUseCase,
		staticUseCase:  staticUseCase,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

func (h *AdvertEndpoints) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/adverts/{advertId}", h.GetAdvertById).Methods("GET")
	router.HandleFunc("/api/v1/adverts/seller/{sellerId}", h.GetAdvertsBySellerId).Methods("GET")
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
// @Summary Retrieve all adverts
// @Description Fetch a list of all adverts with optional pagination.
// @Tags adverts
// @Produce json
// @Param limit query int false "Limit the number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} dto.AdvertResponse "List of adverts"
// @Failure 400 {object} utils.ErrResponse "Invalid limit or offset"
// @Failure 500 {object} utils.ErrResponse "Failed to retrieve adverts"
// @Router /api/v1/adverts [get]
func (h *AdvertEndpoints) GetAdverts(writer http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		h.sendError(writer, http.StatusBadRequest, err, "invalid limit", nil)
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		h.sendError(writer, http.StatusBadRequest, err, "invalid offset", nil)
		return
	}

	adverts, err := h.advertUseCase.GetAdverts(limit, offset)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to get adverts", nil)
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, adverts)
}

// GetAdvertsBySellerId godoc
// @Summary Retrieve adverts by seller ID
// @Description Fetch a list of adverts associated with a specific seller ID.
// @Tags adverts
// @Produce json
// @Param sellerId path string true "Seller ID"
// @Success 200 {array} dto.AdvertResponse "List of adverts"
// @Failure 400 {object} utils.ErrResponse "Invalid seller ID"
// @Failure 403 {object} utils.ErrResponse "Forbidden"
// @Failure 500 {object} utils.ErrResponse "Failed to retrieve adverts by seller ID"
// @Router /api/v1/adverts/seller/{sellerId} [get]
func (h *AdvertEndpoints) GetAdvertsBySellerId(writer http.ResponseWriter, r *http.Request) {
	sellerIdStr := mux.Vars(r)["sellerId"]
	sellerId, err := uuid.Parse(sellerIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, err, "invalid seller ID", nil)
		return
	}

	adverts, err := h.advertUseCase.GetAdvertsByUserId(sellerId)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to get adverts by seller ID", nil)
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, adverts)
}

// GetAdvertsByCartId godoc
// @Summary Retrieve adverts by cart ID
// @Description Fetch a list of adverts in the specified cart.
// @Tags adverts
// @Produce json
// @Param cartId path string true "Cart ID"
// @Success 200 {array} dto.AdvertResponse "List of adverts in cart"
// @Failure 400 {object} utils.ErrResponse "Invalid cart ID"
// @Failure 403 {object} utils.ErrResponse "Forbidden"
// @Failure 500 {object} utils.ErrResponse "Failed to retrieve adverts by cart ID"
// @Router /api/v1/adverts/cart/{cartId} [get]
func (h *AdvertEndpoints) GetAdvertsByCartId(writer http.ResponseWriter, r *http.Request) {
	cartIdStr := mux.Vars(r)["cartId"]
	cartId, err := uuid.Parse(cartIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, err, "invalid cart ID", nil)
		return
	}

	adverts, err := h.advertUseCase.GetAdvertsByCartId(cartId)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to get adverts by cart ID", nil)
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, adverts)
}

// GetAdvertById godoc
// @Summary Retrieve an advert by ID
// @Description Fetch an advert based on its ID.
// @Tags adverts
// @Produce json
// @Param advertId path string true "Advert ID"
// @Success 200 {object} dto.AdvertResponse "Advert details"
// @Failure 400 {object} utils.ErrResponse "Invalid advert ID"
// @Failure 404 {object} utils.ErrResponse "Advert not found"
// @Failure 500 {object} utils.ErrResponse "Failed to retrieve advert by ID"
// @Router /api/v1/adverts/{advertId} [get]
func (h *AdvertEndpoints) GetAdvertById(writer http.ResponseWriter, r *http.Request) {
	advertIdStr := mux.Vars(r)["advertId"]
	advertId, err := uuid.Parse(advertIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, err, "invalid advert ID", nil)
		return
	}

	advert, err := h.advertUseCase.GetAdvertById(advertId)
	if err != nil {
		h.handleError(writer, err, "failed to get advert by ID")
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, advert)
}

// AddAdvert godoc
// @Summary Create a new advert
// @Description Add a new advert to the system.
// @Tags adverts
// @Accept json
// @Produce json
// @Param advert body dto.AdvertRequest true "Advert data"
// @Success 201 {object} dto.AdvertResponse "Advert created"
// @Failure 400 {object} utils.ErrResponse "Invalid advert data"
// @Failure 500 {object} utils.ErrResponse "Failed to create advert"
// @Router /api/v1/adverts [post]
func (h *AdvertEndpoints) AddAdvert(writer http.ResponseWriter, r *http.Request) {
	var advert dto.AdvertRequest
	if err := json.NewDecoder(r.Body).Decode(&advert); err != nil {
		h.sendError(writer, http.StatusBadRequest, err, "invalid advert data", nil)
		return
	}

	userID, err := h.sessionManager.GetUserID(r)
	if err != nil {
		h.sendError(writer, http.StatusUnauthorized, err, "user not found", nil)
		return
	}

	newAdvert, err := h.advertUseCase.AddAdvert(&advert, userID)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to add advert", nil)
		return
	}

	utils.SendJSONResponse(writer, http.StatusCreated, newAdvert)
}

// UpdateAdvert godoc
// @Summary Update an existing advert
// @Description Modify the details of an existing advert.
// @Tags adverts
// @Accept json
// @Produce json
// @Param advertId path string true "Advert ID"
// @Param advert body dto.AdvertRequest true "Updated advert data"
// @Success 200 "Advert updated successfully"
// @Failure 400 {object} utils.ErrResponse "Invalid advert data"
// @Failure 404 {object} utils.ErrResponse "Advert not found"
// @Failure 403 {object} utils.ErrResponse "Forbidden"
// @Failure 500 {object} utils.ErrResponse "Failed to update advert"
// @Router /api/v1/adverts/{advertId} [put]
func (h *AdvertEndpoints) UpdateAdvert(writer http.ResponseWriter, r *http.Request) {
	var advert dto.AdvertRequest
	if err := json.NewDecoder(r.Body).Decode(&advert); err != nil {
		h.sendError(writer, http.StatusBadRequest, err, "invalid advert data", nil)
		return
	}

	userID, err := h.sessionManager.GetUserID(r)
	if err != nil {
		h.sendError(writer, http.StatusUnauthorized, err, "user not found", nil)
		return
	}

	advertIdStr := mux.Vars(r)["advertId"]
	advertId, err := uuid.Parse(advertIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, err, "invalid advert ID", nil)
		return
	}

	if err := h.advertUseCase.UpdateAdvert(&advert, userID, advertId); err != nil {
		h.handleError(writer, err, "failed to update advert")
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, "Advert updated successfully")
}

// DeleteAdvertById godoc
// @Summary Delete an advert by ID
// @Description Remove an advert from the system using its ID.
// @Tags adverts
// @Param advertId path string true "Advert ID"
// @Success 204 "Advert deleted"
// @Failure 400 {object} utils.ErrResponse "Invalid advert ID"
// @Failure 404 {object} utils.ErrResponse "Advert not found"
// @Failure 403 {object} utils.ErrResponse "Forbidden"
// @Failure 500 {object} utils.ErrResponse "Failed to delete advert"
// @Router /api/v1/adverts/{advertId} [delete]
func (h *AdvertEndpoints) DeleteAdvertById(writer http.ResponseWriter, r *http.Request) {
	advertIdStr := mux.Vars(r)["advertId"]
	advertId, err := uuid.Parse(advertIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, err, "invalid advert ID", nil)
		return
	}

	userID, err := h.sessionManager.GetUserID(r)
	if err != nil {
		h.sendError(writer, http.StatusUnauthorized, err, "user not found", nil)
		return
	}

	if err := h.advertUseCase.DeleteAdvertById(advertId, userID); err != nil {
		h.handleError(writer, err, "failed to delete advert")
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, "Advert deleted")
}

// UpdateAdvertStatus godoc
// @Summary Update the status of an advert
// @Description Change the status of an advert by its ID.
// @Tags adverts
// @Param advertId path string true "Advert ID"
// @Param status body string true "New status"
// @Success 200 "Advert status updated"
// @Failure 400 {object} utils.ErrResponse "Invalid advert ID or status"
// @Failure 404 {object} utils.ErrResponse "Advert not found"
// @Failure 403 {object} utils.ErrResponse "Forbidden"
// @Failure 500 {object} utils.ErrResponse "Failed to update advert status"
// @Router /api/v1/adverts/{advertId}/status [put]
func (h *AdvertEndpoints) UpdateAdvertStatus(writer http.ResponseWriter, r *http.Request) {
	advertIdStr := mux.Vars(r)["advertId"]
	advertId, err := uuid.Parse(advertIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, err, "invalid advert ID", nil)
		return
	}

	status := r.FormValue("status")
	if status != string(dto.AdvertStatusActive) && status != string(dto.AdvertStatusInactive) {
		h.sendError(writer, http.StatusBadRequest, err, "invalid advert status", nil)
		return
	}

	userID, err := h.sessionManager.GetUserID(r)
	if err != nil {
		h.sendError(writer, http.StatusUnauthorized, err, "user not found", nil)
		return
	}

	if err := h.advertUseCase.UpdateAdvertStatus(advertId, status, userID); err != nil {
		h.handleError(writer, err, "failed to update advert status")
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, "Advert status updated")
}

// GetAdvertsByCategoryId godoc
// @Summary Retrieve adverts by category ID
// @Description Fetch a list of adverts associated with a specific category ID.
// @Tags adverts
// @Produce json
// @Param categoryId path string true "Category ID"
// @Success 200 {array} dto.AdvertResponse "List of adverts by category ID"
// @Failure 400 {object} utils.ErrResponse "Invalid category ID"
// @Failure 500 {object} utils.ErrResponse "Failed to retrieve adverts by category ID"
// @Router /api/v1/adverts/category/{categoryId} [get]
func (h *AdvertEndpoints) GetAdvertsByCategoryId(writer http.ResponseWriter, r *http.Request) {
	categoryIdStr := mux.Vars(r)["categoryId"]
	categoryId, err := uuid.Parse(categoryIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, err, "invalid category ID", nil)
		return
	}

	adverts, err := h.advertUseCase.GetAdvertsByCategoryId(categoryId)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to get adverts by category ID", nil)
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, adverts)
}

// UploadImage godoc
// @Summary Upload an image for an advert
// @Description Upload an image associated with an advert by its ID.
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
		h.sendError(writer, http.StatusBadRequest, err, "invalid advert ID", nil)
		return
	}

	fileHeader, _, err := r.FormFile("image")
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, err, "file not attached", nil)
		return
	}

	data, err := io.ReadAll(fileHeader)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to read file", nil)
		return
	}

	if err = fileHeader.Close(); err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to close file", nil)
		return
	}

	imageId, err := h.staticUseCase.UploadFile(data)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to upload image", nil)
		return
	}

	userID, err := h.sessionManager.GetUserID(r)	
	if err != nil {
		h.sendError(writer, http.StatusUnauthorized, err, "user not found", nil)
		return
	}

	if err := h.advertUseCase.UploadImage(advertId, imageId, userID); err != nil {
		if errors.Is(err, ErrAdvertNotFound) {
			h.sendError(writer, http.StatusNotFound, err, "advert not found", nil)
		} else if errors.Is(err, ErrForbidden) {
			h.sendError(writer, http.StatusForbidden, err, "forbidden", nil)
		} else {
			h.sendError(writer, http.StatusInternalServerError, err, "failed to upload image", nil)
		}
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, "Image uploaded")
}

func (h *AdvertEndpoints) sendError(w http.ResponseWriter, statusCode int, err error, context string, additionalInfo map[string]string) {
	h.logger.Error(err.Error(), zap.String("context", context), zap.Any("info", additionalInfo))
	utils.SendErrorResponse(w, statusCode, err.Error())
}

func (h *AdvertEndpoints) handleError(writer http.ResponseWriter, err error, context string) {
	switch {
	case errors.Is(err, ErrAdvertNotFound):
		h.sendError(writer, http.StatusNotFound, err, context, nil)
	case errors.Is(err, ErrForbidden):
		h.sendError(writer, http.StatusForbidden, err, context, nil)
	default:
		h.sendError(writer, http.StatusInternalServerError, err, context, nil)
	}
}
