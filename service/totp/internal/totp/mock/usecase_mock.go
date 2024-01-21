// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/GCFactory/dbo-system/service/totp/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockUseCase is a mock of UseCase interface.
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase.
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance.
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// Disable mocks base method.
func (m *MockUseCase) Disable(ctx context.Context, request *models.TOTPRequest) (*models.TOTPBase, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Disable", ctx, request)
	ret0, _ := ret[0].(*models.TOTPBase)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Disable indicates an expected call of Disable.
func (mr *MockUseCaseMockRecorder) Disable(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disable", reflect.TypeOf((*MockUseCase)(nil).Disable), ctx, request)
}

// Enable mocks base method.
func (m *MockUseCase) Enable(ctx context.Context, request *models.TOTPRequest) (*models.TOTPBase, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Enable", ctx, request)
	ret0, _ := ret[0].(*models.TOTPBase)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Enable indicates an expected call of Enable.
func (mr *MockUseCaseMockRecorder) Enable(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enable", reflect.TypeOf((*MockUseCase)(nil).Enable), ctx, request)
}

// Enroll mocks base method.
func (m *MockUseCase) Enroll(ctx context.Context, totp models.TOTPConfig) (*models.TOTPEnroll, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Enroll", ctx, totp)
	ret0, _ := ret[0].(*models.TOTPEnroll)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Enroll indicates an expected call of Enroll.
func (mr *MockUseCaseMockRecorder) Enroll(ctx, totp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enroll", reflect.TypeOf((*MockUseCase)(nil).Enroll), ctx, totp)
}

// Validate mocks base method.
func (m *MockUseCase) Validate(ctx context.Context, request *models.TOTPRequest) (*models.TOTPBase, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", ctx, request)
	ret0, _ := ret[0].(*models.TOTPBase)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Validate indicates an expected call of Validate.
func (mr *MockUseCaseMockRecorder) Validate(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockUseCase)(nil).Validate), ctx, request)
}

// Verify mocks base method.
func (m *MockUseCase) Verify(ctx context.Context, url string) (*models.TOTPVerify, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", ctx, url)
	ret0, _ := ret[0].(*models.TOTPVerify)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Verify indicates an expected call of Verify.
func (mr *MockUseCaseMockRecorder) Verify(ctx, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockUseCase)(nil).Verify), ctx, url)
}
