package http

import (
	"net/http"
)

type PurchaseEndpoints struct {
}

// CreatePurchase Обрабатывает покупку по ID корзины
// @Summary Совершает покупку по ID корзины
// @Description Принимает ID корзины и выполняет процесс покупки
// @Tags Покупки
// @Accept json
// @Produce json
// @Param purchase body dto.PurchaseRequest true "Данные покупки"
// @Success 200 {object} dto.PurchaseResponse "Успешная покупка"
// @Failure 400 {object} utils.ErrResponse "Неверные параметры запроса"
// @Failure 404 {object} utils.ErrResponse "Корзина не найдена"
// @Failure 500 {object} utils.ErrResponse "Внутренняя ошибка сервера"
// @Router /purchase [post]
func (h *PurchaseEndpoints) CreatePurchase(w http.ResponseWriter, r *http.Request) {

}
