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

// GetByID Retrieves the cart by its ID
// @Summary Retrieve cart by ID
// @Description Retrieves detailed information about a cart using its unique identifier
// @Tags Cart
// @Accept json
// @Produce json
// @Param cart_id path string true "Cart ID"
// @Success 200 {object} dto.CartResponse "Successfully retrieved cart"
// @Failure 400 {object} utils.ErrResponse "Invalid cart ID format"
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
	utils.SendJSONResponse(w, http.StatusOK, cart)
}

// GetByUserID Retrieves the cart by the user's ID
// @Summary Retrieve cart by User ID
// @Description Retrieves detailed information about a cart associated with a specific user
// @Tags Cart
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} dto.CartResponse "Successfully retrieved cart"
// @Failure 400 {object} utils.ErrResponse "Invalid user ID format"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /cart/user/{user_id} [get]
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

	utils.SendJSONResponse(w, http.StatusOK, cart)
}

// AddToCart Adds an advert to the user's cart
// @Summary Add advert to user's cart
// @Description Adds a new advert to the cart associated with a user
// @Tags Cart
// @Accept json
// @Produce json
// @Param purchase body dto.AddAdvertToUserCartRequest true "Data to add advert to cart"
// @Success 200 {object} map[string]string "Successfully added advert"
// @Failure 400 {object} utils.ErrResponse "Invalid request data"
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
	case err != nil:
		utils.SendErrorResponse(w, http.StatusInternalServerError, "failed to add advert to user cart")
	}

	logger.Info("advert added to user cart")
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "advert added to user cart"})
}

// DeleteFromCart Deletes an advert from the user's cart
// @Summary Delete advert from user's cart
// @Description Deletes an advert from the cart associated with a user
// @Tags Cart
// @Accept json
// @Produce json
// @Param purchase body dto.DeleteAdvertFromUserCartRequest true "Data to delete advert from cart"
// @Success 200 {object} map[string]string "Successfully deleted advert from user cart"
// @Failure 400 {object} utils.ErrResponse "Invalid request data"
// @Failure 404 {object} utils.ErrResponse "Cart or advert not found"
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/cart/delete [delete]
func (h *CartEndpoint) DeleteFromCart(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())
	logger.Info("delete advert from user cart request")
	var req dto.DeleteAdvertFromUserCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

// CheckExists Checks if a cart exists for a user
// @Summary Check if cart exists for user
// @Description Checks if a cart exists for a user by their ID
// @Tags Cart
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]bool "Cart existence check result"
// @Failure 400 {object} utils.ErrResponse "Invalid user ID format"
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

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"exists": exists.String()})
}
