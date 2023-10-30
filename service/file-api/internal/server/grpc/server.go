package grpc

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	certFile       = "ssl/cert.crt"
	keyFile        = "ssl/private.pem"
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

type Server struct {
	cfg  *config.Config
	echo *echo.Echo

	grpcServer      *grpc.Server
	httpServer      *http.Server
	subscriberCount int64
	logger          logger.Logger
}

func NewServer(cfg *config.Config, logger logger.Logger) *Server {
	return &Server{echo: echo.New(), cfg: cfg, logger: logger}
}

func (s *Server) Run() error {
	grpcOpts := []grpc.ServerOption{}

	if s.cfg.HTTPServer.SSL {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			return err
		}
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	s.grpcServer = grpc.NewServer(grpcOpts...)

	if s.cfg.HTTPServer.SSL {
		if err := s.MapHandlers(s.echo); err != nil {
			return err
		}
		s.echo.Server.ReadTimeout = time.Second * s.cfg.HTTPServer.ReadTimeout
		s.echo.Server.WriteTimeout = time.Second * s.cfg.HTTPServer.WriteTimeout
		s.echo.Server.MaxHeaderBytes = maxHeaderBytes
		go func() {
			s.logger.Infof("Server is listening on PORT: %s", s.cfg.HTTPServer.Port)
			s.echo.Server.MaxHeaderBytes = maxHeaderBytes
			if err := s.echo.StartTLS(s.cfg.HTTPServer.Port, certFile, keyFile); err != nil {
				s.logger.Fatalf("Error starting TLS Server: ", err)
			}
		}()

		if s.cfg.HTTPServer.Debug {
			go func() {
				s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.HTTPServer.PprofPort)
				if err := http.ListenAndServe(s.cfg.HTTPServer.PprofPort, http.DefaultServeMux); err != nil {
					s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
				}
			}()
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		<-quit

		ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
		defer shutdown()
		defer s.grpcServer.Stop()

		s.logger.Info("Server Exited Properly")
		return s.echo.Shutdown(ctx)
	}

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

	if s.cfg.HTTPServer.Debug {
		go func() {
			s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.HTTPServer.PprofPort)
			if err := http.ListenAndServe(s.cfg.HTTPServer.PprofPort, http.DefaultServeMux); err != nil {
				s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
			}
		}()
	}

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()
	defer s.grpcServer.Stop()

	s.logger.Info("Server Exited Properly")
	return s.echo.Shutdown(ctx)
}
