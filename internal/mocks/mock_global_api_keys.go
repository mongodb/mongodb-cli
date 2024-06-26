// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mongodb/mongodb-atlas-cli/mongocli/v2/internal/store (interfaces: GlobalAPIKeyLister,GlobalAPIKeyDescriber,GlobalAPIKeyUpdater,GlobalAPIKeyCreator,GlobalAPIKeyDeleter)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	opsmngr "go.mongodb.org/ops-manager/opsmngr"
)

// MockGlobalAPIKeyLister is a mock of GlobalAPIKeyLister interface.
type MockGlobalAPIKeyLister struct {
	ctrl     *gomock.Controller
	recorder *MockGlobalAPIKeyListerMockRecorder
}

// MockGlobalAPIKeyListerMockRecorder is the mock recorder for MockGlobalAPIKeyLister.
type MockGlobalAPIKeyListerMockRecorder struct {
	mock *MockGlobalAPIKeyLister
}

// NewMockGlobalAPIKeyLister creates a new mock instance.
func NewMockGlobalAPIKeyLister(ctrl *gomock.Controller) *MockGlobalAPIKeyLister {
	mock := &MockGlobalAPIKeyLister{ctrl: ctrl}
	mock.recorder = &MockGlobalAPIKeyListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGlobalAPIKeyLister) EXPECT() *MockGlobalAPIKeyListerMockRecorder {
	return m.recorder
}

// GlobalAPIKeys mocks base method.
func (m *MockGlobalAPIKeyLister) GlobalAPIKeys(arg0 *opsmngr.ListOptions) ([]opsmngr.APIKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GlobalAPIKeys", arg0)
	ret0, _ := ret[0].([]opsmngr.APIKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GlobalAPIKeys indicates an expected call of GlobalAPIKeys.
func (mr *MockGlobalAPIKeyListerMockRecorder) GlobalAPIKeys(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GlobalAPIKeys", reflect.TypeOf((*MockGlobalAPIKeyLister)(nil).GlobalAPIKeys), arg0)
}

// MockGlobalAPIKeyDescriber is a mock of GlobalAPIKeyDescriber interface.
type MockGlobalAPIKeyDescriber struct {
	ctrl     *gomock.Controller
	recorder *MockGlobalAPIKeyDescriberMockRecorder
}

// MockGlobalAPIKeyDescriberMockRecorder is the mock recorder for MockGlobalAPIKeyDescriber.
type MockGlobalAPIKeyDescriberMockRecorder struct {
	mock *MockGlobalAPIKeyDescriber
}

// NewMockGlobalAPIKeyDescriber creates a new mock instance.
func NewMockGlobalAPIKeyDescriber(ctrl *gomock.Controller) *MockGlobalAPIKeyDescriber {
	mock := &MockGlobalAPIKeyDescriber{ctrl: ctrl}
	mock.recorder = &MockGlobalAPIKeyDescriberMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGlobalAPIKeyDescriber) EXPECT() *MockGlobalAPIKeyDescriberMockRecorder {
	return m.recorder
}

// GlobalAPIKey mocks base method.
func (m *MockGlobalAPIKeyDescriber) GlobalAPIKey(arg0 string) (*opsmngr.APIKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GlobalAPIKey", arg0)
	ret0, _ := ret[0].(*opsmngr.APIKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GlobalAPIKey indicates an expected call of GlobalAPIKey.
func (mr *MockGlobalAPIKeyDescriberMockRecorder) GlobalAPIKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GlobalAPIKey", reflect.TypeOf((*MockGlobalAPIKeyDescriber)(nil).GlobalAPIKey), arg0)
}

// MockGlobalAPIKeyUpdater is a mock of GlobalAPIKeyUpdater interface.
type MockGlobalAPIKeyUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockGlobalAPIKeyUpdaterMockRecorder
}

// MockGlobalAPIKeyUpdaterMockRecorder is the mock recorder for MockGlobalAPIKeyUpdater.
type MockGlobalAPIKeyUpdaterMockRecorder struct {
	mock *MockGlobalAPIKeyUpdater
}

// NewMockGlobalAPIKeyUpdater creates a new mock instance.
func NewMockGlobalAPIKeyUpdater(ctrl *gomock.Controller) *MockGlobalAPIKeyUpdater {
	mock := &MockGlobalAPIKeyUpdater{ctrl: ctrl}
	mock.recorder = &MockGlobalAPIKeyUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGlobalAPIKeyUpdater) EXPECT() *MockGlobalAPIKeyUpdaterMockRecorder {
	return m.recorder
}

// UpdateGlobalAPIKey mocks base method.
func (m *MockGlobalAPIKeyUpdater) UpdateGlobalAPIKey(arg0 string, arg1 *opsmngr.APIKeyInput) (*opsmngr.APIKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGlobalAPIKey", arg0, arg1)
	ret0, _ := ret[0].(*opsmngr.APIKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateGlobalAPIKey indicates an expected call of UpdateGlobalAPIKey.
func (mr *MockGlobalAPIKeyUpdaterMockRecorder) UpdateGlobalAPIKey(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGlobalAPIKey", reflect.TypeOf((*MockGlobalAPIKeyUpdater)(nil).UpdateGlobalAPIKey), arg0, arg1)
}

// MockGlobalAPIKeyCreator is a mock of GlobalAPIKeyCreator interface.
type MockGlobalAPIKeyCreator struct {
	ctrl     *gomock.Controller
	recorder *MockGlobalAPIKeyCreatorMockRecorder
}

// MockGlobalAPIKeyCreatorMockRecorder is the mock recorder for MockGlobalAPIKeyCreator.
type MockGlobalAPIKeyCreatorMockRecorder struct {
	mock *MockGlobalAPIKeyCreator
}

// NewMockGlobalAPIKeyCreator creates a new mock instance.
func NewMockGlobalAPIKeyCreator(ctrl *gomock.Controller) *MockGlobalAPIKeyCreator {
	mock := &MockGlobalAPIKeyCreator{ctrl: ctrl}
	mock.recorder = &MockGlobalAPIKeyCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGlobalAPIKeyCreator) EXPECT() *MockGlobalAPIKeyCreatorMockRecorder {
	return m.recorder
}

// CreateGlobalAPIKey mocks base method.
func (m *MockGlobalAPIKeyCreator) CreateGlobalAPIKey(arg0 *opsmngr.APIKeyInput) (*opsmngr.APIKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGlobalAPIKey", arg0)
	ret0, _ := ret[0].(*opsmngr.APIKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGlobalAPIKey indicates an expected call of CreateGlobalAPIKey.
func (mr *MockGlobalAPIKeyCreatorMockRecorder) CreateGlobalAPIKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGlobalAPIKey", reflect.TypeOf((*MockGlobalAPIKeyCreator)(nil).CreateGlobalAPIKey), arg0)
}

// MockGlobalAPIKeyDeleter is a mock of GlobalAPIKeyDeleter interface.
type MockGlobalAPIKeyDeleter struct {
	ctrl     *gomock.Controller
	recorder *MockGlobalAPIKeyDeleterMockRecorder
}

// MockGlobalAPIKeyDeleterMockRecorder is the mock recorder for MockGlobalAPIKeyDeleter.
type MockGlobalAPIKeyDeleterMockRecorder struct {
	mock *MockGlobalAPIKeyDeleter
}

// NewMockGlobalAPIKeyDeleter creates a new mock instance.
func NewMockGlobalAPIKeyDeleter(ctrl *gomock.Controller) *MockGlobalAPIKeyDeleter {
	mock := &MockGlobalAPIKeyDeleter{ctrl: ctrl}
	mock.recorder = &MockGlobalAPIKeyDeleterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGlobalAPIKeyDeleter) EXPECT() *MockGlobalAPIKeyDeleterMockRecorder {
	return m.recorder
}

// DeleteGlobalAPIKey mocks base method.
func (m *MockGlobalAPIKeyDeleter) DeleteGlobalAPIKey(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGlobalAPIKey", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGlobalAPIKey indicates an expected call of DeleteGlobalAPIKey.
func (mr *MockGlobalAPIKeyDeleterMockRecorder) DeleteGlobalAPIKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGlobalAPIKey", reflect.TypeOf((*MockGlobalAPIKeyDeleter)(nil).DeleteGlobalAPIKey), arg0)
}
