package usecase

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/GCFactory/dbo-system/service/totp/internal/totp"
	totpErrors "github.com/GCFactory/dbo-system/service/totp/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
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
