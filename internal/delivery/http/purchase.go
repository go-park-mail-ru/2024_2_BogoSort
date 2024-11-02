package http

import (
	"net/http"
	"encoding/json"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PurchaseEndpoints struct {
	purchaseUC usecase.Purchase
	logger *zap.Logger
}

func NewPurchaseEndpoints(purchaseUC usecase.Purchase, logger *zap.Logger) *PurchaseEndpoints {
	return &PurchaseEndpoints{purchaseUC: purchaseUC, logger: logger}
}

func (h *PurchaseEndpoints) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/purchase", h.AddPurchase).Methods("POST")
}

// AddPurchase processes the addition of a purchase
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
func (h *PurchaseEndpoints) AddPurchase(w http.ResponseWriter, r *http.Request) {
	var purchase dto.PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&purchase); err != nil {
		h.logger.Error("failed to decode purchase request", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, "invalid request parameters")
		return
	}

	purchaseResponse, err := h.purchaseUC.AddPurchase(purchase)
	if err != nil {
		h.logger.Error("failed to add purchase", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, purchaseResponse)
}
