// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mongodb/mongodb-cli/mongocli/v2/internal/store (interfaces: ServiceVersionDescriber)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	opsmngr "go.mongodb.org/ops-manager/opsmngr"
)

// MockServiceVersionDescriber is a mock of ServiceVersionDescriber interface.
type MockServiceVersionDescriber struct {
	ctrl     *gomock.Controller
	recorder *MockServiceVersionDescriberMockRecorder
}

// MockServiceVersionDescriberMockRecorder is the mock recorder for MockServiceVersionDescriber.
type MockServiceVersionDescriberMockRecorder struct {
	mock *MockServiceVersionDescriber
}

// NewMockServiceVersionDescriber creates a new mock instance.
func NewMockServiceVersionDescriber(ctrl *gomock.Controller) *MockServiceVersionDescriber {
	mock := &MockServiceVersionDescriber{ctrl: ctrl}
	mock.recorder = &MockServiceVersionDescriberMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceVersionDescriber) EXPECT() *MockServiceVersionDescriberMockRecorder {
	return m.recorder
}

// ServiceVersion mocks base method.
func (m *MockServiceVersionDescriber) ServiceVersion() (*opsmngr.ServiceVersion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServiceVersion")
	ret0, _ := ret[0].(*opsmngr.ServiceVersion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ServiceVersion indicates an expected call of ServiceVersion.
func (mr *MockServiceVersionDescriberMockRecorder) ServiceVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServiceVersion", reflect.TypeOf((*MockServiceVersionDescriber)(nil).ServiceVersion))
}
