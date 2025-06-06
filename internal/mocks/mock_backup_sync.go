// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mongodb/mongodb-cli/mongocli/v2/internal/store (interfaces: SyncsLister,SyncsDescriber,SyncsCreator,SyncsUpdater,SyncsDeleter)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	opsmngr "go.mongodb.org/ops-manager/opsmngr"
)

// MockSyncsLister is a mock of SyncsLister interface.
type MockSyncsLister struct {
	ctrl     *gomock.Controller
	recorder *MockSyncsListerMockRecorder
}

// MockSyncsListerMockRecorder is the mock recorder for MockSyncsLister.
type MockSyncsListerMockRecorder struct {
	mock *MockSyncsLister
}

// NewMockSyncsLister creates a new mock instance.
func NewMockSyncsLister(ctrl *gomock.Controller) *MockSyncsLister {
	mock := &MockSyncsLister{ctrl: ctrl}
	mock.recorder = &MockSyncsListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSyncsLister) EXPECT() *MockSyncsListerMockRecorder {
	return m.recorder
}

// ListSyncs mocks base method.
func (m *MockSyncsLister) ListSyncs(arg0 *opsmngr.ListOptions) (*opsmngr.BackupStores, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSyncs", arg0)
	ret0, _ := ret[0].(*opsmngr.BackupStores)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSyncs indicates an expected call of ListSyncs.
func (mr *MockSyncsListerMockRecorder) ListSyncs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSyncs", reflect.TypeOf((*MockSyncsLister)(nil).ListSyncs), arg0)
}

// MockSyncsDescriber is a mock of SyncsDescriber interface.
type MockSyncsDescriber struct {
	ctrl     *gomock.Controller
	recorder *MockSyncsDescriberMockRecorder
}

// MockSyncsDescriberMockRecorder is the mock recorder for MockSyncsDescriber.
type MockSyncsDescriberMockRecorder struct {
	mock *MockSyncsDescriber
}

// NewMockSyncsDescriber creates a new mock instance.
func NewMockSyncsDescriber(ctrl *gomock.Controller) *MockSyncsDescriber {
	mock := &MockSyncsDescriber{ctrl: ctrl}
	mock.recorder = &MockSyncsDescriberMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSyncsDescriber) EXPECT() *MockSyncsDescriberMockRecorder {
	return m.recorder
}

// GetSync mocks base method.
func (m *MockSyncsDescriber) GetSync(arg0 string) (*opsmngr.BackupStore, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSync", arg0)
	ret0, _ := ret[0].(*opsmngr.BackupStore)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSync indicates an expected call of GetSync.
func (mr *MockSyncsDescriberMockRecorder) GetSync(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSync", reflect.TypeOf((*MockSyncsDescriber)(nil).GetSync), arg0)
}

// MockSyncsCreator is a mock of SyncsCreator interface.
type MockSyncsCreator struct {
	ctrl     *gomock.Controller
	recorder *MockSyncsCreatorMockRecorder
}

// MockSyncsCreatorMockRecorder is the mock recorder for MockSyncsCreator.
type MockSyncsCreatorMockRecorder struct {
	mock *MockSyncsCreator
}

// NewMockSyncsCreator creates a new mock instance.
func NewMockSyncsCreator(ctrl *gomock.Controller) *MockSyncsCreator {
	mock := &MockSyncsCreator{ctrl: ctrl}
	mock.recorder = &MockSyncsCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSyncsCreator) EXPECT() *MockSyncsCreatorMockRecorder {
	return m.recorder
}

// CreateSync mocks base method.
func (m *MockSyncsCreator) CreateSync(arg0 *opsmngr.BackupStore) (*opsmngr.BackupStore, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSync", arg0)
	ret0, _ := ret[0].(*opsmngr.BackupStore)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSync indicates an expected call of CreateSync.
func (mr *MockSyncsCreatorMockRecorder) CreateSync(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSync", reflect.TypeOf((*MockSyncsCreator)(nil).CreateSync), arg0)
}

// MockSyncsUpdater is a mock of SyncsUpdater interface.
type MockSyncsUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockSyncsUpdaterMockRecorder
}

// MockSyncsUpdaterMockRecorder is the mock recorder for MockSyncsUpdater.
type MockSyncsUpdaterMockRecorder struct {
	mock *MockSyncsUpdater
}

// NewMockSyncsUpdater creates a new mock instance.
func NewMockSyncsUpdater(ctrl *gomock.Controller) *MockSyncsUpdater {
	mock := &MockSyncsUpdater{ctrl: ctrl}
	mock.recorder = &MockSyncsUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSyncsUpdater) EXPECT() *MockSyncsUpdaterMockRecorder {
	return m.recorder
}

// UpdateSync mocks base method.
func (m *MockSyncsUpdater) UpdateSync(arg0 string, arg1 *opsmngr.BackupStore) (*opsmngr.BackupStore, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSync", arg0, arg1)
	ret0, _ := ret[0].(*opsmngr.BackupStore)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSync indicates an expected call of UpdateSync.
func (mr *MockSyncsUpdaterMockRecorder) UpdateSync(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSync", reflect.TypeOf((*MockSyncsUpdater)(nil).UpdateSync), arg0, arg1)
}

// MockSyncsDeleter is a mock of SyncsDeleter interface.
type MockSyncsDeleter struct {
	ctrl     *gomock.Controller
	recorder *MockSyncsDeleterMockRecorder
}

// MockSyncsDeleterMockRecorder is the mock recorder for MockSyncsDeleter.
type MockSyncsDeleterMockRecorder struct {
	mock *MockSyncsDeleter
}

// NewMockSyncsDeleter creates a new mock instance.
func NewMockSyncsDeleter(ctrl *gomock.Controller) *MockSyncsDeleter {
	mock := &MockSyncsDeleter{ctrl: ctrl}
	mock.recorder = &MockSyncsDeleterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSyncsDeleter) EXPECT() *MockSyncsDeleterMockRecorder {
	return m.recorder
}

// DeleteSync mocks base method.
func (m *MockSyncsDeleter) DeleteSync(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSync", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSync indicates an expected call of DeleteSync.
func (mr *MockSyncsDeleterMockRecorder) DeleteSync(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSync", reflect.TypeOf((*MockSyncsDeleter)(nil).DeleteSync), arg0)
}
