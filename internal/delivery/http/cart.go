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

type CartEndpoints struct {
	cartUC usecase.Cart
	logger *zap.Logger
}

func NewCartEndpoints(cartUC usecase.Cart, logger *zap.Logger) *CartEndpoints {
	return &CartEndpoints{
		cartUC: cartUC,
		logger: logger,
	}
}

func (h *CartEndpoints) Configure(router *mux.Router) {
	router.HandleFunc("/cart/{cart_id}", h.GetCartByID).Methods(http.MethodGet)
	router.HandleFunc("/cart/user/{user_id}", h.GetCartByUserID).Methods(http.MethodGet)
	router.HandleFunc("/cart/add", h.AddAdvertToCart).Methods(http.MethodPost)
}

func (h *CartEndpoints) GetCartByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cartIDStr, ok := vars["cart_id"]
	if !ok {
		utils.SendErrorResponse(w, http.StatusBadRequest, "cart_id is required")
		return
	}

	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		h.logger.Error("failed to parse cart_id", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, "invalid cart_id")
		return
	}

	cart, err := h.cartUC.GetCartByID(cartID)
	if err != nil {
		h.logger.Error("failed to get cart", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, "failed to get adverts from cart")
		return
	}

	h.logger.Info("cart", zap.Any("cart", cart))
	utils.SendJSONResponse(w, http.StatusOK, cart)
}

func (h *CartEndpoints) GetCartByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["user_id"]
	if !ok {
		utils.SendErrorResponse(w, http.StatusBadRequest, "user_id is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Error("failed to parse user_id", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	cart, err := h.cartUC.GetCartByUserID(userID)
	if err != nil {
		h.logger.Error("failed to get cart", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, "failed to get cart")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, cart)
}

func (h *CartEndpoints) AddAdvertToCart(w http.ResponseWriter, r *http.Request) {
	var req dto.AddAdvertToUserCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.cartUC.AddAdvertToUserCart(req.UserID, req.AdvertID); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "failed to add advert to user cart")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "advert added to user cart"})
}
