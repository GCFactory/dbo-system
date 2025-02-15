package server

import (
	"context"
	"errors"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/registration/config"
	accounts_api "github.com/GCFactory/dbo-system/service/registration/gen_proto/proto/api/account"
	users_api "github.com/GCFactory/dbo-system/service/registration/gen_proto/proto/api/users"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration/grpc"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration/repository"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration/usecase"
	"github.com/GCFactory/dbo-system/service/registration/pkg/kafka"
	"github.com/IBM/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
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
	grpcH         registration.RegistrationGRPCHandlers
	useCase       registration.UseCase
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
	RepoRegistration := repository.NewRegistrationRepository(
		server.db,
	)
	UCHandlers := usecase.NewRegistrationUseCase(
		cfg,
		RepoRegistration,
		server.logger,
	)
	grpcHandlers := grpc.NewRegistrationGRPCHandlers(
		cfg,
		kProducer,
		UCHandlers,
		server.logger,
	)
	server.grpcH = grpcHandlers
	server.useCase = UCHandlers
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
	// TODO: complete

	switch message.Topic {
	case grpc.ServerTopicUsersProducerRes:
		{
			event_success := &users_api.EventSuccess{}
			err := proto.Unmarshal(message.Value, event_success)
			if err != nil {
				s.logger.Errorf("Error unmarshalling event error: %v", err)
			}

			operation_name := event_success.OperationName

			saga_uuid, err := uuid.Parse(event_success.SagaUuid)
			if err != nil {
				s.logger.Errorf("Error parsing saga uuid: %v", err)
			}
			event_uuid, err := uuid.Parse(event_success.EventUuid)
			if err != nil {
				s.logger.Errorf("Error parsing event uuid: %v", err)
			}

			data := make(map[string]interface{})

			switch operation_name {
			case grpc.OperationCreateUser:
				{
					result := event_success.GetInfo()

					data["user_uuid"] = result

					break
				}
			case grpc.OperationGetUserData:
				{

					result := event_success.GetFullData()
					data["inn"] = result.GetUserInn()

					accounts := result.GetAccounts()
					if accounts != nil {
						data["accounts"] = accounts.Accounts
					}

					break
				}
			case grpc.OperationAddAccountToUser:
				{
					break
				}
			default:
				{
					s.logger.Errorf("Message was gotten from <%v> topic with unknown operation <%v>!", message.Topic, operation_name)
				}
			}

			err = s.grpcH.Process(context.Background(), saga_uuid, nil, event_uuid, nil, data, true)
			if err != nil {
				s.logger.Errorf("Error processing event: %v", err)
			}

			break

		}
	case grpc.ServerTopicUsersProducerErr:
		{
			event_error := &users_api.EventError{}
			err := proto.Unmarshal(message.Value, event_error)
			if err != nil {
				s.logger.Errorf("Error unmarshalling event error: %v", err)
			}

			saga_uuid, err := uuid.Parse(event_error.SagaUuid)
			if err != nil {
				s.logger.Errorf("Error parsing saga uuid: %v", err)
			}
			event_uuid, err := uuid.Parse(event_error.EventUuid)
			if err != nil {
				s.logger.Errorf("Error parsing event uuid: %v", err)
			}

			data := make(map[string]interface{})
			data["info"] = event_error.Info
			data["operation_name"] = event_error.OperationName
			data["status"] = event_error.Status
			err = s.grpcH.Process(context.Background(), saga_uuid, nil, event_uuid, nil, data, false)
			if err != nil {
				s.logger.Errorf("Error processing event: %v", err)
			}

			break
		}
	case grpc.ServerTopicAccountsProducerRes:
		{

			event_success := &accounts_api.EventStatus{}
			err := proto.Unmarshal(message.Value, event_success)
			if err != nil {
				s.logger.Errorf("Error unmarshalling event error: %v", err)
			}

			operation_name := event_success.OperationName

			saga_uuid, err := uuid.Parse(event_success.SagaUuid)
			if err != nil {
				s.logger.Errorf("Error parsing saga uuid: %v", err)
			}
			event_uuid, err := uuid.Parse(event_success.EventUuid)
			if err != nil {
				s.logger.Errorf("Error parsing event uuid: %v", err)
			}

			data := make(map[string]interface{})

			switch operation_name {
			case grpc.OperationReserveAcc:
				{
					data["acc_id"] = event_success.GetInfo()

					break
				}
			case grpc.OperationGetAccountData:
				{
					acc_data := event_success.GetAccData()
					if acc_data != nil {

						data["acc_status"] = acc_data.GetAccStatus()
						data["acc_cache"] = acc_data.GetAccMoneyAmount()

					}

					break
				}
			case grpc.OperationOpenAcc,
				grpc.OperationCreateAcc,
				grpc.OperationCloseAccount,
				grpc.OperationAddAccountCache,
				grpc.OperationWidthAccountCache:
				{
					break
				}
			default:
				{
					s.logger.Errorf("Message was gotten from <%v> topic with unknown operation <%v>!", message.Topic, operation_name)
				}
			}

			err = s.grpcH.Process(context.Background(), saga_uuid, nil, event_uuid, nil, data, true)
			if err != nil {
				s.logger.Errorf("Error processing event: %v", err)
			}

			break
		}
	case grpc.ServerTopicAccountsProducerErr:
		{
			event_error := &accounts_api.EventError{}
			err := proto.Unmarshal(message.Value, event_error)
			if err != nil {
				s.logger.Errorf("Error unmarshalling event error: %v", err)
			}

			saga_uuid, err := uuid.Parse(event_error.SagaUuid)
			if err != nil {
				s.logger.Errorf("Error parsing saga uuid: %v", err)
			}
			event_uuid, err := uuid.Parse(event_error.EventUuid)
			if err != nil {
				s.logger.Errorf("Error parsing event uuid: %v", err)
			}

			data := make(map[string]interface{})
			data["info"] = event_error.Info
			data["operation_name"] = event_error.OperationName
			data["status"] = event_error.Status
			err = s.grpcH.Process(context.Background(), saga_uuid, nil, event_uuid, nil, data, false)
			if err != nil {
				s.logger.Errorf("Error processing event: %v", err)
			}

			break
		}
	default:
		{
			s.logger.Errorf("Message was gotten from wrong topic <%v>!", message.Topic)
		}
	}

	return nil

}
