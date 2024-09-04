package server

func (s *Server) MapHandlers(e *echo.Echo) error {
	//metrics, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)
	//if err != nil {
	//	s.logger.Errorf("CreateMetrics Error: %s", err)
	//}
	//s.logger.Infof(
	//	"Metrics available URL: %s, ServiceName: %s",
	//	s.cfg.Metrics.URL,
	//	s.cfg.Metrics.ServiceName,
	//)
	// todo: сверху раскоментировать

	//sRepo := sessionRepository.NewSessionRepository(s.redisClient, s.cfg)
	//newsRedisRepo := newsRepository.NewNewsRedisRepo(s.redisClient)

	//authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, authRedisRepo, s.logger)

	//authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, sessUC, s.logger)

	////////////
	// handlerWithKakfa := m.NewHandler(s.cfg, usecase1, usecase2, ..., s.kafkaProducer, s.logger)
	////////////

	//	todo: снизу раскоментить
	//mw := apiMiddlewares.NewMiddlewareManager(s.cfg, []string{"*"}, s.logger)
	//
	//e.Use(mw.RequestLoggerMiddleware)
	//if s.cfg.Docs.Enable {
	//	//docs.SwaggerInfo.Title = s.cfg.Docs.Title
	//	//e.GET(fmt.Sprintf("/%s/*", s.cfg.Docs.Prefix), echoSwagger.WrapHandler)
	//}
	//
	//if s.cfg.HTTPServer.SSL {
	//	e.Pre(middleware.HTTPSRedirect())
	//}
	//
	//e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//	AllowOrigins: []string{"*"},
	//	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrf.CSRFHeader},
	//}))
	//e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
	//	StackSize:         1 << 10, // 1 KB
	//	DisablePrintStack: true,
	//	DisableStackAll:   true,
	//}))
	//e.Use(middleware.RequestID())
	//e.Use(mw.MetricsMiddleware(metrics))
	//
	//e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
	//	Level: 5,
	//	Skipper: func(c echo.Context) bool {
	//		return strings.Contains(c.Request().URL.Path, "swagger")
	//	},
	//}))
	//e.Use(middleware.Secure())
	//e.Use(middleware.BodyLimit("3M"))
	//if s.cfg.HTTPServer.Debug {
	//	e.Use(mw.DebugMiddleware)
	//}

	return nil
}
