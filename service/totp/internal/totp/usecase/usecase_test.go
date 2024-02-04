package usecase

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/httpErrors"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/GCFactory/dbo-system/service/totp/internal/totp"
	"github.com/GCFactory/dbo-system/service/totp/internal/totp/mock"
	totpErrors "github.com/GCFactory/dbo-system/service/totp/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var (
	testCfg = &config.Config{
		Env: "Development",
		Logger: config.Logger{
			Development: true,
			Level:       "Debug",
		},
	}
)

func TestVerify(t *testing.T) {
	url := "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&digits=6&issuer=dbo.gcfactory.space&period=30&secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
	var cfg config.Config
	var log logger.Logger
	var repo totp.Repository
	totpUc := NewTOTPUseCase(&cfg, repo, log)
	ctx := context.TODO()
	result, err := totpUc.Verify(ctx, url)
	require.NoError(t, err)
	require.Equal(t, &models.TOTPVerify{Status: "OK"}, result, "Should be correct check")
	url = "otpauth://totp/dbo.gcfactory.space:admin?digits=6&issuer=dbo.gcfactory.space&period=30&secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
	result, err = totpUc.Verify(ctx, url)
	require.Error(t, err)
	require.Equal(t, err, totpErrors.NoAlgorithmField)
	require.Equal(t, result, &models.TOTPVerify{Status: err.Error()})
	url = "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&issuer=dbo.gcfactory.space&period=30&secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
	result, err = totpUc.Verify(ctx, url)
	require.Error(t, err)
	require.Equal(t, err, totpErrors.NoDigitsField)
	require.Equal(t, result, &models.TOTPVerify{Status: err.Error()})
	url = "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&digits=6&period=30&secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
	result, err = totpUc.Verify(ctx, url)
	require.Error(t, err)
	require.Equal(t, err, totpErrors.NoIssuerField)
	require.Equal(t, result, &models.TOTPVerify{Status: err.Error()})
	url = "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&digits=6&issuer=dbo.gcfactory.space&secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
	result, err = totpUc.Verify(ctx, url)
	require.Error(t, err)
	require.Equal(t, err, totpErrors.NoPeriodField)
	require.Equal(t, result, &models.TOTPVerify{Status: err.Error()})
	url = "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&digits=6&issuer=dbo.gcfactory.space&period=30"
	result, err = totpUc.Verify(ctx, url)
	require.Error(t, err)
	require.Equal(t, err, totpErrors.NoSecretField)
	require.Equal(t, result, &models.TOTPVerify{Status: err.Error()})
}

func TestTotpUC_Disable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfg)
	apiLogger.InitLogger()
	mockRepo := mock.NewMockRepository(ctrl)
	totpUC := NewTOTPUseCase(testCfg, mockRepo, apiLogger)

	t.Parallel()
	// Supress output
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	os.Stdout = null
	os.Stderr = null
	defer func() {
		defer null.Close()
		os.Stdout = sout
		os.Stderr = serr
	}()
	t.Run("ByTotpID", func(t *testing.T) {
		totpId := uuid.New()
		t.Run("Valid", func(t *testing.T) {
			configTOTP := models.TOTPConfig{
				IsActive: true,
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authUC.UploadAvatar")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(&configTOTP, nil)
			mockRepo.EXPECT().UpdateTotpActivityByTotpId(ctxWithTrace, gomock.Eq(totpId), gomock.Eq(false)).Return(nil)

			usecaseResult, err := totpUC.Disable(ctxWithTrace, totpId, uuid.Nil)
			require.NoError(t, err)
			require.Nil(t, err)
			require.NotNil(t, usecaseResult)
			require.Equal(t, usecaseResult, &models.TOTPDisable{Status: "OK"})
		})
		t.Run("NotFound", func(t *testing.T) {
			configTOTP := models.TOTPConfig{
				IsActive: true,
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authUC.UploadAvatar")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(&configTOTP, totpErrors.NoTotpId)

			usecaseResult, err := totpUC.Disable(ctxWithTrace, totpId, uuid.Nil)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, totpErrors.NoTotpId)
			require.NotNil(t, usecaseResult)
			require.Equal(t, usecaseResult, &models.TOTPDisable{Status: totpErrors.NoTotpId.Error()})

		})
		t.Run("InActive", func(t *testing.T) {
			configTOTP := models.TOTPConfig{
				IsActive: false,
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authUC.UploadAvatar")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(&configTOTP, nil)

			usecaseResult, err := totpUC.Disable(ctxWithTrace, totpId, uuid.Nil)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, totpErrors.TotpIsDisabled)
			require.NotNil(t, usecaseResult)
			require.Equal(t, usecaseResult, &models.TOTPDisable{Status: totpErrors.TotpIsDisabled.Error()})
		})
		t.Run("FailedToUpdateRepo", func(t *testing.T) {
			configTOTP := models.TOTPConfig{
				IsActive: true,
			}
			expectedRepoError := errors.New("")
			expectedUsecaseError := httpErrors.NewInternalServerError(errors.Wrap(expectedRepoError, "totpUC.Disable.totpRepo.UpdateTotpActivityByTotpId(totpId)"))

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authUC.UploadAvatar")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(&configTOTP, nil)
			mockRepo.EXPECT().UpdateTotpActivityByTotpId(ctxWithTrace, gomock.Eq(totpId), gomock.Eq(false)).Return(expectedRepoError)

			usecaseResult, err := totpUC.Disable(ctxWithTrace, totpId, uuid.Nil)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err.Error(), expectedUsecaseError.Error())
			require.Nil(t, usecaseResult)
		})
	})

	t.Run("ByUserID", func(t *testing.T) {
		//userId := uuid.New()
	})
}
