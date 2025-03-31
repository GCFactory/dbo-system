package http

import (
	"github.com/GCFactory/dbo-system/service/users/internal/middleware"
	"github.com/GCFactory/dbo-system/service/users/internal/users"
	"github.com/labstack/echo/v4"
)

func MapNotificationRoutes(usersGroup *echo.Group, h users.HttpHandlers, mw *middleware.MiddlewareManager) {
	usersGroup.GET("/get_user_data", h.GetUserData())
	usersGroup.GET("/get_user_data_accounts", h.GetUserAccounts())
	usersGroup.GET("/get_user_data_by_login", h.GetUserDataByLogin())
}
