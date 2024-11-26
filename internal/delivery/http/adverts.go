package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/static"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidID           = errors.New("invalid ID")
	ErrAdvertNotFound      = errors.New("advert not found")
	ErrForbidden           = errors.New("forbidden")
	ErrBadRequest          = errors.New("bad request")
	ErrInvalidAdvertData   = errors.New("invalid advert data")
	ErrInvalidAdvertStatus = errors.New("invalid advert status")
	ErrFileNotAttached     = errors.New("file not attached")
	ErrFailedToReadFile    = errors.New("failed to read file")
	ErrFailedToCloseFile   = errors.New("failed to close file")
	ErrFailedToUploadFile  = errors.New("failed to upload file")
	ErrTimeout             = errors.New("timeout exceeded")
	ErrTooLargeFile        = errors.New("file size exceeds limit")
)

type AdvertEndpoint struct {
	advertUC         usecase.AdvertUseCase
	staticGrpcClient static.StaticGrpcClient
	sessionManager   *utils.SessionManager
	policy           *bluemonday.Policy
}

func NewAdvertEndpoint(advertUC usecase.AdvertUseCase,
	staticGrpcClient static.StaticGrpcClient,
	sessionManager *utils.SessionManager,
	policy *bluemonday.Policy) *AdvertEndpoint {
	return &AdvertEndpoint{
		advertUC:         advertUC,
		staticGrpcClient: staticGrpcClient,
		sessionManager:   sessionManager,
		policy:           policy,
	}
}

func (h *AdvertEndpoint) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/adverts/seller/{sellerId}", h.GetBySellerId).Methods("GET")
	router.HandleFunc("/api/v1/adverts/{advertId}", h.GetById).Methods("GET")
	router.HandleFunc("/api/v1/adverts/category/{categoryId}", h.GetByCategoryId).Methods("GET")
	router.HandleFunc("/api/v1/adverts", h.Get).Methods("GET")
	router.HandleFunc("/api/v1/adverts/viewed/{advertId}", h.AddToViewed).Methods("POST")
	router.HandleFunc("/api/v1/search", h.Search).Methods("GET")
}

func (h *AdvertEndpoint) ConfigureProtectedRoutes(router *mux.Router) {
	protected := router.PathPrefix("/api/v1").Subrouter()
	sessionMiddleware := middleware.NewAuthMiddleware(h.sessionManager)
	protected.Use(sessionMiddleware.SessionMiddleware)

	protected.HandleFunc("/adverts", h.Add).Methods("POST")
	protected.HandleFunc("/adverts/my", h.GetByUserId).Methods("GET")
	protected.HandleFunc("/adverts/saved", h.GetSavedByUserId).Methods("GET")
	protected.HandleFunc("/adverts/cart/{cartId}", h.GetByCartId).Methods("GET")
	protected.HandleFunc("/adverts/{advertId}", h.Update).Methods("PUT")
	protected.HandleFunc("/adverts/{advertId}", h.Delete).Methods("DELETE")
	protected.HandleFunc("/adverts/{advertId}/status", h.UpdateStatus).Methods("PUT")
	protected.HandleFunc("/adverts/{advertId}/image", h.UploadImage).Methods("PUT")
	protected.HandleFunc("/adverts/saved/{advertId}", h.AddToSaved).Methods("POST")
	protected.HandleFunc("/adverts/saved/{advertId}", h.RemoveFromSaved).Methods("DELETE")
}

// Get godoc
// @Summary Retrieve all adverts
// @Description Fetch a list of all adverts with optional pagination.
// @Tags adverts
// @Produce json
// @Param limit query int false "Limit the number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} dto.PreviewAdvertCard "List of adverts"
// @Failure 400 {object} utils.ErrResponse "Invalid limit or offset"
// @Failure 500 {object} utils.ErrResponse "Failed to retrieve adverts"
// @Router /api/v1/adverts [get]
func (h *AdvertEndpoint) Get(writer http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	logger.Info("get adverts")
	userId, err := h.sessionManager.GetUserID(r)
	if err != nil {
		logger.Error("user not found", zap.Error(err))
		userId = uuid.Nil
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		h.sendError(writer, http.StatusBadRequest, ErrBadRequest, "invalid limit", nil)
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		h.sendError(writer, http.StatusBadRequest, ErrBadRequest, "invalid offset", nil)
		return
	}

	adverts, err := h.advertUC.Get(limit, offset, userId)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to get adverts", nil)
		return
	}

	for _, advert := range adverts {
		utils.SanitizePreviewAdvert(&advert.Preview, h.policy)
	}

	logger.Info("adverts sent", zap.Any("adverts", adverts))
	utils.SendJSONResponse(writer, http.StatusOK, adverts)
}

// GetBySellerId godoc
// @Summary Retrieve adverts by seller ID
// @Description Fetch a list of adverts associated with a specific seller ID.
// @Tags adverts
// @Produce json
// @Param sellerId path string true "Seller ID"
// @Success 200 {array} dto.PreviewAdvertCard "List of adverts"
// @Failure 400 {object} utils.ErrResponse "Invalid seller ID"
// @Failure 403 {object} utils.ErrResponse "Forbidden"
// @Failure 500 {object} utils.ErrResponse "Failed to retrieve adverts by seller ID"
// @Router /api/v1/adverts/seller/{sellerId} [get]
func (h *AdvertEndpoint) GetBySellerId(writer http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	userId, err := h.sessionManager.GetUserID(r)
	if err != nil {
		h.sendError(writer, http.StatusUnauthorized, err, "user not found", nil)
		return
	}

	sellerIdStr := mux.Vars(r)["sellerId"]
	sellerId, err := uuid.Parse(sellerIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, ErrInvalidID, "invalid seller ID", nil)
		return
	}

	adverts, err := h.advertUC.GetBySellerId(userId, sellerId)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to get adverts by seller ID", nil)
		return
	}

	for _, advert := range adverts {
		utils.SanitizePreviewAdvert(&advert.Preview, h.policy)
	}
	logger.Info("adverts sent", zap.Any("adverts", adverts))
	utils.SendJSONResponse(writer, http.StatusOK, adverts)
}

// GetByCartId godoc
// @Summary Retrieve adverts by cart ID
// @Description Fetch a list of adverts in the specified cart.
// @Tags adverts
// @Produce json
// @Param cartId path string true "Cart ID"
// @Success 200 {array} dto.PreviewAdvertCard "List of adverts in cart"
// @Failure 400 {object} utils.ErrResponse "Invalid cart ID"
// @Failure 403 {object} utils.ErrResponse "Forbidden"
// @Failure 500 {object} utils.ErrResponse "Failed to retrieve adverts by cart ID"
// @Router /api/v1/adverts/cart/{cartId} [get]
func (h *AdvertEndpoint) GetByCartId(writer http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	userId, err := h.sessionManager.GetUserID(r)
	if err != nil {
		h.sendError(writer, http.StatusUnauthorized, err, "user not found", nil)
		return
	}

	cartIdStr := mux.Vars(r)["cartId"]
	cartId, err := uuid.Parse(cartIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, ErrInvalidID, "invalid cart ID", nil)
		return
	}

	adverts, err := h.advertUC.GetByCartId(cartId, userId)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to get adverts by cart ID", nil)
		return
	}

	for _, advert := range adverts {
		utils.SanitizePreviewAdvert(&advert.Preview, h.policy)
	}
	logger.Info("adverts sent", zap.Any("adverts", adverts))
	utils.SendJSONResponse(writer, http.StatusOK, adverts)
}

// GetSavedByUserId godoc
// @Summary Retrieve adverts by user ID
// @Description Fetch a list of adverts saved by the specified user ID.
// @Tags adverts
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {array} dto.PreviewAdvertCard "List of adverts saved by user"
// @Failure 400 {object} utils.ErrResponse "Invalid user ID"
// @Failure 403 {object} utils.ErrResponse "Forbidden"
// @Failure 500 {object} utils.ErrResponse "Failed to retrieve adverts by user ID"
// @Router /api/v1/adverts/saved [get]
func (h *AdvertEndpoint) GetSavedByUserId(writer http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	userId, err := h.sessionManager.GetUserID(r)
	if err != nil {
		h.sendError(writer, http.StatusUnauthorized, err, "user not found", nil)
		return
	}

	adverts, err := h.advertUC.GetSavedByUserId(userId)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to get adverts by user ID", nil)
		return
	}

	for _, advert := range adverts {
		utils.SanitizePreviewAdvert(&advert.Preview, h.policy)
	}
	logger.Info("adverts sent", zap.Any("adverts", adverts))
	utils.SendJSONResponse(writer, http.StatusOK, adverts)
}

// GetById godoc
// @Summary Retrieve an advert by ID
// @Description Fetch an advert based on its ID.
// @Tags adverts
// @Produce json
// @Param advertId path string true "Advert ID"
// @Success 200 {object} dto.AdvertCard "Advert details"
// @Failure 400 {object} utils.ErrResponse "Invalid advert ID"
// @Failure 404 {object} utils.ErrResponse "Advert not found"
// @Failure 500 {object} utils.ErrResponse "Failed to retrieve advert by ID"
// @Router /api/v1/adverts/{advertId} [get]
func (h *AdvertEndpoint) GetById(writer http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	userId, err := h.sessionManager.GetUserID(r)
	if err != nil {
		userId = uuid.Nil
	}

	advertIdStr := mux.Vars(r)["advertId"]
	advertId, err := uuid.Parse(advertIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, ErrInvalidID, "invalid advert ID", nil)
		return
	}

	advert, err := h.advertUC.GetById(advertId, userId)
	if err != nil {
		h.handleError(writer, err, "failed to get advert by ID")
		return
	}
	logger.Info("advert sent", zap.Any("advert", advert))
	utils.SanitizeAdvert(&advert.Advert, h.policy)
	utils.SendJSONResponse(writer, http.StatusOK, advert)
}

// Add godoc
// @Summary Create a new advert
// @Description Add a new advert to the system.
// @Tags adverts
// @Accept json
// @Produce json
// @Param advert body dto.AdvertRequest true "Advert data"
// @Success 201 {object} dto.Advert "Advert created"
// @Failure 400 {object} utils.ErrResponse "Invalid advert data"
// @Failure 500 {object} utils.ErrResponse "Failed to create advert"
// @Router /api/v1/adverts [post]
func (h *AdvertEndpoint) Add(writer http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	var advert dto.AdvertRequest
	if err := json.NewDecoder(r.Body).Decode(&advert); err != nil {
		h.sendError(writer, http.StatusBadRequest, ErrInvalidAdvertData, "invalid advert data", nil)
		return
	}

	utils.SanitizeRequestAdvert(&advert, h.policy)

	userID, err := h.sessionManager.GetUserID(r)
	if err != nil {
		h.sendError(writer, http.StatusUnauthorized, err, "user not found", nil)
		return
	}

	newAdvert, err := h.advertUC.Add(&advert, userID)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to add advert", nil)
		return
	}

	logger.Info("advert created", zap.Any("advert", newAdvert))
	utils.SendJSONResponse(writer, http.StatusCreated, newAdvert)
}

// Update godoc
// @Summary Update an existing advert
// @Description Modify the details of an existing advert.
// @Tags adverts
// @Accept json
// @Produce json
// @Param advertId path string true "Advert ID"
// @Param advert body dto.AdvertRequest true "Updated advert data"
// @Success 200 {string} string "Advert updated successfully"
// @Failure 400 {object} utils.ErrResponse "Invalid advert data"
// @Failure 404 {object} utils.ErrResponse "Advert not found"
// @Failure 403 {object} utils.ErrResponse "Forbidden"
// @Failure 500 {object} utils.ErrResponse "Failed to update advert"
// @Router /api/v1/adverts/{advertId} [put]
func (h *AdvertEndpoint) Update(writer http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())
	var advert dto.AdvertRequest
	if err := json.NewDecoder(r.Body).Decode(&advert); err != nil {
		h.sendError(writer, http.StatusBadRequest, ErrInvalidAdvertData, "invalid advert data", nil)
		return
	}

	utils.SanitizeRequestAdvert(&advert, h.policy)

	userID, err := h.sessionManager.GetUserID(r)
	if err != nil {
		h.sendError(writer, http.StatusUnauthorized, ErrInvalidCredentials, "user not found", nil)
		return
	}

	advertIdStr := mux.Vars(r)["advertId"]
	advertId, err := uuid.Parse(advertIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, ErrInvalidID, "invalid advert ID", nil)
		return
	}

	if err := h.advertUC.Update(&advert, userID, advertId); err != nil {
		h.handleError(writer, err, "failed to update advert")
		return
	}

	logger.Info("advert updated", zap.Any("advert", advert))
	utils.SendJSONResponse(writer, http.StatusOK, "Advert updated successfully")
}

// Delete godoc
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
func (h *AdvertEndpoint) Delete(writer http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	advertIdStr := mux.Vars(r)["advertId"]
	advertId, err := uuid.Parse(advertIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, ErrInvalidID, "invalid advert ID", nil)
		return
	}

	userID, err := h.sessionManager.GetUserID(r)
	if err != nil {
		h.sendError(writer, http.StatusUnauthorized, ErrInvalidCredentials, "user not found", nil)
		return
	}

	if err := h.advertUC.DeleteById(advertId, userID); err != nil {
		h.handleError(writer, err, "failed to delete advert")
		return
	}

	logger.Info("advert deleted")
	utils.SendJSONResponse(writer, http.StatusOK, "Advert deleted")
}

// UpdateStatus godoc
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
func (h *AdvertEndpoint) UpdateStatus(writer http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	advertIdStr := mux.Vars(r)["advertId"]
	advertId, err := uuid.Parse(advertIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, ErrInvalidID, "invalid advert ID", nil)
		return
	}

	status := r.FormValue("status")
	if status != string(dto.AdvertStatusActive) && status != string(dto.AdvertStatusInactive) {
		h.sendError(writer, http.StatusBadRequest, ErrInvalidAdvertStatus, "invalid advert status", nil)
		return
	}

	userID, err := h.sessionManager.GetUserID(r)
	if err != nil {
		h.sendError(writer, http.StatusUnauthorized, ErrInvalidCredentials, "user not found", nil)
		return
	}

	if err := h.advertUC.UpdateStatus(advertId, userID, dto.AdvertStatus(status)); err != nil {
		h.handleError(writer, err, "failed to update advert status")
		return
	}

	logger.Info("advert status updated")
	utils.SendJSONResponse(writer, http.StatusOK, "Advert status updated")
}

// GetByCategoryId godoc
// @Summary Retrieve adverts by category ID
// @Description Fetch a list of adverts associated with a specific category ID.
// @Tags adverts
// @Produce json
// @Param categoryId path string true "Category ID"
// @Success 200 {array} dto.PreviewAdvertCard "List of adverts by category ID"
// @Failure 400 {object} utils.ErrResponse "Invalid category ID"
// @Failure 500 {object} utils.ErrResponse "Failed to retrieve adverts by category ID"
// @Router /api/v1/adverts/category/{categoryId} [get]
func (h *AdvertEndpoint) GetByCategoryId(writer http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	categoryIdStr := mux.Vars(r)["categoryId"]
	categoryId, err := uuid.Parse(categoryIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, ErrInvalidID, "invalid category ID", nil)
		return
	}

	userId, err := h.sessionManager.GetUserID(r)
	if err != nil {
		userId = uuid.Nil
	}

	adverts, err := h.advertUC.GetByCategoryId(categoryId, userId)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to get adverts by category ID", nil)
		return
	}

	for _, advert := range adverts {
		utils.SanitizePreviewAdvert(&advert.Preview, h.policy)
	}
	logger.Info("adverts sent", zap.Any("adverts", adverts))
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
func (h *AdvertEndpoint) UploadImage(writer http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	userID, err := h.sessionManager.GetUserID(r)
	if err != nil {
		h.sendError(writer, http.StatusUnauthorized, ErrInvalidCredentials, "user not found", nil)
		return
	}

	advertIdStr := mux.Vars(r)["advertId"]
	advertId, err := uuid.Parse(advertIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, ErrInvalidID, "invalid advert ID", nil)
		return
	}

	fileHeader, _, err := r.FormFile("image")
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, ErrFileNotAttached, "file not attached or size too large", nil)
		return
	}

	data, err := io.ReadAll(fileHeader)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, ErrFailedToReadFile, "failed to read file", nil)
		return
	}

	if err = fileHeader.Close(); err != nil {
		h.sendError(writer, http.StatusInternalServerError, ErrFailedToCloseFile, "failed to close file", nil)
		return
	}

	imageId, err := h.staticGrpcClient.UploadStatic(bytes.NewReader(data))
	if err != nil {
		if status, ok := status.FromError(err); ok {
			switch status.Code() {
			case codes.DeadlineExceeded:
				h.sendError(writer, http.StatusGatewayTimeout, ErrTimeout, "upload image timeout deadline exceeded", nil)
			case codes.ResourceExhausted:
				h.sendError(writer, http.StatusRequestEntityTooLarge, ErrTooLargeFile, "file size exceeds limit", nil)
			default:
				h.sendError(writer, http.StatusInternalServerError, ErrFailedToUploadFile, "failed to upload image", nil)
			}
		} else {
			h.sendError(writer, http.StatusInternalServerError, ErrFailedToUploadFile, "failed to upload image", nil)
		}
		return
	}

	if err := h.advertUC.UploadImage(advertId, imageId, userID); err != nil {
		if errors.Is(err, ErrAdvertNotFound) {
			h.sendError(writer, http.StatusNotFound, err, "advert not found", nil)
		} else if errors.Is(err, ErrForbidden) {
			h.sendError(writer, http.StatusForbidden, err, "forbidden", nil)
		} else {
			h.sendError(writer, http.StatusInternalServerError, ErrFailedToUploadFile, "failed to upload image", nil)
		}
		return
	}

	logger.Info("image uploaded")
	utils.SendJSONResponse(writer, http.StatusOK, "Image uploaded")
}

// AddToSaved godoc
// @Summary Add an advert to saved
// @Description Add an advert to saved by its ID.
// @Tags adverts
// @Param advertId path string true "Advert ID"
// @Success 200 "Advert added to saved"
// @Failure 400 {object} utils.ErrResponse "Invalid advert ID"
// @Failure 404 {object} utils.ErrResponse "Advert not found"
// @Failure 500 {object} utils.ErrResponse "Failed to add advert to saved"
// @Router /api/v1/adverts/saved/{advertId} [post]
func (h *AdvertEndpoint) AddToSaved(writer http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	userId, err := h.sessionManager.GetUserID(r)
	if err != nil {
		h.sendError(writer, http.StatusUnauthorized, ErrInvalidCredentials, "user not found", nil)
		return
	}

	advertIdStr := mux.Vars(r)["advertId"]
	advertId, err := uuid.Parse(advertIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, ErrInvalidID, "invalid advert ID", nil)
		return
	}

	if err := h.advertUC.AddToSaved(advertId, userId); err != nil {
		h.handleError(writer, err, "failed to add advert to saved")
		return
	}

	logger.Info("advert added to saved")
	utils.SendJSONResponse(writer, http.StatusOK, "Advert added to saved")
}

// RemoveFromSaved godoc
// @Summary Remove an advert from saved
// @Description Remove an advert from saved by its ID.
// @Tags adverts
// @Param advertId path string true "Advert ID"
// @Success 200 "Advert removed from saved"
// @Failure 400 {object} utils.ErrResponse "Invalid advert ID"
// @Failure 404 {object} utils.ErrResponse "Advert not found"
// @Failure 500 {object} utils.ErrResponse "Failed to remove advert from saved"
// @Router /api/v1/adverts/saved/{advertId} [delete]
func (h *AdvertEndpoint) RemoveFromSaved(writer http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	userId, err := h.sessionManager.GetUserID(r)
	if err != nil {
		h.sendError(writer, http.StatusUnauthorized, ErrInvalidCredentials, "user not found", nil)
		return
	}

	advertIdStr := mux.Vars(r)["advertId"]
	advertId, err := uuid.Parse(advertIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, ErrInvalidID, "invalid advert ID", nil)
		return
	}

	if err := h.advertUC.RemoveFromSaved(advertId, userId); err != nil {
		h.handleError(writer, err, "failed to remove advert from saved")
		return
	}

	logger.Info("advert removed from saved")
	utils.SendJSONResponse(writer, http.StatusOK, "Advert removed from saved")
}

func (h *AdvertEndpoint) sendError(w http.ResponseWriter, statusCode int, err error, contextInfo string, additionalInfo map[string]string) {
	logger := middleware.GetLogger(context.Background())

	logger.Error(err.Error(), zap.String("context", contextInfo), zap.Any("info", additionalInfo))
	utils.SendErrorResponse(w, statusCode, err.Error())
}

func (h *AdvertEndpoint) handleError(writer http.ResponseWriter, err error, context string) {
	switch {
	case errors.Is(err, ErrAdvertNotFound):
		h.sendError(writer, http.StatusNotFound, err, context, nil)
	case errors.Is(err, ErrForbidden):
		h.sendError(writer, http.StatusForbidden, err, context, nil)
	default:
		h.sendError(writer, http.StatusInternalServerError, err, context, nil)
	}
}

func (h *AdvertEndpoint) AddToViewed(writer http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	userId, err := h.sessionManager.GetUserID(r)
	if err != nil {
		userId = uuid.Nil
	}

	advertIdStr := mux.Vars(r)["advertId"]
	advertId, err := uuid.Parse(advertIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, ErrInvalidID, "invalid advert ID", nil)
		return
	}

	if err := h.advertUC.AddViewed(advertId, userId); err != nil {
		h.handleError(writer, err, "failed to add advert to viewed")
		return
	}

	logger.Info("advert added to viewed")
	utils.SendJSONResponse(writer, http.StatusOK, "Advert added to viewed")
}

// GetByUserId godoc
// @Summary Retrieve adverts by user ID
// @Description Fetch a list of adverts associated with a specific user ID.
// @Tags adverts
// @Success 200 {array} dto.MyPreviewAdvertCard "List of adverts by user ID"
// @Failure 400 {object} utils.ErrResponse "Invalid user ID"
// @Failure 500 {object} utils.ErrResponse "Failed to retrieve adverts by user ID"
// @Router /api/v1/adverts/my [get]
func (h *AdvertEndpoint) GetByUserId(writer http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	userId, err := h.sessionManager.GetUserID(r)
	if err != nil {
		h.sendError(writer, http.StatusUnauthorized, ErrInvalidCredentials, "user not found", nil)
		return
	}

	adverts, err := h.advertUC.GetByUserId(userId)
	if err != nil {
		h.handleError(writer, err, "failed to get adverts by user ID")
		return
	}

	for _, advert := range adverts {
		utils.SanitizePreviewAdvert(&advert.Preview, h.policy)
	}

	logger.Info("adverts sent", zap.Any("adverts", adverts))
	utils.SendJSONResponse(writer, http.StatusOK, adverts)
}

// Search godoc
// @Summary Поиск объявлений
// @Description Выполняет поиск объявлений по строке запроса с разбивкой на батчи.
// @Tags adverts
// @Produce json
// @Param query query string true "Строка поиска"
// @Param limit query int false "Лимит результатов (по умолчанию 100)"
// @Param offset query int false "Смещение для пагинации (по умолчанию 0)"
// @Success 200 {array} dto.PreviewAdvertCard "Список найденных объявлений"
// @Failure 400 {object} utils.ErrResponse "Неверные параметры запроса"
// @Failure 500 {object} utils.ErrResponse "Ошибка сервера"
// @Router /api/v1/adverts/search [get]
func (h *AdvertEndpoint) Search(writer http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	query := r.URL.Query().Get("query")
	if strings.TrimSpace(query) == "" {
		h.sendError(writer, http.StatusBadRequest, errors.New("search query is empty"), "empty search query", nil)
		return
	}

	batchSize := config.GetSearchBatchSize()

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 100
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	userId, err := h.sessionManager.GetUserID(r)
	if err != nil {
		userId = uuid.Nil
	}

	adverts, err := h.advertUC.Search(query, batchSize, limit, offset, userId)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "error during search execution", nil)
		return
	}

	logger.Info("adverts sent", zap.Any("adverts", adverts))
	utils.SendJSONResponse(writer, http.StatusOK, adverts)
}
