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
			UserId:   "“8e258b27-c787-49ef-9539-11461b251ffad",
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
