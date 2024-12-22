package http

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PromotionEndpoint struct {
	promotionUC usecase.PromotionUseCase
}

func NewPromotionEndpoint(promotionUC usecase.PromotionUseCase) *PromotionEndpoint {
	return &PromotionEndpoint{
		promotionUC: promotionUC,
	}
}

func (e *PromotionEndpoint) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/promotions", e.Get).Methods("GET")
}

// Get godoc
// @Summary Get promotion info
// @Description Retrieve promotion info
// @Tags promotions
// @Accept json
// @Produce json
// @Success 200 {object} entity.Promotion
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/promotions [get]
func (e *PromotionEndpoint) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(context.Background())
	logger.Info("get promotions request")
	promotions, err := e.promotionUC.GetPromotionInfo()
	if err != nil {
		e.sendError(w, http.StatusInternalServerError, err, "error getting promotions", nil)
		return
	}
	logger.Info("get promotions response", zap.Any("promotions", promotions))
	utils.SendJSONResponse(w, http.StatusOK, promotions)
}

func (e *PromotionEndpoint) sendError(w http.ResponseWriter, statusCode int, err error, contextInfo string, additionalInfo map[string]string) {
	logger := middleware.GetLogger(context.Background())
	logger.Error(err.Error(), zap.String("context", contextInfo), zap.Any("info", additionalInfo))
	utils.SendErrorResponse(w, statusCode, err.Error())
}
