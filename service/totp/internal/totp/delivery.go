package totp

import "github.com/labstack/echo/v4"

type Handlers interface {
	Enroll() echo.HandlerFunc
	Verify() echo.HandlerFunc
	Validate() echo.HandlerFunc
	Enable() echo.HandlerFunc
	Disable() echo.HandlerFunc
}
