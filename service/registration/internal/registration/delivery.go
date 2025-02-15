package registration

import "github.com/labstack/echo/v4"

type Handlers interface {
	CreateUser() echo.HandlerFunc
	GetOperationStatus() echo.HandlerFunc
	OpenAccount() echo.HandlerFunc
	AddAccountCache() echo.HandlerFunc
	WidthAccountCache() echo.HandlerFunc
	CloseAccount() echo.HandlerFunc
}
