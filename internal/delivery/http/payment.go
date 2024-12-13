package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PaymentEndpoint struct {
	paymentUC      usecase.PaymentUseCase
	sessionManager *utils.SessionManager
}

func NewPaymentEndpoint(paymentUC usecase.PaymentUseCase, sessionManager *utils.SessionManager) *PaymentEndpoint {
	return &PaymentEndpoint{
		paymentUC:      paymentUC,
		sessionManager: sessionManager,
	}
}

func (h *PaymentEndpoint) ConfigureProtectedRoutes(router *mux.Router) {
	protected := router.PathPrefix("/api/v1").Subrouter()
	sessionMiddleware := middleware.NewAuthMiddleware(h.sessionManager)
	protected.Use(sessionMiddleware.SessionMiddleware)

	router.HandleFunc("/api/v1/payment/init", h.InitPayment).Methods("POST")
}

func (h *PaymentEndpoint) ConfigureUnprotectedRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/payment/callback", h.Callback).Methods("GET")
}

type PaymentRequest struct {
	ItemID string `json:"item_id"`
}

// InitPayment handles the initialization of a payment process.
// @Summary Initialize payment
// @Description Initiates a payment process for a given item ID.
// @Tags payment
// @Accept json
// @Produce json
// @Param paymentRequest body PaymentRequest true "Payment Request"
// @Success 200 {object} map[string]string "Payment URL"
// @Failure 400 {object} utils.ErrResponse "Invalid request or payment service error"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/payment/init [post]
func (h *PaymentEndpoint) InitPayment(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	var paymentReq PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentReq); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	url, err := h.paymentUC.InitPayment(paymentReq.ItemID)
	if err != nil {
		logger.Error("Payment service error", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, "Payment service error")
		return
	}

	if url == nil || *url == "" {
		logger.Error("Payment service returned empty URL")
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Payment service error")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, struct {
		PaymentURL string `json:"payment_url"`
	}{
		PaymentURL: *url,
	})
}

func (h *PaymentEndpoint) Callback(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("Failed to read request body", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to read request body")
		return
	}
	defer r.Body.Close()

	logger.Info("Callback received", zap.String("body", string(body)))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Callback received")
}
