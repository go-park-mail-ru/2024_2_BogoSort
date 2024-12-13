package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
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
// @Failure 400 {string} string "Invalid request body or Payment service error"
// @Failure 500 {string} string "Payment service error"
// @Router /api/v1/payment/init [post]
func (h *PaymentEndpoint) InitPayment(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	var paymentReq PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	url, err := h.paymentUC.InitPayment(paymentReq.ItemID)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Payment service error", http.StatusBadRequest)
		return
	}

	if url == nil || *url == "" {
		http.Error(w, "Payment service error", http.StatusInternalServerError)
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	logger.Info(fmt.Sprintf("Callback received: %s\n", string(body)))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Callback received")
}