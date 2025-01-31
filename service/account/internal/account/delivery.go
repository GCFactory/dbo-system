package account

import "github.com/labstack/echo/v4"

type Handlers interface {
	TestFunc() echo.HandlerFunc
}
