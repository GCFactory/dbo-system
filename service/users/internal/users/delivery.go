package users

import "github.com/labstack/echo/v4"

type HttpHandlers interface {
	GetUserData() echo.HandlerFunc
	GetUserAccounts() echo.HandlerFunc
	GetUserDataByLogin() echo.HandlerFunc
	CheckUserPassw() echo.HandlerFunc
	GetUserTotpInfo() echo.HandlerFunc
	UpdateTotpInfo() echo.HandlerFunc
}
