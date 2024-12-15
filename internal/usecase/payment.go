package usecase

type PaymentUseCase interface {
	// InitPayment инициализация платежа
	InitPayment(itemId string) (*string, error)
}
