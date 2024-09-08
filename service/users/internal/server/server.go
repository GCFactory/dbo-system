package server

import (
	"context"
	"errors"
	platfrom_config "github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/users/config"
	"github.com/GCFactory/dbo-system/service/users/internal/users"
	"github.com/GCFactory/dbo-system/service/users/internal/users/repository"
	"github.com/GCFactory/dbo-system/service/users/internal/users/usecase"
	"github.com/golang/protobuf/proto"
	"github.com/labstack/echo/v4"

	"github.com/GCFactory/dbo-system/service/users/gen_proto/proto/api"
	"github.com/GCFactory/dbo-system/service/users/internal/users/grpc_handlers"
	"github.com/GCFactory/dbo-system/service/users/pkg/kafka"
	"github.com/IBM/sarama"
	"github.com/jmoiron/sqlx"
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
	grpcHandlers      users.GRPCHandlers
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
	usersRepo := repository.NewUserRepository(db)
	usersUC := usecase.NewUsersUseCase(&platfrom_config.Config{
		Env:      cfg.Env,
		Logger:   cfg.Logger,
		App:      cfg.App,
		Postgres: cfg.Postgres,
		Version:  cfg.Version,
	}, usersRepo, logger)
	server.grpcHandlers = grpc_handlers.NewUsersGRPCHandlers(cfg, kProducer, usersUC, logger)
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

	go s.RunKafkaConsumer(ctxWithCancel, s.kafkaConsumerChan)

	for {
		select {
		// If goroutine exit for any reason - restart
		case <-s.kafkaConsumerChan:
			s.logger.Warn("Received end of kafka consumer goroutine")
			go s.RunKafkaConsumer(ctxWithCancel, s.kafkaConsumerChan)
		// Handle system interrupts
		case <-quit:
			ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
			s.kafkaConsumer.Consumer.PauseAll()
			s.kafkaConsumer.Consumer.Close()
			time.Sleep(time.Second * 5)
			close(quit)
			close(s.kafkaConsumerChan)
			defer shutdown()
			defer cancel()

			s.logger.Info("Server Exited Properly")
			return s.echo.Server.Shutdown(ctx)
		}
	}
}

func (s *Server) RunKafkaConsumer(ctx context.Context, quitChan chan<- int) {
	consumer := kafka.Consumer{
		Ready:       make(chan bool),
		HandlerFunc: s.handleData,
	}
	for {
		if err := s.kafkaConsumer.Consumer.Consume(ctx, s.cfg.KafkaConsumer.Topics, &consumer); err != nil {
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
		consumer.Ready = make(chan bool)
	}
}

func (s *Server) handleData(message *sarama.ConsumerMessage) error {

	data := &api.EventData{}
	err := proto.Unmarshal(message.Value, data)

	if err != nil {

		s.logger.Error(grpc_handlers.ErrorInvalidInputData.Error())

		answer := &api.EventError{
			Info:   grpc_handlers.ErrorInvalidInputData.Error(),
			Status: grpc_handlers.GetErrorCode(grpc_handlers.ErrorInvalidInputData),
		}
		answer_data, err := proto.Marshal(answer)
		if err != nil {
			s.logger.Error(err)
			return err
		}
		err = s.kafkaProducer.ProduceRecord(grpc_handlers.TopicError, sarama.ByteEncoder(answer_data))
		if err != nil {
			s.logger.Error(err)
			return err
		}
		return nil
	}
	//
	switch data.GetOperationName() {
	case grpc_handlers.AddUser:
		{
			// Unpack data and handle func
			if extracted_data := data.GetUserInfo(); extracted_data != nil {
				if err = s.grpcHandlers.AddUser(context.Background(), data.GetSagaUuid(), data.GetEventUuid(), extracted_data, s.kafkaProducer); err != nil {
					return err
				}
			} else {
				err = grpc_handlers.ErrorUnkonwnOperationType
			}
		}
	case grpc_handlers.GetUserData:
		{
			// Unpack data and handle func
			if extracted_data := data.GetAdditionalInfo(); extracted_data != nil {
				if err = s.grpcHandlers.GetUserData(context.Background(), data.GetSagaUuid(), data.GetEventUuid(), extracted_data, s.kafkaProducer); err != nil {
					return err
				}
			} else {
				err = grpc_handlers.ErrorUnkonwnOperationType
			}
		}
	case grpc_handlers.UpdateUsersPassport:
		{
			// Unpack data and handle func
			if extracted_data := data.GetAdditionalInfo(); extracted_data != nil {
				if err = s.grpcHandlers.UpdateUsersPassport(context.Background(), data.GetSagaUuid(), data.GetEventUuid(), extracted_data, s.kafkaProducer); err != nil {
					return err
				}
			} else {
				err = grpc_handlers.ErrorUnkonwnOperationType
			}
		}
	case grpc_handlers.AddUserAccounts:
		{
			// Unpack data and handle func
			if extracted_data := data.GetAdditionalInfo(); extracted_data != nil {
				if err = s.grpcHandlers.AddUserAccount(context.Background(), data.GetSagaUuid(), data.GetEventUuid(), extracted_data, s.kafkaProducer); err != nil {
					return err
				}
			} else {
				err = grpc_handlers.ErrorUnkonwnOperationType
			}
		}
	case grpc_handlers.GetUsersAccounts:
		{
			// Unpack data and handle func
			if extracted_data := data.GetAdditionalInfo(); extracted_data != nil {
				if err = s.grpcHandlers.GetUsersAccounts(context.Background(), data.GetSagaUuid(), data.GetEventUuid(), extracted_data, s.kafkaProducer); err != nil {
					return err
				}
			} else {
				err = grpc_handlers.ErrorUnkonwnOperationType
			}
		}
	default:
		{
			answer := &api.EventError{
				SagaUuid:      data.GetSagaUuid(),
				EventUuid:     data.GetEventUuid(),
				Info:          grpc_handlers.ErrorUnkonwnOperationType.Error(),
				Status:        grpc_handlers.GetErrorCode(grpc_handlers.ErrorUnkonwnOperationType),
				OperationName: data.GetOperationName(),
			}
			answer_data, err := proto.Marshal(answer)
			if err != nil {
				s.logger.Error(err)
				return err
			}
			err = s.kafkaProducer.ProduceRecord(grpc_handlers.TopicError, sarama.ByteEncoder(answer_data))
			if err != nil {
				s.logger.Error(err)
				return err
			}
		}
	}

	if err != nil {

		s.logger.Error(err)

		answer := &api.EventError{
			Info:   err.Error(),
			Status: grpc_handlers.GetErrorCode(grpc_handlers.ErrorInvalidInputData),
		}
		answer_data, err := proto.Marshal(answer)
		if err != nil {
			s.logger.Error(err)
			return err
		}
		err = s.kafkaProducer.ProduceRecord(grpc_handlers.TopicError, sarama.ByteEncoder(answer_data))
		if err != nil {
			s.logger.Error(err)
			return err
		}
		return nil
	}

	return nil
}
