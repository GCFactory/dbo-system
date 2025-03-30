package server

import (
	"fmt"
	"github.com/GCFactory/dbo-system/platform/pkg/csrf"
	"github.com/GCFactory/dbo-system/platform/pkg/metric"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/totp/pkg/otp"

	//"github.com/GCFactory/dbo-system/service/totp/docs"
	apiMiddlewares "github.com/GCFactory/dbo-system/service/totp/internal/middleware"
	totpHttp "github.com/GCFactory/dbo-system/service/totp/internal/totp/delivery/http"
	totpRepository "github.com/GCFactory/dbo-system/service/totp/internal/totp/repository"
	totpUsecase "github.com/GCFactory/dbo-system/service/totp/internal/totp/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
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

	// Init repositories
	tRepo := totpRepository.NewTOTPRepository(s.db)
	//sRepo := sessionRepository.NewSessionRepository(s.redisClient, s.cfg)
	//newsRedisRepo := newsRepository.NewNewsRedisRepo(s.redisClient)

	tLogic := otp.NewTOTPStruct(s.cfg, s.logger)

	// Init useCases
	tUC := totpUsecase.NewTOTPUseCase(s.cfg, tRepo, tLogic, s.logger)
	//authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, authRedisRepo, s.logger)

	// Init handlers
	tHandlers := totpHttp.NewTOTPHandlers(s.cfg, tUC, s.logger)
	//authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, sessUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(s.cfg, []string{"*"}, s.logger)

	e.Use(mw.RequestLoggerMiddleware)

	if s.cfg.Docs.Enable {
		//docs.SwaggerInfo.Title = s.cfg.Docs.Title
		e.GET(fmt.Sprintf("/%s/*", s.cfg.Docs.Prefix), echoSwagger.WrapHandler)
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
	tGroup := v1.Group("/totp")
	//authGroup := v1.Group("/auth")

	//authHttp.MapAuthRoutes(authGroup, authHandlers, mw)
	totpHttp.MapTOTPRoutes(tGroup, tHandlers, mw)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
