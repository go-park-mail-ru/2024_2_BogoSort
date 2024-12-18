package repository

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
)

// PaymentRepository интерфейс для работы с платежами
type PaymentRepository interface {
	// InsertOrder создает новую запись о заказе в базе данных
	InsertOrder(orderID, amount, paymentID, status string) (*entity.Order, error)
	// UpdateOrderStatus обновляет статус существующего заказа
	UpdateOrderStatus(orderID string, status string) (*entity.Order, error)
	// GetOrderByID возвращает заказ по его идентификатору
	GetOrderByID(orderID string) (*entity.Order, error)
	// GetOrdersInProcess возвращает массив заказов со статусом 'in_process'
	GetOrdersInProcess() ([]entity.Order, error)
}
