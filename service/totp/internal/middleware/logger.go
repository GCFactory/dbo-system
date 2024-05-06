package middleware

import (
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/labstack/echo/v4"
	"time"
)

// Request logger middleware
func (mw *MiddlewareManager) RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		err := next(ctx)

		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		size := res.Size
		s := time.Since(start).String()
		requestID := utils.GetRequestID(ctx)

		mw.logger.Infof("Version: %v, RequestID: %s, Method: %s, URI: %s, Status: %v, Size: %v, Time: %s",
			req.ProtoMajor, requestID, req.Method, req.URL.String(), status, size, s,
		)
		return err
	}
}
