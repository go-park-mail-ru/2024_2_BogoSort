// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/static.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockStaticRepository is a mock of StaticRepository interface.
type MockStaticRepository struct {
	ctrl     *gomock.Controller
	recorder *MockStaticRepositoryMockRecorder
}

// MockStaticRepositoryMockRecorder is the mock recorder for MockStaticRepository.
type MockStaticRepositoryMockRecorder struct {
	mock *MockStaticRepository
}

// NewMockStaticRepository creates a new mock instance.
func NewMockStaticRepository(ctrl *gomock.Controller) *MockStaticRepository {
	mock := &MockStaticRepository{ctrl: ctrl}
	mock.recorder = &MockStaticRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStaticRepository) EXPECT() *MockStaticRepositoryMockRecorder {
	return m.recorder
}

// GetStatic mocks base method.
func (m *MockStaticRepository) GetStatic(staticID uuid.UUID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatic", staticID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStatic indicates an expected call of GetStatic.
func (mr *MockStaticRepositoryMockRecorder) GetStatic(staticID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatic", reflect.TypeOf((*MockStaticRepository)(nil).GetStatic), staticID)
}

// UploadStatic mocks base method.
func (m *MockStaticRepository) UploadStatic(path, filename string, data []byte) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadStatic", path, filename, data)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadStatic indicates an expected call of UploadStatic.
func (mr *MockStaticRepositoryMockRecorder) UploadStatic(path, filename, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadStatic", reflect.TypeOf((*MockStaticRepository)(nil).UploadStatic), path, filename, data)
}