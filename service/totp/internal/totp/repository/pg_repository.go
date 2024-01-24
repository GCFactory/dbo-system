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
		&totpConfig.URL,
	).StructScan(&s); err != nil {
		return errors.Wrap(err, "totpRepo.CreateConfig.QueryRowxContext")
	}
	return nil
}

func (t totpRepo) GetURL(ctx context.Context, id uuid.UUID) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "totpRepo.Get")
	defer span.Finish()

	type Data struct {
		Url string `json:"url"`
	}

	var tmp Data

	if err := t.db.QueryRowxContext(ctx,
		getURL,
		id,
	).StructScan(&tmp); err != nil {
		return tmp.Url, errors.Wrap(err, "totpRepo.GetSecret.QueryRowxContext")
	}
	return tmp.Url, nil
}

func (t totpRepo) UpdateActive(ctx context.Context, id uuid.UUID, status bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "totpRepo.Update")
	defer span.Finish()

	_, err := t.db.ExecContext(ctx,
		updareActive,
		id,
		status,
	)
	if err != nil {
		return err
	}
	return nil
}

func NewTOTPRepository(db *sqlx.DB) totp.Repository {
	return &totpRepo{db: db}
}
