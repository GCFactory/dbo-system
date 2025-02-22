package grpc_handlers

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/users/config"
	"github.com/GCFactory/dbo-system/service/users/gen_proto/proto/api"
	"github.com/GCFactory/dbo-system/service/users/gen_proto/proto/platform"
	"github.com/GCFactory/dbo-system/service/users/internal/models"
	"github.com/GCFactory/dbo-system/service/users/internal/users"
	"github.com/GCFactory/dbo-system/service/users/pkg/kafka"
	"github.com/IBM/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"time"
)

type UsersGrpcHandlers struct {
	cfg       *config.Config
	kProducer *kafka.ProducerProvider
	usersUC   users.UseCase
	accLog    logger.Logger
}

func (usersGRPC UsersGrpcHandlers) AddUser(ctx context.Context, saga_uuid string, event_uuid string, users_data *api.UserInfo, kProducer *kafka.ProducerProvider) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "usersGRPC.AddUser")
	defer span.Finish()

	var answer_data []byte
	flag_error := false
	var err error
	answer_topic := TopicResult

	answer := &api.EventSuccess{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: AddUser,
	}

	error_answer := &api.EventError{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: AddUser,
	}

	passport := users_data.GetPassport()
	fcs := passport.GetFcs()

	passport_uuid := uuid.New()

	user := &models.User_full_data{
		User: &models.User{
			User_uuid:     uuid.New(),
			User_inn:      users_data.GetUserInn(),
			User_accounts: models.ListOfAccounts{},
			Passport_uuid: passport_uuid,
		},
		Passport: &models.Passport{
			Passport_uuid:       passport_uuid,
			Passport_number:     passport.GetNumber(),
			Passport_series:     passport.GetSeries(),
			Name:                fcs.GetName(),
			Surname:             fcs.GetSurname(),
			Patronimic:          fcs.GetPatronymic(),
			Birth_date:          passport.GetBirthDate().AsTime().Format("2006-01-02 15:04:05 -07:00:00"),
			Birth_location:      passport.GetBirthLocation(),
			Pick_up_point:       passport.GetPickUpPoint(),
			Authority_date:      passport.GetAuthorityDate().AsTime().Format("2006-01-02 15:04:05 -07:00:00"),
			Authority:           passport.GetAuthority(),
			Registration_adress: passport.GetRegistrationAdress(),
		},
	}

	if err = usersGRPC.usersUC.AddUser(ctxWithTrace, user); err != nil {
		usersGRPC.accLog.Error(err)
		flag_error = true
		answer_topic = TopicError

		error_answer.Info = err.Error()
		error_answer.Status = GetErrorCode(err)
	} else {
		answer.Result = &api.EventSuccess_Info{user.User.User_uuid.String()}
	}

	if flag_error {
		answer_data, err = proto.Marshal(error_answer)
	} else {
		answer_data, err = proto.Marshal(answer)
	}

	if err != nil {
		usersGRPC.accLog.Error(err)
		return err
	}

	err = kProducer.ProduceRecord(answer_topic, sarama.ByteEncoder(answer_data))
	if err != nil {
		usersGRPC.accLog.Error(err)
		return err
	}
	usersGRPC.accLog.Info("Success to send answer!")

	return nil
}

func (usersGRPC UsersGrpcHandlers) GetUserData(ctx context.Context, saga_uuid string, event_uuid string, operation_details *api.OperationDetails, kProducer *kafka.ProducerProvider) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "usersGRPC.GetUserData")
	defer span.Finish()

	var answer_data []byte
	flag_error := false
	var err error
	answer_topic := TopicResult

	answer := &api.EventSuccess{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: GetUserData,
	}

	error_answer := &api.EventError{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: GetUserData,
	}

	var result *models.User_full_data

	user_uuid, err := uuid.Parse(operation_details.GetUserUuid())
	if err != nil {
		error_answer.Info = err.Error()
		error_answer.Status = GetErrorCode(ErrorInvalidInputData)

		flag_error = true
	}

	if result, err = usersGRPC.usersUC.GetUserData(ctxWithTrace, user_uuid); err != nil {
		usersGRPC.accLog.Error(err)
		flag_error = true
		answer_topic = TopicError

		error_answer.Info = err.Error()
		error_answer.Status = GetErrorCode(err)
	} else {
		passport := result.Passport
		user := result.User

		var accounts []string

		for i := 0; i < len(user.User_accounts.Data); i++ {
			accounts = append(accounts, user.User_accounts.Data[i].String())
		}

		birth_date, err := time.Parse(time.RFC3339, passport.Birth_date)
		if err != nil {
			return err
		}

		authority_date, err := time.Parse(time.RFC3339, passport.Authority_date)
		if err != nil {
			return err
		}

		answer.Result = &api.EventSuccess_FullData{
			FullData: &api.FullData{
				UserInn: user.User_inn,
				Accounts: &api.ListOfAccounts{
					Accounts: accounts,
				},
				Passport: &platform.Passport{
					Authority: passport.Authority,
					Fcs: &platform.FCs{
						Name:       passport.Name,
						Surname:    passport.Surname,
						Patronymic: passport.Patronimic,
					},
					BirthLocation:      passport.Birth_location,
					Number:             passport.Passport_number,
					Series:             passport.Passport_series,
					RegistrationAdress: passport.Registration_adress,
					PickUpPoint:        passport.Pick_up_point,
					BirthDate: &timestamp.Timestamp{
						Seconds: int64(birth_date.Second()),
						Nanos:   int32(birth_date.Nanosecond()),
					},
					AuthorityDate: &timestamp.Timestamp{
						Seconds: int64(authority_date.Second()),
						Nanos:   int32(authority_date.Nanosecond()),
					},
				},
			},
		}
	}

	if flag_error {
		answer_data, err = proto.Marshal(error_answer)
	} else {
		answer_data, err = proto.Marshal(answer)
	}

	if err != nil {
		usersGRPC.accLog.Error(err)
		return err
	}

	err = kProducer.ProduceRecord(answer_topic, sarama.ByteEncoder(answer_data))
	if err != nil {
		usersGRPC.accLog.Error(err)
		return err
	}
	usersGRPC.accLog.Info("Success to send answer!")

	return nil
}

func (usersGRPC UsersGrpcHandlers) UpdateUsersPassport(ctx context.Context, saga_uuid string, event_uuid string, operation_details *api.OperationDetails, kProducer *kafka.ProducerProvider) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "usersGRPC.UpdateUsersPassport")
	defer span.Finish()

	var answer_data []byte
	flag_error := false
	var err error
	answer_topic := TopicResult

	answer := &api.EventSuccess{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: UpdateUsersPassport,
	}

	error_answer := &api.EventError{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: UpdateUsersPassport,
	}

	user_uuid, err := uuid.Parse(operation_details.GetUserUuid())
	if err != nil {
		error_answer.Info = err.Error()
		error_answer.Status = GetErrorCode(ErrorInvalidInputData)

		flag_error = true
	}

	passport := operation_details.GetPassport()
	if passport == nil {
		return ErrorInvalidInputData
	}

	fcs := passport.GetFcs()

	local_passport := &models.Passport{
		Authority:           passport.GetAuthority(),
		Authority_date:      passport.GetAuthorityDate().AsTime().Format("2006-01-02 15:04:05 -07:00:00"),
		Birth_location:      passport.GetBirthLocation(),
		Name:                fcs.GetName(),
		Surname:             fcs.GetSurname(),
		Patronimic:          fcs.GetPatronymic(),
		Birth_date:          passport.GetBirthDate().AsTime().Format("2006-01-02 15:04:05 -07:00:00"),
		Passport_number:     passport.GetNumber(),
		Passport_series:     passport.GetSeries(),
		Registration_adress: passport.GetRegistrationAdress(),
		Pick_up_point:       passport.GetPickUpPoint(),
	}

	if err = usersGRPC.usersUC.UpdateUsersPassport(ctxWithTrace, user_uuid, local_passport); err != nil {
		usersGRPC.accLog.Error(err)
		flag_error = true
		answer_topic = TopicError

		error_answer.Info = err.Error()
		error_answer.Status = GetErrorCode(err)
	} else {
		answer.Result = &api.EventSuccess_Info{"Success"}
	}

	if flag_error {
		answer_data, err = proto.Marshal(error_answer)
	} else {
		answer_data, err = proto.Marshal(answer)
	}

	if err != nil {
		usersGRPC.accLog.Error(err)
		return err
	}

	err = kProducer.ProduceRecord(answer_topic, sarama.ByteEncoder(answer_data))
	if err != nil {
		usersGRPC.accLog.Error(err)
		return err
	}
	usersGRPC.accLog.Info("Success to send answer!")

	return nil
}

func (usersGRPC UsersGrpcHandlers) AddUserAccount(ctx context.Context, saga_uuid string, event_uuid string, operation_details *api.OperationDetails, kProducer *kafka.ProducerProvider) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "usersGRPC.AddUserAccount")
	defer span.Finish()

	var answer_data []byte
	flag_error := false
	var err error
	answer_topic := TopicResult

	answer := &api.EventSuccess{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: AddUserAccounts,
	}

	error_answer := &api.EventError{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: AddUserAccounts,
	}

	user_uuid, err := uuid.Parse(operation_details.GetUserUuid())
	if err != nil {
		error_answer.Info = err.Error()
		error_answer.Status = GetErrorCode(ErrorInvalidInputData)

		flag_error = true
	}

	account_uuid, err := uuid.Parse(operation_details.GetSomeData())
	if err != nil {
		error_answer.Info = err.Error()
		error_answer.Status = GetErrorCode(ErrorInvalidInputData)

		flag_error = true
	}

	if err = usersGRPC.usersUC.AddUserAccount(ctxWithTrace, user_uuid, account_uuid); err != nil {
		usersGRPC.accLog.Error(err)
		flag_error = true
		answer_topic = TopicError

		error_answer.Info = err.Error()
		error_answer.Status = GetErrorCode(err)
	} else {
		answer.Result = &api.EventSuccess_Info{"Success"}
	}

	if flag_error {
		answer_data, err = proto.Marshal(error_answer)
	} else {
		answer_data, err = proto.Marshal(answer)
	}

	if err != nil {
		usersGRPC.accLog.Error(err)
		return err
	}

	err = kProducer.ProduceRecord(answer_topic, sarama.ByteEncoder(answer_data))
	if err != nil {
		usersGRPC.accLog.Error(err)
		return err
	}
	usersGRPC.accLog.Info("Success to send answer!")

	return nil
}

func (usersGRPC UsersGrpcHandlers) RemoveUserAccount(ctx context.Context, saga_uuid string, event_uuid string, operation_details *api.OperationDetails, kProducer *kafka.ProducerProvider) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "usersGRPC.RemoveUserAccount")
	defer span.Finish()

	var answer_data []byte
	flag_error := false
	var err error
	answer_topic := TopicResult

	answer := &api.EventSuccess{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: RemoveUserAccount,
	}

	error_answer := &api.EventError{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: RemoveUserAccount,
	}

	user_uuid, err := uuid.Parse(operation_details.GetUserUuid())
	if err != nil {
		error_answer.Info = err.Error()
		error_answer.Status = GetErrorCode(ErrorInvalidInputData)

		flag_error = true
	}

	account_uuid, err := uuid.Parse(operation_details.GetSomeData())
	if err != nil {
		error_answer.Info = err.Error()
		error_answer.Status = GetErrorCode(ErrorInvalidInputData)

		flag_error = true
	}

	if err = usersGRPC.usersUC.RemoveUserAccount(ctxWithTrace, user_uuid, account_uuid); err != nil {
		usersGRPC.accLog.Error(err)
		flag_error = true
		answer_topic = TopicError

		error_answer.Info = err.Error()
		error_answer.Status = GetErrorCode(err)
	} else {
		answer.Result = &api.EventSuccess_Info{"Success"}
	}

	if flag_error {
		answer_data, err = proto.Marshal(error_answer)
	} else {
		answer_data, err = proto.Marshal(answer)
	}

	if err != nil {
		usersGRPC.accLog.Error(err)
		return err
	}

	err = kProducer.ProduceRecord(answer_topic, sarama.ByteEncoder(answer_data))
	if err != nil {
		usersGRPC.accLog.Error(err)
		return err
	}
	usersGRPC.accLog.Info("Success to send answer!")

	return nil
}

func (usersGRPC UsersGrpcHandlers) GetUsersAccounts(ctx context.Context, saga_uuid string, event_uuid string, operation_details *api.OperationDetails, kProducer *kafka.ProducerProvider) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "usersGRPC.GetUsersAccounts")
	defer span.Finish()

	var answer_data []byte
	flag_error := false
	var err error
	answer_topic := TopicResult

	answer := &api.EventSuccess{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: GetUsersAccounts,
	}

	error_answer := &api.EventError{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: GetUsersAccounts,
	}

	user_uuid, err := uuid.Parse(operation_details.GetUserUuid())
	if err != nil {
		error_answer.Info = err.Error()
		error_answer.Status = GetErrorCode(ErrorInvalidInputData)

		flag_error = true
	}

	var result *models.ListOfAccounts

	if result, err = usersGRPC.usersUC.GetUserAccounts(ctxWithTrace, user_uuid); err != nil {
		usersGRPC.accLog.Error(err)
		flag_error = true
		answer_topic = TopicError

		error_answer.Info = err.Error()
		error_answer.Status = GetErrorCode(err)
	} else {
		var accounts []string

		for i := 0; i < len(result.Data); i++ {
			accounts = append(accounts, result.Data[i].String())
		}

		answer.Result = &api.EventSuccess_Accounts{
			Accounts: &api.ListOfAccounts{
				Accounts: accounts,
			},
		}

	}

	if flag_error {
		answer_data, err = proto.Marshal(error_answer)
	} else {
		answer_data, err = proto.Marshal(answer)
	}

	if err != nil {
		usersGRPC.accLog.Error(err)
		return err
	}

	err = kProducer.ProduceRecord(answer_topic, sarama.ByteEncoder(answer_data))
	if err != nil {
		usersGRPC.accLog.Error(err)
		return err
	}
	usersGRPC.accLog.Info("Success to send answer!")

	return nil

	return nil
}

func NewUsersGRPCHandlers(cfg *config.Config, kProducer *kafka.ProducerProvider, usersUC users.UseCase, accLog logger.Logger) users.GRPCHandlers {
	return &UsersGrpcHandlers{cfg: cfg, kProducer: kProducer, usersUC: usersUC, accLog: accLog}
}
