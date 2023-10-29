package grpc

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
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
	cfg *config.Config

	Hostname        string
	grpcServer      *grpc.Server
	httpServer      *http.Server
	subscriberCount int64
	logger          logger.Logger
}

func NewServer(cfg *config.Config, logger logger.Logger) *Server {
	return &Server{cfg: cfg, logger: logger}
}

// init check and prepare used config in Server struct
func (s *Server) Init() (err error) {
	if err = s.Close(); err != nil {
		return
	}

	if s.Hostname == "" {
		if s.Hostname, err = os.Hostname(); err != nil {
			return
		}
	}

	return
}

// Close server
func (s *Server) Close() (err error) {
	if s.grpcServer != nil {
		s.grpcServer.Stop()
		s.grpcServer = nil
	}
	if s.httpServer != nil {
		err = s.httpServer.Close()
		s.httpServer = nil
	}
	return
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

	s.httpServer = &http.Server{
		Addr:           s.cfg.HTTPServer.Port,
		ReadTimeout:    time.Second * s.cfg.HTTPServer.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.HTTPServer.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
	s.grpcServer = grpc.NewServer(grpcOpts...)

	if err := s.MapHandlers(); err != nil {
		return err
	}

	if s.cfg.HTTPServer.SSL {
		go func() {
			s.logger.Infof("Server is listening on PORT: %s", s.cfg.HTTPServer.Port)
			if err := s.httpServer.ListenAndServeTLS(certFile, keyFile); err != nil {
				s.logger.Fatalf("Error starting Server: ", err)
			}
		}()
	} else {
		go func() {
			s.logger.Infof("Server is listening on PORT: %s", s.cfg.HTTPServer.Port)
			if err := s.httpServer.ListenAndServe(); err != nil {
				s.logger.Fatalf("Error starting Server: ", err)
			}
		}()
	}

	go func() {
		s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.HTTPServer.PprofPort)
		if err := http.ListenAndServe(s.cfg.HTTPServer.PprofPort, http.DefaultServeMux); err != nil {
			s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.httpServer.Shutdown(ctx)
}
