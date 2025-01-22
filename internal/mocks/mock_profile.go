// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mongodb/mongodb-cli/mongocli/v2/internal/config (interfaces: SetSaver)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSetSaver is a mock of SetSaver interface.
type MockSetSaver struct {
	ctrl     *gomock.Controller
	recorder *MockSetSaverMockRecorder
}

// MockSetSaverMockRecorder is the mock recorder for MockSetSaver.
type MockSetSaverMockRecorder struct {
	mock *MockSetSaver
}

// NewMockSetSaver creates a new mock instance.
func NewMockSetSaver(ctrl *gomock.Controller) *MockSetSaver {
	mock := &MockSetSaver{ctrl: ctrl}
	mock.recorder = &MockSetSaverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSetSaver) EXPECT() *MockSetSaverMockRecorder {
	return m.recorder
}

// Save mocks base method.
func (m *MockSetSaver) Save() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save")
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockSetSaverMockRecorder) Save() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSetSaver)(nil).Save))
}

// Set mocks base method.
func (m *MockSetSaver) Set(arg0 string, arg1 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", arg0, arg1)
}

// Set indicates an expected call of Set.
func (mr *MockSetSaverMockRecorder) Set(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockSetSaver)(nil).Set), arg0, arg1)
}

// SetGlobal mocks base method.
func (m *MockSetSaver) SetGlobal(arg0 string, arg1 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetGlobal", arg0, arg1)
}

// SetGlobal indicates an expected call of SetGlobal.
func (mr *MockSetSaverMockRecorder) SetGlobal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetGlobal", reflect.TypeOf((*MockSetSaver)(nil).SetGlobal), arg0, arg1)
}
