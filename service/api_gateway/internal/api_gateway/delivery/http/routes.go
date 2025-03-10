package http

import (
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/middleware"
	"github.com/labstack/echo/v4"
)

func MapApiGatewayRoutes(apiGatewayGroup *echo.Group, h api_gateway.Handlers, mw *middleware.MiddlewareManager) {
	apiGatewayGroup.GET("/sign_in", h.SignInPage())
	apiGatewayGroup.GET("/sign_up", h.SignUpPage())
}
