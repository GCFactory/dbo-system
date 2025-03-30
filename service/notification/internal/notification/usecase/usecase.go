package usecase

import (
	"context"
	"fmt"
	"github.com/GCFactory/dbo-system/service/notification/internal/models"
	"github.com/GCFactory/dbo-system/service/notification/internal/notification"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/rabbitmq/amqp091-go"
)

type NotificationUseCase struct {
	repo notification.Repository
}

func (uc NotificationUseCase) AddUserSettings(ctx context.Context, user *models.UserNotificationInfo) error {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "NotificationUseCase.AddUserSettings")
	defer span.Finish()

	err := uc.repo.AddUserNotificationSettings(local_ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (uc NotificationUseCase) UpdateUserSettings(ctx context.Context, user *models.UserNotificationInfo) error {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "NotificationUseCase.UpdateUserSettings")
	defer span.Finish()

	err := uc.repo.UpdateUserNotificationSettings(local_ctx, user)
	if err != nil {
		return nil
	}

	return nil
}

func (uc NotificationUseCase) DeleteUserSettings(ctx context.Context, userId uuid.UUID) error {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "NotificationUseCase.DeleteUserSettings")
	defer span.Finish()

	err := uc.repo.DeleteUserNotificationSettings(local_ctx, userId)
	if err != nil {
		return err
	}

	return nil
}

func (uc NotificationUseCase) SendMessage(ctx context.Context, message amqp091.Delivery) error {

	//span, local_ctx := opentracing.StartSpanFromContext(ctx, "NotificationUseCase.SendMessage")
	//defer span.Finish()

	msgBody := message.Body
	msgHeaders := message.Headers

	fmt.Println("Body:", string(msgBody))
	fmt.Println("User_id:", msgHeaders["user_id"])
	fmt.Println("Notification_level :", msgHeaders["notification_level"])

	return nil
}

func NewNotificationUseCase(repo notification.Repository) notification.UseCase {
	return &NotificationUseCase{repo: repo}
}
