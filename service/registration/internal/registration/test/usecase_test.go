package test

//
//import (
//	"encoding/json"
//	"github.com/GCFactory/dbo-system/platform/config"
//	"github.com/GCFactory/dbo-system/platform/pkg/logger"
//	"github.com/GCFactory/dbo-system/service/registration/internal/models"
//	"github.com/GCFactory/dbo-system/service/registration/internal/registration/mock"
//	"github.com/GCFactory/dbo-system/service/registration/internal/registration/repository"
//	"github.com/GCFactory/dbo-system/service/registration/internal/registration/usecase"
//	"github.com/golang/mock/gomock"
//	"github.com/google/uuid"
//	"github.com/opentracing/opentracing-go"
//	"github.com/stretchr/testify/require"
//	"golang.org/x/net/context"
//	"testing"
//)
//
//var (
//	testCfgUC = &config.Config{
//		Env: "Development",
//		Logger: config.Logger{
//			Development: true,
//			Level:       "Debug",
//		},
//	}
//)
//
//func TestRegistrationUC_CreateSaga(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	apiLogger := logger.NewServerLogger(testCfgUC)
//	apiLogger.InitLogger()
//
//	mockRepo := mock.NewMockRepository(ctrl)
//	registrationUC := usecase.NewRegistrationUseCase(testCfgUC, mockRepo, apiLogger)
//
//	t.Parallel()
//
//	saga_uuid := uuid.New()
//	saga := models.Saga{
//		Saga_uuid: saga_uuid,
//		Saga_name: "test saga",
//		Saga_type: usecase.SagaRegistration,
//	}
//
//	t.Run("Success", func(t *testing.T) {
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.CreateSaga")
//		defer span.Finish()
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(nil, repository.ErrorCreateSaga)
//
//		completed_saga := saga
//		completed_saga.Saga_status = usecase.StatusInProcess
//
//		list_of_events := models.SagaListEvents{}
//
//		for i := 0; i < len(usecase.PossibleEventsListForSagaType[saga.Saga_type]); i++ {
//			list_of_events.EventList = append(list_of_events.EventList, usecase.PossibleEventsListForSagaType[saga.Saga_type][i])
//		}
//
//		list_of_events_byte, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_byte)
//
//		completed_saga.Saga_list_of_events = string(list_of_events_byte)
//
//		mockRepo.EXPECT().CreateSaga(ctxWithTrace, gomock.Eq(completed_saga)).Return(nil)
//
//		err = registrationUC.CreateSaga(ctx, saga)
//		require.Nil(t, err)
//
//	})
//	t.Run("Saga already exist", func(t *testing.T) {
//
//		completed_saga := saga
//		completed_saga.Saga_status = usecase.StatusInProcess
//
//		list_of_events := models.SagaListEvents{}
//
//		for i := 0; i < len(usecase.PossibleEventsListForSagaType[saga.Saga_type]); i++ {
//			list_of_events.EventList = append(list_of_events.EventList, usecase.PossibleEventsListForSagaType[saga.Saga_type][i])
//		}
//
//		list_of_events_byte, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_byte)
//
//		completed_saga.Saga_list_of_events = string(list_of_events_byte)
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.CreateSaga")
//		defer span.Finish()
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&completed_saga, nil)
//
//		err = registrationUC.CreateSaga(ctx, saga)
//		require.Equal(t, err, usecase.ErrorSagaExist)
//
//	})
//	t.Run("Wrong saga type", func(t *testing.T) {
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.CreateSaga")
//		defer span.Finish()
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(nil, repository.ErrorGetSaga)
//
//		completed_saga := saga
//		completed_saga.Saga_type = 1000
//
//		err := registrationUC.CreateSaga(ctx, completed_saga)
//		require.Equal(t, err, usecase.ErrorWrongSagaType)
//
//	})
//	t.Run("Error create saga", func(t *testing.T) {
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.CreateSaga")
//		defer span.Finish()
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(nil, repository.ErrorGetSaga)
//
//		completed_saga := saga
//		completed_saga.Saga_status = usecase.StatusInProcess
//
//		list_of_events := models.SagaListEvents{}
//
//		for i := 0; i < len(usecase.PossibleEventsListForSagaType[saga.Saga_type]); i++ {
//			list_of_events.EventList = append(list_of_events.EventList, usecase.PossibleEventsListForSagaType[saga.Saga_type][i])
//		}
//
//		list_of_events_byte, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_byte)
//
//		completed_saga.Saga_list_of_events = string(list_of_events_byte)
//
//		mockRepo.EXPECT().CreateSaga(ctxWithTrace, gomock.Eq(completed_saga)).Return(repository.ErrorCreateSaga)
//
//		err = registrationUC.CreateSaga(ctx, saga)
//		require.Equal(t, err, usecase.ErrorCreateSaga)
//
//	})
//}
//
//func TestRegistrationUC_DeleteSaga(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	apiLogger := logger.NewServerLogger(testCfgUC)
//	apiLogger.InitLogger()
//
//	mockRepo := mock.NewMockRepository(ctrl)
//	registrationUC := usecase.NewRegistrationUseCase(testCfgUC, mockRepo, apiLogger)
//
//	t.Parallel()
//
//	t.Run("Success", func(t *testing.T) {
//		saga_uuid := uuid.New()
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.DeleteSaga")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_name:           "test saga",
//			Saga_type:           usecase.SagaRegistration,
//			Saga_status:         usecase.StatusCompleted,
//			Saga_list_of_events: string(list_of_events_bytes),
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//		mockRepo.EXPECT().DeleteSaga(ctxWithTrace, gomock.Eq(saga_uuid)).Return(nil)
//
//		err = registrationUC.DeleteSaga(ctx, saga_uuid)
//		require.Nil(t, err)
//
//	})
//	t.Run("No saga", func(t *testing.T) {
//		saga_uuid := uuid.New()
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.DeleteSaga")
//		defer span.Finish()
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(nil, repository.ErrorGetSaga)
//
//		err := registrationUC.DeleteSaga(ctx, saga_uuid)
//		require.Equal(t, err, usecase.ErrorNoSaga)
//
//	})
//	t.Run("Bad status", func(t *testing.T) {
//		saga_uuid := uuid.New()
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.DeleteSaga")
//		defer span.Finish()
//
//		saga := models.Saga{
//			Saga_uuid:   saga_uuid,
//			Saga_status: 10000,
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//
//		err := registrationUC.DeleteSaga(ctx, saga_uuid)
//		require.Equal(t, err, usecase.ErrorSagaWasNotCompeted)
//
//	})
//	t.Run("Not empty event list", func(t *testing.T) {
//		saga_uuid := uuid.New()
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.DeleteSaga")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//		list_of_events.EventList = append(list_of_events.EventList, "Event")
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_name:           "test saga",
//			Saga_type:           usecase.SagaRegistration,
//			Saga_status:         usecase.StatusCompleted,
//			Saga_list_of_events: string(list_of_events_bytes),
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//
//		err = registrationUC.DeleteSaga(ctx, saga_uuid)
//		require.Equal(t, err, usecase.ErrorNotEmptyEventList)
//
//	})
//	t.Run("Error delete saga", func(t *testing.T) {
//		saga_uuid := uuid.New()
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.DeleteSaga")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_name:           "test saga",
//			Saga_type:           usecase.SagaRegistration,
//			Saga_status:         usecase.StatusCompleted,
//			Saga_list_of_events: string(list_of_events_bytes),
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//		mockRepo.EXPECT().DeleteSaga(ctxWithTrace, gomock.Eq(saga_uuid)).Return(repository.ErrorDeleteSaga)
//
//		err = registrationUC.DeleteSaga(ctx, saga_uuid)
//		require.Equal(t, err, usecase.ErrorDeleteSaga)
//
//	})
//
//}
//
//func TestRegistrationUC_RemoveSagaEvent(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	apiLogger := logger.NewServerLogger(testCfgUC)
//	apiLogger.InitLogger()
//
//	mockRepo := mock.NewMockRepository(ctrl)
//	registrationUC := usecase.NewRegistrationUseCase(testCfgUC, mockRepo, apiLogger)
//
//	t.Parallel()
//	saga_uuid := uuid.New()
//	event_name := usecase.AntiFrod
//
//	t.Run("Wrong event name", func(t *testing.T) {
//
//		err := registrationUC.RemoveSagaEvent(context.Background(), saga_uuid, "some name")
//		require.Equal(t, err, usecase.ErrorWrongEventName)
//
//	})
//	t.Run("No saga with this id", func(t *testing.T) {
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.RemoveSagaEvent")
//		defer span.Finish()
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(nil, repository.ErrorGetSaga)
//
//		err := registrationUC.RemoveSagaEvent(ctx, saga_uuid, event_name)
//		require.Equal(t, err, usecase.ErrorNoSaga)
//
//	})
//	t.Run("Wrong event name for this saga type", func(t *testing.T) {
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.RemoveSagaEvent")
//		defer span.Finish()
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_name:           "test saga",
//			Saga_type:           usecase.SagaRegistration,
//			Saga_status:         usecase.StatusInProcess,
//			Saga_list_of_events: "",
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//
//		err := registrationUC.RemoveSagaEvent(ctx, saga_uuid, usecase.TestEvent)
//		require.Equal(t, err, usecase.ErrorWrongEventNameForSagaType)
//
//	})
//	t.Run("Saga hasn't such event", func(t *testing.T) {
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.RemoveSagaEvent")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//		list_of_events.EventList = append(list_of_events.EventList, usecase.ReserveNumber)
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_name:           "test saga",
//			Saga_type:           usecase.SagaRegistration,
//			Saga_status:         usecase.StatusInProcess,
//			Saga_list_of_events: string(list_of_events_bytes),
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//
//		err = registrationUC.RemoveSagaEvent(ctx, saga_uuid, usecase.AntiFrod)
//		require.Equal(t, err, usecase.ErrorNoEvent)
//
//	})
//	t.Run("Update saga events error", func(t *testing.T) {
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.RemoveSagaEvent")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//		list_of_events.EventList = append(list_of_events.EventList, usecase.ReserveNumber)
//
//		list_of_events_bytes_end, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes_end)
//
//		list_of_events.EventList = append(list_of_events.EventList, usecase.AntiFrod)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_name:           "test saga",
//			Saga_type:           usecase.SagaRegistration,
//			Saga_status:         usecase.StatusInProcess,
//			Saga_list_of_events: string(list_of_events_bytes),
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//		mockRepo.EXPECT().UpdateSagaEvents(ctxWithTrace, gomock.Eq(saga_uuid), gomock.Eq(string(list_of_events_bytes_end))).Return(repository.ErrorUpdateSagaEvents)
//
//		err = registrationUC.RemoveSagaEvent(ctx, saga_uuid, usecase.AntiFrod)
//		require.Equal(t, err, usecase.ErrorUpdateSagaEvents)
//
//	})
//	t.Run("Update saga status error with empty events list with in process status", func(t *testing.T) {
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.RemoveSagaEvent")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//
//		list_of_events.EventList = append(list_of_events.EventList, usecase.AntiFrod)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes)
//
//		list_of_events.EventList = list_of_events.EventList[:0]
//
//		list_of_events_bytes_end, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes_end)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_name:           "test saga",
//			Saga_type:           usecase.SagaRegistration,
//			Saga_status:         usecase.StatusInProcess,
//			Saga_list_of_events: string(list_of_events_bytes),
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//		mockRepo.EXPECT().UpdateSagaEvents(ctxWithTrace, gomock.Eq(saga_uuid), gomock.Eq(string(list_of_events_bytes_end))).Return(nil)
//		mockRepo.EXPECT().UpdateSagaStatus(ctxWithTrace, gomock.Eq(saga_uuid), gomock.Eq(usecase.StatusCompleted)).Return(repository.ErrorUpdateSagaStatus)
//
//		err = registrationUC.RemoveSagaEvent(ctx, saga_uuid, usecase.AntiFrod)
//		require.Equal(t, err, usecase.ErrorUpdateSagaStatus)
//
//	})
//	t.Run("Update saga status error with empty events list with fall back status", func(t *testing.T) {
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.RemoveSagaEvent")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//
//		list_of_events.EventList = append(list_of_events.EventList, "-"+usecase.AntiFrod)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes)
//
//		list_of_events.EventList = list_of_events.EventList[:0]
//
//		list_of_events_bytes_end, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes_end)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_name:           "test saga",
//			Saga_type:           usecase.SagaRegistration,
//			Saga_status:         usecase.StatusFallBack,
//			Saga_list_of_events: string(list_of_events_bytes),
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//		mockRepo.EXPECT().UpdateSagaEvents(ctxWithTrace, gomock.Eq(saga_uuid), gomock.Eq(string(list_of_events_bytes_end))).Return(nil)
//		mockRepo.EXPECT().UpdateSagaStatus(ctxWithTrace, gomock.Eq(saga_uuid), gomock.Eq(usecase.StatusError)).Return(repository.ErrorUpdateSagaStatus)
//
//		err = registrationUC.RemoveSagaEvent(ctx, saga_uuid, "-"+usecase.AntiFrod)
//		require.Equal(t, err, usecase.ErrorUpdateSagaStatus)
//
//	})
//	t.Run("Success with several usual events", func(t *testing.T) {
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.RemoveSagaEvent")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//
//		list_of_events.EventList = append(list_of_events.EventList, usecase.ReserveNumber)
//
//		list_of_events_bytes_end, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes_end)
//
//		list_of_events.EventList = append(list_of_events.EventList, usecase.AntiFrod)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_name:           "test saga",
//			Saga_type:           usecase.SagaRegistration,
//			Saga_status:         usecase.StatusInProcess,
//			Saga_list_of_events: string(list_of_events_bytes),
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//		mockRepo.EXPECT().UpdateSagaEvents(ctxWithTrace, gomock.Eq(saga_uuid), gomock.Eq(string(list_of_events_bytes_end))).Return(nil)
//
//		err = registrationUC.RemoveSagaEvent(ctx, saga_uuid, usecase.AntiFrod)
//		require.Nil(t, err)
//
//	})
//	t.Run("Success with several inverse events and fall back event", func(t *testing.T) {
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.RemoveSagaEvent")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//
//		list_of_events.EventList = append(list_of_events.EventList, usecase.ReserveNumber)
//
//		list_of_events_bytes_end, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes_end)
//
//		list_of_events.EventList = append(list_of_events.EventList, "-"+usecase.AntiFrod)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_name:           "test saga",
//			Saga_type:           usecase.SagaRegistration,
//			Saga_status:         usecase.StatusInProcess,
//			Saga_list_of_events: string(list_of_events_bytes),
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//		mockRepo.EXPECT().UpdateSagaEvents(ctxWithTrace, gomock.Eq(saga_uuid), gomock.Eq(string(list_of_events_bytes_end))).Return(nil)
//
//		err = registrationUC.RemoveSagaEvent(ctx, saga_uuid, "-"+usecase.AntiFrod)
//		require.Nil(t, err)
//
//	})
//	t.Run("Success with solo usual event", func(t *testing.T) {
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.RemoveSagaEvent")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//
//		list_of_events.EventList = append(list_of_events.EventList, usecase.AntiFrod)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes)
//
//		list_of_events.EventList = list_of_events.EventList[:0]
//
//		list_of_events_bytes_end, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes_end)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_name:           "test saga",
//			Saga_type:           usecase.SagaRegistration,
//			Saga_status:         usecase.StatusInProcess,
//			Saga_list_of_events: string(list_of_events_bytes),
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//		mockRepo.EXPECT().UpdateSagaEvents(ctxWithTrace, gomock.Eq(saga_uuid), gomock.Eq(string(list_of_events_bytes_end))).Return(nil)
//		mockRepo.EXPECT().UpdateSagaStatus(ctxWithTrace, gomock.Eq(saga_uuid), gomock.Eq(usecase.StatusCompleted)).Return(nil)
//
//		err = registrationUC.RemoveSagaEvent(ctx, saga_uuid, usecase.AntiFrod)
//		require.Nil(t, err)
//
//	})
//	t.Run("Success with solo inverse event and fall back event", func(t *testing.T) {
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.RemoveSagaEvent")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//
//		list_of_events.EventList = append(list_of_events.EventList, "-"+usecase.AntiFrod)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes)
//
//		list_of_events.EventList = list_of_events.EventList[:0]
//
//		list_of_events_bytes_end, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes_end)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_name:           "test saga",
//			Saga_type:           usecase.SagaRegistration,
//			Saga_status:         usecase.StatusFallBack,
//			Saga_list_of_events: string(list_of_events_bytes),
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//		mockRepo.EXPECT().UpdateSagaEvents(ctxWithTrace, gomock.Eq(saga_uuid), gomock.Eq(string(list_of_events_bytes_end))).Return(nil)
//		mockRepo.EXPECT().UpdateSagaStatus(ctxWithTrace, gomock.Eq(saga_uuid), gomock.Eq(usecase.StatusError)).Return(nil)
//
//		err = registrationUC.RemoveSagaEvent(ctx, saga_uuid, "-"+usecase.AntiFrod)
//		require.Nil(t, err)
//
//	})
//	t.Run("Success with solo inverse event and fall back event", func(t *testing.T) {
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.RemoveSagaEvent")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//
//		list_of_events.EventList = append(list_of_events.EventList, usecase.AntiFrod)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes)
//
//		list_of_events.EventList[0] = "-" + usecase.AntiFrod
//
//		list_of_events_bytes_end, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes_end)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_name:           "test saga",
//			Saga_type:           usecase.SagaRegistration,
//			Saga_status:         usecase.StatusFallBack,
//			Saga_list_of_events: string(list_of_events_bytes),
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//		mockRepo.EXPECT().UpdateSagaEvents(ctxWithTrace, gomock.Eq(saga_uuid), gomock.Eq(string(list_of_events_bytes_end))).Return(nil)
//
//		err = registrationUC.RemoveSagaEvent(ctx, saga_uuid, usecase.AntiFrod)
//		require.Nil(t, err)
//
//	})
//}
//
//func TestRegistrationUC_FallBackEvent(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	apiLogger := logger.NewServerLogger(testCfgUC)
//	apiLogger.InitLogger()
//
//	mockRepo := mock.NewMockRepository(ctrl)
//	registrationUC := usecase.NewRegistrationUseCase(testCfgUC, mockRepo, apiLogger)
//
//	t.Parallel()
//
//	t.Run("Wrong saga status", func(t *testing.T) {
//		saga := &models.Saga{
//			Saga_status: usecase.StatusUndefined,
//		}
//
//		saga, err := registrationUC.FallBackEvent(context.Background(), saga, usecase.AntiFrod)
//		require.Equal(t, err, usecase.ErrorWrongFallBackStatusForFallBackEvent)
//		require.Nil(t, saga)
//	})
//	t.Run("Wrong event name", func(t *testing.T) {
//		saga := &models.Saga{
//			Saga_status: usecase.StatusFallBack,
//		}
//
//		saga, err := registrationUC.FallBackEvent(context.Background(), saga, "smth")
//		require.Equal(t, err, usecase.ErrorWrongEventName)
//		require.Nil(t, saga)
//	})
//	t.Run("Wrong event name for this saga type", func(t *testing.T) {
//		saga := &models.Saga{
//			Saga_status: usecase.StatusFallBack,
//		}
//
//		saga, err := registrationUC.FallBackEvent(context.Background(), saga, usecase.TestEvent)
//		require.Equal(t, err, usecase.ErrorWrongEventNameForSagaType)
//		require.Nil(t, saga)
//	})
//	t.Run("Get inverse event", func(t *testing.T) {
//		saga := &models.Saga{
//			Saga_status: usecase.StatusFallBack,
//		}
//
//		saga, err := registrationUC.FallBackEvent(context.Background(), saga, "-"+usecase.AntiFrod)
//		require.Equal(t, err, usecase.ErrorInversedEvent)
//		require.Nil(t, saga)
//	})
//	t.Run("Get usual event and saga already has this inversed saga too", func(t *testing.T) {
//		list_of_events := models.SagaListEvents{}
//		list_of_events.EventList = append(list_of_events.EventList, "-"+usecase.AntiFrod)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes)
//
//		saga := &models.Saga{
//			Saga_status:         usecase.StatusFallBack,
//			Saga_list_of_events: string(list_of_events_bytes),
//		}
//
//		saga, err = registrationUC.FallBackEvent(context.Background(), saga, usecase.AntiFrod)
//		require.Equal(t, err, usecase.ErrorEventExist)
//		require.Nil(t, saga)
//	})
//	t.Run("No such event in current saga", func(t *testing.T) {
//		list_of_events := models.SagaListEvents{}
//		list_of_events.EventList = append(list_of_events.EventList, usecase.AntiFrod)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes)
//
//		saga := &models.Saga{
//			Saga_status:         usecase.StatusFallBack,
//			Saga_list_of_events: string(list_of_events_bytes),
//		}
//
//		saga, err = registrationUC.FallBackEvent(context.Background(), saga, usecase.ReserveNumber)
//		require.Equal(t, err, usecase.ErrorNoEvent)
//		require.Nil(t, saga)
//	})
//	t.Run("Success", func(t *testing.T) {
//		list_of_events := models.SagaListEvents{}
//		list_of_events.EventList = append(list_of_events.EventList, usecase.AntiFrod)
//		list_of_events.EventList = append(list_of_events.EventList, usecase.ReserveNumber)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_bytes)
//
//		list_of_events.EventList = list_of_events.EventList[:1]
//		list_of_events.EventList = append(list_of_events.EventList, "-"+usecase.ReserveNumber)
//		list_of_events_end_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_end_bytes)
//
//		saga := &models.Saga{
//			Saga_status:         usecase.StatusFallBack,
//			Saga_list_of_events: string(list_of_events_bytes),
//		}
//
//		saga_end := &models.Saga{
//			Saga_status:         saga.Saga_status,
//			Saga_list_of_events: string(list_of_events_end_bytes),
//		}
//
//		result, err := registrationUC.FallBackEvent(context.Background(), saga, usecase.ReserveNumber)
//		require.Nil(t, err)
//		require.Equal(t, result, saga_end)
//	})
//}
//
//func TestRegistrationUC_FallBackSaga(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	apiLogger := logger.NewServerLogger(testCfgUC)
//	apiLogger.InitLogger()
//
//	mockRepo := mock.NewMockRepository(ctrl)
//	registrationUC := usecase.NewRegistrationUseCase(testCfgUC, mockRepo, apiLogger)
//
//	t.Parallel()
//
//	t.Run("No saga", func(t *testing.T) {
//		saga_uuid := uuid.New()
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.FallBackSaga")
//		defer span.Finish()
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(nil, repository.ErrorGetSaga)
//
//		err := registrationUC.FallBackSaga(ctx, saga_uuid, "smth")
//		require.Equal(t, err, usecase.ErrorNoSaga)
//	})
//	t.Run("Saga already fall back", func(t *testing.T) {
//		saga_uuid := uuid.New()
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.FallBackSaga")
//		defer span.Finish()
//
//		saga := models.Saga{
//			Saga_uuid:   saga_uuid,
//			Saga_status: usecase.StatusFallBack,
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//
//		err := registrationUC.FallBackSaga(ctx, saga_uuid, "smth")
//		require.Equal(t, err, usecase.ErrorSagaAlreadyFallBack)
//	})
//	t.Run("Saga has wrong status for beginning fall back", func(t *testing.T) {
//		saga_uuid := uuid.New()
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.FallBackSaga")
//		defer span.Finish()
//
//		saga := models.Saga{
//			Saga_uuid:   saga_uuid,
//			Saga_status: usecase.StatusCreated,
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//
//		err := registrationUC.FallBackSaga(ctx, saga_uuid, "smth")
//		require.Equal(t, err, usecase.ErrorWrongFallBackStatus)
//	})
//	t.Run("Wrong event name", func(t *testing.T) {
//		saga_uuid := uuid.New()
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.FallBackSaga")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//		list_of_events.EventList = append(list_of_events.EventList, usecase.AntiFrod)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotNil(t, list_of_events_bytes)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_status:         usecase.StatusInProcess,
//			Saga_list_of_events: string(list_of_events_bytes),
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//
//		err = registrationUC.FallBackSaga(ctx, saga_uuid, "smth")
//		require.Equal(t, err, usecase.ErrorWrongEventName)
//	})
//	t.Run("Wrong event name for this saga type", func(t *testing.T) {
//		saga_uuid := uuid.New()
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.FallBackSaga")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//		list_of_events.EventList = append(list_of_events.EventList, usecase.AntiFrod)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotNil(t, list_of_events_bytes)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_status:         usecase.StatusInProcess,
//			Saga_list_of_events: string(list_of_events_bytes),
//			Saga_type:           usecase.SagaRegistration,
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//
//		err = registrationUC.FallBackSaga(ctx, saga_uuid, usecase.TestEvent)
//		require.Equal(t, err, usecase.ErrorWrongEventNameForSagaType)
//	})
//	t.Run("Saga has inverse event", func(t *testing.T) {
//		saga_uuid := uuid.New()
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.FallBackSaga")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//		list_of_events.EventList = append(list_of_events.EventList, "-"+usecase.AntiFrod)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotNil(t, list_of_events_bytes)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_status:         usecase.StatusInProcess,
//			Saga_list_of_events: string(list_of_events_bytes),
//			Saga_type:           usecase.SagaRegistration,
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//
//		err = registrationUC.FallBackSaga(ctx, saga_uuid, usecase.AntiFrod)
//		require.Equal(t, err, usecase.ErrorEventAlreadyFallBack)
//	})
//	t.Run("Error update saga events", func(t *testing.T) {
//		saga_uuid := uuid.New()
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.FallBackSaga")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//		list_of_events.EventList = append(list_of_events.EventList, usecase.AntiFrod)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotNil(t, list_of_events_bytes)
//
//		list_of_events.EventList[0] = "-" + list_of_events.EventList[0]
//		list_of_events.EventList = append(list_of_events.EventList, "-"+usecase.ReserveNumber)
//		list_of_events_end_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_end_bytes)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_status:         usecase.StatusInProcess,
//			Saga_list_of_events: string(list_of_events_bytes),
//			Saga_type:           usecase.SagaRegistration,
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//		mockRepo.EXPECT().UpdateSagaEvents(ctxWithTrace, gomock.Eq(saga_uuid), string(list_of_events_end_bytes)).Return(repository.ErrorUpdateSagaEvents)
//
//		err = registrationUC.FallBackSaga(ctx, saga_uuid, usecase.AntiFrod)
//		require.Equal(t, err, usecase.ErrorUpdateSagaEvents)
//	})
//	t.Run("Error update saga status", func(t *testing.T) {
//		saga_uuid := uuid.New()
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.FallBackSaga")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//		list_of_events.EventList = append(list_of_events.EventList, usecase.AntiFrod)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotNil(t, list_of_events_bytes)
//
//		list_of_events.EventList[0] = "-" + list_of_events.EventList[0]
//		list_of_events.EventList = append(list_of_events.EventList, "-"+usecase.ReserveNumber)
//		list_of_events_end_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_end_bytes)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_status:         usecase.StatusInProcess,
//			Saga_list_of_events: string(list_of_events_bytes),
//			Saga_type:           usecase.SagaRegistration,
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//		mockRepo.EXPECT().UpdateSagaEvents(ctxWithTrace, gomock.Eq(saga_uuid), string(list_of_events_end_bytes)).Return(nil)
//		mockRepo.EXPECT().UpdateSagaStatus(ctxWithTrace, gomock.Eq(saga_uuid), gomock.Eq(usecase.StatusFallBack)).Return(repository.ErrorUpdateSagaStatus)
//
//		err = registrationUC.FallBackSaga(ctx, saga_uuid, usecase.AntiFrod)
//		require.Equal(t, err, usecase.ErrorUpdateSagaStatus)
//	})
//	t.Run("Success", func(t *testing.T) {
//		saga_uuid := uuid.New()
//
//		ctx := context.Background()
//		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.FallBackSaga")
//		defer span.Finish()
//
//		list_of_events := models.SagaListEvents{}
//		list_of_events.EventList = append(list_of_events.EventList, usecase.AntiFrod)
//
//		list_of_events_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotNil(t, list_of_events_bytes)
//
//		list_of_events.EventList[0] = "-" + list_of_events.EventList[0]
//		list_of_events.EventList = append(list_of_events.EventList, "-"+usecase.ReserveNumber)
//		list_of_events_end_bytes, err := json.Marshal(list_of_events)
//		require.Nil(t, err)
//		require.NotEmpty(t, list_of_events_end_bytes)
//
//		saga := models.Saga{
//			Saga_uuid:           saga_uuid,
//			Saga_status:         usecase.StatusInProcess,
//			Saga_list_of_events: string(list_of_events_bytes),
//			Saga_type:           usecase.SagaRegistration,
//		}
//
//		mockRepo.EXPECT().GetSagaById(ctxWithTrace, gomock.Eq(saga_uuid)).Return(&saga, nil)
//		mockRepo.EXPECT().UpdateSagaEvents(ctxWithTrace, gomock.Eq(saga_uuid), string(list_of_events_end_bytes)).Return(nil)
//		mockRepo.EXPECT().UpdateSagaStatus(ctxWithTrace, gomock.Eq(saga_uuid), gomock.Eq(usecase.StatusFallBack)).Return(nil)
//
//		err = registrationUC.FallBackSaga(ctx, saga_uuid, usecase.AntiFrod)
//		require.Nil(t, err)
//	})
//}
