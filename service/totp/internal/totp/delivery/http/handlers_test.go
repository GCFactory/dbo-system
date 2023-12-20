package http

import (
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/converter"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/GCFactory/dbo-system/service/totp/internal/totp/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTOTPHandlers_Enroll(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTOTPUsecase := mock.NewMockUseCase(ctrl)

	cfg := &config.Config{
		Logger: config.Logger{
			Development: true,
		},
	}

	apiLogger := logger.NewServerLogger(cfg)
	totpHandlers := NewTOTPHandlers(cfg, mockTOTPUsecase, apiLogger)

	type User struct {
		UserId uuid.UUID `json:"user_id"`
	}
	user := User{
		UserId: uuid.New(),
	}
	totpConfigReq := models.TOTPConfig{
		UserId: user.UserId,
	}

	buff, err := converter.AnyToBytesBuffer(user)

	require.NoError(t, err)
	require.NotNil(t, buff)
	require.Nil(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/totp/enroll", strings.NewReader(buff.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	ctx := utils.GetRequestCtx(c)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totp.Enroll")
	defer span.Finish()

	totpConfigResponse := &models.TOTPEnroll{
		Base32:    "",
		OTPathURL: "",
	}
	handlerFunc := totpHandlers.Enroll()

	mockTOTPUsecase.EXPECT().Enroll(ctxWithTrace, gomock.Eq(totpConfigReq)).Return(totpConfigResponse, nil)

	err = handlerFunc(c)
	require.NoError(t, err)
	require.Nil(t, err)
}
