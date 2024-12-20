// Code generated by MockGen. DO NOT EDIT.
// Source: internal/usecase/static.go

// Package mocks is a generated GoMock package.
package mocks

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockStaticUseCase is a mock of StaticUseCase interface.
type MockStaticUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockStaticUseCaseMockRecorder
}

// MockStaticUseCaseMockRecorder is the mock recorder for MockStaticUseCase.
type MockStaticUseCaseMockRecorder struct {
	mock *MockStaticUseCase
}

// NewMockStaticUseCase creates a new mock instance.
func NewMockStaticUseCase(ctrl *gomock.Controller) *MockStaticUseCase {
	mock := &MockStaticUseCase{ctrl: ctrl}
	mock.recorder = &MockStaticUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStaticUseCase) EXPECT() *MockStaticUseCaseMockRecorder {
	return m.recorder
}

// GetAvatar mocks base method.
func (m *MockStaticUseCase) GetAvatar(staticID uuid.UUID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvatar", staticID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvatar indicates an expected call of GetAvatar.
func (mr *MockStaticUseCaseMockRecorder) GetAvatar(staticID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvatar", reflect.TypeOf((*MockStaticUseCase)(nil).GetAvatar), staticID)
}

// GetStatic mocks base method.
func (m *MockStaticUseCase) GetStatic(id uuid.UUID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatic", id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStatic indicates an expected call of GetStatic.
func (mr *MockStaticUseCaseMockRecorder) GetStatic(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatic", reflect.TypeOf((*MockStaticUseCase)(nil).GetStatic), id)
}

// GetStaticFile mocks base method.
func (m *MockStaticUseCase) GetStaticFile(uri string) (io.ReadSeeker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStaticFile", uri)
	ret0, _ := ret[0].(io.ReadSeeker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStaticFile indicates an expected call of GetStaticFile.
func (mr *MockStaticUseCaseMockRecorder) GetStaticFile(uri interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStaticFile", reflect.TypeOf((*MockStaticUseCase)(nil).GetStaticFile), uri)
}

// UploadStatic mocks base method.
func (m *MockStaticUseCase) UploadStatic(data io.ReadSeeker) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadStatic", data)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadStatic indicates an expected call of UploadStatic.
func (mr *MockStaticUseCaseMockRecorder) UploadStatic(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadStatic", reflect.TypeOf((*MockStaticUseCase)(nil).UploadStatic), data)
}
