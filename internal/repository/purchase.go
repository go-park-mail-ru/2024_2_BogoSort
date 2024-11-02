package repository

type PurchaseRepository interface {
	// CreatePurchase создает запись о покупке
	CreatePurchase() error
}
