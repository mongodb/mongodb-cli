// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mongodb/mongodb-cli/mongocli/v2/internal/store (interfaces: ProjectServerTypeGetter,ProjectServerTypeUpdater,OrganizationServerTypeGetter,OrganizationServerTypeUpdater,ProjectHostAssignmentLister,OrganizationHostAssignmentLister,SnapshotGenerator,ServerUsageReportDownloader)

// Package mocks is a generated GoMock package.
package mocks

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	opsmngr "go.mongodb.org/ops-manager/opsmngr"
)

// MockProjectServerTypeGetter is a mock of ProjectServerTypeGetter interface.
type MockProjectServerTypeGetter struct {
	ctrl     *gomock.Controller
	recorder *MockProjectServerTypeGetterMockRecorder
}

// MockProjectServerTypeGetterMockRecorder is the mock recorder for MockProjectServerTypeGetter.
type MockProjectServerTypeGetterMockRecorder struct {
	mock *MockProjectServerTypeGetter
}

// NewMockProjectServerTypeGetter creates a new mock instance.
func NewMockProjectServerTypeGetter(ctrl *gomock.Controller) *MockProjectServerTypeGetter {
	mock := &MockProjectServerTypeGetter{ctrl: ctrl}
	mock.recorder = &MockProjectServerTypeGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectServerTypeGetter) EXPECT() *MockProjectServerTypeGetterMockRecorder {
	return m.recorder
}

// ProjectServerType mocks base method.
func (m *MockProjectServerTypeGetter) ProjectServerType(arg0 string) (*opsmngr.ServerType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProjectServerType", arg0)
	ret0, _ := ret[0].(*opsmngr.ServerType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProjectServerType indicates an expected call of ProjectServerType.
func (mr *MockProjectServerTypeGetterMockRecorder) ProjectServerType(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProjectServerType", reflect.TypeOf((*MockProjectServerTypeGetter)(nil).ProjectServerType), arg0)
}

// MockProjectServerTypeUpdater is a mock of ProjectServerTypeUpdater interface.
type MockProjectServerTypeUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockProjectServerTypeUpdaterMockRecorder
}

// MockProjectServerTypeUpdaterMockRecorder is the mock recorder for MockProjectServerTypeUpdater.
type MockProjectServerTypeUpdaterMockRecorder struct {
	mock *MockProjectServerTypeUpdater
}

// NewMockProjectServerTypeUpdater creates a new mock instance.
func NewMockProjectServerTypeUpdater(ctrl *gomock.Controller) *MockProjectServerTypeUpdater {
	mock := &MockProjectServerTypeUpdater{ctrl: ctrl}
	mock.recorder = &MockProjectServerTypeUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectServerTypeUpdater) EXPECT() *MockProjectServerTypeUpdaterMockRecorder {
	return m.recorder
}

// UpdateProjectServerType mocks base method.
func (m *MockProjectServerTypeUpdater) UpdateProjectServerType(arg0 string, arg1 *opsmngr.ServerTypeRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProjectServerType", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProjectServerType indicates an expected call of UpdateProjectServerType.
func (mr *MockProjectServerTypeUpdaterMockRecorder) UpdateProjectServerType(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProjectServerType", reflect.TypeOf((*MockProjectServerTypeUpdater)(nil).UpdateProjectServerType), arg0, arg1)
}

// MockOrganizationServerTypeGetter is a mock of OrganizationServerTypeGetter interface.
type MockOrganizationServerTypeGetter struct {
	ctrl     *gomock.Controller
	recorder *MockOrganizationServerTypeGetterMockRecorder
}

// MockOrganizationServerTypeGetterMockRecorder is the mock recorder for MockOrganizationServerTypeGetter.
type MockOrganizationServerTypeGetterMockRecorder struct {
	mock *MockOrganizationServerTypeGetter
}

// NewMockOrganizationServerTypeGetter creates a new mock instance.
func NewMockOrganizationServerTypeGetter(ctrl *gomock.Controller) *MockOrganizationServerTypeGetter {
	mock := &MockOrganizationServerTypeGetter{ctrl: ctrl}
	mock.recorder = &MockOrganizationServerTypeGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrganizationServerTypeGetter) EXPECT() *MockOrganizationServerTypeGetterMockRecorder {
	return m.recorder
}

// OrganizationServerType mocks base method.
func (m *MockOrganizationServerTypeGetter) OrganizationServerType(arg0 string) (*opsmngr.ServerType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OrganizationServerType", arg0)
	ret0, _ := ret[0].(*opsmngr.ServerType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OrganizationServerType indicates an expected call of OrganizationServerType.
func (mr *MockOrganizationServerTypeGetterMockRecorder) OrganizationServerType(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrganizationServerType", reflect.TypeOf((*MockOrganizationServerTypeGetter)(nil).OrganizationServerType), arg0)
}

// MockOrganizationServerTypeUpdater is a mock of OrganizationServerTypeUpdater interface.
type MockOrganizationServerTypeUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockOrganizationServerTypeUpdaterMockRecorder
}

// MockOrganizationServerTypeUpdaterMockRecorder is the mock recorder for MockOrganizationServerTypeUpdater.
type MockOrganizationServerTypeUpdaterMockRecorder struct {
	mock *MockOrganizationServerTypeUpdater
}

// NewMockOrganizationServerTypeUpdater creates a new mock instance.
func NewMockOrganizationServerTypeUpdater(ctrl *gomock.Controller) *MockOrganizationServerTypeUpdater {
	mock := &MockOrganizationServerTypeUpdater{ctrl: ctrl}
	mock.recorder = &MockOrganizationServerTypeUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrganizationServerTypeUpdater) EXPECT() *MockOrganizationServerTypeUpdaterMockRecorder {
	return m.recorder
}

// UpdateOrganizationServerType mocks base method.
func (m *MockOrganizationServerTypeUpdater) UpdateOrganizationServerType(arg0 string, arg1 *opsmngr.ServerTypeRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrganizationServerType", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrganizationServerType indicates an expected call of UpdateOrganizationServerType.
func (mr *MockOrganizationServerTypeUpdaterMockRecorder) UpdateOrganizationServerType(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrganizationServerType", reflect.TypeOf((*MockOrganizationServerTypeUpdater)(nil).UpdateOrganizationServerType), arg0, arg1)
}

// MockProjectHostAssignmentLister is a mock of ProjectHostAssignmentLister interface.
type MockProjectHostAssignmentLister struct {
	ctrl     *gomock.Controller
	recorder *MockProjectHostAssignmentListerMockRecorder
}

// MockProjectHostAssignmentListerMockRecorder is the mock recorder for MockProjectHostAssignmentLister.
type MockProjectHostAssignmentListerMockRecorder struct {
	mock *MockProjectHostAssignmentLister
}

// NewMockProjectHostAssignmentLister creates a new mock instance.
func NewMockProjectHostAssignmentLister(ctrl *gomock.Controller) *MockProjectHostAssignmentLister {
	mock := &MockProjectHostAssignmentLister{ctrl: ctrl}
	mock.recorder = &MockProjectHostAssignmentListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectHostAssignmentLister) EXPECT() *MockProjectHostAssignmentListerMockRecorder {
	return m.recorder
}

// ProjectHostAssignments mocks base method.
func (m *MockProjectHostAssignmentLister) ProjectHostAssignments(arg0 string, arg1 *opsmngr.ServerTypeOptions) (*opsmngr.HostAssignments, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProjectHostAssignments", arg0, arg1)
	ret0, _ := ret[0].(*opsmngr.HostAssignments)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProjectHostAssignments indicates an expected call of ProjectHostAssignments.
func (mr *MockProjectHostAssignmentListerMockRecorder) ProjectHostAssignments(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProjectHostAssignments", reflect.TypeOf((*MockProjectHostAssignmentLister)(nil).ProjectHostAssignments), arg0, arg1)
}

// MockOrganizationHostAssignmentLister is a mock of OrganizationHostAssignmentLister interface.
type MockOrganizationHostAssignmentLister struct {
	ctrl     *gomock.Controller
	recorder *MockOrganizationHostAssignmentListerMockRecorder
}

// MockOrganizationHostAssignmentListerMockRecorder is the mock recorder for MockOrganizationHostAssignmentLister.
type MockOrganizationHostAssignmentListerMockRecorder struct {
	mock *MockOrganizationHostAssignmentLister
}

// NewMockOrganizationHostAssignmentLister creates a new mock instance.
func NewMockOrganizationHostAssignmentLister(ctrl *gomock.Controller) *MockOrganizationHostAssignmentLister {
	mock := &MockOrganizationHostAssignmentLister{ctrl: ctrl}
	mock.recorder = &MockOrganizationHostAssignmentListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrganizationHostAssignmentLister) EXPECT() *MockOrganizationHostAssignmentListerMockRecorder {
	return m.recorder
}

// OrganizationHostAssignments mocks base method.
func (m *MockOrganizationHostAssignmentLister) OrganizationHostAssignments(arg0 string, arg1 *opsmngr.ServerTypeOptions) (*opsmngr.HostAssignments, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OrganizationHostAssignments", arg0, arg1)
	ret0, _ := ret[0].(*opsmngr.HostAssignments)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OrganizationHostAssignments indicates an expected call of OrganizationHostAssignments.
func (mr *MockOrganizationHostAssignmentListerMockRecorder) OrganizationHostAssignments(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrganizationHostAssignments", reflect.TypeOf((*MockOrganizationHostAssignmentLister)(nil).OrganizationHostAssignments), arg0, arg1)
}

// MockSnapshotGenerator is a mock of SnapshotGenerator interface.
type MockSnapshotGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockSnapshotGeneratorMockRecorder
}

// MockSnapshotGeneratorMockRecorder is the mock recorder for MockSnapshotGenerator.
type MockSnapshotGeneratorMockRecorder struct {
	mock *MockSnapshotGenerator
}

// NewMockSnapshotGenerator creates a new mock instance.
func NewMockSnapshotGenerator(ctrl *gomock.Controller) *MockSnapshotGenerator {
	mock := &MockSnapshotGenerator{ctrl: ctrl}
	mock.recorder = &MockSnapshotGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSnapshotGenerator) EXPECT() *MockSnapshotGeneratorMockRecorder {
	return m.recorder
}

// GenerateSnapshot mocks base method.
func (m *MockSnapshotGenerator) GenerateSnapshot() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateSnapshot")
	ret0, _ := ret[0].(error)
	return ret0
}

// GenerateSnapshot indicates an expected call of GenerateSnapshot.
func (mr *MockSnapshotGeneratorMockRecorder) GenerateSnapshot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateSnapshot", reflect.TypeOf((*MockSnapshotGenerator)(nil).GenerateSnapshot))
}

// MockServerUsageReportDownloader is a mock of ServerUsageReportDownloader interface.
type MockServerUsageReportDownloader struct {
	ctrl     *gomock.Controller
	recorder *MockServerUsageReportDownloaderMockRecorder
}

// MockServerUsageReportDownloaderMockRecorder is the mock recorder for MockServerUsageReportDownloader.
type MockServerUsageReportDownloaderMockRecorder struct {
	mock *MockServerUsageReportDownloader
}

// NewMockServerUsageReportDownloader creates a new mock instance.
func NewMockServerUsageReportDownloader(ctrl *gomock.Controller) *MockServerUsageReportDownloader {
	mock := &MockServerUsageReportDownloader{ctrl: ctrl}
	mock.recorder = &MockServerUsageReportDownloaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServerUsageReportDownloader) EXPECT() *MockServerUsageReportDownloaderMockRecorder {
	return m.recorder
}

// DownloadServerUsageReport mocks base method.
func (m *MockServerUsageReportDownloader) DownloadServerUsageReport(arg0 *opsmngr.ServerTypeOptions, arg1 io.Writer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadServerUsageReport", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DownloadServerUsageReport indicates an expected call of DownloadServerUsageReport.
func (mr *MockServerUsageReportDownloaderMockRecorder) DownloadServerUsageReport(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadServerUsageReport", reflect.TypeOf((*MockServerUsageReportDownloader)(nil).DownloadServerUsageReport), arg0, arg1)
}
