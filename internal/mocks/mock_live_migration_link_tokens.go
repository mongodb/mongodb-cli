// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mongodb/mongodb-cli/mongocli/v2/internal/store (interfaces: LinkTokenDeleter)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLinkTokenDeleter is a mock of LinkTokenDeleter interface.
type MockLinkTokenDeleter struct {
	ctrl     *gomock.Controller
	recorder *MockLinkTokenDeleterMockRecorder
}

// MockLinkTokenDeleterMockRecorder is the mock recorder for MockLinkTokenDeleter.
type MockLinkTokenDeleterMockRecorder struct {
	mock *MockLinkTokenDeleter
}

// NewMockLinkTokenDeleter creates a new mock instance.
func NewMockLinkTokenDeleter(ctrl *gomock.Controller) *MockLinkTokenDeleter {
	mock := &MockLinkTokenDeleter{ctrl: ctrl}
	mock.recorder = &MockLinkTokenDeleterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLinkTokenDeleter) EXPECT() *MockLinkTokenDeleterMockRecorder {
	return m.recorder
}

// DeleteLinkToken mocks base method.
func (m *MockLinkTokenDeleter) DeleteLinkToken(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLinkToken", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLinkToken indicates an expected call of DeleteLinkToken.
func (mr *MockLinkTokenDeleterMockRecorder) DeleteLinkToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLinkToken", reflect.TypeOf((*MockLinkTokenDeleter)(nil).DeleteLinkToken), arg0)
}
