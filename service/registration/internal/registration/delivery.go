package registration

import "github.com/labstack/echo/v4"

type Handlers interface {
	CreateUser() echo.HandlerFunc
	GetOperationStatus() echo.HandlerFunc
}
