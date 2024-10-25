// Code generated by MockGen. DO NOT EDIT.
// Source: ./user/repository.go

// Package user is a generated GoMock package.
package user

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// FindByUsername mocks base method.
func (m *MockRepository) FindByUsername(id string) (User, error, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUsername", id)
	ret0, _ := ret[0].(User)
	ret1, _ := ret[1].(error)
	ret2, _ := ret[2].(bool)
	return ret0, ret1, ret2
}

// FindByUsername indicates an expected call of FindByUsername.
func (mr *MockRepositoryMockRecorder) FindByUsername(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUsername", reflect.TypeOf((*MockRepository)(nil).FindByUsername), id)
}

// doSomething mocks base method.
func (m *MockRepository) doSomething() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "doSomething")
	ret0, _ := ret[0].(string)
	return ret0
}

// doSomething indicates an expected call of doSomething.
func (mr *MockRepositoryMockRecorder) doSomething() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "doSomething", reflect.TypeOf((*MockRepository)(nil).doSomething))
}