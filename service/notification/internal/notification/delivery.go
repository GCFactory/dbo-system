package notification

import "github.com/labstack/echo/v4"

type HttpHandlers interface {
	AddUserSettings() echo.HandlerFunc
	UpdateUserSettings() echo.HandlerFunc
	DeleteUserSettings() echo.HandlerFunc
}
