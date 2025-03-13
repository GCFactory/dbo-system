package api_gateway

import "github.com/labstack/echo/v4"

type Handlers interface {
	SignInPage() echo.HandlerFunc
	SignIn() echo.HandlerFunc
	SignUpPage() echo.HandlerFunc
	SignUp() echo.HandlerFunc
	SignOut() echo.HandlerFunc
	HomePage() echo.HandlerFunc
	OpenAccountPage() echo.HandlerFunc
	AccountCreditsPage() echo.HandlerFunc
	CloseAccountPage() echo.HandlerFunc
	AddAccountCachePage() echo.HandlerFunc
	WidthAccountCachePage() echo.HandlerFunc
}
