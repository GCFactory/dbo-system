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

func (t totpUC) Verify(ctx context.Context, request *models.TOTPRequest) (*models.TOTPConfig, error) {
	//TODO implement me
	panic("implement me")
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
