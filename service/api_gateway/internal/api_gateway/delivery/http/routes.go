package http

import (
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/middleware"
	"github.com/labstack/echo/v4"
)

func MapApiGatewayRoutes(RegistrationGroup *echo.Group, h api_gateway.Handlers, mw *middleware.MiddlewareManager) {
	RegistrationGroup.POST("/sign_in", h.SignInPage())

}
