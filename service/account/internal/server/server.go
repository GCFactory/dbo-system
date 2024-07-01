package server

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
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
	echo   *echo.Echo
	cfg    *config.Config
	db     *sqlx.DB
	logger logger.Logger
	// Channel to control goroutines
	kafkaProducerChan chan int
	kafkaConsumerChan chan int
}

func NewServer(cfg *config.Config, db *sqlx.DB, logger logger.Logger) *Server {
	return &Server{
		echo:              echo.New(),
		cfg:               cfg,
		db:                db,
		logger:            logger,
		kafkaProducerChan: make(chan int, 3),
		kafkaConsumerChan: make(chan int, 3),
	}
}

func (s *Server) Run() error {
	ctxWithCancel, cancel := context.WithCancel(context.Background())

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select { // If goroutine exit for any reason - restart
		case <-s.kafkaProducerChan:
			s.logger.Warn("Received end of kafka producer goroutine")
			go s.RunKafkaProducer(&ctxWithCancel, &s.kafkaProducerChan)
		// If goroutine exit for any reason - restart
		case <-s.kafkaConsumerChan:
			s.logger.Warn("Received end of kafka consumer goroutine")
			go s.RunKafkaProducer(&ctxWithCancel, &s.kafkaConsumerChan)
		// Handle system interrupts
		case <-quit:
			ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
			cancel()
			defer close(quit)
			defer close(s.kafkaProducerChan)
			defer shutdown()

			s.logger.Info("Server Exited Properly")
			return s.echo.Server.Shutdown(ctx)
		}
	}
}

func (s *Server) RunKafkaProducer(ctx *context.Context, quitChan *chan int) {

}

func (s *Server) RunKafkaConsumer(ctx *context.Context, quitChan *chan int) {

}
