package server

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/api_gateway/config"
	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server struct
type Server struct {
	echo     *echo.Echo
	cfg      *config.Config
	redis    *redis.Client
	logger   logger.Logger
	rmqChan  *amqp091.Channel
	rmqQueue amqp091.Queue
}

// NewServer New Server constructor
func NewServer(cfg *config.Config, redis *redis.Client, rmqChan *amqp091.Channel, rmqQueue amqp091.Queue, logger logger.Logger) *Server {
	return &Server{echo: echo.New(), cfg: cfg, redis: redis, logger: logger, rmqChan: rmqChan, rmqQueue: rmqQueue}
}

const (
	certFile       = "ssl/Server.crt"
	keyFile        = "ssl/Server.pem"
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

func (s *Server) Run() error {
	//Сервер с включённым TLS (SSL)
	if s.cfg.HTTPServer.SSL {
		if err := s.MapHandlers(s.echo); err != nil {
			return err
		}

		s.echo.Server.ReadTimeout = time.Second * s.cfg.HTTPServer.ReadTimeout
		s.echo.Server.WriteTimeout = time.Second * s.cfg.HTTPServer.WriteTimeout

		//Разворачиваем TLS
		go func() {
			s.logger.Infof("Server is listening on PORT: %s", s.cfg.HTTPServer.Port)
			s.echo.Server.MaxHeaderBytes = maxHeaderBytes
			if err := s.echo.StartTLS(s.cfg.HTTPServer.Port, certFile, keyFile); err != nil {
				s.logger.Fatalf("Error starting TLS Server: ", err)
			}
		}()

		//Наинаем слушать запросы
		go func() {
			s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.HTTPServer.PprofPort)
			if err := http.ListenAndServe(s.cfg.HTTPServer.PprofPort, http.DefaultServeMux); err != nil {
				s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
			}
		}()

		//Подвесили обработчик прерывания
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		<-quit

		ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
		defer shutdown()

		s.logger.Info("Server Exited Properly")
		return s.echo.Server.Shutdown(ctx)
	}
	//Аналог сервера, только без SSL
	server := &http.Server{
		Addr:           s.cfg.HTTPServer.Port,
		ReadTimeout:    time.Second * s.cfg.HTTPServer.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.HTTPServer.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.HTTPServer.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.logger.Fatalf("Error starting Server: ", err)
		}
	}()

	go func() {
		s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.HTTPServer.PprofPort)
		if err := http.ListenAndServe(s.cfg.HTTPServer.PprofPort, http.DefaultServeMux); err != nil {
			s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
		}
	}()

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
