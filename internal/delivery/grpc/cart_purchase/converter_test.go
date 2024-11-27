package cart_purchase

import (
	"testing"

	proto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/cart_purchase/proto"
	"github.com/pkg/errors"
)

func TestConvertDBPurchaseStatusToEnum(t *testing.T) {
	tests := []struct {
		input    string
		expected proto.PurchaseStatus
		err      error
	}{
		{"PURCHASE_STATUS_PENDING", proto.PurchaseStatus_PURCHASE_STATUS_PENDING, nil},
		{"PURCHASE_STATUS_IN_PROGRESS", proto.PurchaseStatus_PURCHASE_STATUS_IN_PROGRESS, nil},
		{"PURCHASE_STATUS_COMPLETED", proto.PurchaseStatus_PURCHASE_STATUS_COMPLETED, nil},
		{"PURCHASE_STATUS_CANCELED", proto.PurchaseStatus_PURCHASE_STATUS_CANCELED, nil},
		{"unknown", proto.PurchaseStatus_PURCHASE_STATUS_PENDING, errors.New("unknown purchase status")},
	}

	for _, test := range tests {
		result, err := ConvertDBPurchaseStatusToEnum(test.input)
		if result != test.expected || (err != nil && err.Error() != test.err.Error()) {
			t.Errorf("ConvertDBPurchaseStatusToEnum(%q) = %v, %v; want %v, %v", test.input, result, err, test.expected, test.err)
		}
	}
}

func TestConvertDBPaymentMethodToEnum(t *testing.T) {
	tests := []struct {
		input    string
		expected proto.PaymentMethod
		err      error
	}{
		{"card", proto.PaymentMethod_PAYMENT_METHOD_CARD, nil},
		{"cash", proto.PaymentMethod_PAYMENT_METHOD_CASH, nil},
		{"unknown", proto.PaymentMethod_PAYMENT_METHOD_CARD, errors.New("unknown payment method")},
	}

	for _, test := range tests {
		result, err := ConvertDBPaymentMethodToEnum(test.input)
		if result != test.expected || (err != nil && err.Error() != test.err.Error()) {
			t.Errorf("ConvertDBPaymentMethodToEnum(%q) = %v, %v; want %v, %v", test.input, result, err, test.expected, test.err)
		}
	}
}

func TestConvertDBDeliveryMethodToEnum(t *testing.T) {
	tests := []struct {
		input    string
		expected proto.DeliveryMethod
		err      error
	}{
		{"pickup", proto.DeliveryMethod_DELIVERY_METHOD_PICKUP, nil},
		{"delivery", proto.DeliveryMethod_DELIVERY_METHOD_DELIVERY, nil},
		{"unknown", proto.DeliveryMethod_DELIVERY_METHOD_PICKUP, errors.New("unknown delivery method")},
	}

	for _, test := range tests {
		result, err := ConvertDBDeliveryMethodToEnum(test.input)
		if result != test.expected || (err != nil && err.Error() != test.err.Error()) {
			t.Errorf("ConvertDBDeliveryMethodToEnum(%q) = %v, %v; want %v, %v", test.input, result, err, test.expected, test.err)
		}
	}
}

func TestConvertDBCartStatusToEnum(t *testing.T) {
	tests := []struct {
		input    string
		expected proto.CartStatus
		err      error
	}{
		{"active", proto.CartStatus_CART_STATUS_ACTIVE, nil},
		{"inactive", proto.CartStatus_CART_STATUS_INACTIVE, nil},
		{"deleted", proto.CartStatus_CART_STATUS_DELETED, nil},
		{"unknown", proto.CartStatus_CART_STATUS_ACTIVE, errors.New("unknown cart status")},
	}

	for _, test := range tests {
		result, err := ConvertDBCartStatusToEnum(test.input)
		if result != test.expected || (err != nil && err.Error() != test.err.Error()) {
			t.Errorf("ConvertDBCartStatusToEnum(%q) = %v, %v; want %v, %v", test.input, result, err, test.expected, test.err)
		}
	}
}

func TestConvertDBAdvertStatusToEnum(t *testing.T) {
	tests := []struct {
		input    string
		expected proto.AdvertStatus
		err      error
	}{
		{"active", proto.AdvertStatus_ADVERT_STATUS_ACTIVE, nil},
		{"inactive", proto.AdvertStatus_ADVERT_STATUS_INACTIVE, nil},
		{"reserved", proto.AdvertStatus_ADVERT_STATUS_RESERVED, nil},
		{"unknown", proto.AdvertStatus_ADVERT_STATUS_ACTIVE, errors.New("unknown advert status")},
	}

	for _, test := range tests {
		result, err := ConvertDBAdvertStatusToEnum(test.input)
		if result != test.expected || (err != nil && err.Error() != test.err.Error()) {
			t.Errorf("ConvertDBAdvertStatusToEnum(%q) = %v, %v; want %v, %v", test.input, result, err, test.expected, test.err)
		}
	}
}
