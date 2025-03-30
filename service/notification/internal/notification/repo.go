package notification

import (
	"context"
	"github.com/GCFactory/dbo-system/service/notification/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	AddUserNotificationSettings(ctx context.Context, user *models.UserNotificationInfo) error
	UpdateUserNotificationSettings(ctx context.Context, user *models.UserNotificationInfo) error
	GetUserNotificationSettings(ctx context.Context, userId uuid.UUID) (*models.UserNotificationInfo, error)
	DeleteUserNotificationSettings(ctx context.Context, userId uuid.UUID) error
}
