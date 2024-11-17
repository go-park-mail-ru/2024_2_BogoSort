package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PurchaseEndpoint struct {
	purchaseUC usecase.Purchase
	logger     *zap.Logger
}

func NewPurchaseEndpoint(purchaseUC usecase.Purchase, logger *zap.Logger) *PurchaseEndpoint {
	return &PurchaseEndpoint{purchaseUC: purchaseUC, logger: logger}
}

func (h *PurchaseEndpoint) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/purchase", h.Add).Methods("POST")
	router.HandleFunc("/api/v1/purchase/{user_id}", h.GetByUserID).Methods("GET")
}

// Add processes the addition of a purchase
// @Summary Adds a purchase
// @Description Accepts purchase data, validates it, and adds it to the system. Returns a response with purchase data or an error.
// @Tags Purchases
// @Accept json
// @Produce json
// @Param purchase body dto.PurchaseRequest true "Purchase request"
// @Success 201 {object} dto.PurchaseResponse "Successful purchase"
// @Failure 400 {object} utils.ErrResponse "Invalid request parameters"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/purchase [post]
func (h *PurchaseEndpoint) Add(w http.ResponseWriter, r *http.Request) {
	var purchase dto.PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&purchase); err != nil {
		h.logger.Error("failed to decode purchase request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, "invalid request parameters")
		return
	}

	purchaseResponse, err := h.purchaseUC.Add(purchase)
	if err != nil {
		h.logger.Error("failed to add purchase", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, purchaseResponse)
}

// GetByUserID processes the retrieval of purchases by user ID
// @Summary Retrieves purchases by user ID
// @Description Accepts a user ID, validates it, and retrieves purchases from the system. Returns a response with purchase data or an error.
// @Tags Purchases
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {array} dto.PurchaseResponse "Successful purchase"
// @Failure 400 {object} utils.ErrResponse "Invalid user ID"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/purchase/{user_id} [get]
func (h *PurchaseEndpoint) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["user_id"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err, "invalid user ID", nil)
		return
	}

	purchases, err := h.purchaseUC.GetByUserId(userID)
	if err != nil {
		h.handleError(w, err, "failed to get purchases")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, purchases)
}

func (h *PurchaseEndpoint) sendError(w http.ResponseWriter, status int, err error, message string, data interface{}) {
	utils.SendErrorResponse(w, status, message)
}

func (h *PurchaseEndpoint) handleError(w http.ResponseWriter, err error, message string) {
	h.logger.Error(message, zap.Error(err))
	utils.SendErrorResponse(w, http.StatusInternalServerError, "internal server error")
}
