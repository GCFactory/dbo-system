package http

import (
	"github.com/GCFactory/dbo-system/service/account/internal/account"
	"github.com/GCFactory/dbo-system/service/account/internal/middleware"
	"github.com/labstack/echo/v4"
)

func MapACCRoutes(accountGroup *echo.Group, h account.Handlers, mw *middleware.MiddlewareManager) {
	accountGroup.POST("/test_func", h.TestFunc())
}
