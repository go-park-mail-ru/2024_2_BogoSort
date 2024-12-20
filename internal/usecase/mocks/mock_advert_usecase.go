// Code generated by MockGen. DO NOT EDIT.
// Source: internal/usecase/adverts.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	dto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockAdvertUseCase is a mock of AdvertUseCase interface.
type MockAdvertUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockAdvertUseCaseMockRecorder
}

// MockAdvertUseCaseMockRecorder is the mock recorder for MockAdvertUseCase.
type MockAdvertUseCaseMockRecorder struct {
	mock *MockAdvertUseCase
}

// NewMockAdvertUseCase creates a new mock instance.
func NewMockAdvertUseCase(ctrl *gomock.Controller) *MockAdvertUseCase {
	mock := &MockAdvertUseCase{ctrl: ctrl}
	mock.recorder = &MockAdvertUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdvertUseCase) EXPECT() *MockAdvertUseCaseMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockAdvertUseCase) Add(advert *dto.AdvertRequest, userId uuid.UUID) (*dto.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", advert, userId)
	ret0, _ := ret[0].(*dto.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockAdvertUseCaseMockRecorder) Add(advert, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockAdvertUseCase)(nil).Add), advert, userId)
}

// AddToSaved mocks base method.
func (m *MockAdvertUseCase) AddToSaved(advertId, userId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToSaved", advertId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToSaved indicates an expected call of AddToSaved.
func (mr *MockAdvertUseCaseMockRecorder) AddToSaved(advertId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToSaved", reflect.TypeOf((*MockAdvertUseCase)(nil).AddToSaved), advertId, userId)
}

// AddViewed mocks base method.
func (m *MockAdvertUseCase) AddViewed(advertId, userId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddViewed", advertId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddViewed indicates an expected call of AddViewed.
func (mr *MockAdvertUseCaseMockRecorder) AddViewed(advertId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddViewed", reflect.TypeOf((*MockAdvertUseCase)(nil).AddViewed), advertId, userId)
}

// DeleteById mocks base method.
func (m *MockAdvertUseCase) DeleteById(advertId, userId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", advertId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockAdvertUseCaseMockRecorder) DeleteById(advertId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockAdvertUseCase)(nil).DeleteById), advertId, userId)
}

// Get mocks base method.
func (m *MockAdvertUseCase) Get(limit, offset int, userId uuid.UUID) ([]*dto.PreviewAdvertCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", limit, offset, userId)
	ret0, _ := ret[0].([]*dto.PreviewAdvertCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockAdvertUseCaseMockRecorder) Get(limit, offset, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAdvertUseCase)(nil).Get), limit, offset, userId)
}

// GetByCartId mocks base method.
func (m *MockAdvertUseCase) GetByCartId(cartId, userId uuid.UUID) ([]*dto.PreviewAdvertCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByCartId", cartId, userId)
	ret0, _ := ret[0].([]*dto.PreviewAdvertCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByCartId indicates an expected call of GetByCartId.
func (mr *MockAdvertUseCaseMockRecorder) GetByCartId(cartId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByCartId", reflect.TypeOf((*MockAdvertUseCase)(nil).GetByCartId), cartId, userId)
}

// GetByCategoryId mocks base method.
func (m *MockAdvertUseCase) GetByCategoryId(categoryId, userId uuid.UUID) ([]*dto.PreviewAdvertCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByCategoryId", categoryId, userId)
	ret0, _ := ret[0].([]*dto.PreviewAdvertCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByCategoryId indicates an expected call of GetByCategoryId.
func (mr *MockAdvertUseCaseMockRecorder) GetByCategoryId(categoryId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByCategoryId", reflect.TypeOf((*MockAdvertUseCase)(nil).GetByCategoryId), categoryId, userId)
}

// GetById mocks base method.
func (m *MockAdvertUseCase) GetById(advertId, userId uuid.UUID) (*dto.AdvertCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", advertId, userId)
	ret0, _ := ret[0].(*dto.AdvertCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockAdvertUseCaseMockRecorder) GetById(advertId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockAdvertUseCase)(nil).GetById), advertId, userId)
}

// GetBySellerId mocks base method.
func (m *MockAdvertUseCase) GetBySellerId(userId, sellerId uuid.UUID) ([]*dto.PreviewAdvertCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBySellerId", userId, sellerId)
	ret0, _ := ret[0].([]*dto.PreviewAdvertCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBySellerId indicates an expected call of GetBySellerId.
func (mr *MockAdvertUseCaseMockRecorder) GetBySellerId(userId, sellerId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBySellerId", reflect.TypeOf((*MockAdvertUseCase)(nil).GetBySellerId), userId, sellerId)
}

// GetByUserId mocks base method.
func (m *MockAdvertUseCase) GetByUserId(userId uuid.UUID) ([]*dto.MyPreviewAdvertCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserId", userId)
	ret0, _ := ret[0].([]*dto.MyPreviewAdvertCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserId indicates an expected call of GetByUserId.
func (mr *MockAdvertUseCaseMockRecorder) GetByUserId(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserId", reflect.TypeOf((*MockAdvertUseCase)(nil).GetByUserId), userId)
}

// GetSavedByUserId mocks base method.
func (m *MockAdvertUseCase) GetSavedByUserId(userId uuid.UUID) ([]*dto.PreviewAdvertCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSavedByUserId", userId)
	ret0, _ := ret[0].([]*dto.PreviewAdvertCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSavedByUserId indicates an expected call of GetSavedByUserId.
func (mr *MockAdvertUseCaseMockRecorder) GetSavedByUserId(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSavedByUserId", reflect.TypeOf((*MockAdvertUseCase)(nil).GetSavedByUserId), userId)
}

// RemoveFromSaved mocks base method.
func (m *MockAdvertUseCase) RemoveFromSaved(advertId, userId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFromSaved", advertId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFromSaved indicates an expected call of RemoveFromSaved.
func (mr *MockAdvertUseCaseMockRecorder) RemoveFromSaved(advertId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFromSaved", reflect.TypeOf((*MockAdvertUseCase)(nil).RemoveFromSaved), advertId, userId)
}

// Search mocks base method.
func (m *MockAdvertUseCase) Search(query string, batchSize, limit, offset int, userId uuid.UUID) ([]*dto.PreviewAdvertCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", query, batchSize, limit, offset, userId)
	ret0, _ := ret[0].([]*dto.PreviewAdvertCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockAdvertUseCaseMockRecorder) Search(query, batchSize, limit, offset, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockAdvertUseCase)(nil).Search), query, batchSize, limit, offset, userId)
}

// Update mocks base method.
func (m *MockAdvertUseCase) Update(advert *dto.AdvertRequest, userId, advertId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", advert, userId, advertId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockAdvertUseCaseMockRecorder) Update(advert, userId, advertId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAdvertUseCase)(nil).Update), advert, userId, advertId)
}

// UpdateStatus mocks base method.
func (m *MockAdvertUseCase) UpdateStatus(advertId, userId uuid.UUID, status dto.AdvertStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", advertId, userId, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockAdvertUseCaseMockRecorder) UpdateStatus(advertId, userId, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockAdvertUseCase)(nil).UpdateStatus), advertId, userId, status)
}

// UploadImage mocks base method.
func (m *MockAdvertUseCase) UploadImage(advertId, imageId, userId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadImage", advertId, imageId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadImage indicates an expected call of UploadImage.
func (mr *MockAdvertUseCaseMockRecorder) UploadImage(advertId, imageId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadImage", reflect.TypeOf((*MockAdvertUseCase)(nil).UploadImage), advertId, imageId, userId)
}
