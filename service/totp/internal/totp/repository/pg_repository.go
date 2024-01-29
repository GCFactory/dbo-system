package repository

import (
	"context"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/GCFactory/dbo-system/service/totp/internal/totp"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type totpRepo struct {
	db *sqlx.DB
}

func (t totpRepo) CreateConfig(ctx context.Context, totpConfig models.TOTPConfig) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "totpRepo.CreateConfig")
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
		&totpConfig.URL,
	).StructScan(&s); err != nil {
		return errors.Wrap(err, "totpRepo.CreateConfig.QueryRowxContext")
	}
	return nil
}

func (t totpRepo) GetActiveConfig(ctx context.Context, userId uuid.UUID) (*models.TOTPConfig, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "totpRepo.GetActiveConfig")
	defer span.Finish()

	var s models.TOTPConfig

	if err := t.db.QueryRowxContext(ctx,
		getActiveConfig,
		&userId,
	).StructScan(&s); err != nil {
		return nil, errors.Wrap(err, "totpRepo.GetActiveConfig.QueryRowContext")
	}
	return &s, nil
}

func (t totpRepo) GetConfigByTotpId(ctx context.Context, totpId uuid.UUID) (*models.TOTPConfig, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "totpRepo.GetConfigByTotpId")
	defer span.Finish()

	var s models.TOTPConfig

	if err := t.db.QueryRowxContext(ctx,
		getConfigByTotpId,
		&totpId,
	).StructScan(&s); err != nil {
		return nil, errors.Wrap(err, "totpRepo.GetConfigByTotpId.QueryRowContext")
	}
	return &s, nil
}

func (t totpRepo) UpdateTotpActivityByTotpId(ctx context.Context, totpId uuid.UUID, status bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "totpRepo.UpdateActive")
	defer span.Finish()

	_, err := t.db.ExecContext(ctx,
		updateTotpActivityByTotpId,
		totpId,
		status,
	)
	if err != nil {
		return err
	}
	return nil
}

func (t totpRepo) GetConfigByUserId(ctx context.Context, userId uuid.UUID) (*models.TOTPConfig, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "totpRepo.GetConfigByUserId")
	defer span.Finish()

	var s models.TOTPConfig

	if err := t.db.QueryRowxContext(ctx,
		getConfigByUserId,
		&userId,
	).StructScan(&s); err != nil {
		return nil, errors.Wrap(err, "totpRepo.GetConfigByUserId.QueryRowContext")
	}
	return &s, nil
}

func (t totpRepo) GetLastDisabledConfig(ctx context.Context, userId uuid.UUID) (*models.TOTPConfig, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "totpRepo.GetConfigByUserId")
	defer span.Finish()

	var s models.TOTPConfig

	if err := t.db.QueryRowxContext(ctx,
		getLastDisabledConfig,
		&userId,
	).StructScan(&s); err != nil {
		return nil, errors.Wrap(err, "totpRepo.GetLastDisabledConfig.QueryRowContext")
	}
	return &s, nil
}

func NewTOTPRepository(db *sqlx.DB) totp.Repository {
	return &totpRepo{db: db}
}
