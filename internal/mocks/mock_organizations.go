// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mongodb/mongodb-cli/mongocli/v2/internal/store (interfaces: OrganizationLister,OrganizationCreator,OrganizationDeleter,OrganizationDescriber)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	opsmngr "go.mongodb.org/ops-manager/opsmngr"
)

// MockOrganizationLister is a mock of OrganizationLister interface.
type MockOrganizationLister struct {
	ctrl     *gomock.Controller
	recorder *MockOrganizationListerMockRecorder
}

// MockOrganizationListerMockRecorder is the mock recorder for MockOrganizationLister.
type MockOrganizationListerMockRecorder struct {
	mock *MockOrganizationLister
}

// NewMockOrganizationLister creates a new mock instance.
func NewMockOrganizationLister(ctrl *gomock.Controller) *MockOrganizationLister {
	mock := &MockOrganizationLister{ctrl: ctrl}
	mock.recorder = &MockOrganizationListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrganizationLister) EXPECT() *MockOrganizationListerMockRecorder {
	return m.recorder
}

// Organizations mocks base method.
func (m *MockOrganizationLister) Organizations(arg0 *opsmngr.OrganizationsListOptions) (*opsmngr.Organizations, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Organizations", arg0)
	ret0, _ := ret[0].(*opsmngr.Organizations)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Organizations indicates an expected call of Organizations.
func (mr *MockOrganizationListerMockRecorder) Organizations(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Organizations", reflect.TypeOf((*MockOrganizationLister)(nil).Organizations), arg0)
}

// MockOrganizationCreator is a mock of OrganizationCreator interface.
type MockOrganizationCreator struct {
	ctrl     *gomock.Controller
	recorder *MockOrganizationCreatorMockRecorder
}

// MockOrganizationCreatorMockRecorder is the mock recorder for MockOrganizationCreator.
type MockOrganizationCreatorMockRecorder struct {
	mock *MockOrganizationCreator
}

// NewMockOrganizationCreator creates a new mock instance.
func NewMockOrganizationCreator(ctrl *gomock.Controller) *MockOrganizationCreator {
	mock := &MockOrganizationCreator{ctrl: ctrl}
	mock.recorder = &MockOrganizationCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrganizationCreator) EXPECT() *MockOrganizationCreatorMockRecorder {
	return m.recorder
}

// CreateOrganization mocks base method.
func (m *MockOrganizationCreator) CreateOrganization(arg0 string) (*opsmngr.Organization, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrganization", arg0)
	ret0, _ := ret[0].(*opsmngr.Organization)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrganization indicates an expected call of CreateOrganization.
func (mr *MockOrganizationCreatorMockRecorder) CreateOrganization(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrganization", reflect.TypeOf((*MockOrganizationCreator)(nil).CreateOrganization), arg0)
}

// MockOrganizationDeleter is a mock of OrganizationDeleter interface.
type MockOrganizationDeleter struct {
	ctrl     *gomock.Controller
	recorder *MockOrganizationDeleterMockRecorder
}

// MockOrganizationDeleterMockRecorder is the mock recorder for MockOrganizationDeleter.
type MockOrganizationDeleterMockRecorder struct {
	mock *MockOrganizationDeleter
}

// NewMockOrganizationDeleter creates a new mock instance.
func NewMockOrganizationDeleter(ctrl *gomock.Controller) *MockOrganizationDeleter {
	mock := &MockOrganizationDeleter{ctrl: ctrl}
	mock.recorder = &MockOrganizationDeleterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrganizationDeleter) EXPECT() *MockOrganizationDeleterMockRecorder {
	return m.recorder
}

// DeleteOrganization mocks base method.
func (m *MockOrganizationDeleter) DeleteOrganization(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrganization", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOrganization indicates an expected call of DeleteOrganization.
func (mr *MockOrganizationDeleterMockRecorder) DeleteOrganization(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrganization", reflect.TypeOf((*MockOrganizationDeleter)(nil).DeleteOrganization), arg0)
}

// MockOrganizationDescriber is a mock of OrganizationDescriber interface.
type MockOrganizationDescriber struct {
	ctrl     *gomock.Controller
	recorder *MockOrganizationDescriberMockRecorder
}

// MockOrganizationDescriberMockRecorder is the mock recorder for MockOrganizationDescriber.
type MockOrganizationDescriberMockRecorder struct {
	mock *MockOrganizationDescriber
}

// NewMockOrganizationDescriber creates a new mock instance.
func NewMockOrganizationDescriber(ctrl *gomock.Controller) *MockOrganizationDescriber {
	mock := &MockOrganizationDescriber{ctrl: ctrl}
	mock.recorder = &MockOrganizationDescriberMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrganizationDescriber) EXPECT() *MockOrganizationDescriberMockRecorder {
	return m.recorder
}

// Organization mocks base method.
func (m *MockOrganizationDescriber) Organization(arg0 string) (*opsmngr.Organization, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Organization", arg0)
	ret0, _ := ret[0].(*opsmngr.Organization)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Organization indicates an expected call of Organization.
func (mr *MockOrganizationDescriberMockRecorder) Organization(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Organization", reflect.TypeOf((*MockOrganizationDescriber)(nil).Organization), arg0)
}
