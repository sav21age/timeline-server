// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"
	domain "timeline/internal/domain"

	gomock "github.com/golang/mock/gomock"
)

// MockUserInterface is a mock of UserInterface interface.
type MockUserInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUserInterfaceMockRecorder
}

// MockUserInterfaceMockRecorder is the mock recorder for MockUserInterface.
type MockUserInterfaceMockRecorder struct {
	mock *MockUserInterface
}

// NewMockUserInterface creates a new mock instance.
func NewMockUserInterface(ctrl *gomock.Controller) *MockUserInterface {
	mock := &MockUserInterface{ctrl: ctrl}
	mock.recorder = &MockUserInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserInterface) EXPECT() *MockUserInterfaceMockRecorder {
	return m.recorder
}

// AccountActivateByCode mocks base method.
func (m *MockUserInterface) AccountActivateByCode(ctx context.Context, code string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccountActivateByCode", ctx, code)
	ret0, _ := ret[0].(error)
	return ret0
}

// AccountActivateByCode indicates an expected call of AccountActivateByCode.
func (mr *MockUserInterfaceMockRecorder) AccountActivateByCode(ctx, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountActivateByCode", reflect.TypeOf((*MockUserInterface)(nil).AccountActivateByCode), ctx, code)
}

// AccountActivateResendCode mocks base method.
func (m *MockUserInterface) AccountActivateResendCode(ctx context.Context, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccountActivateResendCode", ctx, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// AccountActivateResendCode indicates an expected call of AccountActivateResendCode.
func (mr *MockUserInterfaceMockRecorder) AccountActivateResendCode(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountActivateResendCode", reflect.TypeOf((*MockUserInterface)(nil).AccountActivateResendCode), ctx, email)
}

// CreateSession mocks base method.
func (m *MockUserInterface) CreateSession(ctx context.Context, user domain.User, rememberMe bool) (domain.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, user, rememberMe)
	ret0, _ := ret[0].(domain.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockUserInterfaceMockRecorder) CreateSession(ctx, user, rememberMe interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockUserInterface)(nil).CreateSession), ctx, user, rememberMe)
}

// RecoveryPasswordByEmail mocks base method.
func (m *MockUserInterface) RecoveryPasswordByEmail(ctx context.Context, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecoveryPasswordByEmail", ctx, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecoveryPasswordByEmail indicates an expected call of RecoveryPasswordByEmail.
func (mr *MockUserInterfaceMockRecorder) RecoveryPasswordByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecoveryPasswordByEmail", reflect.TypeOf((*MockUserInterface)(nil).RecoveryPasswordByEmail), ctx, email)
}

// RecoveryPasswordByUsername mocks base method.
func (m *MockUserInterface) RecoveryPasswordByUsername(ctx context.Context, username string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecoveryPasswordByUsername", ctx, username)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecoveryPasswordByUsername indicates an expected call of RecoveryPasswordByUsername.
func (mr *MockUserInterfaceMockRecorder) RecoveryPasswordByUsername(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecoveryPasswordByUsername", reflect.TypeOf((*MockUserInterface)(nil).RecoveryPasswordByUsername), ctx, username)
}

// RefreshSession mocks base method.
func (m *MockUserInterface) RefreshSession(ctx context.Context, refreshToken string) (domain.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshSession", ctx, refreshToken)
	ret0, _ := ret[0].(domain.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshSession indicates an expected call of RefreshSession.
func (mr *MockUserInterfaceMockRecorder) RefreshSession(ctx, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshSession", reflect.TypeOf((*MockUserInterface)(nil).RefreshSession), ctx, refreshToken)
}

// SetNewPassword mocks base method.
func (m *MockUserInterface) SetNewPassword(ctx context.Context, password, code string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetNewPassword", ctx, password, code)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetNewPassword indicates an expected call of SetNewPassword.
func (mr *MockUserInterfaceMockRecorder) SetNewPassword(ctx, password, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNewPassword", reflect.TypeOf((*MockUserInterface)(nil).SetNewPassword), ctx, password, code)
}

// SetRecoveryCode mocks base method.
func (m *MockUserInterface) SetRecoveryCode(ctx context.Context, user domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetRecoveryCode", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRecoveryCode indicates an expected call of SetRecoveryCode.
func (mr *MockUserInterfaceMockRecorder) SetRecoveryCode(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRecoveryCode", reflect.TypeOf((*MockUserInterface)(nil).SetRecoveryCode), ctx, user)
}

// SignInByEmail mocks base method.
func (m *MockUserInterface) SignInByEmail(ctx context.Context, input domain.UserSignInEmailInput) (domain.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignInByEmail", ctx, input)
	ret0, _ := ret[0].(domain.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignInByEmail indicates an expected call of SignInByEmail.
func (mr *MockUserInterfaceMockRecorder) SignInByEmail(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignInByEmail", reflect.TypeOf((*MockUserInterface)(nil).SignInByEmail), ctx, input)
}

// SignInByUsername mocks base method.
func (m *MockUserInterface) SignInByUsername(ctx context.Context, input domain.UserSignInUsernameInput) (domain.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignInByUsername", ctx, input)
	ret0, _ := ret[0].(domain.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignInByUsername indicates an expected call of SignInByUsername.
func (mr *MockUserInterfaceMockRecorder) SignInByUsername(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignInByUsername", reflect.TypeOf((*MockUserInterface)(nil).SignInByUsername), ctx, input)
}

// SignOut mocks base method.
func (m *MockUserInterface) SignOut(ctx context.Context, refreshToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignOut", ctx, refreshToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignOut indicates an expected call of SignOut.
func (mr *MockUserInterfaceMockRecorder) SignOut(ctx, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignOut", reflect.TypeOf((*MockUserInterface)(nil).SignOut), ctx, refreshToken)
}

// SignUp mocks base method.
func (m *MockUserInterface) SignUp(ctx context.Context, input domain.UserSignUpInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignUp indicates an expected call of SignUp.
func (mr *MockUserInterfaceMockRecorder) SignUp(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockUserInterface)(nil).SignUp), ctx, input)
}

// ValidateAccessToken mocks base method.
func (m *MockUserInterface) ValidateAccessToken(accessToken string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAccessToken", accessToken)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateAccessToken indicates an expected call of ValidateAccessToken.
func (mr *MockUserInterfaceMockRecorder) ValidateAccessToken(accessToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAccessToken", reflect.TypeOf((*MockUserInterface)(nil).ValidateAccessToken), accessToken)
}

// VerifyRecoveryCode mocks base method.
func (m *MockUserInterface) VerifyRecoveryCode(ctx context.Context, code string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyRecoveryCode", ctx, code)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyRecoveryCode indicates an expected call of VerifyRecoveryCode.
func (mr *MockUserInterfaceMockRecorder) VerifyRecoveryCode(ctx, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyRecoveryCode", reflect.TypeOf((*MockUserInterface)(nil).VerifyRecoveryCode), ctx, code)
}