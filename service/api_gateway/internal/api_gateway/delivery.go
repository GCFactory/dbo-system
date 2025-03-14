package api_gateway

import "github.com/labstack/echo/v4"

type Handlers interface {
	SignInPage() echo.HandlerFunc
	SignUpPage() echo.HandlerFunc
	HomePage() echo.HandlerFunc
	OpenAccountPage() echo.HandlerFunc
	AccountCreditsPage() echo.HandlerFunc
	CloseAccountPage() echo.HandlerFunc
	AddAccountCachePage() echo.HandlerFunc
	WidthAccountCachePage() echo.HandlerFunc
	AdminPage() echo.HandlerFunc
	SignIn() echo.HandlerFunc
	SignUp() echo.HandlerFunc
	SignOut() echo.HandlerFunc
	OpenAccount() echo.HandlerFunc
	CloseAccount() echo.HandlerFunc
	AddAccountCache() echo.HandlerFunc
	WidthAccountCache() echo.HandlerFunc
}
