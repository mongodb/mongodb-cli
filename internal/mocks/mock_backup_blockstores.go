// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mongodb/mongodb-cli/mongocli/v2/internal/store (interfaces: BlockstoresLister,BlockstoresDescriber,BlockstoresCreator,BlockstoresUpdater,BlockstoresDeleter)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	opsmngr "go.mongodb.org/ops-manager/opsmngr"
)

// MockBlockstoresLister is a mock of BlockstoresLister interface.
type MockBlockstoresLister struct {
	ctrl     *gomock.Controller
	recorder *MockBlockstoresListerMockRecorder
}

// MockBlockstoresListerMockRecorder is the mock recorder for MockBlockstoresLister.
type MockBlockstoresListerMockRecorder struct {
	mock *MockBlockstoresLister
}

// NewMockBlockstoresLister creates a new mock instance.
func NewMockBlockstoresLister(ctrl *gomock.Controller) *MockBlockstoresLister {
	mock := &MockBlockstoresLister{ctrl: ctrl}
	mock.recorder = &MockBlockstoresListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlockstoresLister) EXPECT() *MockBlockstoresListerMockRecorder {
	return m.recorder
}

// ListBlockstores mocks base method.
func (m *MockBlockstoresLister) ListBlockstores(arg0 *opsmngr.ListOptions) (*opsmngr.BackupStores, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBlockstores", arg0)
	ret0, _ := ret[0].(*opsmngr.BackupStores)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListBlockstores indicates an expected call of ListBlockstores.
func (mr *MockBlockstoresListerMockRecorder) ListBlockstores(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBlockstores", reflect.TypeOf((*MockBlockstoresLister)(nil).ListBlockstores), arg0)
}

// MockBlockstoresDescriber is a mock of BlockstoresDescriber interface.
type MockBlockstoresDescriber struct {
	ctrl     *gomock.Controller
	recorder *MockBlockstoresDescriberMockRecorder
}

// MockBlockstoresDescriberMockRecorder is the mock recorder for MockBlockstoresDescriber.
type MockBlockstoresDescriberMockRecorder struct {
	mock *MockBlockstoresDescriber
}

// NewMockBlockstoresDescriber creates a new mock instance.
func NewMockBlockstoresDescriber(ctrl *gomock.Controller) *MockBlockstoresDescriber {
	mock := &MockBlockstoresDescriber{ctrl: ctrl}
	mock.recorder = &MockBlockstoresDescriberMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlockstoresDescriber) EXPECT() *MockBlockstoresDescriberMockRecorder {
	return m.recorder
}

// DescribeBlockstore mocks base method.
func (m *MockBlockstoresDescriber) DescribeBlockstore(arg0 string) (*opsmngr.BackupStore, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeBlockstore", arg0)
	ret0, _ := ret[0].(*opsmngr.BackupStore)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeBlockstore indicates an expected call of DescribeBlockstore.
func (mr *MockBlockstoresDescriberMockRecorder) DescribeBlockstore(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeBlockstore", reflect.TypeOf((*MockBlockstoresDescriber)(nil).DescribeBlockstore), arg0)
}

// MockBlockstoresCreator is a mock of BlockstoresCreator interface.
type MockBlockstoresCreator struct {
	ctrl     *gomock.Controller
	recorder *MockBlockstoresCreatorMockRecorder
}

// MockBlockstoresCreatorMockRecorder is the mock recorder for MockBlockstoresCreator.
type MockBlockstoresCreatorMockRecorder struct {
	mock *MockBlockstoresCreator
}

// NewMockBlockstoresCreator creates a new mock instance.
func NewMockBlockstoresCreator(ctrl *gomock.Controller) *MockBlockstoresCreator {
	mock := &MockBlockstoresCreator{ctrl: ctrl}
	mock.recorder = &MockBlockstoresCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlockstoresCreator) EXPECT() *MockBlockstoresCreatorMockRecorder {
	return m.recorder
}

// CreateBlockstore mocks base method.
func (m *MockBlockstoresCreator) CreateBlockstore(arg0 *opsmngr.BackupStore) (*opsmngr.BackupStore, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBlockstore", arg0)
	ret0, _ := ret[0].(*opsmngr.BackupStore)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBlockstore indicates an expected call of CreateBlockstore.
func (mr *MockBlockstoresCreatorMockRecorder) CreateBlockstore(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBlockstore", reflect.TypeOf((*MockBlockstoresCreator)(nil).CreateBlockstore), arg0)
}

// MockBlockstoresUpdater is a mock of BlockstoresUpdater interface.
type MockBlockstoresUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockBlockstoresUpdaterMockRecorder
}

// MockBlockstoresUpdaterMockRecorder is the mock recorder for MockBlockstoresUpdater.
type MockBlockstoresUpdaterMockRecorder struct {
	mock *MockBlockstoresUpdater
}

// NewMockBlockstoresUpdater creates a new mock instance.
func NewMockBlockstoresUpdater(ctrl *gomock.Controller) *MockBlockstoresUpdater {
	mock := &MockBlockstoresUpdater{ctrl: ctrl}
	mock.recorder = &MockBlockstoresUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlockstoresUpdater) EXPECT() *MockBlockstoresUpdaterMockRecorder {
	return m.recorder
}

// UpdateBlockstore mocks base method.
func (m *MockBlockstoresUpdater) UpdateBlockstore(arg0 *opsmngr.BackupStore) (*opsmngr.BackupStore, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBlockstore", arg0)
	ret0, _ := ret[0].(*opsmngr.BackupStore)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBlockstore indicates an expected call of UpdateBlockstore.
func (mr *MockBlockstoresUpdaterMockRecorder) UpdateBlockstore(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBlockstore", reflect.TypeOf((*MockBlockstoresUpdater)(nil).UpdateBlockstore), arg0)
}

// MockBlockstoresDeleter is a mock of BlockstoresDeleter interface.
type MockBlockstoresDeleter struct {
	ctrl     *gomock.Controller
	recorder *MockBlockstoresDeleterMockRecorder
}

// MockBlockstoresDeleterMockRecorder is the mock recorder for MockBlockstoresDeleter.
type MockBlockstoresDeleterMockRecorder struct {
	mock *MockBlockstoresDeleter
}

// NewMockBlockstoresDeleter creates a new mock instance.
func NewMockBlockstoresDeleter(ctrl *gomock.Controller) *MockBlockstoresDeleter {
	mock := &MockBlockstoresDeleter{ctrl: ctrl}
	mock.recorder = &MockBlockstoresDeleterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlockstoresDeleter) EXPECT() *MockBlockstoresDeleterMockRecorder {
	return m.recorder
}

// DeleteBlockstore mocks base method.
func (m *MockBlockstoresDeleter) DeleteBlockstore(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBlockstore", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBlockstore indicates an expected call of DeleteBlockstore.
func (mr *MockBlockstoresDeleterMockRecorder) DeleteBlockstore(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBlockstore", reflect.TypeOf((*MockBlockstoresDeleter)(nil).DeleteBlockstore), arg0)
}
