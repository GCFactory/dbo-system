package tests

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/users/internal/models"
	"github.com/GCFactory/dbo-system/service/users/internal/users/mock"
	"github.com/GCFactory/dbo-system/service/users/internal/users/usecase"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	usecaseTestCfg = &config.Config{
		Env: "Development",
		Logger: config.Logger{
			Development: true,
			Level:       "Debug",
		},
	}
)

func TestUsersUC_AddUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(usecaseTestCfg)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	usersUC := usecase.NewUsersUseCase(usecaseTestCfg, mockRepo, apiLogger)

	passport := &models.Passport{
		Passport_uuid:       uuid.New(),
		Registration_adress: "...",
		Authority_date:      time.DateTime,
		Authority:           "5",
		Pick_up_point:       "000-111",
		Birth_location:      "4",
		Birth_date:          time.DateTime,
		Patronimic:          "3",
		Name:                "1",
		Passport_number:     "123456",
		Passport_series:     "1234",
		Surname:             "2",
	}

	user := &models.User{
		User_uuid:     uuid.New(),
		Passport_uuid: passport.Passport_uuid,
		User_inn:      "01234567890123456789",
		User_accounts: "{\"accounts\": [] }",
	}

	user_full := &models.User_full_data{
		User:     user,
		Passport: passport,
	}

	ctx := context.Background()

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userUsecase.AddUser")
	defer span.Finish()

	var err error

	t.Run("Success", func(t *testing.T) {

		mockRepo.EXPECT().GetUserData(gomock.Eq(ctxWithTrace), gomock.Eq(user.User_uuid)).Return(nil, err)
		mockRepo.EXPECT().AddUser(gomock.Eq(ctxWithTrace), gomock.Eq(user_full)).Return(nil)

		err = usersUC.AddUser(ctx, user_full)
		require.Nil(t, err)

	})

	t.Run("Error user exist", func(t *testing.T) {

		mockRepo.EXPECT().GetUserData(gomock.Eq(ctxWithTrace), gomock.Eq(user.User_uuid)).Return(user_full, nil)

		err = usersUC.AddUser(ctx, user_full)
		require.Equal(t, err, usecase.ErrorUserExists)

	})

	t.Run("Error create user", func(t *testing.T) {

		mockRepo.EXPECT().GetUserData(gomock.Eq(ctxWithTrace), gomock.Eq(user.User_uuid)).Return(nil, err)
		mockRepo.EXPECT().AddUser(gomock.Eq(ctxWithTrace), gomock.Eq(user_full)).Return(err)

		tmp_err := err

		err = usersUC.AddUser(ctx, user_full)
		require.Equal(t, err, tmp_err)

	})
}

func TestUsersUC_GetUserData(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(usecaseTestCfg)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	usersUC := usecase.NewUsersUseCase(usecaseTestCfg, mockRepo, apiLogger)

	passport := &models.Passport{
		Passport_uuid:       uuid.New(),
		Registration_adress: "...",
		Authority_date:      time.DateTime,
		Authority:           "5",
		Pick_up_point:       "000-111",
		Birth_location:      "4",
		Birth_date:          time.DateTime,
		Patronimic:          "3",
		Name:                "1",
		Passport_number:     "123456",
		Passport_series:     "1234",
		Surname:             "2",
	}

	user := &models.User{
		User_uuid:     uuid.New(),
		Passport_uuid: passport.Passport_uuid,
		User_inn:      "01234567890123456789",
		User_accounts: "{\"accounts\": [] }",
	}

	user_full := &models.User_full_data{
		User:     user,
		Passport: passport,
	}

	ctx := context.Background()

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userUsecase.GetUserData")
	defer span.Finish()

	t.Run("Success", func(t *testing.T) {

		mockRepo.EXPECT().GetUserData(gomock.Eq(ctxWithTrace), gomock.Eq(user.User_uuid)).Return(user_full, nil)

		result, err := usersUC.GetUserData(ctx, user.User_uuid)
		require.Nil(t, err)
		require.Equal(t, result, user_full)

	})

	t.Run("Error user not found", func(t *testing.T) {

		mockRepo.EXPECT().GetUserData(gomock.Eq(ctxWithTrace), gomock.Eq(user.User_uuid)).Return(nil, errors.New("Some error"))

		result, err := usersUC.GetUserData(ctx, user.User_uuid)
		require.Equal(t, err, usecase.ErrorUserNotFound)
		require.Nil(t, result)

	})
}

func TestUsersUC_UpdateUsersPassport(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(usecaseTestCfg)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	usersUC := usecase.NewUsersUseCase(usecaseTestCfg, mockRepo, apiLogger)

	passport := &models.Passport{
		Passport_uuid:       uuid.New(),
		Registration_adress: "...",
		Authority_date:      time.DateTime,
		Authority:           "5",
		Pick_up_point:       "000-111",
		Birth_location:      "4",
		Birth_date:          time.DateTime,
		Patronimic:          "3",
		Name:                "1",
		Passport_number:     "123456",
		Passport_series:     "1234",
		Surname:             "2",
	}

	new_passport := &models.Passport{
		Passport_uuid:       passport.Passport_uuid,
		Registration_adress: "...",
		Authority_date:      time.DateTime,
		Authority:           "5",
		Pick_up_point:       "000-111",
		Birth_location:      "4",
		Birth_date:          time.DateTime,
		Patronimic:          "3",
		Name:                "1",
		Passport_number:     "123456",
		Passport_series:     "1234",
		Surname:             "2",
	}

	user := &models.User{
		User_uuid:     uuid.New(),
		Passport_uuid: passport.Passport_uuid,
		User_inn:      "01234567890123456789",
		User_accounts: "{\"accounts\": [] }",
	}

	user_full := &models.User_full_data{
		User:     user,
		Passport: passport,
	}

	ctx := context.Background()

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userUsecase.UpdateUsersPassport")
	defer span.Finish()

	t.Run("Success", func(t *testing.T) {

		mockRepo.EXPECT().GetUserData(gomock.Eq(ctxWithTrace), gomock.Eq(user.User_uuid)).Return(user_full, nil)
		mockRepo.EXPECT().UpdatePassport(gomock.Eq(ctxWithTrace), gomock.Eq(new_passport)).Return(nil)

		err := usersUC.UpdateUsersPassport(ctx, user.User_uuid, new_passport)
		require.Nil(t, err)

	})

	t.Run("Error user not found", func(t *testing.T) {

		mockRepo.EXPECT().GetUserData(gomock.Eq(ctxWithTrace), gomock.Eq(user.User_uuid)).Return(nil, errors.New("Some error"))

		err := usersUC.UpdateUsersPassport(ctx, user.User_uuid, new_passport)
		require.Equal(t, err, usecase.ErrorUserNotFound)

	})

	t.Run("Error update passport", func(t *testing.T) {

		local_err := errors.New("Some error")

		mockRepo.EXPECT().GetUserData(gomock.Eq(ctxWithTrace), gomock.Eq(user.User_uuid)).Return(user_full, nil)
		mockRepo.EXPECT().UpdatePassport(gomock.Eq(ctxWithTrace), gomock.Eq(new_passport)).Return(local_err)

		err := usersUC.UpdateUsersPassport(ctx, user.User_uuid, new_passport)
		require.Equal(t, err, local_err)

	})
}

func TestUsersUC_AddUserAccount(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(usecaseTestCfg)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	usersUC := usecase.NewUsersUseCase(usecaseTestCfg, mockRepo, apiLogger)

	passport := &models.Passport{
		Passport_uuid:       uuid.New(),
		Registration_adress: "...",
		Authority_date:      time.DateTime,
		Authority:           "5",
		Pick_up_point:       "000-111",
		Birth_location:      "4",
		Birth_date:          time.DateTime,
		Patronimic:          "3",
		Name:                "1",
		Passport_number:     "123456",
		Passport_series:     "1234",
		Surname:             "2",
	}

	user := &models.User{
		User_uuid:     uuid.New(),
		Passport_uuid: passport.Passport_uuid,
		User_inn:      "01234567890123456789",
		User_accounts: "{\"accounts\": [" +
			"\"" + uuid.New().String() + "\"," +
			"\"" + uuid.New().String() + "\"" +
			"]" +
			"}",
	}

	new_account := uuid.New()

	user_full := &models.User_full_data{
		User:     user,
		Passport: passport,
	}

	ctx := context.Background()

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userUsecase.AddUserAccount")
	defer span.Finish()

	t.Run("Success", func(t *testing.T) {

		var tmp models.ListOfAccounts

		err := json.Unmarshal([]byte(user.User_accounts), &tmp)
		require.Nil(t, err)

		accounts := tmp.Data

		accounts = append(accounts, new_account)

		tmp.Data = accounts

		bytes, err := json.Marshal(tmp)
		require.Nil(t, err)

		mockRepo.EXPECT().GetUserData(gomock.Eq(ctxWithTrace), gomock.Eq(user.User_uuid)).Return(user_full, nil)
		mockRepo.EXPECT().UpdateUserAccount(gomock.Eq(ctxWithTrace), gomock.Eq(user.User_uuid), gomock.Eq(string(bytes))).Return(nil)

		err = usersUC.AddUserAccount(ctx, user.User_uuid, new_account)
		require.Nil(t, err)

	})

	t.Run("Error user not found", func(t *testing.T) {

		mockRepo.EXPECT().GetUserData(gomock.Eq(ctxWithTrace), gomock.Eq(user.User_uuid)).Return(nil, errors.New("Some error"))

		err := usersUC.AddUserAccount(ctx, user.User_uuid, new_account)
		require.Equal(t, err, usecase.ErrorUserNotFound)

	})

	t.Run("Error account already exists", func(t *testing.T) {

		save := user.User_accounts

		var tmp models.ListOfAccounts

		err := json.Unmarshal([]byte(user.User_accounts), &tmp)
		require.Nil(t, err)

		accounts := tmp.Data

		accounts = append(accounts, new_account)

		tmp.Data = accounts

		bytes, err := json.Marshal(tmp)
		require.Nil(t, err)

		user.User_accounts = string(bytes)

		mockRepo.EXPECT().GetUserData(gomock.Eq(ctxWithTrace), gomock.Eq(user.User_uuid)).Return(user_full, nil)

		err = usersUC.AddUserAccount(ctx, user.User_uuid, new_account)
		require.Equal(t, err, usecase.ErrorAccountAlreadyExists)

		user.User_accounts = save

	})

	t.Run("Error update accounts", func(t *testing.T) {

		var tmp models.ListOfAccounts

		err := json.Unmarshal([]byte(user.User_accounts), &tmp)
		require.Nil(t, err)

		accounts := tmp.Data

		accounts = append(accounts, new_account)

		tmp.Data = accounts

		bytes, err := json.Marshal(tmp)
		require.Nil(t, err)

		local_error := errors.New("Some error")

		mockRepo.EXPECT().GetUserData(gomock.Eq(ctxWithTrace), gomock.Eq(user.User_uuid)).Return(user_full, nil)
		mockRepo.EXPECT().UpdateUserAccount(gomock.Eq(ctxWithTrace), gomock.Eq(user.User_uuid), gomock.Eq(string(bytes))).Return(local_error)

		err = usersUC.AddUserAccount(ctx, user.User_uuid, new_account)
		require.Equal(t, err, local_error)

	})
}

func TestUsersUC_GetUserAccounts(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(usecaseTestCfg)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	usersUC := usecase.NewUsersUseCase(usecaseTestCfg, mockRepo, apiLogger)

	passport := &models.Passport{
		Passport_uuid:       uuid.New(),
		Registration_adress: "...",
		Authority_date:      time.DateTime,
		Authority:           "5",
		Pick_up_point:       "000-111",
		Birth_location:      "4",
		Birth_date:          time.DateTime,
		Patronimic:          "3",
		Name:                "1",
		Passport_number:     "123456",
		Passport_series:     "1234",
		Surname:             "2",
	}

	user := &models.User{
		User_uuid:     uuid.New(),
		Passport_uuid: passport.Passport_uuid,
		User_inn:      "01234567890123456789",
		User_accounts: "{\"accounts\": [" +
			"\"" + uuid.New().String() + "\"," +
			"\"" + uuid.New().String() + "\"" +
			"]" +
			"}",
	}

	ctx := context.Background()

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userUsecase.GetUserAccounts")
	defer span.Finish()

	t.Run("Success", func(t *testing.T) {

		mockRepo.EXPECT().GetUsersAccounts(gomock.Eq(ctxWithTrace), gomock.Eq(user.User_uuid)).Return(user.User_accounts, nil)

		result, err := usersUC.GetUserAccounts(ctx, user.User_uuid)
		require.Nil(t, err)
		require.Equal(t, user.User_accounts, result)

	})

	t.Run("Error user not found", func(t *testing.T) {

		mockRepo.EXPECT().GetUsersAccounts(gomock.Eq(ctxWithTrace), gomock.Eq(user.User_uuid)).Return("", errors.New("Some error"))

		result, err := usersUC.GetUserAccounts(ctx, user.User_uuid)
		require.Equal(t, err, usecase.ErrorUserNotFound)
		require.Equal(t, result, "")

	})
}
