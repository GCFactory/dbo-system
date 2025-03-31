package notification

import "github.com/labstack/echo/v4"

type HttpHandlers interface {
	UpdateUserSettings() echo.HandlerFunc
}
