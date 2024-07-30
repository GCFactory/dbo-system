package server

import (
	"context"
	"errors"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/registration/config"
	"github.com/GCFactory/dbo-system/service/registration/pkg/kafka"
	"github.com/IBM/sarama"
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
	echo          *echo.Echo
	kafkaProducer *kafka.ProducerProvider
	kafkaConsumer *kafka.ConsumerGroup
	cfg           *config.Config
	db            *sqlx.DB
	logger        logger.Logger
	// Channel to control goroutines
	kafkaConsumerChan chan int
}

func NewServer(cfg *config.Config, kConsumer *kafka.ConsumerGroup, kProducer *kafka.ProducerProvider, db *sqlx.DB, logger logger.Logger) *Server {
	server := Server{
		echo:              echo.New(),
		cfg:               cfg,
		db:                db,
		logger:            logger,
		kafkaConsumer:     kConsumer,
		kafkaConsumerChan: make(chan int, 3),
		kafkaProducer:     kProducer,
	}
	server.echo.HidePort = true
	server.echo.HideBanner = true
	return &server
}

func (s *Server) Run() error {
	ctxWithCancel, cancel := context.WithCancel(context.Background())
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go s.RunKafkaConsumer(ctxWithCancel, s.kafkaConsumerChan, s.kafkaConsumer.Consumer)
	//// Remove example
	go func() {
		for i := 0; i < 10; i++ {
			err := s.kafkaProducer.ProduceRecord("test", []byte("test message"))
			if err != nil {
				s.logger.Warnf("Error on produce: %v", err)
			}
			time.Sleep(time.Second * 10)
		}
	}()
	//// Remove example

	for {
		select {
		// If goroutine exit for any reason - restart
		case <-s.kafkaConsumerChan:
			s.logger.Warn("Received end of kafka consumer goroutine")
			go s.RunKafkaConsumer(ctxWithCancel, s.kafkaConsumerChan, s.kafkaConsumer.Consumer)
		// Handle system interrupts
		case <-quit:
			ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
			s.kafkaConsumer.Consumer.PauseAll()
			s.kafkaConsumer.Consumer.Close()
			time.Sleep(time.Second * 15)
			cancel()
			close(quit)
			close(s.kafkaConsumerChan)
			defer shutdown()

			s.logger.Info("Server Exited Properly")
			return s.echo.Server.Shutdown(ctx)
		}
	}
}

func (s *Server) RunKafkaConsumer(ctx context.Context, quitChan chan<- int, client sarama.ConsumerGroup) {
	consumer := Consumer{
		ready: make(chan bool),
		handlerFunc: func(message *sarama.ConsumerMessage) error {
			/// DO IT
			return nil
		},
	}
	for {
		if err := client.Consume(ctx, s.cfg.KafkaConsumer.Topics, &consumer); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return
			}
			quitChan <- 1
			s.logger.Errorf("Error from consumer: %v", err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			s.logger.Infof("Stopping kafka consumer: context close: %v", ctx.Err())
			return
		}
		s.logger.Info("Stopping kafka consumer...")
		consumer.ready = make(chan bool)
	}
}
