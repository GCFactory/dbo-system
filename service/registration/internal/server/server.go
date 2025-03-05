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
	kafkaConsumerChan  chan int
	sagaSyncChan       chan bool
	sagaProcessingSaga map[uuid.UUID][]*sarama.ConsumerMessage
}

func NewServer(cfg *config.Config, kConsumer *kafka.ConsumerGroup, kProducer *kafka.ProducerProvider, db *sqlx.DB, logger logger.Logger) *Server {
	server := Server{
		echo:               echo.New(),
		cfg:                cfg,
		db:                 db,
		logger:             logger,
		kafkaConsumer:      kConsumer,
		kafkaConsumerChan:  make(chan int, 3),
		kafkaProducer:      kProducer,
		sagaSyncChan:       make(chan bool, 1),
		sagaProcessingSaga: make(map[uuid.UUID][]*sarama.ConsumerMessage),
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
	server.sagaSyncChan <- true
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

func (s *Server) handleData(message *sarama.ConsumerMessage) (err error) {

	err = nil

	var saga_uuid uuid.UUID = uuid.Nil
	var event_uuid uuid.UUID = uuid.Nil
	var success bool = false
	var data map[string]interface{} = make(map[string]interface{})

	switch message.Topic {
	case grpc.ServerTopicUsersProducerRes:
		{
			event_success := &users_api.EventSuccess{}
			err := proto.Unmarshal(message.Value, event_success)
			if err != nil {
				s.logger.Errorf("Error unmarshalling event error: %v", err)
			}

			operation_name := event_success.OperationName

			saga_uuid, err = uuid.Parse(event_success.SagaUuid)
			if err != nil {
				s.logger.Errorf("Error parsing saga uuid: %v", err)
			}
			event_uuid, err = uuid.Parse(event_success.EventUuid)
			if err != nil {
				s.logger.Errorf("Error parsing event uuid: %v", err)
			}

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
					data["user_id"] = result.GetUserId()
					data["user_login"] = result.GetUserLogin()

					accounts := result.GetAccounts()
					if accounts != nil {
						data["accounts"] = accounts.Accounts
					}

					passport := result.GetPassport()
					if passport != nil {
						data["passport_series"] = passport.GetSeries()
						data["passport_number"] = passport.GetNumber()

						fcs := passport.GetFcs()
						if fcs != nil {
							data["passport_first_name"] = fcs.GetName()
							data["passport_first_surname"] = fcs.GetSurname()
							data["passport_first_patronimic"] = fcs.GetPatronymic()
						}

						data["passport_birth_date"] = passport.GetBirthDate().AsTime().Format("2006-01-02 15:04:05 -07:00:00")
						data["passport_birth_location"] = passport.GetBirthLocation()
						data["passport_pick_up_point"] = passport.GetPickUpPoint()
						data["passport_authority"] = passport.GetAuthority()
						data["passport_authority_date"] = passport.GetAuthorityDate().AsTime().Format("2006-01-02 15:04:05 -07:00:00")
						data["passport_registration_address"] = passport.GetRegistrationAdress()

					}

					break
				}
			case grpc.OperationGetUserDataByLogin:
				{

					result := event_success.GetFullData()
					data["inn"] = result.GetUserInn()
					data["user_id"] = result.GetUserId()
					data["user_login"] = result.GetUserLogin()

					accounts := result.GetAccounts()
					if accounts != nil {
						data["accounts"] = accounts.Accounts
					}

				}
			case grpc.OperationAddAccountToUser,
				grpc.OperationRemoveUserAccount,
				grpc.OperationCheckUSerPassword:
				{
					break
				}
			default:
				{
					s.logger.Errorf("Message was gotten from <%v> topic with unknown operation <%v>!", message.Topic, operation_name)
				}
			}

			success = true

			break

		}
	case grpc.ServerTopicUsersProducerErr:
		{
			event_error := &users_api.EventError{}
			err := proto.Unmarshal(message.Value, event_error)
			if err != nil {
				s.logger.Errorf("Error unmarshalling event error: %v", err)
			}

			saga_uuid, err = uuid.Parse(event_error.SagaUuid)
			if err != nil {
				s.logger.Errorf("Error parsing saga uuid: %v", err)
			}
			event_uuid, err = uuid.Parse(event_error.EventUuid)
			if err != nil {
				s.logger.Errorf("Error parsing event uuid: %v", err)
			}

			data["info"] = event_error.Info
			data["operation_name"] = event_error.OperationName
			data["status"] = event_error.Status

			success = false

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

			saga_uuid, err = uuid.Parse(event_success.SagaUuid)
			if err != nil {
				s.logger.Errorf("Error parsing saga uuid: %v", err)
			}
			event_uuid, err = uuid.Parse(event_success.EventUuid)
			if err != nil {
				s.logger.Errorf("Error parsing event uuid: %v", err)
			}

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
						data["acc_cache_value"] = acc_data.GetAccMoneyValue()

						account_details := acc_data.GetAccDetails()

						if account_details != nil {
							data["acc_culc_number"] = account_details.GetCulcNumber()
							data["acc_corr_number"] = account_details.GetCorrNumber()
							data["acc_bic"] = account_details.GetBic()
							data["acc_cio"] = account_details.GetCio()
							data["acc_reserve_reason"] = account_details.GetReserveReason()
						}

					}

					break
				}
			case grpc.OperationOpenAcc,
				grpc.OperationCreateAcc,
				grpc.OperationCloseAccount,
				grpc.OperationAddAccountCache,
				grpc.OperationWidthAccountCache,
				grpc.OperationRemoveAccount:
				{
					break
				}
			default:
				{
					s.logger.Errorf("Message was gotten from <%v> topic with unknown operation <%v>!", message.Topic, operation_name)
				}
			}

			success = true

			break
		}
	case grpc.ServerTopicAccountsProducerErr:
		{
			event_error := &accounts_api.EventError{}
			err := proto.Unmarshal(message.Value, event_error)
			if err != nil {
				s.logger.Errorf("Error unmarshalling event error: %v", err)
			}

			saga_uuid, err = uuid.Parse(event_error.SagaUuid)
			if err != nil {
				s.logger.Errorf("Error parsing saga uuid: %v", err)
			}
			event_uuid, err = uuid.Parse(event_error.EventUuid)
			if err != nil {
				s.logger.Errorf("Error parsing event uuid: %v", err)
			}

			data["info"] = event_error.Info
			data["operation_name"] = event_error.OperationName
			data["status"] = event_error.Status

			success = false

			break
		}
	default:
		{
			s.logger.Errorf("Message was gotten from wrong topic <%v>!", message.Topic)
		}
	}

	s.logger.Debug("INCOMING:", message.Topic, "|", event_uuid, "|", success)

	<-s.sagaSyncChan

	if s.checkSagaProcessing(saga_uuid) {
		s.logger.Errorf("Error processing event: %v", errors.New("Saga already in process, event uuid:"), event_uuid.String())
		s.addSagaToProcess(saga_uuid, message)
		s.sagaSyncChan <- true
		return
	}

	s.sagaSyncChan <- true

	s.logger.Debug("Processing event:", event_uuid)
	err = s.grpcH.Process(context.Background(), saga_uuid, nil, event_uuid, nil, data, success)
	if err != nil {
		s.logger.Errorf("Error processing event: %v", err)
	}

	<-s.sagaSyncChan

	s.removeSagaFromProcess(saga_uuid)

	s.sagaSyncChan <- true

	s.startSagaProcessAgain(saga_uuid)

	s.logger.Debug("End processing:", event_uuid)

	return nil

}

func (s *Server) addSagaToProcess(saga_uuid uuid.UUID, message *sarama.ConsumerMessage) {

	if !s.checkSagaProcessing(saga_uuid) {

		s.sagaProcessingSaga[saga_uuid] = append(s.sagaProcessingSaga[saga_uuid], message)

	}

}

func (s *Server) checkSagaProcessing(saga_uuid uuid.UUID) bool {

	_, exist := s.sagaProcessingSaga[saga_uuid]

	return exist

}

func (s *Server) removeSagaFromProcess(saga_uuid uuid.UUID) {

	if s.checkSagaProcessing(saga_uuid) {

		delete(s.sagaProcessingSaga, saga_uuid)

	}

}

func (s *Server) startSagaProcessAgain(saga_uuid uuid.UUID) {

	if s.checkSagaProcessing(saga_uuid) {

		messages, ok := s.sagaProcessingSaga[saga_uuid]
		if ok {
			_ = s.handleData(messages[0])
		}

	}

}
