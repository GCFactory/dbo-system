package usecase

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/service/notification/internal/models"
	"github.com/GCFactory/dbo-system/service/notification/internal/notification"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/rabbitmq/amqp091-go"
	"net/smtp"
	"slices"
)

type NotificationUseCase struct {
	repo       notification.Repository
	smtpClient *smtp.Client
	smtpCfg    config.Smtp
}

func (uc NotificationUseCase) AddUserSettings(ctx context.Context, user *models.UserNotificationInfo) error {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "NotificationUseCase.AddUserSettings")
	defer span.Finish()

	userD, err := uc.repo.GetUserNotificationSettings(local_ctx, user.UserUuid)
	if err == nil && userD != nil {
		return ErrorUserSettingsAlreadyExist
	}

	err = uc.repo.AddUserNotificationSettings(local_ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (uc NotificationUseCase) GetUserSettings(ctx context.Context, userId uuid.UUID) (*models.UserNotificationInfo, error) {
	span, local_ctx := opentracing.StartSpanFromContext(ctx, "NotificationUseCase.GetUserSettings")
	defer span.Finish()

	result, err := uc.repo.GetUserNotificationSettings(local_ctx, userId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (uc NotificationUseCase) UpdateUserSettings(ctx context.Context, user *models.UserNotificationInfo) error {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "NotificationUseCase.UpdateUserSettings")
	defer span.Finish()

	userD, err := uc.repo.GetUserNotificationSettings(local_ctx, user.UserUuid)
	if err != nil && userD == nil {
		return ErrorUserNotFound
	}

	err = uc.repo.UpdateUserNotificationSettings(local_ctx, user)
	if err != nil {
		return nil
	}

	return nil
}

func (uc NotificationUseCase) DeleteUserSettings(ctx context.Context, userId uuid.UUID) error {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "NotificationUseCase.DeleteUserSettings")
	defer span.Finish()

	userD, err := uc.repo.GetUserNotificationSettings(local_ctx, userId)
	if err != nil && userD == nil {
		return ErrorUserNotFound
	}

	err = uc.repo.DeleteUserNotificationSettings(local_ctx, userId)
	if err != nil {
		return err
	}

	return nil
}

func (uc NotificationUseCase) SendMessage(ctx context.Context, message amqp091.Delivery) error {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "NotificationUseCase.SendMessage")
	defer span.Finish()

	msgHeaders := message.Headers
	ok, err := uc.checkRuquiredHeaders(local_ctx, msgHeaders)
	if !ok {
		return err
	}

	userIdStr := msgHeaders[HeaderUserId]
	userId, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		return err
	}

	notificationLvl := msgHeaders[HeaderNotificationLvl].(string)
	ok, err = uc.validateNotificationLvl(local_ctx, notificationLvl)
	if !ok {
		return err
	}

	userNotificationSettings, err := uc.repo.GetUserNotificationSettings(local_ctx, userId)
	if err != nil {
		return err
	}

	sendEmail := false

	if notificationLvl == NotificationLvlEmail {
		sendEmail = userNotificationSettings.EmailUsage
	} else if notificationLvl == NotificationLvlAll {
		sendEmail = userNotificationSettings.EmailUsage
	}

	if sendEmail {
		// 1. Указываем отправителя
		if err = uc.smtpClient.Mail(uc.smtpCfg.From); err != nil {
			return err
		}

		// 2. Указываем получателя
		if err = uc.smtpClient.Rcpt(userNotificationSettings.Email); err != nil {
			return err
		}

		// 3. Отправляем тело письма
		dataWriter, err := uc.smtpClient.Data()
		if err != nil {
			return err
		}

		// Формируем письмо (заголовки + тело)
		messageBytes := []byte(
			"Subject: Notification message\r\n" +
				"From: " + uc.smtpCfg.User + "\r\n" +
				"To: " + userNotificationSettings.Email + "\r\n" +
				"\r\n" + // Пустая строка — разделитель заголовков и тела
				string(message.Body) + "\r\n",
		)

		_, err = dataWriter.Write(messageBytes)
		if err != nil {
			return err
		}

		// Закрываем writer
		err = dataWriter.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc NotificationUseCase) checkRuquiredHeaders(ctx context.Context, headers map[string]interface{}) (bool, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "NotificationUseCase.checkRuquiredHeaders")
	defer span.Finish()

	_, ok := headers[HeaderUserId]
	if !ok {
		return false, ErrorNoUserIdHeader
	}

	_, ok = headers[HeaderNotificationLvl]
	if !ok {
		return false, ErrorNoNotificationLvlHeader
	}

	return true, nil

}

func (uc NotificationUseCase) validateNotificationLvl(ctx context.Context, notificationLvl string) (bool, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "NotificationUseCase.validateNotificationLvl")
	defer span.Finish()

	exists := slices.Contains(PossibleNotificationLvl, notificationLvl)
	if !exists {
		return exists, ErrorInvalidNotificationLvl
	}

	return exists, nil

}

func NewNotificationUseCase(repo notification.Repository, smtpClient *smtp.Client, smtpCfg config.Smtp) notification.UseCase {
	return &NotificationUseCase{repo: repo, smtpClient: smtpClient, smtpCfg: smtpCfg}
}
