// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cetteup/joinme.click-launcher/pkg/game_launcher (interfaces: FileRepository)
//
// Generated by this command:
//
//	mockgen -destination=mock_test.go -package=internal -write_package_comment=false github.com/cetteup/joinme.click-launcher/pkg/game_launcher FileRepository
package internal

import (
	fs "io/fs"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockFileRepository is a mock of FileRepository interface.
type MockFileRepository struct {
	ctrl     *gomock.Controller
	recorder *MockFileRepositoryMockRecorder
}

// MockFileRepositoryMockRecorder is the mock recorder for MockFileRepository.
type MockFileRepositoryMockRecorder struct {
	mock *MockFileRepository
}

// NewMockFileRepository creates a new mock instance.
func NewMockFileRepository(ctrl *gomock.Controller) *MockFileRepository {
	mock := &MockFileRepository{ctrl: ctrl}
	mock.recorder = &MockFileRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileRepository) EXPECT() *MockFileRepositoryMockRecorder {
	return m.recorder
}

// DirExists mocks base method.
func (m *MockFileRepository) DirExists(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DirExists", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DirExists indicates an expected call of DirExists.
func (mr *MockFileRepositoryMockRecorder) DirExists(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DirExists", reflect.TypeOf((*MockFileRepository)(nil).DirExists), arg0)
}

// FileExists mocks base method.
func (m *MockFileRepository) FileExists(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileExists", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FileExists indicates an expected call of FileExists.
func (mr *MockFileRepositoryMockRecorder) FileExists(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileExists", reflect.TypeOf((*MockFileRepository)(nil).FileExists), arg0)
}

// Glob mocks base method.
func (m *MockFileRepository) Glob(arg0 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Glob", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Glob indicates an expected call of Glob.
func (mr *MockFileRepositoryMockRecorder) Glob(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Glob", reflect.TypeOf((*MockFileRepository)(nil).Glob), arg0)
}

// ReadDir mocks base method.
func (m *MockFileRepository) ReadDir(arg0 string) ([]fs.DirEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadDir", arg0)
	ret0, _ := ret[0].([]fs.DirEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadDir indicates an expected call of ReadDir.
func (mr *MockFileRepositoryMockRecorder) ReadDir(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadDir", reflect.TypeOf((*MockFileRepository)(nil).ReadDir), arg0)
}

// ReadFile mocks base method.
func (m *MockFileRepository) ReadFile(arg0 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadFile", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadFile indicates an expected call of ReadFile.
func (mr *MockFileRepositoryMockRecorder) ReadFile(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFile", reflect.TypeOf((*MockFileRepository)(nil).ReadFile), arg0)
}

// RemoveAll mocks base method.
func (m *MockFileRepository) RemoveAll(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAll", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAll indicates an expected call of RemoveAll.
func (mr *MockFileRepositoryMockRecorder) RemoveAll(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAll", reflect.TypeOf((*MockFileRepository)(nil).RemoveAll), arg0)
}

// WriteFile mocks base method.
func (m *MockFileRepository) WriteFile(arg0 string, arg1 []byte, arg2 fs.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteFile", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteFile indicates an expected call of WriteFile.
func (mr *MockFileRepositoryMockRecorder) WriteFile(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteFile", reflect.TypeOf((*MockFileRepository)(nil).WriteFile), arg0, arg1, arg2)
}
