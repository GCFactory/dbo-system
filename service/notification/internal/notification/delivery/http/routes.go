package http

import (
	"github.com/GCFactory/dbo-system/service/notification/internal/middleware"
	"github.com/GCFactory/dbo-system/service/notification/internal/notification"
	"github.com/labstack/echo/v4"
)

func MapNotificationRoutes(notificationGroup *echo.Group, h notification.HttpHandlers, mw *middleware.MiddlewareManager) {
	notificationGroup.POST("/add_user_notification_settings", h.AddUserSettings())
	notificationGroup.POST("/update_user_notification_settings", h.UpdateUserSettings())
	notificationGroup.POST("/delete_user_notification_settings", h.DeleteUserSettings())
}
