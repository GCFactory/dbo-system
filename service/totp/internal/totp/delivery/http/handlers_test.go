package http

import (
	"encoding/json"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/converter"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/GCFactory/dbo-system/service/totp/internal/totp/mock"
	totpRepo "github.com/GCFactory/dbo-system/service/totp/internal/totp/repository"
	totpUseCase "github.com/GCFactory/dbo-system/service/totp/internal/totp/usecase"
	totpErrors "github.com/GCFactory/dbo-system/service/totp/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

func TestTOTPHandlers_Enroll(t *testing.T) {
	t.Parallel()
	type inputStruct struct {
		UserName string `json:"user_name"`
		UserId   string `json:"user_id"`
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLoggger := logger.NewServerLogger(testCfg)
	apiLoggger.InitLogger()

	mockUC := mock.NewMockUseCase(ctrl)
	tHandlers := NewTOTPHandlers(testCfg, mockUC, apiLoggger)

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
	userName := "Rueie"
	t.Run("EmptyBody", func(t *testing.T) {
		inputBody := ""

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enroll", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handlerFunc := tHandlers.Enroll()

		err = handlerFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Equal(t, rec.Body.String(), "{\"status\":500,\"error\":\"Internal Server Error\"}\n")
	})
	t.Run("NoUserId", func(t *testing.T) {
		inputBody := inputStruct{}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enroll", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Enroll()

		err = handleFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusBadRequest)
		require.Equal(t, rec.Body.String(), "{\"status\":400,\"error\":\"No user_id field\"}\n")
	})
	t.Run("NoUserName", func(t *testing.T) {
		inputBody := inputStruct{
			UserId: uuid.New().String(),
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enroll", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Enroll()

		err = handleFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusBadRequest)
		require.Equal(t, rec.Body.String(), "{\"status\":400,\"error\":\"No user_name field\"}\n")
	})
	t.Run("InternalUuidParseError", func(t *testing.T) {
		inputBody := inputStruct{
			UserId:   "8e258b27-c787-49ef-9539-11461b251ffad",
			UserName: userName,
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enroll", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Enroll()

		err = handleFunc(c)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusInternalServerError)
		require.Equal(t, rec.Body.String(), "{\"status\":500,\"error\":\"Internal Server Error\"}\n")
	})
	t.Run("EnrollError", func(t *testing.T) {
		inputBody := inputStruct{
			UserId:   uuid.New().String(),
			UserName: userName,
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enroll", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Enroll()

		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpH.Enroll")
		defer span.Finish()

		userId, err := uuid.Parse(inputBody.UserId)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, userId)

		mockUC.EXPECT().Enroll(ctx, gomock.Eq(models.TOTPConfig{
			UserId:      userId,
			AccountName: userName,
		})).Return(nil, totpErrors.ActiveTotp)

		err = handleFunc(c)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusForbidden)
		require.Equal(t, rec.Body.String(), "{\"status\":403,\"error\":\"User's totp is active\"}\n")
	})
	t.Run("OK", func(t *testing.T) {
		secret := "JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
		url := "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&digits=6&issuer=dbo.gcfactory.space&period=30&secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"

		inputBody := inputStruct{
			UserId:   uuid.New().String(),
			UserName: userName,
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enroll", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Enroll()

		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpH.Enroll")
		defer span.Finish()

		userId, err := uuid.Parse(inputBody.UserId)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, userId)

		enrollResult := models.TOTPEnroll{
			TotpId:     uuid.New().String(),
			TotpSecret: secret,
			TotpUrl:    url,
		}

		mockUC.EXPECT().Enroll(ctx, gomock.Eq(models.TOTPConfig{
			UserId:      userId,
			AccountName: userName,
		})).Return(&enrollResult, nil)

		err = handleFunc(c)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusCreated)

		buf, err = converter.AnyToBytesBuffer(enrollResult)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		require.Equal(t, rec.Body.String(), buf.String())
	})
	t.Run("BadRequest", func(t *testing.T) {
		inputBody := inputStruct{
			UserId:   uuid.New().String(),
			UserName: userName,
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enroll", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Enroll()

		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpH.Enroll")
		defer span.Finish()

		userId, err := uuid.Parse(inputBody.UserId)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, userId)

		mockUC.EXPECT().Enroll(ctx, gomock.Eq(models.TOTPConfig{
			UserId:      userId,
			AccountName: userName,
		})).Return(nil, totpErrors.NoUserName)

		err = handleFunc(c)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusBadRequest)
		require.Equal(t, rec.Body.String(), "{\"status\":400,\"error\":\"No user name\"}\n")
	})
	t.Run("InteralServerError", func(t *testing.T) {
		inputBody := inputStruct{
			UserId:   uuid.New().String(),
			UserName: userName,
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enroll", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Enroll()

		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpH.Enroll")
		defer span.Finish()

		userId, err := uuid.Parse(inputBody.UserId)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, userId)

		mockUC.EXPECT().Enroll(ctx, gomock.Eq(models.TOTPConfig{
			UserId:      userId,
			AccountName: userName,
		})).Return(nil, totpRepo.ErrorCreateConfig)

		err = handleFunc(c)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusInternalServerError)
		require.Equal(t, rec.Body.String(), "{\"status\":500,\"error\":\"Internal Server Error\"}\n")
	})
}

func TestTotpHandlers_Verify(t *testing.T) {
	t.Parallel()
	type inputStruct struct {
		TotpUrl string `json:"totp_url"`
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLoggger := logger.NewServerLogger(testCfg)
	apiLoggger.InitLogger()

	mockUC := mock.NewMockUseCase(ctrl)
	tHandlers := NewTOTPHandlers(testCfg, mockUC, apiLoggger)

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
	url := "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&digits=6&issuer=dbo.gcfactory.space&period=30&secret=JOMS6CZZZ4L7S4F6CADFZWMZRJAB5WPP"
	t.Run("EmptyBody", func(t *testing.T) {
		inputBody := ""

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/verify", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handlerFunc := tHandlers.Verify()

		err = handlerFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Equal(t, rec.Body.String(), "{\"status\":500,\"error\":\"Internal Server Error\"}\n")
	})
	t.Run("NoUrl", func(t *testing.T) {
		inputBody := inputStruct{}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/verify", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Verify()

		err = handleFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusBadRequest)
		require.Equal(t, rec.Body.String(), "{\"status\":400,\"error\":\"No totp_url field\"}\n")
	})
	t.Run("OK", func(t *testing.T) {
		inputBody := inputStruct{
			TotpUrl: url,
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/verify", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Verify()

		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpUC.Verify")
		defer span.Finish()

		verifyResult := models.TOTPVerify{
			Status: "OK",
		}

		mockUC.EXPECT().Verify(ctx, gomock.Eq(inputBody.TotpUrl)).Return(&verifyResult, nil)

		err = handleFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusOK)
		require.Equal(t, rec.Body.String(), "{\"status\":\"OK\"}\n")
	})
	t.Run("BadRequest", func(t *testing.T) {
		inputBody := inputStruct{
			TotpUrl: "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&digits=6&issuer=dbo.gcfactory.space&period=30",
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/verify", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Verify()

		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpUC.Verify")
		defer span.Finish()

		verifyResult := models.TOTPVerify{
			Status: totpErrors.NoSecretField.Error(),
		}

		mockUC.EXPECT().Verify(ctx, gomock.Eq(inputBody.TotpUrl)).Return(&verifyResult, totpErrors.NoSecretField)

		err = handleFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusBadRequest)
		require.Equal(t, rec.Body.String(), "{\"status\":\"No secret field\"}\n")
	})
	t.Run("InternalServerError", func(t *testing.T) {
		inputBody := inputStruct{
			TotpUrl: "otpauth://totp/dbo.gcfactory.space:admin?algorithm=SHA1&digits=6&issuer=dbo.gcfactory.space&period=30",
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/verify", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Verify()

		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpUC.Verify")
		defer span.Finish()

		mockUC.EXPECT().Verify(ctx, gomock.Eq(inputBody.TotpUrl)).Return(nil, totpUseCase.ErrorRegexCompile)

		err = handleFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusInternalServerError)
		require.Equal(t, rec.Body.String(), "{\"status\":500,\"error\":\"Internal Server Error\"}\n")
	})
}

func TestTotpHandlers_Validate(t *testing.T) {
	t.Parallel()
	type inputStruct struct {
		UserId   string `json:"user_id"`
		TotpCode string `json:"totp_code"`
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLoggger := logger.NewServerLogger(testCfg)
	apiLoggger.InitLogger()

	mockUC := mock.NewMockUseCase(ctrl)
	tHandlers := NewTOTPHandlers(testCfg, mockUC, apiLoggger)

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
	t.Run("EmptyBody", func(t *testing.T) {
		inputBody := ""

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/validate", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handlerFunc := tHandlers.Validate()

		err = handlerFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Equal(t, rec.Body.String(), "{\"status\":500,\"error\":\"Internal Server Error\"}\n")
	})
	t.Run("NoUserId", func(t *testing.T) {
		inputBody := inputStruct{}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/validate", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Validate()

		err = handleFunc(c)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusBadRequest)
		require.Equal(t, rec.Body.String(), "{\"status\":400,\"error\":\"No user_id field\"}\n")
	})
	t.Run("NoTotpCode", func(t *testing.T) {
		inputBody := inputStruct{
			UserId: uuid.New().String(),
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/validate", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Validate()

		err = handleFunc(c)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusBadRequest)
		require.Equal(t, rec.Body.String(), "{\"status\":400,\"error\":\"No totp_code field\"}\n")
	})
	t.Run("InternalErrorParseUserId", func(t *testing.T) {
		inputBody := inputStruct{
			UserId:   "8e258b27-c787-49ef-9539-11461b251ffad",
			TotpCode: "123456",
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/validate", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Validate()

		err = handleFunc(c)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusInternalServerError)
		require.Equal(t, rec.Body.String(), "{\"status\":500,\"error\":\"Internal Server Error\"}\n")
	})
	t.Run("OK", func(t *testing.T) {
		inputBody := inputStruct{
			UserId:   uuid.New().String(),
			TotpCode: "123456",
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/validate", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Validate()

		userId, err := uuid.Parse(inputBody.UserId)
		require.Nil(t, err)
		require.NotNil(t, userId)
		require.NotEqual(t, userId, uuid.Nil)

		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpUC.Validate")
		defer span.Finish()

		validateResult := models.TOTPValidate{
			Status: "OK",
		}

		mockUC.EXPECT().Validate(ctx, gomock.Eq(userId), gomock.Eq(inputBody.TotpCode), gomock.Any()).Return(&validateResult, nil)

		err = handleFunc(c)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusOK)
		require.Equal(t, rec.Body.String(), "{\"status\":\"OK\"}\n")
	})
	t.Run("UserNotFound", func(t *testing.T) {
		inputBody := inputStruct{
			UserId:   uuid.New().String(),
			TotpCode: "123456",
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/validate", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Validate()

		userId, err := uuid.Parse(inputBody.UserId)
		require.Nil(t, err)
		require.NotNil(t, userId)
		require.NotEqual(t, userId, uuid.Nil)

		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpUC.Validate")
		defer span.Finish()

		validateResult := models.TOTPValidate{
			Status: totpErrors.NoUserId.Error(),
		}

		mockUC.EXPECT().Validate(ctx, gomock.Eq(userId), gomock.Eq(inputBody.TotpCode), gomock.Any()).Return(&validateResult, totpErrors.NoUserId)

		err = handleFunc(c)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusNotFound)
		require.Equal(t, rec.Body.String(), "{\"status\":\"No totp for this user id\"}\n")
	})
	t.Run("WrongTotpCode", func(t *testing.T) {
		inputBody := inputStruct{
			UserId:   uuid.New().String(),
			TotpCode: "123456",
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/validate", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Validate()

		userId, err := uuid.Parse(inputBody.UserId)
		require.Nil(t, err)
		require.NotNil(t, userId)
		require.NotEqual(t, userId, uuid.Nil)

		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpUC.Validate")
		defer span.Finish()

		validateResult := models.TOTPValidate{
			Status: totpErrors.WrongTotpCode.Error(),
		}

		mockUC.EXPECT().Validate(ctx, gomock.Eq(userId), gomock.Eq(inputBody.TotpCode), gomock.Any()).Return(&validateResult, totpErrors.WrongTotpCode)

		err = handleFunc(c)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusBadRequest)
		require.Equal(t, rec.Body.String(), "{\"status\":\"Wrong totp code\"}\n")
	})
	t.Run("InternalValidateError", func(t *testing.T) {
		inputBody := inputStruct{
			UserId:   uuid.New().String(),
			TotpCode: "123456",
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/validate", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handleFunc := tHandlers.Validate()

		userId, err := uuid.Parse(inputBody.UserId)
		require.Nil(t, err)
		require.NotNil(t, userId)
		require.NotEqual(t, userId, uuid.Nil)

		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpUC.Validate")
		defer span.Finish()

		mockUC.EXPECT().Validate(ctx, gomock.Eq(userId), gomock.Eq(inputBody.TotpCode), gomock.Any()).Return(nil, totpUseCase.ErrorRegexCompile)

		err = handleFunc(c)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusInternalServerError)
		require.Equal(t, rec.Body.String(), "{\"status\":500,\"error\":\"Internal Server Error\"}\n")
	})
}

func TestTotpHandlers_Enable(t *testing.T) {
	t.Parallel()
	type inputStruct struct {
		UserId string `json:"user_id"`
		TotpId string `json:"totp_id"`
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLoggger := logger.NewServerLogger(testCfg)
	apiLoggger.InitLogger()

	mockUC := mock.NewMockUseCase(ctrl)
	tHandlers := NewTOTPHandlers(testCfg, mockUC, apiLoggger)

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
	t.Run("EmptyBody", func(t *testing.T) {
		inputBody := ""

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enable", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handlerFunc := tHandlers.Enable()

		err = handlerFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusInternalServerError)
		require.Equal(t, rec.Body.String(), "{\"status\":500,\"error\":\"Internal Server Error\"}\n")
	})
	t.Run("Empty both id", func(t *testing.T) {
		inputBody := inputStruct{
			UserId: "",
			TotpId: "",
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enable", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handlerFunc := tHandlers.Enable()

		err = handlerFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusBadRequest)
		require.Equal(t, rec.Body.String(), "{\"status\":400,\"error\":\"No user_id and totp_id fields\"}\n")
	})
	t.Run("InternalTotpIdParse", func(t *testing.T) {
		inputBody := inputStruct{
			UserId: "",
			TotpId: "0f14d0ab-9605-4a62-a9e4-5ed26688389!",
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enable", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handlerFunc := tHandlers.Enable()

		err = handlerFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusInternalServerError)
		require.Equal(t, rec.Body.String(), "{\"status\":500,\"error\":\"Internal Server Error\"}\n")
	})
	t.Run("InternalUserIdParse", func(t *testing.T) {
		inputBody := inputStruct{
			UserId: "0f14d0ab-9605-4a62-a9e4-5ed26688389!",
			TotpId: "",
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enable", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handlerFunc := tHandlers.Enable()

		err = handlerFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusInternalServerError)
		require.Equal(t, rec.Body.String(), "{\"status\":500,\"error\":\"Internal Server Error\"}\n")
	})
	t.Run("OK", func(t *testing.T) {
		inputBody := inputStruct{
			UserId: uuid.New().String(),
			TotpId: uuid.New().String(),
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enable", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handlerFunc := tHandlers.Enable()

		userId, err := uuid.Parse(inputBody.UserId)
		require.Nil(t, err)
		require.NotNil(t, userId)
		require.NotEqual(t, userId, uuid.Nil)

		totpId, err := uuid.Parse(inputBody.TotpId)
		require.Nil(t, err)
		require.NotNil(t, totpId)
		require.NotEqual(t, totpId, uuid.Nil)

		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpUC.Enable")
		defer span.Finish()

		enableResult := models.TOTPEnable{
			Status: "OK",
		}

		mockUC.EXPECT().Enable(ctx, gomock.Eq(totpId), gomock.Eq(userId)).Return(&enableResult, nil)

		err = handlerFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusOK)
		require.Equal(t, rec.Body.String(), "{\"status\":\"OK\"}\n")
	})
	t.Run("NoTotpId", func(t *testing.T) {
		inputBody := inputStruct{
			UserId: "",
			TotpId: uuid.New().String(),
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enable", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handlerFunc := tHandlers.Enable()

		totpId, err := uuid.Parse(inputBody.TotpId)
		require.Nil(t, err)
		require.NotNil(t, totpId)
		require.NotEqual(t, totpId, uuid.Nil)

		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpUC.Enable")
		defer span.Finish()

		enableResult := models.TOTPEnable{
			Status: totpErrors.NoTotpId.Error(),
		}

		mockUC.EXPECT().Enable(ctx, gomock.Eq(totpId), gomock.Eq(uuid.Nil)).Return(&enableResult, totpErrors.NoTotpId)

		err = handlerFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusNotFound)
		require.Equal(t, rec.Body.String(), "{\"status\":\"No totp with this totp id\"}\n")
	})
	t.Run("TotpIsActive", func(t *testing.T) {
		inputBody := inputStruct{
			UserId: "",
			TotpId: uuid.New().String(),
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enable", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handlerFunc := tHandlers.Enable()

		totpId, err := uuid.Parse(inputBody.TotpId)
		require.Nil(t, err)
		require.NotNil(t, totpId)
		require.NotEqual(t, totpId, uuid.Nil)

		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpUC.Enable")
		defer span.Finish()

		enableResult := models.TOTPEnable{
			Status: totpErrors.TotpIsActive.Error(),
		}

		mockUC.EXPECT().Enable(ctx, gomock.Eq(totpId), gomock.Eq(uuid.Nil)).Return(&enableResult, totpErrors.TotpIsActive)

		err = handlerFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusBadRequest)
		require.Equal(t, rec.Body.String(), "{\"status\":\"Totp is active yet\"}\n")
	})
	t.Run("InternalEnableError", func(t *testing.T) {
		inputBody := inputStruct{
			UserId: "",
			TotpId: uuid.New().String(),
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, buf)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enable", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		handlerFunc := tHandlers.Enable()

		totpId, err := uuid.Parse(inputBody.TotpId)
		require.Nil(t, err)
		require.NotNil(t, totpId)
		require.NotEqual(t, totpId, uuid.Nil)

		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpUC.Enable")
		defer span.Finish()

		mockUC.EXPECT().Enable(ctx, gomock.Eq(totpId), gomock.Eq(uuid.Nil)).Return(nil, totpUseCase.ErrorUpdateActivityByTotpId)

		err = handlerFunc(c)
		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, rec.Code, http.StatusInternalServerError)
		require.Equal(t, rec.Body.String(), "{\"status\":500,\"error\":\"Internal Server Error\"}\n")
	})
}

func TestTotpHandlers_Disable(t *testing.T) {
	type inputStruct struct {
		UserId string `json:"user_id"`
		TotpId string `json:"totp_id"`
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfg)
	apiLogger.InitLogger()

	mockTOTP := mock.NewMockUseCase(ctrl)
	tHandlers := NewTOTPHandlers(testCfg, mockTOTP, apiLogger)

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

	t.Run("EmptyBody", func(t *testing.T) {
		inputBody := ""

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.NotNil(t, buf)
		require.Nil(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/disable", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		handlerFunc := tHandlers.Disable()

		err = handlerFunc(c)
		require.Nil(t, err)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Equal(t, rec.Body.String(), "{\"status\":500,\"error\":\"Internal Server Error\"}\n")
	})

	t.Run("EmptyJSON", func(t *testing.T) {
		inputBody := inputStruct{
			UserId: "",
			TotpId: "",
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.NotNil(t, buf)
		require.Nil(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/disable", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		handlerFunc := tHandlers.Disable()

		err = handlerFunc(c)
		require.Nil(t, err)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, rec.Body.String(), "{\"status\":400,\"error\":\"No user_id and totp_id fields\"}\n")
	})

	t.Run("BothID", func(t *testing.T) {
		inputBody := inputStruct{
			UserId: "00000000-0000-0000-1234-000000000000",
			TotpId: "00000000-0000-1234-0000-000000000000",
		}

		totpDisable := models.TOTPDisable{
			Status: "",
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.NotNil(t, buf)
		require.Nil(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/disable", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		ctx := utils.GetRequestCtx(c)
		handlerFunc := tHandlers.Disable()

		totpId, err := uuid.Parse(inputBody.TotpId)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, totpId)

		userId, err := uuid.Parse(inputBody.UserId)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, totpId)

		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "auth.Login")
		defer span.Finish()
		mockTOTP.EXPECT().Disable(ctxWithTrace, gomock.Eq(totpId), gomock.Eq(userId)).Return(&totpDisable, nil)

		expectedBody, _ := json.Marshal(totpDisable)
		// Echo outputs JSON with \n
		expectedBody = append(expectedBody, '\n')

		err = handlerFunc(c)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, rec.Body.Bytes(), expectedBody)
	})

	t.Run("OnlyUserID", func(t *testing.T) {
		inputBody := inputStruct{
			UserId: "00000000-0000-0000-1234-000000000000",
			TotpId: "",
		}

		totpDisable := models.TOTPDisable{
			Status: "",
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.NotNil(t, buf)
		require.Nil(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/disable", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		ctx := utils.GetRequestCtx(c)
		handlerFunc := tHandlers.Disable()

		totpId := uuid.Nil

		userId, err := uuid.Parse(inputBody.UserId)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, userId)

		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "auth.Login")
		defer span.Finish()
		mockTOTP.EXPECT().Disable(ctxWithTrace, gomock.Eq(totpId), gomock.Eq(userId)).Return(&totpDisable, nil)

		expectedBody, _ := json.Marshal(totpDisable)
		// Echo outputs JSON with \n
		expectedBody = append(expectedBody, '\n')

		err = handlerFunc(c)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, rec.Body.Bytes(), expectedBody)
	})

	t.Run("OnlyTotpID", func(t *testing.T) {
		inputBody := inputStruct{
			UserId: "",
			TotpId: "00000000-0000-0000-1234-000000000000",
		}

		totpDisable := models.TOTPDisable{
			Status: "",
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.NotNil(t, buf)
		require.Nil(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/disable", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		ctx := utils.GetRequestCtx(c)
		handlerFunc := tHandlers.Disable()

		totpId, err := uuid.Parse(inputBody.TotpId)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, totpId)

		userId := uuid.Nil

		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "auth.Login")
		defer span.Finish()
		mockTOTP.EXPECT().Disable(ctxWithTrace, gomock.Eq(totpId), gomock.Eq(userId)).Return(&totpDisable, nil)

		expectedBody, _ := json.Marshal(totpDisable)
		// Echo outputs JSON with \n
		expectedBody = append(expectedBody, '\n')

		err = handlerFunc(c)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, rec.Body.Bytes(), expectedBody)
	})

	t.Run("NotFoundTotpId", func(t *testing.T) {
		inputBody := inputStruct{
			UserId: "00000000-0000-0000-1234-000000000000",
			TotpId: "00000000-0000-1234-0000-000000000000",
		}

		totpDisable := models.TOTPDisable{
			Status: "",
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.NotNil(t, buf)
		require.Nil(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/disable", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		ctx := utils.GetRequestCtx(c)
		handlerFunc := tHandlers.Disable()

		totpId, err := uuid.Parse(inputBody.TotpId)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, totpId)

		userId, err := uuid.Parse(inputBody.UserId)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, totpId)

		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "auth.Login")
		defer span.Finish()
		mockTOTP.EXPECT().Disable(ctxWithTrace, gomock.Eq(totpId), gomock.Eq(userId)).Return(&totpDisable, totpErrors.NoTotpId)

		expectedBody, _ := json.Marshal(totpDisable)
		// Echo outputs JSON with \n
		expectedBody = append(expectedBody, '\n')

		err = handlerFunc(c)
		require.Equal(t, http.StatusNotFound, rec.Code)
		require.Equal(t, rec.Body.Bytes(), expectedBody)
	})

	t.Run("DisabledTotp", func(t *testing.T) {
		inputBody := inputStruct{
			UserId: "00000000-0000-0000-1234-000000000000",
			TotpId: "00000000-0000-1234-0000-000000000000",
		}

		totpDisable := models.TOTPDisable{
			Status: "",
		}

		buf, err := converter.AnyToBytesBuffer(inputBody)
		require.NoError(t, err)
		require.NotNil(t, buf)
		require.Nil(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/disable", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		ctx := utils.GetRequestCtx(c)
		handlerFunc := tHandlers.Disable()

		totpId, err := uuid.Parse(inputBody.TotpId)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, totpId)

		userId, err := uuid.Parse(inputBody.UserId)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, totpId)

		span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "auth.Login")
		defer span.Finish()
		mockTOTP.EXPECT().Disable(ctxWithTrace, gomock.Eq(totpId), gomock.Eq(userId)).Return(&totpDisable, totpErrors.TotpIsDisabled)

		expectedBody, _ := json.Marshal(totpDisable)
		// Echo outputs JSON with \n
		expectedBody = append(expectedBody, '\n')

		err = handlerFunc(c)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, rec.Body.Bytes(), expectedBody)
	})
}
