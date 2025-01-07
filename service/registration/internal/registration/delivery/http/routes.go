package http

import (
	"github.com/GCFactory/dbo-system/service/registration/internal/middleware"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration"
	"github.com/labstack/echo/v4"
)

func MapRegistrationRoutes(RegistrationGroup *echo.Group, h registration.Handlers, mw *middleware.MiddlewareManager) {
	RegistrationGroup.POST("/create_user", h.CreateUser())
	RegistrationGroup.GET("/get_operation_status", h.GetOperationStatus())
}
