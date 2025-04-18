// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	models "github.com/GCFactory/dbo-system/service/account/internal/models"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	context "golang.org/x/net/context"
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

// AddingAcc mocks base method.
func (m *MockUseCase) AddingAcc(ctx context.Context, acc_uuid uuid.UUID, add_value float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddingAcc", ctx, acc_uuid, add_value)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddingAcc indicates an expected call of AddingAcc.
func (mr *MockUseCaseMockRecorder) AddingAcc(ctx, acc_uuid, add_value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddingAcc", reflect.TypeOf((*MockUseCase)(nil).AddingAcc), ctx, acc_uuid, add_value)
}

// BlockAcc mocks base method.
func (m *MockUseCase) BlockAcc(ctx context.Context, acc_uuid uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockAcc", ctx, acc_uuid)
	ret0, _ := ret[0].(error)
	return ret0
}

// BlockAcc indicates an expected call of BlockAcc.
func (mr *MockUseCaseMockRecorder) BlockAcc(ctx, acc_uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockAcc", reflect.TypeOf((*MockUseCase)(nil).BlockAcc), ctx, acc_uuid)
}

// CloseAcc mocks base method.
func (m *MockUseCase) CloseAcc(ctx context.Context, acc_uuid uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseAcc", ctx, acc_uuid)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseAcc indicates an expected call of CloseAcc.
func (mr *MockUseCaseMockRecorder) CloseAcc(ctx, acc_uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseAcc", reflect.TypeOf((*MockUseCase)(nil).CloseAcc), ctx, acc_uuid)
}

// CreateAcc mocks base method.
func (m *MockUseCase) CreateAcc(ctx context.Context, acc_uuid uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAcc", ctx, acc_uuid)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAcc indicates an expected call of CreateAcc.
func (mr *MockUseCaseMockRecorder) CreateAcc(ctx, acc_uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAcc", reflect.TypeOf((*MockUseCase)(nil).CreateAcc), ctx, acc_uuid)
}

// GetAccInfo mocks base method.
func (m *MockUseCase) GetAccInfo(ctx context.Context, acc_uuid uuid.UUID) (*models.FullAccountData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccInfo", ctx, acc_uuid)
	ret0, _ := ret[0].(*models.FullAccountData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccInfo indicates an expected call of GetAccInfo.
func (mr *MockUseCaseMockRecorder) GetAccInfo(ctx, acc_uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccInfo", reflect.TypeOf((*MockUseCase)(nil).GetAccInfo), ctx, acc_uuid)
}

// OpenAcc mocks base method.
func (m *MockUseCase) OpenAcc(ctx context.Context, acc_uuid uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenAcc", ctx, acc_uuid)
	ret0, _ := ret[0].(error)
	return ret0
}

// OpenAcc indicates an expected call of OpenAcc.
func (mr *MockUseCaseMockRecorder) OpenAcc(ctx, acc_uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenAcc", reflect.TypeOf((*MockUseCase)(nil).OpenAcc), ctx, acc_uuid)
}

// ReservAcc mocks base method.
func (m *MockUseCase) ReservAcc(ctx context.Context, acc_data *models.FullAccountData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReservAcc", ctx, acc_data)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReservAcc indicates an expected call of ReservAcc.
func (mr *MockUseCaseMockRecorder) ReservAcc(ctx, acc_data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReservAcc", reflect.TypeOf((*MockUseCase)(nil).ReservAcc), ctx, acc_data)
}

// ValidateAccBankNumber mocks base method.
func (m *MockUseCase) ValidateAccBankNumber(ctx context.Context, acc_bank_number string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAccBankNumber", ctx, acc_bank_number)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateAccBankNumber indicates an expected call of ValidateAccBankNumber.
func (mr *MockUseCaseMockRecorder) ValidateAccBankNumber(ctx, acc_bank_number interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAccBankNumber", reflect.TypeOf((*MockUseCase)(nil).ValidateAccBankNumber), ctx, acc_bank_number)
}

// ValidateAccCorrNumberOwner mocks base method.
func (m *MockUseCase) ValidateAccCorrNumberOwner(ctx context.Context, acc_corr_number_owner string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAccCorrNumberOwner", ctx, acc_corr_number_owner)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateAccCorrNumberOwner indicates an expected call of ValidateAccCorrNumberOwner.
func (mr *MockUseCaseMockRecorder) ValidateAccCorrNumberOwner(ctx, acc_corr_number_owner interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAccCorrNumberOwner", reflect.TypeOf((*MockUseCase)(nil).ValidateAccCorrNumberOwner), ctx, acc_corr_number_owner)
}

// ValidateAccCountry mocks base method.
func (m *MockUseCase) ValidateAccCountry(ctx context.Context, acc_country string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAccCountry", ctx, acc_country)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateAccCountry indicates an expected call of ValidateAccCountry.
func (mr *MockUseCaseMockRecorder) ValidateAccCountry(ctx, acc_country interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAccCountry", reflect.TypeOf((*MockUseCase)(nil).ValidateAccCountry), ctx, acc_country)
}

// ValidateAccCountryRegion mocks base method.
func (m *MockUseCase) ValidateAccCountryRegion(ctx context.Context, acc_country, acc_country_region string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAccCountryRegion", ctx, acc_country, acc_country_region)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateAccCountryRegion indicates an expected call of ValidateAccCountryRegion.
func (mr *MockUseCaseMockRecorder) ValidateAccCountryRegion(ctx, acc_country, acc_country_region interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAccCountryRegion", reflect.TypeOf((*MockUseCase)(nil).ValidateAccCountryRegion), ctx, acc_country, acc_country_region)
}

// ValidateAccMainOffice mocks base method.
func (m *MockUseCase) ValidateAccMainOffice(ctx context.Context, acc_country, acc_country_region, acc_main_office string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAccMainOffice", ctx, acc_country, acc_country_region, acc_main_office)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateAccMainOffice indicates an expected call of ValidateAccMainOffice.
func (mr *MockUseCaseMockRecorder) ValidateAccMainOffice(ctx, acc_country, acc_country_region, acc_main_office interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAccMainOffice", reflect.TypeOf((*MockUseCase)(nil).ValidateAccMainOffice), ctx, acc_country, acc_country_region, acc_main_office)
}

// ValidateAccMoneyValue mocks base method.
func (m *MockUseCase) ValidateAccMoneyValue(ctx context.Context, money_value uint8) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAccMoneyValue", ctx, money_value)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateAccMoneyValue indicates an expected call of ValidateAccMoneyValue.
func (mr *MockUseCaseMockRecorder) ValidateAccMoneyValue(ctx, money_value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAccMoneyValue", reflect.TypeOf((*MockUseCase)(nil).ValidateAccMoneyValue), ctx, money_value)
}

// ValidateAccStatus mocks base method.
func (m *MockUseCase) ValidateAccStatus(ctx context.Context, status uint8) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAccStatus", ctx, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateAccStatus indicates an expected call of ValidateAccStatus.
func (mr *MockUseCaseMockRecorder) ValidateAccStatus(ctx, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAccStatus", reflect.TypeOf((*MockUseCase)(nil).ValidateAccStatus), ctx, status)
}

// ValidateActivity mocks base method.
func (m *MockUseCase) ValidateActivity(ctx context.Context, owner, activity string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateActivity", ctx, owner, activity)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateActivity indicates an expected call of ValidateActivity.
func (mr *MockUseCaseMockRecorder) ValidateActivity(ctx, owner, activity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateActivity", reflect.TypeOf((*MockUseCase)(nil).ValidateActivity), ctx, owner, activity)
}

// ValidateBIC mocks base method.
func (m *MockUseCase) ValidateBIC(ctx context.Context, BIC string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateBIC", ctx, BIC)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateBIC indicates an expected call of ValidateBIC.
func (mr *MockUseCaseMockRecorder) ValidateBIC(ctx, BIC interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateBIC", reflect.TypeOf((*MockUseCase)(nil).ValidateBIC), ctx, BIC)
}

// ValidateCorrNumber mocks base method.
func (m *MockUseCase) ValidateCorrNumber(ctx context.Context, acc_corr_number, acc_bic string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateCorrNumber", ctx, acc_corr_number, acc_bic)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateCorrNumber indicates an expected call of ValidateCorrNumber.
func (mr *MockUseCaseMockRecorder) ValidateCorrNumber(ctx, acc_corr_number, acc_bic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateCorrNumber", reflect.TypeOf((*MockUseCase)(nil).ValidateCorrNumber), ctx, acc_corr_number, acc_bic)
}

// ValidateCulcNumber mocks base method.
func (m *MockUseCase) ValidateCulcNumber(ctx context.Context, culc_number string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateCulcNumber", ctx, culc_number)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateCulcNumber indicates an expected call of ValidateCulcNumber.
func (mr *MockUseCaseMockRecorder) ValidateCulcNumber(ctx, culc_number interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateCulcNumber", reflect.TypeOf((*MockUseCase)(nil).ValidateCulcNumber), ctx, culc_number)
}

// ValidateCurrency mocks base method.
func (m *MockUseCase) ValidateCurrency(ctx context.Context, currency string) (uint8, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateCurrency", ctx, currency)
	ret0, _ := ret[0].(uint8)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateCurrency indicates an expected call of ValidateCurrency.
func (mr *MockUseCaseMockRecorder) ValidateCurrency(ctx, currency interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateCurrency", reflect.TypeOf((*MockUseCase)(nil).ValidateCurrency), ctx, currency)
}

// ValidateKPP mocks base method.
func (m *MockUseCase) ValidateKPP(ctx context.Context, acc_kpp, acc_bic string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateKPP", ctx, acc_kpp, acc_bic)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateKPP indicates an expected call of ValidateKPP.
func (mr *MockUseCaseMockRecorder) ValidateKPP(ctx, acc_kpp, acc_bic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateKPP", reflect.TypeOf((*MockUseCase)(nil).ValidateKPP), ctx, acc_kpp, acc_bic)
}

// ValidateOwner mocks base method.
func (m *MockUseCase) ValidateOwner(ctx context.Context, owner string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateOwner", ctx, owner)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateOwner indicates an expected call of ValidateOwner.
func (mr *MockUseCaseMockRecorder) ValidateOwner(ctx, owner interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateOwner", reflect.TypeOf((*MockUseCase)(nil).ValidateOwner), ctx, owner)
}

// ValidatePaymentSystem mocks base method.
func (m *MockUseCase) ValidatePaymentSystem(ctx context.Context, payment_system string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidatePaymentSystem", ctx, payment_system)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidatePaymentSystem indicates an expected call of ValidatePaymentSystem.
func (mr *MockUseCaseMockRecorder) ValidatePaymentSystem(ctx, payment_system interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidatePaymentSystem", reflect.TypeOf((*MockUseCase)(nil).ValidatePaymentSystem), ctx, payment_system)
}

// WidthAcc mocks base method.
func (m *MockUseCase) WidthAcc(ctx context.Context, acc_uuid uuid.UUID, width_value float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WidthAcc", ctx, acc_uuid, width_value)
	ret0, _ := ret[0].(error)
	return ret0
}

// WidthAcc indicates an expected call of WidthAcc.
func (mr *MockUseCaseMockRecorder) WidthAcc(ctx, acc_uuid, width_value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WidthAcc", reflect.TypeOf((*MockUseCase)(nil).WidthAcc), ctx, acc_uuid, width_value)
}
