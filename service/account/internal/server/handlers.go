package server

import (
	"github.com/GCFactory/dbo-system/platform/pkg/csrf"
	"github.com/GCFactory/dbo-system/platform/pkg/metric"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	delivery "github.com/GCFactory/dbo-system/service/account/internal/account/delivery/http"
	"github.com/GCFactory/dbo-system/service/account/internal/account/repository"
	"github.com/GCFactory/dbo-system/service/account/internal/account/usecase"
	accountMiddleware "github.com/GCFactory/dbo-system/service/account/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strings"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	metrics, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics Error: %s", err)
	}
	s.logger.Infof(
		"Metrics available URL: %s, ServiceName: %s",
		s.cfg.Metrics.URL,
		s.cfg.Metrics.ServiceName,
	)

	accRepo := repository.NewAccountRepository(s.db)

	accUC := usecase.NewAccountUseCase(s.cfg, accRepo, s.logger)

	accHandlers := delivery.NewACCHandlers(s.cfg, accUC, s.logger)

	////////////
	// handlerWithKakfa := m.NewHandler(s.cfg, usecase1, usecase2, ..., s.kafkaProducer, s.logger)
	////////////

	mw := accountMiddleware.NewMiddlewareManager(s.cfg, []string{"*"}, s.logger)

	e.Use(mw.RequestLoggerMiddleware)
	if s.cfg.Docs.Enable {
		//docs.SwaggerInfo.Title = s.cfg.Docs.Title
		//e.GET(fmt.Sprintf("/%s/*", s.cfg.Docs.Prefix), echoSwagger.WrapHandler)
	}

	if s.cfg.HTTPServer.SSL {
		e.Pre(middleware.HTTPSRedirect())
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrf.CSRFHeader},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())
	e.Use(mw.MetricsMiddleware(metrics))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("3M"))
	if s.cfg.HTTPServer.Debug {
		e.Use(mw.DebugMiddleware)
	}

	v1 := e.Group("/api/v1")

	health := e.Group("/health/ready")
	apiGatewayGroup := v1.Group("/account")

	delivery.MapACCRoutes(apiGatewayGroup, accHandlers, mw)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
