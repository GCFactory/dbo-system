package usecase

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/httpErrors"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/GCFactory/dbo-system/service/totp/internal/totp"
	totpPkg "github.com/GCFactory/dbo-system/service/totp/pkg/otp"
	totpPkgConfig "github.com/GCFactory/dbo-system/service/totp/pkg/otp/config"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type totpUC struct {
	cfg      *config.Config
	totpRepo totp.Repository
	logger   logger.Logger
}

func (t totpUC) Enroll(ctx context.Context, totpConfig models.TOTPConfig) (*models.TOTPEnroll, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "totpUC.Enroll")
	defer span.Finish()

	totpConfig.Id = uuid.New()
	totpConfig.AccountName = "admin"
	totpConfig.IsActive = true
	totpConfig.Issuer = "dbo.gcfactory.space"

	secret, url, err := totpPkg.Generate(totpPkgConfig.GenerateOpts{
		Issuer:      totpConfig.Issuer,
		AccountName: totpConfig.AccountName,
		SecretSize:  20,
		Algorithm:   totpPkgConfig.AlgorithmSHA256,
	})

	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "totpUC.Enroll.Generate"))
	}
	totpConfig.Secret = *secret
	totpConfig.URL = *url

	err = t.totpRepo.CreateConfig(ctx, totpConfig)
	if err != nil {
		return nil, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)
	}

	totpEnroll := models.TOTPEnroll{
		Base32:    totpConfig.Secret,
		OTPathURL: *url,
	}

	return &totpEnroll, nil
}

func (t totpUC) Verify(ctx context.Context, url string) (*models.TOTPVerify, error) {
	if len(url) == 0 {
		return &models.TOTPVerify{Status: "Empty TOTP URL"}, errors.New("Empty TOTP URL")
	}
	algorithm := strings.Contains(url, "algorithm")
	if !algorithm {
		return &models.TOTPVerify{Status: "No algorithm field"}, errors.New("URL hadn't algorithm field")
	}
	digits := strings.Contains(url, "digits")
	if !digits {
		return &models.TOTPVerify{Status: "No digits field"}, errors.New("URL hadn't digits field")
	}
	issuer := strings.Contains(url, "issuer")
	if !issuer {
		return &models.TOTPVerify{Status: "No issuer field"}, errors.New("URL hadn't issuer field")
	}
	period := strings.Contains(url, "period")
	if !period {
		return &models.TOTPVerify{Status: "No period field"}, errors.New("URL hadn't period field")
	}
	secret := strings.Contains(url, "secret")
	if !secret {
		return &models.TOTPVerify{Status: "No secret field"}, errors.New("URL hadn't secret field")
	}
	return &models.TOTPVerify{Status: "OK"}, nil
}

func (t totpUC) Validate(ctx context.Context, request *models.TOTPRequest) (*models.TOTPBase, error) {
	//TODO implement me
	panic("implement me")
}

func (t totpUC) Enable(ctx context.Context, request *models.TOTPRequest) (*models.TOTPBase, error) {
	//TODO implement me
	panic("implement me")
}

func (t totpUC) Disable(ctx context.Context, request *models.TOTPRequest) (*models.TOTPBase, error) {
	//TODO implement me
	panic("implement me")
}

func NewTOTPUseCase(cfg *config.Config, totpRepo totp.Repository, log logger.Logger) totp.UseCase {
	return &totpUC{cfg: cfg, totpRepo: totpRepo, logger: log}
}
