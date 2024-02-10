package usecase

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/GCFactory/dbo-system/service/totp/internal/totp/mock"
	totpRepo "github.com/GCFactory/dbo-system/service/totp/internal/totp/repository"
	totpErrors "github.com/GCFactory/dbo-system/service/totp/pkg/errors"
	otpPkg "github.com/GCFactory/dbo-system/service/totp/pkg/otp"
	totpPkgConfig "github.com/GCFactory/dbo-system/service/totp/pkg/otp/config"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
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

func TestTotpUC_Enroll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfg)
	apiLogger.InitLogger()
	mockRepo := mock.NewMockRepository(ctrl)
	mockTotp := mock.NewMockTotp(ctrl)
	totpUC := NewTOTPUseCase(testCfg, mockRepo, mockTotp, apiLogger)

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
	userId := uuid.New()
	userName := "Rueie"
	issuer := "dbo.gcfactory.space"
	t.Run("Valid", func(t *testing.T) {
		inputCfg := models.TOTPConfig{
			UserId:      userId,
			AccountName: userName,
		}
		genOpts := totpPkgConfig.GenerateOpts{
			Issuer:      issuer,
			AccountName: userName,
			SecretSize:  totpPkgConfig.DefaultSecretLength,
			Algorithm:   totpPkgConfig.DefaultAlgorithm,
		}
		tmp := otpPkg.TotpStruct{}
		secret, url, err := tmp.Generate(genOpts)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, secret)
		require.NotNil(t, url)

		ctx := context.Background()
		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enroll")
		defer span.Finish()

		mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(nil, totpRepo.ErrorGetActiveConfig)
		mockTotp.EXPECT().Generate(gomock.Eq(genOpts)).Return(secret, url, nil)
		mockRepo.EXPECT().CreateConfig(ctxWithTrace, gomock.Any()).Return(nil)

		result, err := totpUC.Enroll(ctx, inputCfg)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotEqual(t, result.TotpId, "")
		require.Equal(t, *secret, result.TotpSecret)
		require.Equal(t, *url, result.TotpUrl)
	})
	t.Run("UserIsActive", func(t *testing.T) {
		totpCfg := models.TOTPConfig{
			UserId:      userId,
			AccountName: userName,
			IsActive:    true,
		}

		ctx := context.Background()
		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enroll")
		defer span.Finish()

		mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(&totpCfg, nil)

		result, err := totpUC.Enroll(ctx, models.TOTPConfig{UserId: userId, AccountName: userName})
		require.NotNil(t, err)
		require.Error(t, err)
		require.Equal(t, err, totpErrors.ActiveTotp)
		require.Nil(t, result)
	})
	t.Run("NoUserName", func(t *testing.T) {
		ctx := context.Background()
		result, err := totpUC.Enroll(ctx, models.TOTPConfig{UserId: userId})
		require.NotNil(t, err)
		require.Error(t, err)
		require.Equal(t, err, totpErrors.NoUserName)
		require.Nil(t, result)
	})
	t.Run("EmptyUserId", func(t *testing.T) {
		ctx := context.Background()
		result, err := totpUC.Enroll(ctx, models.TOTPConfig{AccountName: userName})
		require.NotNil(t, err)
		require.Error(t, err)
		require.Equal(t, err, totpErrors.EmptyUserId)
		require.Nil(t, result)
	})
	t.Run("ErrorGenTotp", func(t *testing.T) {
		inputCfg := models.TOTPConfig{
			UserId:      userId,
			AccountName: userName,
		}
		genOpts := totpPkgConfig.GenerateOpts{
			Issuer:      issuer,
			AccountName: userName,
			SecretSize:  totpPkgConfig.DefaultSecretLength,
			Algorithm:   totpPkgConfig.DefaultAlgorithm,
		}
		tmp := otpPkg.TotpStruct{}
		secret, url, err := tmp.Generate(genOpts)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, secret)
		require.NotNil(t, url)

		ctx := context.Background()
		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enroll")
		defer span.Finish()

		mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(nil, totpRepo.ErrorGetActiveConfig)
		mockTotp.EXPECT().Generate(gomock.Eq(genOpts)).Return(nil, nil, totpPkgConfig.ErrGenerateMissingIssuer)

		result, err := totpUC.Enroll(ctx, inputCfg)
		require.NotNil(t, err)
		require.Error(t, err)
		require.Equal(t, err, ErrorGenTotp)
		require.Nil(t, result)
	})
	t.Run("ErrorCreateConfig", func(t *testing.T) {
		inputCfg := models.TOTPConfig{
			UserId:      userId,
			AccountName: userName,
		}
		genOpts := totpPkgConfig.GenerateOpts{
			Issuer:      issuer,
			AccountName: userName,
			SecretSize:  totpPkgConfig.DefaultSecretLength,
			Algorithm:   totpPkgConfig.DefaultAlgorithm,
		}
		tmp := otpPkg.TotpStruct{}
		secret, url, err := tmp.Generate(genOpts)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, secret)
		require.NotNil(t, url)

		ctx := context.Background()
		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enroll")
		defer span.Finish()

		mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(nil, totpRepo.ErrorGetActiveConfig)
		mockTotp.EXPECT().Generate(gomock.Eq(genOpts)).Return(secret, url, nil)
		mockRepo.EXPECT().CreateConfig(ctxWithTrace, gomock.Any()).Return(totpRepo.ErrorCreateConfig)

		result, err := totpUC.Enroll(ctx, inputCfg)
		require.NotNil(t, err)
		require.Error(t, err)
		require.Equal(t, err, ErrorCreateConfig)
		require.Nil(t, result)
	})
}

func TestTotpUC_Verify(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfg)
	apiLogger.InitLogger()
	mockRepo := mock.NewMockRepository(ctrl)
	mockTotp := mock.NewMockTotp(ctrl)
	totpUC := NewTOTPUseCase(testCfg, mockRepo, mockTotp, apiLogger)

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
	t.Run("Valid", func(t *testing.T) {
		ctx := context.Background()
		url := "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&digits=6&issuer=dbo.gcfactory.space&period=30&secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
		result, err := totpUC.Verify(ctx, url)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, result, &models.TOTPVerify{Status: "OK"})
	})
	t.Run("NoAlgorithm", func(t *testing.T) {
		ctx := context.Background()
		url := "otpauth://totp/dbo.gcfactory.space:admin?digits=6&issuer=dbo.gcfactory.space&period=30&secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
		result, err := totpUC.Verify(ctx, url)
		require.NotNil(t, err)
		require.Error(t, err)
		require.Equal(t, err, totpErrors.NoAlgorithmField)
		require.NotNil(t, result)
		require.Equal(t, result, &models.TOTPVerify{Status: err.Error()})
	})
	t.Run("NoDigits", func(t *testing.T) {
		ctx := context.Background()
		url := "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&issuer=dbo.gcfactory.space&period=30&secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
		result, err := totpUC.Verify(ctx, url)
		require.NotNil(t, err)
		require.Error(t, err)
		require.Equal(t, err, totpErrors.NoDigitsField)
		require.NotNil(t, result)
		require.Equal(t, result, &models.TOTPVerify{Status: err.Error()})
	})
	t.Run("NoIssuer", func(t *testing.T) {
		ctx := context.Background()
		url := "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&digits=6&period=30&secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
		result, err := totpUC.Verify(ctx, url)
		require.NotNil(t, err)
		require.Error(t, err)
		require.Equal(t, err, totpErrors.NoIssuerField)
		require.NotNil(t, result)
		require.Equal(t, result, &models.TOTPVerify{Status: err.Error()})
	})
	t.Run("NoPeriod", func(t *testing.T) {
		ctx := context.Background()
		url := "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&digits=6&issuer=dbo.gcfactory.space&secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
		result, err := totpUC.Verify(ctx, url)
		require.NotNil(t, err)
		require.Error(t, err)
		require.Equal(t, err, totpErrors.NoPeriodField)
		require.NotNil(t, result)
		require.Equal(t, result, &models.TOTPVerify{Status: err.Error()})
	})
	t.Run("NoSecret", func(t *testing.T) {
		ctx := context.Background()
		url := "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&digits=6&issuer=dbo.gcfactory.space&period=30"
		result, err := totpUC.Verify(ctx, url)
		require.NotNil(t, err)
		require.Error(t, err)
		require.Equal(t, err, totpErrors.NoSecretField)
		require.NotNil(t, result)
		require.Equal(t, result, &models.TOTPVerify{Status: err.Error()})
	})
}

func TestTotpUC_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfg)
	apiLogger.InitLogger()
	mockRepo := mock.NewMockRepository(ctrl)
	mockTotp := mock.NewMockTotp(ctrl)
	totpUC := NewTOTPUseCase(testCfg, mockRepo, mockTotp, apiLogger)

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
	userId := uuid.New()
	t.Run("Valid", func(t *testing.T) {
		secret := "JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
		time := time.Now()
		code, err := otpPkg.TotpStruct{}.GenerateCode(secret, time)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, code)
		totpCfg := models.TOTPConfig{
			UserId: userId,
			URL:    "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&digits=6&issuer=dbo.gcfactory.space&period=30&secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP",
			Secret: secret,
		}
		validOpts := totpPkgConfig.ValidateOpts{
			Digits:    totpPkgConfig.DefaultDigits,
			Algorithm: totpPkgConfig.DefaultAlgorithm,
			Period:    totpPkgConfig.DefaultPeriod,
		}

		ctx := context.Background()
		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Validate")
		defer span.Finish()

		mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(&totpCfg, nil)
		mockTotp.EXPECT().GenerateCodeCustom(gomock.Eq(secret), gomock.Eq(time), gomock.Eq(validOpts)).Return(code, nil)

		result, err := totpUC.Validate(ctx, userId, code, time)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, result, &models.TOTPValidate{Status: "OK"})
	})
	t.Run("NotFound", func(t *testing.T) {
		ctx := context.Background()
		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Validate")
		defer span.Finish()

		mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(nil, totpRepo.ErrorGetActiveConfig)

		result, err := totpUC.Validate(ctx, userId, "", time.Now())
		require.NotNil(t, err)
		require.Error(t, err)
		require.Equal(t, err, totpErrors.NoUserId)
		require.NotNil(t, result)
		require.Equal(t, result, &models.TOTPValidate{Status: err.Error()})
	})
	t.Run("ErrorGenerateCode", func(t *testing.T) {
		secret := "JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
		time := time.Now()
		totpCfg := models.TOTPConfig{
			UserId: userId,
			URL:    "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&digits=6&issuer=dbo.gcfactory.space&period=30&secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP",
			Secret: secret,
		}
		validOpts := totpPkgConfig.ValidateOpts{
			Digits:    totpPkgConfig.DefaultDigits,
			Algorithm: totpPkgConfig.DefaultAlgorithm,
			Period:    totpPkgConfig.DefaultPeriod,
		}

		ctx := context.Background()
		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Validate")
		defer span.Finish()

		mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(&totpCfg, nil)
		mockTotp.EXPECT().GenerateCodeCustom(gomock.Eq(secret), gomock.Eq(time), gomock.Eq(validOpts)).Return("", totpPkgConfig.ErrValidateSecretInvalidBase32)

		result, err := totpUC.Validate(ctx, userId, "", time)
		require.NotNil(t, err)
		require.Error(t, err)
		require.Equal(t, err, ErrorGenCodeCustom)
		require.Nil(t, result)
	})
	t.Run("WrongTotp", func(t *testing.T) {
		secret := "JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
		wrongSecret := "JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPA"
		time := time.Now()
		code, err := otpPkg.TotpStruct{}.GenerateCode(secret, time)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, code)
		wrongCode, err := otpPkg.TotpStruct{}.GenerateCode(wrongSecret, time)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, wrongCode)
		totpCfg := models.TOTPConfig{
			UserId: userId,
			URL:    "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&digits=6&issuer=dbo.gcfactory.space&period=30&secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP",
			Secret: secret,
		}
		validOpts := totpPkgConfig.ValidateOpts{
			Digits:    totpPkgConfig.DefaultDigits,
			Algorithm: totpPkgConfig.DefaultAlgorithm,
			Period:    totpPkgConfig.DefaultPeriod,
		}

		ctx := context.Background()
		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Validate")
		defer span.Finish()

		mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(&totpCfg, nil)
		mockTotp.EXPECT().GenerateCodeCustom(gomock.Eq(secret), gomock.Eq(time), gomock.Eq(validOpts)).Return(code, nil)

		result, err := totpUC.Validate(ctx, userId, wrongCode, time)
		require.NotNil(t, err)
		require.Error(t, err)
		require.Equal(t, err, totpErrors.WrongTotpCode)
		require.NotNil(t, result)
		require.Equal(t, result, &models.TOTPValidate{Status: err.Error()})
	})
	t.Run("OnlySecretInUrl", func(t *testing.T) {
		secret := "JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
		time := time.Now()
		code, err := otpPkg.TotpStruct{}.GenerateCode(secret, time)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, code)
		totpCfg := models.TOTPConfig{
			UserId: userId,
			URL:    "otpauth://totp/dbo.gcfactory.space:admin?secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP",
			Secret: secret,
		}
		validOpts := totpPkgConfig.ValidateOpts{
			Digits:    totpPkgConfig.DefaultDigits,
			Algorithm: totpPkgConfig.DefaultAlgorithm,
			Period:    totpPkgConfig.DefaultPeriod,
		}

		ctx := context.Background()
		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Validate")
		defer span.Finish()

		mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(&totpCfg, nil)
		mockTotp.EXPECT().GenerateCodeCustom(gomock.Eq(secret), gomock.Eq(time), gomock.Eq(validOpts)).Return(code, nil)

		result, err := totpUC.Validate(ctx, userId, code, time)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, result, &models.TOTPValidate{Status: "OK"})
	})
	t.Run("UseMD5", func(t *testing.T) {
		secret := "JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
		time := time.Now()
		code, err := otpPkg.TotpStruct{}.GenerateCode(secret, time)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, code)
		totpCfg := models.TOTPConfig{
			UserId: userId,
			URL:    "otpauth://totp/dbo.gcfactory.space:admin?secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP&algorithm=MD5",
			Secret: secret,
		}
		validOpts := totpPkgConfig.ValidateOpts{
			Digits:    totpPkgConfig.DefaultDigits,
			Algorithm: totpPkgConfig.AlgorithmMD5,
			Period:    totpPkgConfig.DefaultPeriod,
		}

		ctx := context.Background()
		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Validate")
		defer span.Finish()

		mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(&totpCfg, nil)
		mockTotp.EXPECT().GenerateCodeCustom(gomock.Eq(secret), gomock.Eq(time), gomock.Eq(validOpts)).Return(code, nil)

		result, err := totpUC.Validate(ctx, userId, code, time)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, result, &models.TOTPValidate{Status: "OK"})
	})
	t.Run("UseSHA256", func(t *testing.T) {
		secret := "JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
		time := time.Now()
		code, err := otpPkg.TotpStruct{}.GenerateCode(secret, time)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, code)
		totpCfg := models.TOTPConfig{
			UserId: userId,
			URL:    "otpauth://totp/dbo.gcfactory.space:admin?secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP&algorithm=SHA256",
			Secret: secret,
		}
		validOpts := totpPkgConfig.ValidateOpts{
			Digits:    totpPkgConfig.DefaultDigits,
			Algorithm: totpPkgConfig.AlgorithmSHA256,
			Period:    totpPkgConfig.DefaultPeriod,
		}

		ctx := context.Background()
		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Validate")
		defer span.Finish()

		mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(&totpCfg, nil)
		mockTotp.EXPECT().GenerateCodeCustom(gomock.Eq(secret), gomock.Eq(time), gomock.Eq(validOpts)).Return(code, nil)

		result, err := totpUC.Validate(ctx, userId, code, time)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, result, &models.TOTPValidate{Status: "OK"})
	})
	t.Run("UseSHA512", func(t *testing.T) {
		secret := "JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
		time := time.Now()
		code, err := otpPkg.TotpStruct{}.GenerateCode(secret, time)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, code)
		totpCfg := models.TOTPConfig{
			UserId: userId,
			URL:    "otpauth://totp/dbo.gcfactory.space:admin?secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP&algorithm=SHA512",
			Secret: secret,
		}
		validOpts := totpPkgConfig.ValidateOpts{
			Digits:    totpPkgConfig.DefaultDigits,
			Algorithm: totpPkgConfig.AlgorithmSHA512,
			Period:    totpPkgConfig.DefaultPeriod,
		}

		ctx := context.Background()
		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Validate")
		defer span.Finish()

		mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(&totpCfg, nil)
		mockTotp.EXPECT().GenerateCodeCustom(gomock.Eq(secret), gomock.Eq(time), gomock.Eq(validOpts)).Return(code, nil)

		result, err := totpUC.Validate(ctx, userId, code, time)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, result, &models.TOTPValidate{Status: "OK"})
	})
}

func TestTotpUC_Enable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfg)
	apiLogger.InitLogger()
	mockRepo := mock.NewMockRepository(ctrl)
	mockTotp := mock.NewMockTotp(ctrl)
	totpUC := NewTOTPUseCase(testCfg, mockRepo, mockTotp, apiLogger)

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
	t.Run("By totp id", func(t *testing.T) {
		totpId := uuid.New()
		t.Run("Valid", func(t *testing.T) {
			totpCfg := models.TOTPConfig{
				Id:       totpId,
				IsActive: false,
				UserId:   uuid.New(),
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(&totpCfg, nil)
			mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(totpCfg.UserId)).Return(nil, totpRepo.ErrorGetActiveConfig)
			mockRepo.EXPECT().UpdateTotpActivityByTotpId(ctxWithTrace, gomock.Eq(totpId), gomock.Eq(true)).Return(nil)

			result, err := totpUC.Enable(ctx, totpId, uuid.Nil)
			require.Nil(t, err)
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, result, &models.TOTPEnable{Status: "OK"})
		})
		t.Run("Not found", func(t *testing.T) {
			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(nil, totpRepo.ErrorGetConfigByTotpId)

			result, err := totpUC.Enable(ctx, totpId, uuid.Nil)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, totpErrors.NoTotpId)
			require.NotNil(t, result)
			require.Equal(t, result, &models.TOTPEnable{Status: err.Error()})
		})
		t.Run("Is active", func(t *testing.T) {
			totpCfg := models.TOTPConfig{
				Id:       totpId,
				IsActive: true,
			}
			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(&totpCfg, nil)

			result, err := totpUC.Enable(ctx, totpId, uuid.Nil)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, totpErrors.TotpIsActive)
			require.NotNil(t, result)
			require.Equal(t, result, &models.TOTPEnable{Status: err.Error()})
		})
		t.Run("Has another active config", func(t *testing.T) {
			totpCfg := models.TOTPConfig{
				Id:       totpId,
				UserId:   uuid.New(),
				IsActive: false,
			}
			activeCfg := models.TOTPConfig{
				Id:       uuid.New(),
				UserId:   totpCfg.UserId,
				IsActive: true,
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(&totpCfg, nil)
			mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(totpCfg.UserId)).Return(&activeCfg, nil)

			result, err := totpUC.Enable(ctx, totpId, uuid.Nil)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, totpErrors.TotpIsActive)
			require.NotNil(t, result)
			require.Equal(t, result, &models.TOTPEnable{Status: err.Error()})
		})
		t.Run("ErrorUpdateRepo", func(t *testing.T) {
			totpCfg := models.TOTPConfig{
				Id:       totpId,
				UserId:   uuid.New(),
				IsActive: false,
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(&totpCfg, nil)
			mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(totpCfg.UserId)).Return(nil, totpRepo.ErrorGetActiveConfig)
			mockRepo.EXPECT().UpdateTotpActivityByTotpId(ctxWithTrace, gomock.Eq(totpId), gomock.Eq(true)).Return(totpRepo.ErrorUpdateTotpActivityByTotpId)

			result, err := totpUC.Enable(ctx, totpId, uuid.Nil)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, ErrorUpdateActivityByTotpId)
			require.Nil(t, result)
		})
	})
	t.Run("By user id", func(t *testing.T) {
		userId := uuid.New()
		t.Run("Valid", func(t *testing.T) {
			totpCfg := models.TOTPConfig{
				UserId: userId,
			}
			disabledCfg := models.TOTPConfig{
				UserId:   userId,
				Id:       uuid.New(),
				IsActive: false,
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByUserId(ctxWithTrace, gomock.Eq(userId)).Return(&totpCfg, nil)
			mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(nil, totpRepo.ErrorGetActiveConfig)
			mockRepo.EXPECT().GetLastDisabledConfig(ctxWithTrace, gomock.Eq(userId)).Return(&disabledCfg, nil)
			mockRepo.EXPECT().UpdateTotpActivityByTotpId(ctxWithTrace, gomock.Eq(disabledCfg.Id), gomock.Eq(true)).Return(nil)

			result, err := totpUC.Enable(ctx, uuid.Nil, userId)
			require.Nil(t, err)
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, result, &models.TOTPEnable{Status: "OK"})
		})
		t.Run("Not found", func(t *testing.T) {
			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByUserId(ctxWithTrace, gomock.Eq(userId)).Return(nil, totpRepo.ErrorGetConfiByUserId)

			result, err := totpUC.Enable(ctx, uuid.Nil, userId)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, totpErrors.NoUserId)
			require.NotNil(t, result)
			require.Equal(t, result, &models.TOTPEnable{Status: err.Error()})
		})
		t.Run("Is active", func(t *testing.T) {
			totpCfg := models.TOTPConfig{
				UserId: userId,
			}
			activeCfg := models.TOTPConfig{
				UserId:   userId,
				IsActive: true,
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByUserId(ctxWithTrace, gomock.Eq(userId)).Return(&totpCfg, nil)
			mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(&activeCfg, nil)

			result, err := totpUC.Enable(ctx, uuid.Nil, userId)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, totpErrors.TotpIsActive)
			require.NotNil(t, result)
			require.Equal(t, result, &models.TOTPEnable{Status: err.Error()})
		})
		t.Run("No disabled totp", func(t *testing.T) {
			totpCfg := models.TOTPConfig{
				UserId: userId,
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByUserId(ctxWithTrace, gomock.Eq(userId)).Return(&totpCfg, nil)
			mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(nil, totpRepo.ErrorGetActiveConfig)
			mockRepo.EXPECT().GetLastDisabledConfig(ctxWithTrace, gomock.Eq(userId)).Return(nil, totpRepo.ErrorGetLastDisabledConfig)

			result, err := totpUC.Enable(ctx, uuid.Nil, userId)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, totpErrors.TotpIsActive)
			require.NotNil(t, result)
			require.Equal(t, result, &models.TOTPEnable{Status: err.Error()})
		})
		t.Run("ErrorUpdateRepo", func(t *testing.T) {
			totpCfg := models.TOTPConfig{
				UserId: userId,
			}
			disabledCfg := models.TOTPConfig{
				UserId:   userId,
				Id:       uuid.New(),
				IsActive: false,
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByUserId(ctxWithTrace, gomock.Eq(userId)).Return(&totpCfg, nil)
			mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(nil, totpRepo.ErrorGetActiveConfig)
			mockRepo.EXPECT().GetLastDisabledConfig(ctxWithTrace, gomock.Eq(userId)).Return(&disabledCfg, nil)
			mockRepo.EXPECT().UpdateTotpActivityByTotpId(ctxWithTrace, gomock.Eq(disabledCfg.Id), gomock.Eq(true)).Return(totpRepo.ErrorUpdateTotpActivityByTotpId)

			result, err := totpUC.Enable(ctx, uuid.Nil, userId)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, ErrorUpdateActivityByTotpId)
			require.Nil(t, result)
		})
	})
	t.Run("Both id", func(t *testing.T) {
		userId := uuid.New()
		totpId := uuid.New()
		t.Run("Valid", func(t *testing.T) {
			totpCfg := models.TOTPConfig{
				UserId:   userId,
				Id:       totpId,
				IsActive: false,
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(&totpCfg, nil)
			mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(nil, totpRepo.ErrorGetActiveConfig)
			mockRepo.EXPECT().UpdateTotpActivityByTotpId(ctxWithTrace, gomock.Eq(totpId), gomock.Eq(true)).Return(nil)

			result, err := totpUC.Enable(ctx, totpId, userId)
			require.Nil(t, err)
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, result, &models.TOTPEnable{Status: "OK"})
		})
		t.Run("Not found", func(t *testing.T) {
			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(nil, totpRepo.ErrorGetConfigByTotpId)

			result, err := totpUC.Enable(ctx, totpId, userId)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, totpErrors.NoId)
			require.NotNil(t, result)
			require.Equal(t, result, &models.TOTPEnable{Status: err.Error()})
		})
		t.Run("Incorrect ids", func(t *testing.T) {
			totpCfg := models.TOTPConfig{
				UserId: uuid.New(),
				Id:     totpId,
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(&totpCfg, nil)

			result, err := totpUC.Enable(ctx, totpId, userId)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, totpErrors.NoId)
			require.NotNil(t, result)
			require.Equal(t, result, &models.TOTPEnable{Status: err.Error()})
		})
	})
	t.Run("No id", func(t *testing.T) {
		ctx := context.Background()
		result, err := totpUC.Enable(ctx, uuid.Nil, uuid.Nil)
		require.NotNil(t, err)
		require.Error(t, err)
		require.Equal(t, err, totpErrors.NoId)
		require.NotNil(t, result)
		require.Equal(t, result, &models.TOTPEnable{Status: err.Error()})
	})
}

func TestTotpUC_Disable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfg)
	apiLogger.InitLogger()
	mockRepo := mock.NewMockRepository(ctrl)
	mockTotp := mock.NewMockTotp(ctrl)
	totpUC := NewTOTPUseCase(testCfg, mockRepo, mockTotp, apiLogger)

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
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Disable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(&configTOTP, nil)
			mockRepo.EXPECT().UpdateTotpActivityByTotpId(ctxWithTrace, gomock.Eq(totpId), gomock.Eq(false)).Return(nil)

			usecaseResult, err := totpUC.Disable(ctx, totpId, uuid.Nil)
			require.NoError(t, err)
			require.Nil(t, err)
			require.NotNil(t, usecaseResult)
			require.Equal(t, usecaseResult, &models.TOTPDisable{Status: "OK"})
		})
		t.Run("NotFound", func(t *testing.T) {
			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Disable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(nil, totpRepo.ErrorGetConfigByTotpId)

			usecaseResult, err := totpUC.Disable(ctx, totpId, uuid.Nil)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, totpErrors.NoTotpId)
			require.NotNil(t, usecaseResult)
			require.Equal(t, usecaseResult, &models.TOTPDisable{Status: totpErrors.NoTotpId.Error()})

		})
		t.Run("IsDisabled", func(t *testing.T) {
			configTOTP := models.TOTPConfig{
				IsActive: false,
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Disable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(&configTOTP, nil)

			usecaseResult, err := totpUC.Disable(ctx, totpId, uuid.Nil)
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

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Disable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(&configTOTP, nil)
			mockRepo.EXPECT().UpdateTotpActivityByTotpId(ctxWithTrace, gomock.Eq(totpId), gomock.Eq(false)).Return(totpRepo.ErrorUpdateTotpActivityByTotpId)

			usecaseResult, err := totpUC.Disable(ctx, totpId, uuid.Nil)
			require.NotNil(t, err)
			require.Nil(t, usecaseResult)
			require.Error(t, err)
			require.Equal(t, err, ErrorUpdateActivityByTotpId)
		})
	})

	t.Run("ByUserID", func(t *testing.T) {
		userId := uuid.New()
		t.Run("Valid", func(t *testing.T) {
			userConfigTOTP := models.TOTPConfig{
				UserId: userId,
			}
			userActiveCfg := models.TOTPConfig{
				UserId:   userId,
				IsActive: true,
				Id:       uuid.New(),
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Disable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByUserId(ctxWithTrace, gomock.Eq(userId)).Return(&userConfigTOTP, nil)
			mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(&userActiveCfg, nil)
			mockRepo.EXPECT().UpdateTotpActivityByTotpId(ctxWithTrace, gomock.Eq(userActiveCfg.Id), gomock.Eq(false)).Return(nil)

			usecaseResult, err := totpUC.Disable(ctx, uuid.Nil, userId)
			require.NoError(t, err)
			require.Nil(t, err)
			require.NotNil(t, usecaseResult)
			require.Equal(t, usecaseResult, &models.TOTPDisable{Status: "OK"})
		})
		t.Run("NotFound", func(t *testing.T) {
			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Disable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByUserId(ctxWithTrace, gomock.Eq(userId)).Return(nil, totpRepo.ErrorGetConfiByUserId)

			usecaseResult, err := totpUC.Disable(ctx, uuid.Nil, userId)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, totpErrors.NoTotpId)
			require.NotNil(t, usecaseResult)
			require.Equal(t, usecaseResult, &models.TOTPDisable{Status: totpErrors.NoTotpId.Error()})

		})
		t.Run("IsDisabled", func(t *testing.T) {
			userConfigTOTP := models.TOTPConfig{
				UserId: userId,
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Disable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByUserId(ctxWithTrace, gomock.Eq(userId)).Return(&userConfigTOTP, nil)
			mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(nil, totpRepo.ErrorGetActiveConfig)

			result, err := totpUC.Disable(ctx, uuid.Nil, userId)
			require.NotNil(t, result)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, totpErrors.TotpIsDisabled)
			require.Equal(t, result, &models.TOTPDisable{Status: err.Error()})
		})
		t.Run("FailedToUpdateRepo", func(t *testing.T) {
			userConfigTOTP := models.TOTPConfig{
				UserId: userId,
			}
			userActiveCfg := models.TOTPConfig{
				UserId:   userId,
				IsActive: true,
				Id:       uuid.New(),
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Disable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByUserId(ctxWithTrace, gomock.Eq(userId)).Return(&userConfigTOTP, nil)
			mockRepo.EXPECT().GetActiveConfig(ctxWithTrace, gomock.Eq(userId)).Return(&userActiveCfg, nil)
			mockRepo.EXPECT().UpdateTotpActivityByTotpId(ctxWithTrace, gomock.Eq(userActiveCfg.Id), gomock.Eq(false)).Return(totpRepo.ErrorUpdateTotpActivityByTotpId)

			result, err := totpUC.Disable(ctx, uuid.Nil, userId)
			require.Nil(t, result)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, ErrorUpdateActivityByTotpId)
		})
	})

	t.Run("No Id", func(t *testing.T) {
		ctx := context.Background()

		result, err := totpUC.Disable(ctx, uuid.Nil, uuid.Nil)
		require.NotNil(t, result)
		require.NotNil(t, err)
		require.Error(t, err)
		require.Equal(t, err, totpErrors.NoId)
		require.Equal(t, result, &models.TOTPDisable{Status: totpErrors.NoId.Error()})
	})

	t.Run("Both id", func(t *testing.T) {
		userId := uuid.New()
		totpId := uuid.New()
		t.Run("Valid", func(t *testing.T) {
			totpCfg := models.TOTPConfig{
				UserId:   userId,
				Id:       totpId,
				IsActive: true,
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Disable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(&totpCfg, nil)
			mockRepo.EXPECT().UpdateTotpActivityByTotpId(ctxWithTrace, gomock.Eq(totpId), gomock.Eq(false)).Return(nil)

			result, err := totpUC.Disable(ctx, totpId, userId)
			require.NoError(t, err)
			require.Nil(t, err)
			require.NotNil(t, result)
			require.Equal(t, result, &models.TOTPDisable{Status: "OK"})
		})
		t.Run("Not found", func(t *testing.T) {
			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Disable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(nil, totpRepo.ErrorGetConfigByTotpId)

			result, err := totpUC.Disable(ctx, totpId, userId)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, totpErrors.NoId)
			require.NotNil(t, result)
			require.Equal(t, result, &models.TOTPDisable{Status: totpErrors.NoId.Error()})
		})
		t.Run("IsDisabled", func(t *testing.T) {
			totpCfg := models.TOTPConfig{
				UserId:   userId,
				Id:       totpId,
				IsActive: false,
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Disable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(&totpCfg, nil)

			result, err := totpUC.Disable(ctx, totpId, userId)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, totpErrors.TotpIsDisabled)
			require.NotNil(t, result)
			require.Equal(t, result, &models.TOTPDisable{Status: err.Error()})
		})
		t.Run("FailedToUpdateRepo", func(t *testing.T) {
			totpCfg := models.TOTPConfig{
				UserId:   userId,
				Id:       totpId,
				IsActive: true,
			}

			ctx := context.Background()
			span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Disable")
			defer span.Finish()

			mockRepo.EXPECT().GetConfigByTotpId(ctxWithTrace, gomock.Eq(totpId)).Return(&totpCfg, nil)
			mockRepo.EXPECT().UpdateTotpActivityByTotpId(ctxWithTrace, gomock.Eq(totpId), gomock.Eq(false)).Return(ErrorUpdateActivityByTotpId)

			result, err := totpUC.Disable(ctx, totpId, userId)
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, ErrorUpdateActivityByTotpId)
			require.Nil(t, result)
		})
	})
}
