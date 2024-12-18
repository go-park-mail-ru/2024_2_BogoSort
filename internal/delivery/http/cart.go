package http

import (
	"context"
	"encoding/json"
	"errors"
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

type CartEndpoint struct {
	cartClient *cart_purchase.CartPurchaseClient
}

func NewCartEndpoint(cartClient *cart_purchase.CartPurchaseClient) *CartEndpoint {
	return &CartEndpoint{
		cartClient: cartClient,
	}
}

func (h *CartEndpoint) Configure(router *mux.Router) {
	router.HandleFunc("/api/v1/cart/{cart_id}", h.GetByID).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/cart/user/{user_id}", h.GetByUserID).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/cart/add", h.AddToCart).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/cart/delete", h.DeleteFromCart).Methods(http.MethodDelete)
	router.HandleFunc("/api/v1/cart/exists/{user_id}", h.CheckExists).Methods(http.MethodGet)
}

// GetByID godoc
// @Summary Get cart by ID
// @Description Retrieves a cart by its unique identifier
// @Tags Cart
// @Accept json
// @Produce json
// @Param cart_id path string true "Cart ID (UUID format)"
// @Success 200 {object} dto.CartResponse "Cart details successfully retrieved"
// @Failure 400 {object} utils.ErrResponse "Invalid cart ID format or missing cart_id"
// @Failure 404 {object} utils.ErrResponse "Cart not found"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/cart/{cart_id} [get]
func (h *CartEndpoint) GetByID(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())
	logger.Info("get cart by id request")
	vars := mux.Vars(r)
	cartIDStr, ok := vars["cart_id"]
	if !ok {
		utils.SendErrorResponse(w, http.StatusBadRequest, "cart_id is required")
		return
	}

	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		logger.Error("failed to parse cart_id", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, "invalid cart_id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cart, err := h.cartClient.GetCartByID(ctx, cartID)
	if errors.Is(err, cart_purchase.ErrCartNotFound) {
		logger.Error("cart not found", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusNotFound, "cart not found")
		return
	}
	if err != nil {
		logger.Error("failed to get cart", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, "failed to get adverts from cart")
		return
	}

	logger.Info("cart", zap.Any("cart", cart))
	if err := utils.WriteJSON(w, cart, http.StatusOK); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "failed to send cart")
		return
	}
}

// GetByUserID godoc
// @Summary Get cart by user ID
// @Description Retrieves a cart associated with a specific user
// @Tags Cart
// @Accept json
// @Produce json
// @Param user_id path string true "User ID (UUID format)"
// @Success 200 {object} dto.CartResponse "Cart details successfully retrieved"
// @Failure 400 {object} utils.ErrResponse "Invalid user ID format or missing user_id"
// @Failure 404 {object} utils.ErrResponse "Cart not found for user"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/cart/user/{user_id} [get]
func (h *CartEndpoint) GetByUserID(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())
	logger.Info("get cart by user id request")
	vars := mux.Vars(r)
	userIDStr, ok := vars["user_id"]
	if !ok {
		utils.SendErrorResponse(w, http.StatusBadRequest, "user_id is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Error("failed to parse user_id", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cart, err := h.cartClient.GetCartByUserID(ctx, userID)
	if errors.Is(err, cart_purchase.ErrCartNotFound) {
		logger.Error("cart not found", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusNotFound, "cart not found")
		return
	}
	if err != nil {
		logger.Error("failed to get cart", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, "failed to get cart")
		return
	}

	if err := utils.WriteJSON(w, cart, http.StatusOK); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "failed to send cart")
		return
	}
}

// AddToCart godoc
// @Summary Add item to cart
// @Description Adds an advertisement to a user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param request body dto.AddAdvertToUserCartRequest true "Add to cart request containing user_id and advert_id"
// @Success 200 {object} map[string]string "Successfully added item to cart"
// @Failure 400 {object} utils.ErrResponse "Invalid request body"
// @Failure 404 {object} utils.ErrResponse "Cart not found"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/cart/add [post]
func (h *CartEndpoint) AddToCart(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())
	logger.Info("add advert to user cart request")
	var req dto.AddAdvertToUserCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := h.cartClient.AddAdvertToCart(ctx, req.UserID, req.AdvertID)

	switch {
	case errors.Is(err, cart_purchase.ErrCartNotFound):
		utils.SendErrorResponse(w, http.StatusNotFound, "cart not found")
		return
	case err != nil:
		utils.SendErrorResponse(w, http.StatusInternalServerError, "failed to add advert to user cart")
		return
	}

	logger.Info("advert added to user cart")
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "advert added to user cart"})
}

// DeleteFromCart godoc
// @Summary Remove item from cart
// @Description Removes an advertisement from a user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param request body dto.DeleteAdvertFromUserCartRequest true "Delete from cart request containing cart_id and advert_id"
// @Success 200 {object} map[string]string "Successfully removed item from cart"
// @Failure 400 {object} utils.ErrResponse "Invalid request body"
// @Failure 404 {object} utils.ErrResponse "Cart or item not found"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/cart/delete [delete]
func (h *CartEndpoint) DeleteFromCart(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())
	logger.Info("delete advert from user cart request")
	var req dto.DeleteAdvertFromUserCartRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := h.cartClient.DeleteAdvertFromCart(ctx, req.CartID, req.AdvertID)
	switch {
	case errors.Is(err, cart_purchase.ErrCartNotFound):
		utils.SendErrorResponse(w, http.StatusNotFound, "cart or advert not found")
	case err != nil:
		utils.SendErrorResponse(w, http.StatusInternalServerError, "failed to delete advert from user cart")
	}

	logger.Info("advert deleted from user cart")
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "advert deleted from user cart"})
}

// CheckExists godoc
// @Summary Check cart existence
// @Description Checks if a cart exists for a specific user
// @Tags Cart
// @Accept json
// @Produce json
// @Param user_id path string true "User ID (UUID format)"
// @Success 200 {object} map[string]string "Cart existence status with cart_id if exists"
// @Failure 400 {object} utils.ErrResponse "Invalid user ID format or missing user_id"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/cart/exists/{user_id} [get]
func (h *CartEndpoint) CheckExists(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())
	logger.Info("check cart existence request")
	vars := mux.Vars(r)
	userIDStr, ok := vars["user_id"]
	if !ok {
		utils.SendErrorResponse(w, http.StatusBadRequest, "user_id is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exists, err := h.cartClient.CheckCartExists(ctx, userID)
	if err != nil {
		logger.Error("failed to check cart existence", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusInternalServerError, "failed to check cart existence")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"cart_id": exists.String()})
}
