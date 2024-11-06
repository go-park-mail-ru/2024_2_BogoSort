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

// AddAdvert mocks base method.
func (m *MockAdvertUseCase) AddAdvert(advert *dto.AdvertRequest, userId uuid.UUID) (*dto.AdvertResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAdvert", advert, userId)
	ret0, _ := ret[0].(*dto.AdvertResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddAdvert indicates an expected call of AddAdvert.
func (mr *MockAdvertUseCaseMockRecorder) AddAdvert(advert, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAdvert", reflect.TypeOf((*MockAdvertUseCase)(nil).AddAdvert), advert, userId)
}

// DeleteAdvertById mocks base method.
func (m *MockAdvertUseCase) DeleteAdvertById(advertId, userId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAdvertById", advertId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAdvertById indicates an expected call of DeleteAdvertById.
func (mr *MockAdvertUseCaseMockRecorder) DeleteAdvertById(advertId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAdvertById", reflect.TypeOf((*MockAdvertUseCase)(nil).DeleteAdvertById), advertId, userId)
}

// GetAdvertById mocks base method.
func (m *MockAdvertUseCase) GetAdvertById(advertId uuid.UUID) (*dto.AdvertResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdvertById", advertId)
	ret0, _ := ret[0].(*dto.AdvertResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdvertById indicates an expected call of GetAdvertById.
func (mr *MockAdvertUseCaseMockRecorder) GetAdvertById(advertId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdvertById", reflect.TypeOf((*MockAdvertUseCase)(nil).GetAdvertById), advertId)
}

// GetAdverts mocks base method.
func (m *MockAdvertUseCase) GetAdverts(limit, offset int) ([]*dto.AdvertResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdverts", limit, offset)
	ret0, _ := ret[0].([]*dto.AdvertResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdverts indicates an expected call of GetAdverts.
func (mr *MockAdvertUseCaseMockRecorder) GetAdverts(limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdverts", reflect.TypeOf((*MockAdvertUseCase)(nil).GetAdverts), limit, offset)
}

// GetAdvertsByCartId mocks base method.
func (m *MockAdvertUseCase) GetAdvertsByCartId(cartId uuid.UUID) ([]*dto.AdvertResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdvertsByCartId", cartId)
	ret0, _ := ret[0].([]*dto.AdvertResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdvertsByCartId indicates an expected call of GetAdvertsByCartId.
func (mr *MockAdvertUseCaseMockRecorder) GetAdvertsByCartId(cartId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdvertsByCartId", reflect.TypeOf((*MockAdvertUseCase)(nil).GetAdvertsByCartId), cartId)
}

// GetAdvertsByCategoryId mocks base method.
func (m *MockAdvertUseCase) GetAdvertsByCategoryId(categoryId uuid.UUID) ([]*dto.AdvertResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdvertsByCategoryId", categoryId)
	ret0, _ := ret[0].([]*dto.AdvertResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdvertsByCategoryId indicates an expected call of GetAdvertsByCategoryId.
func (mr *MockAdvertUseCaseMockRecorder) GetAdvertsByCategoryId(categoryId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdvertsByCategoryId", reflect.TypeOf((*MockAdvertUseCase)(nil).GetAdvertsByCategoryId), categoryId)
}

// GetAdvertsByUserId mocks base method.
func (m *MockAdvertUseCase) GetAdvertsByUserId(userId uuid.UUID) ([]*dto.AdvertResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdvertsByUserId", userId)
	ret0, _ := ret[0].([]*dto.AdvertResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdvertsByUserId indicates an expected call of GetAdvertsByUserId.
func (mr *MockAdvertUseCaseMockRecorder) GetAdvertsByUserId(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdvertsByUserId", reflect.TypeOf((*MockAdvertUseCase)(nil).GetAdvertsByUserId), userId)
}

// UpdateAdvert mocks base method.
func (m *MockAdvertUseCase) UpdateAdvert(advert *dto.AdvertRequest, userId, advertId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAdvert", advert, userId, advertId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAdvert indicates an expected call of UpdateAdvert.
func (mr *MockAdvertUseCaseMockRecorder) UpdateAdvert(advert, userId, advertId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAdvert", reflect.TypeOf((*MockAdvertUseCase)(nil).UpdateAdvert), advert, userId, advertId)
}

// UpdateAdvertStatus mocks base method.
func (m *MockAdvertUseCase) UpdateAdvertStatus(advertId uuid.UUID, status dto.AdvertStatus, userId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAdvertStatus", advertId, status, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAdvertStatus indicates an expected call of UpdateAdvertStatus.
func (mr *MockAdvertUseCaseMockRecorder) UpdateAdvertStatus(advertId, status, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAdvertStatus", reflect.TypeOf((*MockAdvertUseCase)(nil).UpdateAdvertStatus), advertId, status, userId)
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
