package http

import (
	"github.com/GCFactory/dbo-system/service/totp/internal/middleware"
	"github.com/GCFactory/dbo-system/service/totp/internal/totp"
	"github.com/labstack/echo/v4"
)

func MapTOTPRoutes(totpGroup *echo.Group, h totp.Handlers, mw *middleware.MiddlewareManager) {
	totpGroup.POST("/enroll", h.Enroll())
	totpGroup.POST("/verify", h.Verify())
	totpGroup.POST("/validate", h.Validate())
	//totpGroup.POST("/enable", h.Enable())
	//totpGroup.POST("/disable", h.Disable())
}
