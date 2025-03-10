package api_gateway

import "github.com/labstack/echo/v4"

type Handlers interface {
	SignInPage() echo.HandlerFunc
}
