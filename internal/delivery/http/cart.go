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
	router.HandleFunc("api/v1/cart/{cart_id}", h.GetCartByID).Methods(http.MethodGet)
	router.HandleFunc("api/v1/cart/user/{user_id}", h.GetCartByUserID).Methods(http.MethodGet)
	router.HandleFunc("api/v1/cart/add", h.AddAdvertToCart).Methods(http.MethodPost)
}

// GetCartByID Retrieves the cart by its ID
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

// GetCartByUserID Retrieves the cart by the user's ID
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

// AddAdvertToCart Adds an advert to the user's cart
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
