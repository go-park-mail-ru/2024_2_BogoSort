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

// Add mocks base method.
func (m *MockAdvertRepository) Add(advert *entity.Advert) (*entity.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", advert)
	ret0, _ := ret[0].(*entity.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockAdvertRepositoryMockRecorder) Add(advert interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockAdvertRepository)(nil).Add), advert)
}

// AddToSaved mocks base method.
func (m *MockAdvertRepository) AddToSaved(advertId, userId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToSaved", advertId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToSaved indicates an expected call of AddToSaved.
func (mr *MockAdvertRepositoryMockRecorder) AddToSaved(advertId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToSaved", reflect.TypeOf((*MockAdvertRepository)(nil).AddToSaved), advertId, userId)
}

// AddViewed mocks base method.
func (m *MockAdvertRepository) AddViewed(userId, advertId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddViewed", userId, advertId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddViewed indicates an expected call of AddViewed.
func (mr *MockAdvertRepositoryMockRecorder) AddViewed(userId, advertId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddViewed", reflect.TypeOf((*MockAdvertRepository)(nil).AddViewed), userId, advertId)
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

// CheckIfExists mocks base method.
func (m *MockAdvertRepository) CheckIfExists(advertId uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckIfExists", advertId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckIfExists indicates an expected call of CheckIfExists.
func (mr *MockAdvertRepositoryMockRecorder) CheckIfExists(advertId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIfExists", reflect.TypeOf((*MockAdvertRepository)(nil).CheckIfExists), advertId)
}

// Count mocks base method.
func (m *MockAdvertRepository) Count() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockAdvertRepositoryMockRecorder) Count() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockAdvertRepository)(nil).Count))
}

// DeleteById mocks base method.
func (m *MockAdvertRepository) DeleteById(advertId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", advertId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockAdvertRepositoryMockRecorder) DeleteById(advertId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockAdvertRepository)(nil).DeleteById), advertId)
}

// DeleteFromSaved mocks base method.
func (m *MockAdvertRepository) DeleteFromSaved(userId, advertId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFromSaved", userId, advertId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFromSaved indicates an expected call of DeleteFromSaved.
func (mr *MockAdvertRepositoryMockRecorder) DeleteFromSaved(userId, advertId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFromSaved", reflect.TypeOf((*MockAdvertRepository)(nil).DeleteFromSaved), userId, advertId)
}

// Get mocks base method.
func (m *MockAdvertRepository) Get(limit, offset int, userId uuid.UUID) ([]*entity.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", limit, offset, userId)
	ret0, _ := ret[0].([]*entity.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockAdvertRepositoryMockRecorder) Get(limit, offset, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAdvertRepository)(nil).Get), limit, offset, userId)
}

// GetByCartId mocks base method.
func (m *MockAdvertRepository) GetByCartId(cartId, userId uuid.UUID) ([]*entity.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByCartId", cartId, userId)
	ret0, _ := ret[0].([]*entity.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByCartId indicates an expected call of GetByCartId.
func (mr *MockAdvertRepositoryMockRecorder) GetByCartId(cartId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByCartId", reflect.TypeOf((*MockAdvertRepository)(nil).GetByCartId), cartId, userId)
}

// GetByCategoryId mocks base method.
func (m *MockAdvertRepository) GetByCategoryId(categoryId, userId uuid.UUID) ([]*entity.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByCategoryId", categoryId, userId)
	ret0, _ := ret[0].([]*entity.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByCategoryId indicates an expected call of GetByCategoryId.
func (mr *MockAdvertRepositoryMockRecorder) GetByCategoryId(categoryId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByCategoryId", reflect.TypeOf((*MockAdvertRepository)(nil).GetByCategoryId), categoryId, userId)
}

// GetById mocks base method.
func (m *MockAdvertRepository) GetById(advertId, userId uuid.UUID) (*entity.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", advertId, userId)
	ret0, _ := ret[0].(*entity.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockAdvertRepositoryMockRecorder) GetById(advertId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockAdvertRepository)(nil).GetById), advertId, userId)
}

// GetBySellerId mocks base method.
func (m *MockAdvertRepository) GetBySellerId(sellerId, userId uuid.UUID) ([]*entity.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBySellerId", sellerId, userId)
	ret0, _ := ret[0].([]*entity.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBySellerId indicates an expected call of GetBySellerId.
func (mr *MockAdvertRepositoryMockRecorder) GetBySellerId(sellerId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBySellerId", reflect.TypeOf((*MockAdvertRepository)(nil).GetBySellerId), sellerId, userId)
}

// GetByUserId mocks base method.
func (m *MockAdvertRepository) GetByUserId(sellerId, userId uuid.UUID) ([]*entity.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserId", sellerId, userId)
	ret0, _ := ret[0].([]*entity.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserId indicates an expected call of GetByUserId.
func (mr *MockAdvertRepositoryMockRecorder) GetByUserId(sellerId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserId", reflect.TypeOf((*MockAdvertRepository)(nil).GetByUserId), sellerId, userId)
}

// GetSavedByUserId mocks base method.
func (m *MockAdvertRepository) GetSavedByUserId(userId uuid.UUID) ([]*entity.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSavedByUserId", userId)
	ret0, _ := ret[0].([]*entity.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSavedByUserId indicates an expected call of GetSavedByUserId.
func (mr *MockAdvertRepositoryMockRecorder) GetSavedByUserId(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSavedByUserId", reflect.TypeOf((*MockAdvertRepository)(nil).GetSavedByUserId), userId)
}

// Search mocks base method.
func (m *MockAdvertRepository) Search(query string, limit, offset int, userId uuid.UUID) ([]*entity.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", query, limit, offset, userId)
	ret0, _ := ret[0].([]*entity.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockAdvertRepositoryMockRecorder) Search(query, limit, offset, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockAdvertRepository)(nil).Search), query, limit, offset, userId)
}

// Update mocks base method.
func (m *MockAdvertRepository) Update(advert *entity.Advert) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", advert)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockAdvertRepositoryMockRecorder) Update(advert interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAdvertRepository)(nil).Update), advert)
}

// UpdateStatus mocks base method.
func (m *MockAdvertRepository) UpdateStatus(tx pgx.Tx, advertId uuid.UUID, status entity.AdvertStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", tx, advertId, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockAdvertRepositoryMockRecorder) UpdateStatus(tx, advertId, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockAdvertRepository)(nil).UpdateStatus), tx, advertId, status)
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
