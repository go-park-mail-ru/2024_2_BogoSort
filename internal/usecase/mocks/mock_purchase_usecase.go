// Code generated by MockGen. DO NOT EDIT.
// Source: internal/usecase/purchase.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	dto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockPurchase is a mock of Purchase interface.
type MockPurchase struct {
	ctrl     *gomock.Controller
	recorder *MockPurchaseMockRecorder
}

// MockPurchaseMockRecorder is the mock recorder for MockPurchase.
type MockPurchaseMockRecorder struct {
	mock *MockPurchase
}

// NewMockPurchase creates a new mock instance.
func NewMockPurchase(ctrl *gomock.Controller) *MockPurchase {
	mock := &MockPurchase{ctrl: ctrl}
	mock.recorder = &MockPurchaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPurchase) EXPECT() *MockPurchaseMockRecorder {
	return m.recorder
}

// AddPurchase mocks base method.
func (m *MockPurchase) AddPurchase(purchaseRequest dto.PurchaseRequest) (*dto.PurchaseResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPurchase", purchaseRequest)
	ret0, _ := ret[0].(*dto.PurchaseResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddPurchase indicates an expected call of AddPurchase.
func (mr *MockPurchaseMockRecorder) AddPurchase(purchaseRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPurchase", reflect.TypeOf((*MockPurchase)(nil).AddPurchase), purchaseRequest)
}

// GetPurchasesByUserID mocks base method.
func (m *MockPurchase) GetPurchasesByUserID(userID uuid.UUID) ([]*dto.PurchaseResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPurchasesByUserID", userID)
	ret0, _ := ret[0].([]*dto.PurchaseResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPurchasesByUserID indicates an expected call of GetPurchasesByUserID.
func (mr *MockPurchaseMockRecorder) GetPurchasesByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPurchasesByUserID", reflect.TypeOf((*MockPurchase)(nil).GetPurchasesByUserID), userID)
}