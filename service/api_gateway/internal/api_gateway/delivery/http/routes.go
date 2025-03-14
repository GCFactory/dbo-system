package http

import (
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/middleware"
	"github.com/labstack/echo/v4"
)

func MapApiGatewayRoutes(apiGatewayGroup *echo.Group, h api_gateway.Handlers, mw *middleware.MiddlewareManager) {
	apiGatewayGroup.GET("/sign_in", h.SignInPage())
	apiGatewayGroup.GET("/sign_up", h.SignUpPage())
	apiGatewayGroup.GET("/main_page", h.HomePage())
	apiGatewayGroup.GET("/admin", h.AdminPage())
	apiGatewayGroup.GET("/open_account", h.OpenAccountPage())
	apiGatewayGroup.GET("/get_account_info", h.AccountCreditsPage())
	apiGatewayGroup.GET("/adding_account", h.AddAccountCachePage())
	apiGatewayGroup.GET("/width_account", h.WidthAccountCachePage())
	apiGatewayGroup.GET("/close_account", h.CloseAccountPage())
	apiGatewayGroup.POST("/sign_in/sign_in", h.SignIn())
	apiGatewayGroup.POST("/sign_up/sign_up", h.SignUp())
	apiGatewayGroup.POST("/sign_out", h.SignOut())
	apiGatewayGroup.POST("/open_account/open_account", h.OpenAccount())
	apiGatewayGroup.POST("/close_account/close_account", h.CloseAccount())
	apiGatewayGroup.POST("/adding_account/adding_account", h.AddAccountCache())
	apiGatewayGroup.POST("/width_account/width_account", h.WidthAccountCache())
}
