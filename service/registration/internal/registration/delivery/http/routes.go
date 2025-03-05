package http

import (
	"github.com/GCFactory/dbo-system/service/registration/internal/middleware"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration"
	"github.com/labstack/echo/v4"
)

func MapRegistrationRoutes(RegistrationGroup *echo.Group, h registration.Handlers, mw *middleware.MiddlewareManager) {
	RegistrationGroup.POST("/create_user", h.CreateUser())
	RegistrationGroup.GET("/get_operation_status", h.GetOperationStatus())
	RegistrationGroup.POST("/open_account", h.OpenAccount())
	RegistrationGroup.POST("/add_account_cache", h.AddAccountCache())
	RegistrationGroup.POST("/width_account_cache", h.WidthAccountCache())
	RegistrationGroup.POST("/close_acc", h.CloseAccount())
	RegistrationGroup.POST("/get_user_data", h.GetUserData())
	RegistrationGroup.POST("/get_account_data", h.GetAccountData())
	RegistrationGroup.POST("/update_password", h.UpdateUserPassword())
	RegistrationGroup.POST("/user_data_by_login", h.GetUserDataByLogin())
}
