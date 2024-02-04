package http

import (
	"encoding/json"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/converter"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/GCFactory/dbo-system/service/totp/internal/totp/mock"
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

//func TestTOTPHandlers_Enroll(t *testing.T) {
//	t.Parallel()
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockTOTPUsecase := mock.NewMockUseCase(ctrl)
//
//	testCfg := &config.Config{
//		Logger: config.Logger{
//			Development: true,
//		},
//	}
//
//	apiLogger := logger.NewServerLogger(testCfg)
//	totpHandlers := NewTOTPHandlers(testCfg, mockTOTPUsecase, apiLogger)
//
//	type User struct {
//		UserId uuid.UUID `json:"user_id"`
//	}
//	user := User{
//		UserId: uuid.New(),
//	}
//	totpConfigReq := models.TOTPConfig{
//		UserId: user.UserId,
//	}
//
//	buff, err := converter.AnyToBytesBuffer(user)
//
//	require.NoError(t, err)
//	require.NotNil(t, buff)
//	require.Nil(t, err)
//
//	e := echo.New()
//	req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enroll", strings.NewReader(buff.String()))
//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//	rec := httptest.NewRecorder()
//
//	c := e.NewContext(req, rec)
//	ctx := utils.GetRequestCtx(c)
//	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totp.Enroll")
//	defer span.Finish()
//
//	totpConfigResponse := &models.TOTPEnroll{
//		Base32:    "",
//		OTPathURL: "",
//	}
//	handlerFunc := totpHandlers.Enroll()
//
//	mockTOTPUsecase.EXPECT().Enroll(ctxWithTrace, gomock.Eq(totpConfigReq)).Return(totpConfigResponse, nil)
//
//	err = handlerFunc(c)
//	require.NoError(t, err)
//	require.Nil(t, err)
//}

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
