package server

import (
	"fmt"
	"github.com/GCFactory/dbo-system/platform/pkg/csrf"
	"github.com/GCFactory/dbo-system/platform/pkg/metric"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	delivery "github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway/delivery/http"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway/repository"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway/usecase"
	apiMiddlewares "github.com/GCFactory/dbo-system/service/api_gateway/internal/middleware"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/models"
	"io/fs"
	"os"
	"path/filepath"
	"time"

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
	apiGatewayRepo := repository.NewApiGatewayRepository(s.redis)
	// Init useCases
	registrationServerInfo := &models.InternalServerInfo{}
	usersServerInfo := &models.InternalServerInfo{}
	accountsServerInfo := &models.InternalServerInfo{}
	notificationServerInfo := &models.InternalServerInfo{}
	totpServerInfo := &models.InternalServerInfo{}

	if regCfg, ok := s.cfg.InternalServices[ServerNameRegistration]; ok {
		registrationServerInfo.Port = regCfg.Port
		registrationServerInfo.Host = regCfg.Host
		registrationServerInfo.NumRetry = regCfg.Retry
		registrationServerInfo.WaitTimeRetry = time.Duration(time.Millisecond.Nanoseconds() * int64(regCfg.TimeWaitRetry))
		registrationServerInfo.TimeWaitResponse = time.Duration(time.Millisecond.Nanoseconds() * int64(regCfg.TimeWaitResponse))
	}

	if usCfg, ok := s.cfg.InternalServices[ServerNameUsers]; ok {
		usersServerInfo.Port = usCfg.Port
		usersServerInfo.Host = usCfg.Host
		usersServerInfo.NumRetry = usCfg.Retry
		usersServerInfo.WaitTimeRetry = time.Duration(time.Millisecond.Nanoseconds() * int64(usCfg.TimeWaitRetry))
		usersServerInfo.TimeWaitResponse = time.Duration(time.Millisecond.Nanoseconds() * int64(usCfg.TimeWaitResponse))
	}

	if accCfg, ok := s.cfg.InternalServices[ServerNameAccounts]; ok {
		accountsServerInfo.Port = accCfg.Port
		accountsServerInfo.Host = accCfg.Host
		accountsServerInfo.NumRetry = accCfg.Retry
		accountsServerInfo.WaitTimeRetry = time.Duration(time.Millisecond.Nanoseconds() * int64(accCfg.TimeWaitRetry))
		accountsServerInfo.TimeWaitResponse = time.Duration(time.Millisecond.Nanoseconds() * int64(accCfg.TimeWaitResponse))
	}

	if notifCfg, ok := s.cfg.InternalServices[ServerNameNotification]; ok {
		notificationServerInfo.Port = notifCfg.Port
		notificationServerInfo.Host = notifCfg.Host
		notificationServerInfo.NumRetry = notifCfg.Retry
		notificationServerInfo.WaitTimeRetry = time.Duration(time.Millisecond.Nanoseconds() * int64(notifCfg.TimeWaitRetry))
		notificationServerInfo.TimeWaitResponse = time.Duration(time.Millisecond.Nanoseconds() * int64(notifCfg.TimeWaitResponse))
	}

	if totpCfg, ok := s.cfg.InternalServices[ServerNameTotp]; ok {
		totpServerInfo.Port = totpCfg.Port
		totpServerInfo.Host = totpCfg.Host
		totpServerInfo.NumRetry = totpCfg.Retry
		totpServerInfo.WaitTimeRetry = time.Duration(time.Millisecond.Nanoseconds() * int64(totpCfg.TimeWaitRetry))
		totpServerInfo.TimeWaitResponse = time.Duration(time.Millisecond.Nanoseconds() * int64(totpCfg.TimeWaitResponse))
	}

	//executable, err := os.Executable()
	//if err != nil {
	//	panic(err)
	//}
	//exPath := filepath.Dir(executable)

	folderGrapthImagesPath := "./graph"

	ex, err := exists(folderGrapthImagesPath)
	if err != nil {
		return err
	}
	if ex {
		err = os.RemoveAll(folderGrapthImagesPath)
		if err != nil {
			return err
		}
	}
	err = os.MkdirAll(folderGrapthImagesPath, fs.ModePerm)
	if err != nil {
		return err
	}
	fmt.Println(filepath.Abs(folderGrapthImagesPath))

	folderQrImagesPath := "./qr"
	ex, err = exists(folderQrImagesPath)
	if err != nil {
		return err
	}
	if ex {
		err = os.RemoveAll(folderQrImagesPath)
		if err != nil {
			return err
		}
	}
	err = os.MkdirAll(folderQrImagesPath, fs.ModePerm)
	if err != nil {
		return err
	}
	fmt.Println(filepath.Abs(folderQrImagesPath))

	apiGatewayUsecase := usecase.NewApiGatewayUseCase(s.cfg, apiGatewayRepo, registrationServerInfo, usersServerInfo,
		accountsServerInfo, notificationServerInfo, totpServerInfo, folderGrapthImagesPath, folderQrImagesPath, s.rmqChan, s.rmqQueue)
	// Init handlers
	apiGatewayHalndlers := delivery.NewApiGatewayHandlers(s.cfg, s.logger, folderGrapthImagesPath, folderQrImagesPath, apiGatewayUsecase)

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
	apiGatewayGroup := v1.Group("/api_gateway")

	delivery.MapApiGatewayRoutes(apiGatewayGroup, apiGatewayHalndlers, mw)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}

func exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
