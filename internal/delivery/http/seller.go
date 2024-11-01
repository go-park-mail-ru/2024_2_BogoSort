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

type SellerEndpoints struct {
	sellerRepo repository.Seller
	logger     *zap.Logger
}

func NewSellerEndpoints(sellerRepo repository.Seller, logger *zap.Logger) *SellerEndpoints {
	return &SellerEndpoints{
		sellerRepo: sellerRepo,
		logger:     logger,
	}
}

func (s *SellerEndpoints) Configure(router *mux.Router) {
	router.HandleFunc("api/v1/seller/{seller_id}", s.GetSellerByID).Methods(http.MethodGet)
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
// @Router /seller/{seller_id} [get]
func (s *SellerEndpoints) GetSellerByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sellerID, err := uuid.Parse(vars["seller_id"])
	if err != nil {
		s.logger.Error("error parsing seller_id", zap.Error(err))
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid seller_id")
		return
	}

	seller, err := s.sellerRepo.GetSellerByID(sellerID)
	if err != nil {
		s.handleError(w, err, "GetSellerByID")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, seller)
}

func (s *SellerEndpoints) handleError(w http.ResponseWriter, err error, context string) {
	switch {
	case errors.Is(err, repository.ErrSellerNotFound):
		s.sendError(w, http.StatusNotFound, ErrSellerNotFound, context, nil)
	case errors.Is(err, repository.ErrSellerAlreadyExists):
		s.sendError(w, http.StatusBadRequest, ErrSellerAlreadyExists, context, nil)
	case err != nil:
		s.sendError(w, http.StatusInternalServerError, err, context, nil)
	}
}

func (s *SellerEndpoints) sendError(w http.ResponseWriter, statusCode int, err error, context string, additionalInfo map[string]string) {
	s.logger.Error(err.Error(), zap.String("context", context), zap.Any("info", additionalInfo))
	utils.SendErrorResponse(w, statusCode, err.Error())
}
