package http

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	ErrInvalidRequestBody  = errors.New("invalid request body")
	ErrSellerNotFound      = errors.New("seller not found")
	ErrSellerAlreadyExists = errors.New("seller already exists")
	ErrInvalidSellerData   = errors.New("invalid seller data")
)

type SellerEndpoint struct {
	sellerRepo repository.Seller
	logger     *zap.Logger
}

func NewSellerEndpoint(sellerRepo repository.Seller, logger *zap.Logger) *SellerEndpoint {
	return &SellerEndpoint{
		sellerRepo: sellerRepo,
		logger:     logger,
	}
}

func (s *SellerEndpoint) Configure(router *mux.Router) {
	router.HandleFunc("/api/v1/seller/{seller_id}", s.GetByID).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/seller/user/{user_id}", s.GetByUserID).Methods(http.MethodGet)
}

// GetSellerByID
// @Summary Получение продавца по ID
// @Description Возвращает информацию о продавце по его ID
// @Tags Продавцы
// @Accept json
// @Produce json
// @Param seller_id path string true "ID продавца"
// @Success 200 {object} entity.Seller "Информация о продавце"
// @Failure 400 {object} utils.ErrResponse "Некорректный запрос"
// @Failure 404 {object} utils.ErrResponse "Продавец не найден"
// @Failure 500 {object} utils.ErrResponse "Внутренняя ошибка сервера"
// @Router /api/v1/seller/{seller_id} [get]
func (s *SellerEndpoint) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sellerID, err := uuid.Parse(vars["seller_id"])
	if err != nil {
		s.handleError(w, err, "error parsing seller_id")
		return
	}

	seller, err := s.sellerRepo.GetById(sellerID)
	switch {
	case errors.Is(err, repository.ErrSellerNotFound):
		s.handleError(w, err, "error getting seller by id")
	case err != nil:
		s.handleError(w, err, "error getting seller by id")
	}

	s.logger.Info("seller found", zap.String("seller_id", sellerID.String()))
	utils.SendJSONResponse(w, http.StatusOK, seller)
}

// GetSellerByUserID Получение продавца по ID пользователя
// @Summary Получить продавца по ID пользователя
// @Description Возвращает информацию о продавце, связанном с указанным ID пользователя
// @Tags Продавцы
// @Accept json
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Success 200 {object} entity.Seller "Информация о продавце"
// @Failure 400 {object} utils.ErrResponse "Некорректный ID пользователя"
// @Failure 404 {object} utils.ErrResponse "Продавец не найден"
// @Failure 500 {object} utils.ErrResponse "Внутренняя ошибка сервера"
// @Router /api/v1/seller/user/{user_id} [get]
func (s *SellerEndpoint) GetByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["user_id"])
	if err != nil {
		s.handleError(w, err, "error parsing user_id")
		return
	}

	seller, err := s.sellerRepo.GetByUserId(userID)
	switch {
	case errors.Is(err, repository.ErrSellerNotFound):
		s.handleError(w, err, "error getting seller by user_id")
	case err != nil:
		s.handleError(w, err, "error getting seller by user_id")
	}

	s.logger.Info("seller found", zap.String("user_id", userID.String()))
	utils.SendJSONResponse(w, http.StatusOK, seller)
}

func (s *SellerEndpoint) handleError(w http.ResponseWriter, err error, context string) {
	switch {
	case errors.Is(err, repository.ErrSellerNotFound):
		s.sendError(w, http.StatusNotFound, ErrSellerNotFound, context, nil)
	case errors.Is(err, repository.ErrSellerAlreadyExists):
		s.sendError(w, http.StatusBadRequest, ErrSellerAlreadyExists, context, nil)
	case err != nil:
		s.sendError(w, http.StatusInternalServerError, err, context, nil)
	}
}

func (s *SellerEndpoint) sendError(w http.ResponseWriter, statusCode int, err error, context string, additionalInfo map[string]string) {
	s.logger.Error(err.Error(), zap.String("context", context), zap.Any("info", additionalInfo))
	utils.SendErrorResponse(w, statusCode, err.Error())
}
