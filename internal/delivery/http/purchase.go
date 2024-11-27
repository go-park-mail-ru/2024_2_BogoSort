package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/cart_purchase"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	proto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/cart_purchase/proto"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PurchaseEndpoint struct {
	purchaseServer *cart_purchase.GrpcServer
	logger         *zap.Logger
}

func NewPurchaseEndpoint(purchaseServer *cart_purchase.GrpcServer, logger *zap.Logger) *PurchaseEndpoint {
	return &PurchaseEndpoint{purchaseServer: purchaseServer, logger: logger}
}

func (h *PurchaseEndpoint) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/purchase/{user_id}", h.Add).Methods("POST")
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
// @Router /api/v1/purchase/{user_id} [post]
func (h *PurchaseEndpoint) Add(w http.ResponseWriter, r *http.Request) {
	var purchase dto.PurchaseRequest

	userIDStr := mux.Vars(r)["user_id"]
	_, err := uuid.Parse(userIDStr)
	if err != nil {
		h.handleError(w, err, "invalid user ID")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&purchase); err != nil {
		h.logger.Error("failed to decode purchase request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, "invalid request parameters")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Hour)
	defer cancel()

	purchaseResponse, err := h.purchaseServer.AddPurchase(ctx, &proto.AddPurchaseRequest{CartId: purchase.CartID.String(), Address: purchase.Address, PaymentMethod: proto.PaymentMethod(proto.PaymentMethod_value[string(purchase.PaymentMethod)]), DeliveryMethod: proto.DeliveryMethod(proto.DeliveryMethod_value[string(purchase.DeliveryMethod)]), UserId: purchase.UserID.String()})
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
		h.handleError(w, err, "invalid user ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	purchases, err := h.purchaseServer.GetPurchasesByUserID(ctx, &proto.GetPurchasesByUserIDRequest{UserId: userID.String()})
	if err != nil {
		h.handleError(w, err, "failed to get purchases")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, purchases)
}

func (h *PurchaseEndpoint) handleError(w http.ResponseWriter, err error, message string) {
	h.logger.Error(message, zap.Error(err))
	utils.SendErrorResponse(w, http.StatusInternalServerError, "internal server error")
}
