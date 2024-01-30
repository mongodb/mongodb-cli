// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mongodb/mongodb-atlas-cli/internal/store/atlas (interfaces: OperatorOrgStore)

// Package atlas is a generated GoMock package.
package atlas

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	admin "go.mongodb.org/atlas-sdk/v20231115004/admin"
)

// MockOperatorOrgStore is a mock of OperatorOrgStore interface.
type MockOperatorOrgStore struct {
	ctrl     *gomock.Controller
	recorder *MockOperatorOrgStoreMockRecorder
}

// MockOperatorOrgStoreMockRecorder is the mock recorder for MockOperatorOrgStore.
type MockOperatorOrgStoreMockRecorder struct {
	mock *MockOperatorOrgStore
}

// NewMockOperatorOrgStore creates a new mock instance.
func NewMockOperatorOrgStore(ctrl *gomock.Controller) *MockOperatorOrgStore {
	mock := &MockOperatorOrgStore{ctrl: ctrl}
	mock.recorder = &MockOperatorOrgStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOperatorOrgStore) EXPECT() *MockOperatorOrgStoreMockRecorder {
	return m.recorder
}

// AssignProjectAPIKey mocks base method.
func (m *MockOperatorOrgStore) AssignProjectAPIKey(arg0, arg1 string, arg2 *admin.UpdateAtlasProjectApiKey) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignProjectAPIKey", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssignProjectAPIKey indicates an expected call of AssignProjectAPIKey.
func (mr *MockOperatorOrgStoreMockRecorder) AssignProjectAPIKey(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignProjectAPIKey", reflect.TypeOf((*MockOperatorOrgStore)(nil).AssignProjectAPIKey), arg0, arg1, arg2)
}

// CreateOrganizationAPIKey mocks base method.
func (m *MockOperatorOrgStore) CreateOrganizationAPIKey(arg0 string, arg1 *admin.CreateAtlasOrganizationApiKey) (*admin.ApiKeyUserDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrganizationAPIKey", arg0, arg1)
	ret0, _ := ret[0].(*admin.ApiKeyUserDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrganizationAPIKey indicates an expected call of CreateOrganizationAPIKey.
func (mr *MockOperatorOrgStoreMockRecorder) CreateOrganizationAPIKey(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrganizationAPIKey", reflect.TypeOf((*MockOperatorOrgStore)(nil).CreateOrganizationAPIKey), arg0, arg1)
}
