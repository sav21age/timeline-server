// Code generated by MockGen. DO NOT EDIT.
// Source: email.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"
	domain "timeline/internal/domain"

	gomock "github.com/golang/mock/gomock"
)

// MockEmailInterface is a mock of EmailInterface interface.
type MockEmailInterface struct {
	ctrl     *gomock.Controller
	recorder *MockEmailInterfaceMockRecorder
}

// MockEmailInterfaceMockRecorder is the mock recorder for MockEmailInterface.
type MockEmailInterfaceMockRecorder struct {
	mock *MockEmailInterface
}

// NewMockEmailInterface creates a new mock instance.
func NewMockEmailInterface(ctrl *gomock.Controller) *MockEmailInterface {
	mock := &MockEmailInterface{ctrl: ctrl}
	mock.recorder = &MockEmailInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailInterface) EXPECT() *MockEmailInterfaceMockRecorder {
	return m.recorder
}

// AccountActivateEmail mocks base method.
func (m *MockEmailInterface) AccountActivateEmail(arg0 domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccountActivateEmail", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AccountActivateEmail indicates an expected call of AccountActivateEmail.
func (mr *MockEmailInterfaceMockRecorder) AccountActivateEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountActivateEmail", reflect.TypeOf((*MockEmailInterface)(nil).AccountActivateEmail), arg0)
}

// PasswordRecoveryEmail mocks base method.
func (m *MockEmailInterface) PasswordRecoveryEmail(arg0 domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PasswordRecoveryEmail", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// PasswordRecoveryEmail indicates an expected call of PasswordRecoveryEmail.
func (mr *MockEmailInterfaceMockRecorder) PasswordRecoveryEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PasswordRecoveryEmail", reflect.TypeOf((*MockEmailInterface)(nil).PasswordRecoveryEmail), arg0)
}