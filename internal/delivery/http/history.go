package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HistoryEndpoint struct {
	historyRepo repository.HistoryRepository
}

func NewHistoryEndpoint(historyRepo repository.HistoryRepository) *HistoryEndpoint {
	return &HistoryEndpoint{historyRepo: historyRepo}
}

func (h *HistoryEndpoint) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/history/{advert_id}", h.GetAdvertPriceHistory).Methods(http.MethodGet)
}

// @Summary Get Advert Price History
// @Description Получает историю изменения цены для указанного объявления
// @Tags History
// @Accept json
// @Produce json
// @Param advert_id path string true "Advert ID"
// @Success 200 {array} dto.PriceHistoryResponse
// @Failure 500 {object} utils.ErrResponse "Internal server error"
// @Router /api/v1/history/{advert_id} [get]
func (h *HistoryEndpoint) GetAdvertPriceHistory(w http.ResponseWriter, r *http.Request) {
	advertID := mux.Vars(r)["advert_id"]
	history, err := h.historyRepo.GetAdvertPriceHistory(uuid.MustParse(advertID))
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, history)
}
