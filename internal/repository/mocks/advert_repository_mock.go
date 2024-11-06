// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/advert.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	entity "github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
)

// MockAdvertRepository is a mock of AdvertRepository interface.
type MockAdvertRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAdvertRepositoryMockRecorder
}

// MockAdvertRepositoryMockRecorder is the mock recorder for MockAdvertRepository.
type MockAdvertRepositoryMockRecorder struct {
	mock *MockAdvertRepository
}

// NewMockAdvertRepository creates a new mock instance.
func NewMockAdvertRepository(ctrl *gomock.Controller) *MockAdvertRepository {
	mock := &MockAdvertRepository{ctrl: ctrl}
	mock.recorder = &MockAdvertRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdvertRepository) EXPECT() *MockAdvertRepositoryMockRecorder {
	return m.recorder
}

// AddAdvert mocks base method.
func (m *MockAdvertRepository) AddAdvert(advert *entity.Advert) (*entity.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAdvert", advert)
	ret0, _ := ret[0].(*entity.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddAdvert indicates an expected call of AddAdvert.
func (mr *MockAdvertRepositoryMockRecorder) AddAdvert(advert interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAdvert", reflect.TypeOf((*MockAdvertRepository)(nil).AddAdvert), advert)
}

// BeginTransaction mocks base method.
func (m *MockAdvertRepository) BeginTransaction() (pgx.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTransaction")
	ret0, _ := ret[0].(pgx.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTransaction indicates an expected call of BeginTransaction.
func (mr *MockAdvertRepositoryMockRecorder) BeginTransaction() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTransaction", reflect.TypeOf((*MockAdvertRepository)(nil).BeginTransaction))
}

// DeleteAdvertById mocks base method.
func (m *MockAdvertRepository) DeleteAdvertById(advertId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAdvertById", advertId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAdvertById indicates an expected call of DeleteAdvertById.
func (mr *MockAdvertRepositoryMockRecorder) DeleteAdvertById(advertId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAdvertById", reflect.TypeOf((*MockAdvertRepository)(nil).DeleteAdvertById), advertId)
}

// GetAdvertById mocks base method.
func (m *MockAdvertRepository) GetAdvertById(advertId uuid.UUID) (*entity.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdvertById", advertId)
	ret0, _ := ret[0].(*entity.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdvertById indicates an expected call of GetAdvertById.
func (mr *MockAdvertRepositoryMockRecorder) GetAdvertById(advertId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdvertById", reflect.TypeOf((*MockAdvertRepository)(nil).GetAdvertById), advertId)
}

// GetAdverts mocks base method.
func (m *MockAdvertRepository) GetAdverts(limit, offset int) ([]*entity.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdverts", limit, offset)
	ret0, _ := ret[0].([]*entity.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdverts indicates an expected call of GetAdverts.
func (mr *MockAdvertRepositoryMockRecorder) GetAdverts(limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdverts", reflect.TypeOf((*MockAdvertRepository)(nil).GetAdverts), limit, offset)
}

// GetAdvertsByCartId mocks base method.
func (m *MockAdvertRepository) GetAdvertsByCartId(cartId uuid.UUID) ([]*entity.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdvertsByCartId", cartId)
	ret0, _ := ret[0].([]*entity.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdvertsByCartId indicates an expected call of GetAdvertsByCartId.
func (mr *MockAdvertRepositoryMockRecorder) GetAdvertsByCartId(cartId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdvertsByCartId", reflect.TypeOf((*MockAdvertRepository)(nil).GetAdvertsByCartId), cartId)
}

// GetAdvertsByCategoryId mocks base method.
func (m *MockAdvertRepository) GetAdvertsByCategoryId(categoryId uuid.UUID) ([]*entity.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdvertsByCategoryId", categoryId)
	ret0, _ := ret[0].([]*entity.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdvertsByCategoryId indicates an expected call of GetAdvertsByCategoryId.
func (mr *MockAdvertRepositoryMockRecorder) GetAdvertsByCategoryId(categoryId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdvertsByCategoryId", reflect.TypeOf((*MockAdvertRepository)(nil).GetAdvertsByCategoryId), categoryId)
}

// GetAdvertsBySellerId mocks base method.
func (m *MockAdvertRepository) GetAdvertsBySellerId(sellerId uuid.UUID) ([]*entity.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdvertsBySellerId", sellerId)
	ret0, _ := ret[0].([]*entity.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdvertsBySellerId indicates an expected call of GetAdvertsBySellerId.
func (mr *MockAdvertRepositoryMockRecorder) GetAdvertsBySellerId(sellerId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdvertsBySellerId", reflect.TypeOf((*MockAdvertRepository)(nil).GetAdvertsBySellerId), sellerId)
}

// UpdateAdvert mocks base method.
func (m *MockAdvertRepository) UpdateAdvert(advert *entity.Advert) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAdvert", advert)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAdvert indicates an expected call of UpdateAdvert.
func (mr *MockAdvertRepositoryMockRecorder) UpdateAdvert(advert interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAdvert", reflect.TypeOf((*MockAdvertRepository)(nil).UpdateAdvert), advert)
}

// UpdateAdvertStatus mocks base method.
func (m *MockAdvertRepository) UpdateAdvertStatus(tx pgx.Tx, advertId uuid.UUID, status entity.AdvertStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAdvertStatus", tx, advertId, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAdvertStatus indicates an expected call of UpdateAdvertStatus.
func (mr *MockAdvertRepositoryMockRecorder) UpdateAdvertStatus(tx, advertId, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAdvertStatus", reflect.TypeOf((*MockAdvertRepository)(nil).UpdateAdvertStatus), tx, advertId, status)
}

// UploadImage mocks base method.
func (m *MockAdvertRepository) UploadImage(advertId, imageId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadImage", advertId, imageId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadImage indicates an expected call of UploadImage.
func (mr *MockAdvertRepositoryMockRecorder) UploadImage(advertId, imageId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadImage", reflect.TypeOf((*MockAdvertRepository)(nil).UploadImage), advertId, imageId)
}