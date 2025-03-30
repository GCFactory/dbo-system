package notification

import (
	"context"
	"github.com/GCFactory/dbo-system/service/notification/internal/models"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type UseCase interface {
	AddUserSettings(ctx context.Context, user *models.UserNotificationInfo) error
	UpdateUserSettings(ctx context.Context, user *models.UserNotificationInfo) error
	DeleteUserSettings(ctx context.Context, userId uuid.UUID) error
	SendMessage(ctx context.Context, message amqp091.Delivery) error
}
