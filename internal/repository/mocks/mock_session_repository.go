package mocks

import (
	reflect "reflect"

	"github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

type SessionRepository interface {
	Create(userID uuid.UUID) (string, error)
	Delete(sessionID string) error
}

// MockSession is a mock of Session interface
type MockSession struct {
	ctrl     *gomock.Controller
	recorder *MockSessionMockRecorder
}

type MockSessionMockRecorder struct {
	mock *MockSession
}

func (m *MockSession) EXPECT() *MockSessionMockRecorder {
	return m.recorder
}

func NewMockSession(ctrl *gomock.Controller) *MockSession {
	mock := &MockSession{ctrl: ctrl}
	mock.recorder = &MockSessionMockRecorder{mock}
	return mock
}

func (m *MockSession) Delete(sessionID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", sessionID)
	ret0, _ := ret[0].(error)
	return ret0
}
func (mr *MockSessionMockRecorder) Delete(sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSession)(nil).Delete), sessionID)
}

func (m *MockSession) Create(userID uuid.UUID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockSessionMockRecorder) Create(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSession)(nil).Create), userID)
}

func (m *MockSession) Get(sessionID string) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", sessionID)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockSessionMockRecorder) Get(sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSession)(nil).Get), sessionID)
}
