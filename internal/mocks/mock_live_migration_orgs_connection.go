// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mongodb/mongodb-cli/mongocli/v2/internal/store (interfaces: OrganizationsConnector,OrganizationsDescriber)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	opsmngr "go.mongodb.org/ops-manager/opsmngr"
)

// MockOrganizationsConnector is a mock of OrganizationsConnector interface.
type MockOrganizationsConnector struct {
	ctrl     *gomock.Controller
	recorder *MockOrganizationsConnectorMockRecorder
}

// MockOrganizationsConnectorMockRecorder is the mock recorder for MockOrganizationsConnector.
type MockOrganizationsConnectorMockRecorder struct {
	mock *MockOrganizationsConnector
}

// NewMockOrganizationsConnector creates a new mock instance.
func NewMockOrganizationsConnector(ctrl *gomock.Controller) *MockOrganizationsConnector {
	mock := &MockOrganizationsConnector{ctrl: ctrl}
	mock.recorder = &MockOrganizationsConnectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrganizationsConnector) EXPECT() *MockOrganizationsConnectorMockRecorder {
	return m.recorder
}

// ConnectOrganizations mocks base method.
func (m *MockOrganizationsConnector) ConnectOrganizations(arg0 string, arg1 *opsmngr.LinkToken) (*opsmngr.ConnectionStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnectOrganizations", arg0, arg1)
	ret0, _ := ret[0].(*opsmngr.ConnectionStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConnectOrganizations indicates an expected call of ConnectOrganizations.
func (mr *MockOrganizationsConnectorMockRecorder) ConnectOrganizations(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectOrganizations", reflect.TypeOf((*MockOrganizationsConnector)(nil).ConnectOrganizations), arg0, arg1)
}

// MockOrganizationsDescriber is a mock of OrganizationsDescriber interface.
type MockOrganizationsDescriber struct {
	ctrl     *gomock.Controller
	recorder *MockOrganizationsDescriberMockRecorder
}

// MockOrganizationsDescriberMockRecorder is the mock recorder for MockOrganizationsDescriber.
type MockOrganizationsDescriberMockRecorder struct {
	mock *MockOrganizationsDescriber
}

// NewMockOrganizationsDescriber creates a new mock instance.
func NewMockOrganizationsDescriber(ctrl *gomock.Controller) *MockOrganizationsDescriber {
	mock := &MockOrganizationsDescriber{ctrl: ctrl}
	mock.recorder = &MockOrganizationsDescriberMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrganizationsDescriber) EXPECT() *MockOrganizationsDescriberMockRecorder {
	return m.recorder
}

// OrganizationConnectionStatus mocks base method.
func (m *MockOrganizationsDescriber) OrganizationConnectionStatus(arg0 string) (*opsmngr.ConnectionStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OrganizationConnectionStatus", arg0)
	ret0, _ := ret[0].(*opsmngr.ConnectionStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OrganizationConnectionStatus indicates an expected call of OrganizationConnectionStatus.
func (mr *MockOrganizationsDescriberMockRecorder) OrganizationConnectionStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrganizationConnectionStatus", reflect.TypeOf((*MockOrganizationsDescriber)(nil).OrganizationConnectionStatus), arg0)
}
