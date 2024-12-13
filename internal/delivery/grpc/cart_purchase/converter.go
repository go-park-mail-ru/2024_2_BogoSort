package cart_purchase

import (
	proto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/cart_purchase/proto"
	"github.com/pkg/errors"
)

const (
	StatusPending          = "pending"
	StatusInProgress       = "in_progress"
	StatusCompleted        = "completed"
	StatusCanceled         = "canceled"
	StatusActive           = "active"
	StatusInactive         = "inactive"
	StatusReserved         = "reserved"
	StatusUnknown          = "unknown"
	PaymentMethodCard      = "card"
	PaymentMethodCash      = "cash"
	DeliveryMethodPickup   = "pickup"
	DeliveryMethodDelivery = "delivery"
)

func ConvertDBPurchaseStatusToEnum(dbStatus string) (proto.PurchaseStatus, error) {
	switch dbStatus {
	case StatusPending:
		return proto.PurchaseStatus_PURCHASE_STATUS_PENDING, nil
	case StatusInProgress:
		return proto.PurchaseStatus_PURCHASE_STATUS_IN_PROGRESS, nil
	case StatusCompleted:
		return proto.PurchaseStatus_PURCHASE_STATUS_COMPLETED, nil
	case StatusCanceled:
		return proto.PurchaseStatus_PURCHASE_STATUS_CANCELED, nil
	default:
		return proto.PurchaseStatus_PURCHASE_STATUS_PENDING, errors.New(StatusUnknown)
	}
}

func ConvertPurchaseStatusToDB(status proto.PurchaseStatus) string {
	switch status {
	case proto.PurchaseStatus_PURCHASE_STATUS_PENDING:
		return StatusPending
	case proto.PurchaseStatus_PURCHASE_STATUS_IN_PROGRESS:
		return StatusInProgress
	case proto.PurchaseStatus_PURCHASE_STATUS_COMPLETED:
		return StatusCompleted
	case proto.PurchaseStatus_PURCHASE_STATUS_CANCELED:
		return StatusCanceled
	default:
		return StatusUnknown
	}
}

func ConvertDBPaymentMethodToEnum(dbStatus string) (proto.PaymentMethod, error) {
	switch dbStatus {
	case "card":
		return proto.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case "cash":
		return proto.PaymentMethod_PAYMENT_METHOD_CASH, nil
	default:
		return proto.PaymentMethod_PAYMENT_METHOD_CARD, errors.New("unknown payment method")
	}
}

func ConvertPaymentMethodToDB(status proto.PaymentMethod) string {
	switch status {
	case proto.PaymentMethod_PAYMENT_METHOD_CARD:
		return "card"
	case proto.PaymentMethod_PAYMENT_METHOD_CASH:
		return "cash"
	default:
		return "unknown"
	}
}

func ConvertDBDeliveryMethodToEnum(dbStatus string) (proto.DeliveryMethod, error) {
	switch dbStatus {
	case "pickup":
		return proto.DeliveryMethod_DELIVERY_METHOD_PICKUP, nil
	case "delivery":
		return proto.DeliveryMethod_DELIVERY_METHOD_DELIVERY, nil
	default:
		return proto.DeliveryMethod_DELIVERY_METHOD_PICKUP, errors.New("unknown delivery method")
	}
}

func ConvertDeliveryMethodToDB(status proto.DeliveryMethod) string {
	switch status {
	case proto.DeliveryMethod_DELIVERY_METHOD_PICKUP:
		return "pickup"
	case proto.DeliveryMethod_DELIVERY_METHOD_DELIVERY:
		return "delivery"
	default:
		return "unknown"
	}
}

func ConvertDBCartStatusToEnum(dbStatus string) (proto.CartStatus, error) {
	switch dbStatus {
	case StatusActive:
		return proto.CartStatus_CART_STATUS_ACTIVE, nil
	case StatusInactive:
		return proto.CartStatus_CART_STATUS_INACTIVE, nil
	default:
		return proto.CartStatus_CART_STATUS_ACTIVE, errors.New(StatusUnknown)
	}
}

func ConvertCartStatusToDB(status proto.CartStatus) string {
	switch status {
	case proto.CartStatus_CART_STATUS_ACTIVE:
		return StatusActive
	case proto.CartStatus_CART_STATUS_INACTIVE:
		return StatusInactive
	default:
		return StatusUnknown
	}
}

func ConvertDBAdvertStatusToEnum(dbStatus string) (proto.AdvertStatus, error) {
	switch dbStatus {
	case StatusActive:
		return proto.AdvertStatus_ADVERT_STATUS_ACTIVE, nil
	case StatusInactive:
		return proto.AdvertStatus_ADVERT_STATUS_INACTIVE, nil
	case StatusReserved:
		return proto.AdvertStatus_ADVERT_STATUS_RESERVED, nil
	default:
		return proto.AdvertStatus_ADVERT_STATUS_ACTIVE, errors.New(StatusUnknown)
	}
}

func ConvertAdvertStatusToDB(status proto.AdvertStatus) string {
	switch status {
	case proto.AdvertStatus_ADVERT_STATUS_ACTIVE:
		return StatusActive
	case proto.AdvertStatus_ADVERT_STATUS_INACTIVE:
		return StatusInactive
	case proto.AdvertStatus_ADVERT_STATUS_RESERVED:
		return StatusReserved
	default:
		return StatusUnknown
	}
}
