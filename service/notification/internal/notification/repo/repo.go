package repo

import (
	"context"
	"github.com/GCFactory/dbo-system/service/notification/internal/models"
	"github.com/GCFactory/dbo-system/service/notification/internal/notification"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

type NotificationRepository struct {
	db *sqlx.DB
}

func (repo NotificationRepository) AddUserNotificationSettings(ctx context.Context, user *models.UserNotificationInfo) error {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "NotificationRepository.AddUserNotificationSettings")
	defer span.Finish()

	if _, err := repo.db.ExecContext(local_ctx,
		AddUserSettings,
		user.UserUuid,
		user.EmailUsage,
		user.Email,
	); err != nil {
		return err
	}

	return nil
}

func (repo NotificationRepository) UpdateUserNotificationSettings(ctx context.Context, user *models.UserNotificationInfo) error {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "NotificationRepository.UpdateUserNotificationSettings")
	defer span.Finish()

	res, err := repo.db.ExecContext(local_ctx,
		UpdateUserSettings,
		user.UserUuid,
		user.EmailUsage,
		user.Email,
	)

	if err != nil {
		return err
	} else if count, err := res.RowsAffected(); err != nil || count == 0 {
		return ErrorNoUserSettingsFound
	}

	return nil
}

func (repo NotificationRepository) GetUserNotificationSettings(ctx context.Context, userId uuid.UUID) (*models.UserNotificationInfo, error) {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "NotificationRepository.GetUserNotificationSettings")
	defer span.Finish()

	var userSettings = &models.UserNotificationInfo{}

	if err := repo.db.QueryRowxContext(local_ctx,
		GetUserSettings,
		&userId,
	).Scan(&userSettings.UserUuid, &userSettings.EmailUsage, &userSettings.Email); err != nil {
		return nil, err
	}

	return userSettings, nil
}

func (repo NotificationRepository) DeleteUserNotificationSettings(ctx context.Context, userId uuid.UUID) error {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "NotificationRepository.DeleteUserNotificationSettings")
	defer span.Finish()

	res, err := repo.db.ExecContext(local_ctx,
		DeleteUserSettings,
		userId,
	)

	if err != nil {
		return err
	} else if count, err := res.RowsAffected(); err != nil || count == 0 {
		return ErrorNoUserSettingsFound
	}

	return nil
}

func NewNotificationRepository(db *sqlx.DB) notification.Repository {
	return &NotificationRepository{db: db}
}
