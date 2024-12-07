package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/cart_purchase"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PurchaseEndpoint struct {
	purchaseClient *cart_purchase.CartPurchaseClient
}

func NewPurchaseEndpoint(purchaseClient *cart_purchase.CartPurchaseClient) *PurchaseEndpoint {
	return &PurchaseEndpoint{purchaseClient: purchaseClient}
}

func (h *PurchaseEndpoint) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/purchase/{user_id}", h.Add).Methods("POST")
	router.HandleFunc("/api/v1/purchase/{user_id}", h.GetByUserID).Methods("GET")
}

// Add processes the addition of a purchase
// @Summary Create new purchase
// @Description Creates a new purchase record for a user
// @Tags Purchases
// @Accept json
// @Produce json
// @Param user_id path string true "User ID" format(uuid)
// @Param purchase body dto.PurchaseRequest true "Purchase details"
// @Success 201 {object} dto.PurchaseResponse "Purchase created successfully"
// @Failure 400 {object} utils.ErrResponse "Invalid request parameters or user ID"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/purchase/{user_id} [post]
func (h *PurchaseEndpoint) Add(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	logger.Info("add purchase request")

	var purchase dto.PurchaseRequest

	userIDStr := mux.Vars(r)["user_id"]
	_, err := uuid.Parse(userIDStr)
	if err != nil {
		h.handleError(w, err, "invalid user ID")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&purchase); err != nil {
		logger.Error("failed to decode purchase request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, "invalid request parameters")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Hour)
	defer cancel()

	purchaseResponse, err := h.purchaseClient.AddPurchase(ctx, purchase)
	if err != nil {
		logger.Error("failed to add purchase", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	logger.Info("purchase added", zap.Any("purchase", purchaseResponse))
	utils.SendJSONResponse(w, http.StatusCreated, purchaseResponse)
}

// GetByUserID processes the retrieval of purchases by user ID
// @Summary Get user purchases
// @Description Retrieves all purchases associated with a specific user ID
// @Tags Purchases
// @Accept json
// @Produce json
// @Param user_id path string true "User ID" format(uuid)
// @Success 200 {array} dto.PurchaseResponse "List of user purchases"
// @Failure 400 {object} utils.ErrResponse "Invalid user ID"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/purchase/{user_id} [get]
func (h *PurchaseEndpoint) GetByUserID(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	logger.Info("get purchases by user id request")

	userIDStr := mux.Vars(r)["user_id"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.handleError(w, err, "invalid user ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	purchases, err := h.purchaseClient.GetPurchasesByUserID(ctx, userID)
	if err != nil {
		h.handleError(w, err, "failed to get purchases")
		return
	}

	logger.Info("purchases found", zap.Any("purchases", purchases))
	utils.SendJSONResponse(w, http.StatusOK, purchases)
}

func (h *PurchaseEndpoint) handleError(w http.ResponseWriter, err error, message string) {
	logger := middleware.GetLogger(context.Background())

	logger.Error(message, zap.Error(err))
	utils.SendErrorResponse(w, http.StatusInternalServerError, "internal server error")
}
