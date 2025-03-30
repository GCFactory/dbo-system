package server

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/notification/internal/notification"
	"github.com/GCFactory/dbo-system/service/notification/internal/notification/repo"
	"github.com/GCFactory/dbo-system/service/notification/internal/notification/usecase"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/smtp"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"net/http"
	"net/http/pprof"
	"time"
)

const (
	certFile       = "ssl/Server.crt"
	keyFile        = "ssl/Server.pem"
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

// Server struct
type Server struct {
	echo       *echo.Echo
	cfg        *config.Config
	db         *sqlx.DB
	logger     logger.Logger
	msg        <-chan amqp.Delivery
	useCase    notification.UseCase
	smtpClient *smtp.Client
}

func NewServer(cfg *config.Config, db *sqlx.DB, msgChan <-chan amqp.Delivery, smtpClient *smtp.Client, logger logger.Logger) *Server {
	server := Server{
		echo:       echo.New(),
		cfg:        cfg,
		db:         db,
		msg:        msgChan,
		logger:     logger,
		smtpClient: smtpClient,
	}
	server.echo.HidePort = true
	server.echo.HideBanner = true

	serverRepo := repo.NewNotificationRepository(server.db)
	server.useCase = usecase.NewNotificationUseCase(serverRepo, server.smtpClient, server.cfg.NotificationSmtp)

	return &server
}

func (s *Server) Run() error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := &http.Server{
		Addr:           s.cfg.HTTPServer.Port,
		ReadTimeout:    time.Second * s.cfg.HTTPServer.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.HTTPServer.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.HTTPServer.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.logger.Fatalf("Error starting Server: %s", err)
		}
	}()

	go func() {
		r := http.NewServeMux()
		// Регистрация pprof-обработчиков
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)

		s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.HTTPServer.PprofPort)
		if err := http.ListenAndServe(s.cfg.HTTPServer.PprofPort, r); err != nil {
			s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
		}
	}()

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	go func() {
		for d := range s.msg {
			err := s.useCase.SendMessage(context.Background(), d)
			if err != nil {
				s.logger.Errorf("Error send msg: %s", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
