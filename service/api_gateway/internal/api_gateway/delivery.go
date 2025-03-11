package api_gateway

import "github.com/labstack/echo/v4"

type Handlers interface {
	SignInPage() echo.HandlerFunc
	SignIn() echo.HandlerFunc
	SignUpPage() echo.HandlerFunc
	SignUp() echo.HandlerFunc
}
