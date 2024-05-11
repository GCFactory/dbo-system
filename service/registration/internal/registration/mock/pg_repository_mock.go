// Code generated by MockGen. DO NOT EDIT.
// Source: pg_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/GCFactory/dbo-system/service/registration/internal/models"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
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

// AddEventInToSaga mocks base method.
func (m *MockRepository) AddEventInToSaga(ctx context.Context, saga_uuid uuid.UUID, event_name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEventInToSaga", ctx, saga_uuid, event_name)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddEventInToSaga indicates an expected call of AddEventInToSaga.
func (mr *MockRepositoryMockRecorder) AddEventInToSaga(ctx, saga_uuid, event_name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEventInToSaga", reflect.TypeOf((*MockRepository)(nil).AddEventInToSaga), ctx, saga_uuid, event_name)
}

// ChangeSagaStatus mocks base method.
func (m *MockRepository) ChangeSagaStatus(ctx context.Context, saga_uuid uuid.UUID, saga_status uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeSagaStatus", ctx, saga_uuid, saga_status)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeSagaStatus indicates an expected call of ChangeSagaStatus.
func (mr *MockRepositoryMockRecorder) ChangeSagaStatus(ctx, saga_uuid, saga_status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeSagaStatus", reflect.TypeOf((*MockRepository)(nil).ChangeSagaStatus), ctx, saga_uuid, saga_status)
}

// CreateSaga mocks base method.
func (m *MockRepository) CreateSaga(ctx context.Context, saga models.Saga) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSaga", ctx, saga)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSaga indicates an expected call of CreateSaga.
func (mr *MockRepositoryMockRecorder) CreateSaga(ctx, saga interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSaga", reflect.TypeOf((*MockRepository)(nil).CreateSaga), ctx, saga)
}

// GetSagaById mocks base method.
func (m *MockRepository) GetSagaById(ctx context.Context, saga_uuid uuid.UUID) (models.Saga, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSagaById", ctx, saga_uuid)
	ret0, _ := ret[0].(models.Saga)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSagaById indicates an expected call of GetSagaById.
func (mr *MockRepositoryMockRecorder) GetSagaById(ctx, saga_uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSagaById", reflect.TypeOf((*MockRepository)(nil).GetSagaById), ctx, saga_uuid)
}

// RemoveEventFromSaga mocks base method.
func (m *MockRepository) RemoveEventFromSaga(ctx context.Context, saga_uuid uuid.UUID, event_name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveEventFromSaga", ctx, saga_uuid, event_name)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveEventFromSaga indicates an expected call of RemoveEventFromSaga.
func (mr *MockRepositoryMockRecorder) RemoveEventFromSaga(ctx, saga_uuid, event_name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveEventFromSaga", reflect.TypeOf((*MockRepository)(nil).RemoveEventFromSaga), ctx, saga_uuid, event_name)
}
