package test

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration/mock"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration/usecase"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	testCfgUC = &config.Config{
		Env: "Development",
		Logger: config.Logger{
			Development: true,
			Level:       "Debug",
		},
	}
)

func TestRegistrationUC_ProcessingSagaAndEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfgUC)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	regUC := usecase.NewRegistrationUseCase(testCfgUC, mockRepo, apiLogger)

	test_event := &models.Event{
		Event_uuid:          uuid.New(),
		Saga_uuid:           uuid.Nil,
		Event_status:        usecase.EventStatusUndefined,
		Event_is_roll_back:  false,
		Event_rollback_uuid: uuid.Nil,
		Event_name:          usecase.EventTypeCreateUser,
		Event_result:        "",
	}

	ctx := context.Background()

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.ProcessingSagaAndEvents")
	defer span.Finish()

	t.Parallel()

	t.Run("Events", func(t *testing.T) {

		t.Run("Event wrong status", func(t *testing.T) {

			mockRepo.EXPECT().GetEvent(gomock.Eq(ctxWithTrace), gomock.Eq(test_event.Event_uuid)).Return(test_event, nil)

			result, err := regUC.ProcessingSagaAndEvents(ctx, uuid.Nil, test_event.Event_uuid, true)
			require.Nil(t, result)
			require.Error(t, err)
			require.Equal(t, err, usecase.ErrorInvalidEventStatus)

		})

		t.Run("Success to created", func(t *testing.T) {

			old_status := test_event.Event_status
			test_event.Event_status = usecase.EventStatusCreated

			mockRepo.EXPECT().GetEvent(gomock.Eq(ctxWithTrace), gomock.Eq(test_event.Event_uuid)).Return(test_event, nil)

			test_event.Event_status = usecase.EventStatusInProgress
			mockRepo.EXPECT().UpdateEvent(gomock.Eq(ctxWithTrace), gomock.Eq(test_event)).Return(nil)

			test_event.Event_status = usecase.EventStatusCreated
			result, err := regUC.ProcessingSagaAndEvents(ctx, uuid.Nil, test_event.Event_uuid, true)
			require.Nil(t, err)
			require.Equal(t, []*models.Event{test_event}, result)

			test_event.Event_status = old_status
		})

		t.Run("Success to in process", func(t *testing.T) {
			t.Run("Success completed", func(t *testing.T) {

				old_status := test_event.Event_status
				test_event.Event_status = usecase.EventStatusInProgress

				mockRepo.EXPECT().GetEvent(gomock.Eq(ctxWithTrace), gomock.Eq(test_event.Event_uuid)).Return(test_event, nil)
				test_event.Event_status = usecase.EventStatusCompleted
				mockRepo.EXPECT().UpdateEvent(gomock.Eq(ctxWithTrace), gomock.Eq(test_event)).Return(nil)

				test_event.Event_status = usecase.EventStatusInProgress
				result, err := regUC.ProcessingSagaAndEvents(ctx, uuid.Nil, test_event.Event_uuid, true)
				require.Nil(t, err)
				require.Nil(t, result)

				test_event.Event_status = old_status
			})
			t.Run("Wrong compketed", func(t *testing.T) {

				old_status := test_event.Event_status
				test_event.Event_status = usecase.EventStatusInProgress

				mockRepo.EXPECT().GetEvent(gomock.Eq(ctxWithTrace), gomock.Eq(test_event.Event_uuid)).Return(test_event, nil)
				test_event.Event_status = usecase.EventStatusError
				mockRepo.EXPECT().UpdateEvent(gomock.Eq(ctxWithTrace), gomock.Eq(test_event)).Return(nil)

				test_event.Event_status = usecase.EventStatusInProgress
				result, err := regUC.ProcessingSagaAndEvents(ctx, uuid.Nil, test_event.Event_uuid, false)
				require.Nil(t, err)
				require.Nil(t, result)

				test_event.Event_status = old_status
			})
		})

		t.Run("Success to completed", func(t *testing.T) {

			old_status := test_event.Event_status
			test_event.Event_status = usecase.EventStatusCompleted

			mockRepo.EXPECT().GetEvent(gomock.Eq(ctxWithTrace), gomock.Eq(test_event.Event_uuid)).Return(test_event, nil)

			result, err := regUC.ProcessingSagaAndEvents(ctx, uuid.Nil, test_event.Event_uuid, true)
			require.Nil(t, err)
			require.Nil(t, result)

			test_event.Event_status = old_status
		})

		t.Run("Success to fallback in process", func(t *testing.T) {
			t.Run("Success completed", func(t *testing.T) {

				old_status := test_event.Event_status
				test_event.Event_status = usecase.EventStatusFallBackInProcess

				mockRepo.EXPECT().GetEvent(gomock.Eq(ctxWithTrace), gomock.Eq(test_event.Event_uuid)).Return(test_event, nil)
				test_event.Event_status = usecase.EventStatusFallBackCompleted
				mockRepo.EXPECT().UpdateEvent(gomock.Eq(ctxWithTrace), gomock.Eq(test_event)).Return(nil)

				test_event.Event_status = usecase.EventStatusFallBackInProcess
				result, err := regUC.ProcessingSagaAndEvents(ctx, uuid.Nil, test_event.Event_uuid, true)
				require.Nil(t, err)
				require.Nil(t, result)

				test_event.Event_status = old_status
			})
			t.Run("Wrong completed", func(t *testing.T) {

				old_status := test_event.Event_status
				test_event.Event_status = usecase.EventStatusFallBackInProcess

				mockRepo.EXPECT().GetEvent(gomock.Eq(ctxWithTrace), gomock.Eq(test_event.Event_uuid)).Return(test_event, nil)
				test_event.Event_status = usecase.EventStatusFallBackError
				mockRepo.EXPECT().UpdateEvent(gomock.Eq(ctxWithTrace), gomock.Eq(test_event)).Return(nil)

				test_event.Event_status = usecase.EventStatusFallBackInProcess
				result, err := regUC.ProcessingSagaAndEvents(ctx, uuid.Nil, test_event.Event_uuid, false)
				require.Nil(t, err)
				require.Nil(t, result)

				test_event.Event_status = old_status
			})
		})

		t.Run("Success fallback complete", func(t *testing.T) {

			old_status := test_event.Event_status
			test_event.Event_status = usecase.EventStatusFallBackCompleted

			mockRepo.EXPECT().GetEvent(gomock.Eq(ctxWithTrace), gomock.Eq(test_event.Event_uuid)).Return(test_event, nil)

			result, err := regUC.ProcessingSagaAndEvents(ctx, uuid.Nil, test_event.Event_uuid, true)
			require.Nil(t, err)
			require.Nil(t, result)

			test_event.Event_status = old_status

		})

		t.Run("Success fallback error", func(t *testing.T) {

			old_status := test_event.Event_status
			test_event.Event_status = usecase.EventStatusFallBackError

			mockRepo.EXPECT().GetEvent(gomock.Eq(ctxWithTrace), gomock.Eq(test_event.Event_uuid)).Return(test_event, nil)

			result, err := regUC.ProcessingSagaAndEvents(ctx, uuid.Nil, test_event.Event_uuid, true)
			require.Nil(t, err)
			require.Nil(t, result)

			test_event.Event_status = old_status

		})

		t.Run("Success error", func(t *testing.T) {

			old_status := test_event.Event_status
			test_event.Event_status = usecase.EventStatusError

			mockRepo.EXPECT().GetEvent(gomock.Eq(ctxWithTrace), gomock.Eq(test_event.Event_uuid)).Return(test_event, nil)

			result, err := regUC.ProcessingSagaAndEvents(ctx, uuid.Nil, test_event.Event_uuid, true)
			require.Nil(t, err)
			require.Nil(t, result)

			test_event.Event_status = old_status

		})
	})

	t.Run("Sagas", func(t *testing.T) {
		// TODO: complete
	})
}
