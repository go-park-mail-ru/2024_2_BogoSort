package cart_purchase

import (
	proto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/cart_purchase/proto"
	"github.com/pkg/errors"
)

func ConvertDBPurchaseStatusToEnum(dbStatus string) (proto.PurchaseStatus, error) {
	switch dbStatus {
	case "PURCHASE_STATUS_PENDING":
		return proto.PurchaseStatus_PURCHASE_STATUS_PENDING, nil
	case "PURCHASE_STATUS_IN_PROGRESS":
		return proto.PurchaseStatus_PURCHASE_STATUS_IN_PROGRESS, nil
	case "PURCHASE_STATUS_COMPLETED":
		return proto.PurchaseStatus_PURCHASE_STATUS_COMPLETED, nil
	case "PURCHASE_STATUS_CANCELED":
		return proto.PurchaseStatus_PURCHASE_STATUS_CANCELED, nil
	default:
		return proto.PurchaseStatus_PURCHASE_STATUS_PENDING, errors.New("unknown purchase status")
	}
}

func ConvertPurchaseStatusToDB(status proto.PurchaseStatus) string {
	switch status {
	case proto.PurchaseStatus_PURCHASE_STATUS_PENDING:
		return "pending"
	case proto.PurchaseStatus_PURCHASE_STATUS_IN_PROGRESS:
		return "in_progress"
	case proto.PurchaseStatus_PURCHASE_STATUS_COMPLETED:
		return "completed"
	case proto.PurchaseStatus_PURCHASE_STATUS_CANCELED:
		return "canceled"
	default:
		return "unknown"
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
	case "active":
		return proto.CartStatus_CART_STATUS_ACTIVE, nil
	case "inactive":
		return proto.CartStatus_CART_STATUS_INACTIVE, nil
	case "deleted":
		return proto.CartStatus_CART_STATUS_DELETED, nil
	default:
		return proto.CartStatus_CART_STATUS_ACTIVE, errors.New("unknown cart status")
	}
}

func ConvertCartStatusToDB(status proto.CartStatus) string {
	switch status {
	case proto.CartStatus_CART_STATUS_ACTIVE:
		return "active"
	case proto.CartStatus_CART_STATUS_INACTIVE:
		return "inactive"
	case proto.CartStatus_CART_STATUS_DELETED:
		return "deleted"
	default:
		return "unknown"
	}
}

func ConvertDBAdvertStatusToEnum(dbStatus string) (proto.AdvertStatus, error) {
	switch dbStatus {
	case "active":
		return proto.AdvertStatus_ADVERT_STATUS_ACTIVE, nil
	case "inactive":
		return proto.AdvertStatus_ADVERT_STATUS_INACTIVE, nil
	case "reserved":
		return proto.AdvertStatus_ADVERT_STATUS_RESERVED, nil
	default:
		return proto.AdvertStatus_ADVERT_STATUS_ACTIVE, errors.New("unknown advert status")
	}
}

func ConvertAdvertStatusToDB(status proto.AdvertStatus) string {
	switch status {
	case proto.AdvertStatus_ADVERT_STATUS_ACTIVE:
		return "active"
	case proto.AdvertStatus_ADVERT_STATUS_INACTIVE:
		return "inactive"
	case proto.AdvertStatus_ADVERT_STATUS_RESERVED:
		return "reserved"
	default:
		return "unknown"
	}
}
