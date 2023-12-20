package repository

import (
	"context"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/GCFactory/dbo-system/service/totp/internal/totp"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type totpRepo struct {
	db *sqlx.DB
}

func (t totpRepo) CreateConfig(ctx context.Context, totpConfig models.TOTPConfig) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "totpRepo.Save")
	defer span.Finish()

	var s models.TOTPConfig

	if err := t.db.QueryRowxContext(ctx,
		createConfig,
		&totpConfig.Id,
		&totpConfig.UserId,
		&totpConfig.IsActive,
		&totpConfig.Issuer,
		&totpConfig.AccountName,
		&totpConfig.Secret,
	).StructScan(&s); err != nil {
		return errors.Wrap(err, "totpRepo.CreateConfig.QueryRowxContext")
	}
	return nil
}

func (t totpRepo) Verify(ctx context.Context, request models.TOTPRequest) (*models.TOTPConfig, error) {
	//TODO implement me
	panic("implement me")
}

func (t totpRepo) Validate(ctx context.Context, request models.TOTPRequest) (*models.TOTPBase, error) {
	//TODO implement me
	panic("implement me")
}

func (t totpRepo) Enable(ctx context.Context, request models.TOTPRequest) (*models.TOTPBase, error) {
	//TODO implement me
	panic("implement me")
}

func (t totpRepo) Disable(ctx context.Context, request models.TOTPRequest) (*models.TOTPBase, error) {
	//TODO implement me
	panic("implement me")
}

func NewTOTPRepository(db *sqlx.DB) totp.Repository {
	return &totpRepo{db: db}
}
