package http

import (
	"fmt"
	"github.com/GCFactory/dbo-system/platform/pkg/httpErrors"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/registration/config"
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"time"
)

type RegistrationHandlers struct {
	cfg              *config.Config
	registrationGRPC registration.RegistrationGRPCHandlers
	registrationUC   registration.UseCase
	logger           logger.Logger
}

func (h RegistrationHandlers) safeReadRequest(c echo.Context, v interface{}) error {
	var err error = nil
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("Validation error: %v", r)
				h.logger.Error(err.Error())
			}
		}()
		utils.ReadRequest(c, v)
	}()

	validate := validator.New()
	if validation_err := validate.Struct(v); validation_err != nil {
		if validationErrors, ok := validation_err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				err = fmt.Errorf("Ошибка в поле '%s': условие '%s' не выполнено\n", fieldError.Field(), fieldError.Tag())
				return err
			}
		} else {
			err = validation_err
			return err
		}

	}

	return err
}

func (h RegistrationHandlers) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctxWithTrace := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "RegistrationHandlers.CreateUser")
		defer span.Finish()

		user_info := &models.RegistrationUserInfo{}
		if err := h.safeReadRequest(c, user_info); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		}

		data := make(map[string]interface{})

		data["user_inn"] = user_info.User_INN
		data["passport_number"] = user_info.Passport.Number
		data["passport_series"] = user_info.Passport.Series
		data["name"] = user_info.Passport.Name
		data["surname"] = user_info.Passport.Surname
		data["patronimic"] = user_info.Passport.Patronymic
		data["birth_date"] = user_info.Passport.Birth_date
		data["birth_location"] = user_info.Passport.Birth_location
		data["pick_up_point"] = user_info.Passport.Pick_up_point
		data["authority"] = user_info.Passport.Authority
		data["authority_date"] = user_info.Passport.Authority_date
		data["registration_adress"] = user_info.Passport.Registration_address
		data["login"] = user_info.User.Login
		data["password"] = user_info.User.Password
		data["email"] = user_info.UserEmail
		data["email_notification"] = true

		operation_uuid, err := h.registrationGRPC.StartOperation(ctxWithTrace, usecase.OperationCreateUser, data)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
		}

		response := make(map[string]interface{})
		response["info"] = operation_uuid.String()

		return c.JSON(http.StatusAccepted, response)
	}
}

func (h RegistrationHandlers) GetOperationStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctxWithTrace := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "RegistrationHandlers.GetOperationStatus")
		defer span.Finish()

		operation_info := &models.OperationID{}
		if err := h.safeReadRequest(c, operation_info); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		}

		operation_id := uuid.MustParse(operation_info.Operation_ID)

		operation_status, err := h.registrationUC.GetOperationStatus(ctxWithTrace, operation_id)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
		}

		response := make(map[string]interface{})
		response["info"] = operation_status

		operation_data, err := h.registrationUC.GetOperationResultData(ctxWithTrace, operation_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
		}
		response["additional_info"] = operation_data

		return c.JSON(http.StatusOK, response)
	}
}

func (h RegistrationHandlers) OpenAccount() echo.HandlerFunc {
	return func(c echo.Context) error {

		span, ctxWithTrace := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "RegistrationHandlers.OpenAccount")
		defer span.Finish()

		operation_info := &models.OpenAccountInfo{}
		if err := h.safeReadRequest(c, operation_info); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		}

		data := make(map[string]interface{})

		data["user_id"] = operation_info.User_ID
		data["acc_name"] = operation_info.Acc_name
		data["culc_number"] = operation_info.Culc_number
		data["corr_number"] = operation_info.Corr_number
		data["bic"] = operation_info.BIC
		data["cio"] = operation_info.CIO
		data["reserve_reason"] = operation_info.Reserve_reason

		operation_uuid, err := h.registrationGRPC.StartOperation(ctxWithTrace, usecase.OperationAddAccount, data)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
		}

		response := make(map[string]interface{})
		response["info"] = operation_uuid.String()

		return c.JSON(http.StatusAccepted, response)
	}
}

func (h RegistrationHandlers) AddAccountCache() echo.HandlerFunc {
	return func(c echo.Context) error {

		span, ctxWithTrace := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "RegistrationHandlers.AddAccountCache")
		defer span.Finish()

		operation_info := &models.AddAccountCache{}
		if err := h.safeReadRequest(c, operation_info); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		}

		data := make(map[string]interface{})

		data["user_id"] = operation_info.User_ID
		data["acc_id"] = operation_info.Account_ID
		data["cache_diff"] = operation_info.Cache_diff

		operation_uuid, err := h.registrationGRPC.StartOperation(ctxWithTrace, usecase.OperationAddAccountCache, data)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
		}

		response := make(map[string]interface{})
		response["info"] = operation_uuid.String()

		return c.JSON(http.StatusAccepted, response)
	}
}

func (h RegistrationHandlers) WidthAccountCache() echo.HandlerFunc {
	return func(c echo.Context) error {

		span, ctxWithTrace := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "RegistrationHandlers.WidthAccountCache")
		defer span.Finish()

		operation_info := &models.WidthAccountCache{}
		if err := h.safeReadRequest(c, operation_info); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		}

		data := make(map[string]interface{})

		data["user_id"] = operation_info.User_ID
		data["acc_id"] = operation_info.Account_ID
		data["cache_diff"] = operation_info.Cache_diff

		operation_uuid, err := h.registrationGRPC.StartOperation(ctxWithTrace, usecase.OperationWidthAccountCache, data)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
		}

		response := make(map[string]interface{})
		response["info"] = operation_uuid.String()

		return c.JSON(http.StatusAccepted, response)
	}

}

func (h RegistrationHandlers) CloseAccount() echo.HandlerFunc {
	return func(c echo.Context) error {

		span, ctxWithTrace := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "RegistrationHandlers.WidthAccountCache")
		defer span.Finish()

		operation_info := &models.CloseAccount{}
		if err := h.safeReadRequest(c, operation_info); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		}

		data := make(map[string]interface{})

		data["user_id"] = operation_info.User_ID
		data["acc_id"] = operation_info.Account_ID

		operation_uuid, err := h.registrationGRPC.StartOperation(ctxWithTrace, usecase.OperationCloseAccount, data)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
		}

		response := make(map[string]interface{})
		response["info"] = operation_uuid.String()

		return c.JSON(http.StatusAccepted, response)
	}

}

func (h RegistrationHandlers) GetUserData() echo.HandlerFunc {
	return func(c echo.Context) error {

		span, ctxWithTrace := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "RegistrationHandlers.WidthAccountCache")
		defer span.Finish()

		operation_info := &models.GetUserData{}
		if err := h.safeReadRequest(c, operation_info); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		}

		data := make(map[string]interface{})

		data["user_id"] = operation_info.User_ID

		operation_uuid, err := h.registrationGRPC.StartOperation(ctxWithTrace, usecase.OperationGetUserData, data)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
		}

		response := make(map[string]interface{})
		response["info"] = operation_uuid.String()

		return c.JSON(http.StatusAccepted, response)
	}

}

func (h RegistrationHandlers) GetAccountData() echo.HandlerFunc {
	return func(c echo.Context) error {

		span, ctxWithTrace := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "RegistrationHandlers.WidthAccountCache")
		defer span.Finish()

		operation_info := &models.GetAccountData{}
		if err := h.safeReadRequest(c, operation_info); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		}

		data := make(map[string]interface{})

		data["user_id"] = operation_info.User_ID
		data["acc_id"] = operation_info.Account_ID

		operation_uuid, err := h.registrationGRPC.StartOperation(ctxWithTrace, usecase.OperationGetAccountData, data)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
		}

		response := make(map[string]interface{})
		response["info"] = operation_uuid.String()

		return c.JSON(http.StatusAccepted, response)
	}

}

func (h RegistrationHandlers) UpdateUserPassword() echo.HandlerFunc {

	return func(c echo.Context) error {

		span, ctxWithTrace := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "RegistrationHandlers.UpdateUserPassword")
		defer span.Finish()

		operation_info := &models.UpdateUserPassword{}
		if err := h.safeReadRequest(c, operation_info); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		}

		data := make(map[string]interface{})

		data["user_id"] = operation_info.User_ID
		data["new_password"] = operation_info.New_password

		operation_uuid, err := h.registrationGRPC.StartOperation(ctxWithTrace, usecase.OperationGroupUpdateUserPassword, data)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
		}

		response := make(map[string]interface{})
		response["info"] = operation_uuid.String()

		return c.JSON(http.StatusAccepted, response)
	}

}

func (h RegistrationHandlers) GetUserDataByLogin() echo.HandlerFunc {

	return func(c echo.Context) error {

		span, ctxWithTrace := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "RegistrationHandlers.UpdateUserPassword")
		defer span.Finish()

		operation_info := &models.GetUserDataByLogin{}
		if err := h.safeReadRequest(c, operation_info); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		}

		data := make(map[string]interface{})

		data["user_login"] = operation_info.User_login

		operation_uuid, err := h.registrationGRPC.StartOperation(ctxWithTrace, usecase.OperationGetUserDataByLogin, data)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
		}

		response := make(map[string]interface{})
		response["info"] = operation_uuid.String()

		return c.JSON(http.StatusAccepted, response)
	}

}

func (h RegistrationHandlers) CheckUserPassword() echo.HandlerFunc {

	return func(c echo.Context) error {

		span, ctxWithTrace := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "RegistrationHandlers.UpdateUserPassword")
		defer span.Finish()

		operation_info := &models.CheckUserPassword{}
		if err := h.safeReadRequest(c, operation_info); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		}

		data := make(map[string]interface{})

		data["user_id"] = operation_info.User_ID
		data["password"] = operation_info.Password

		operation_uuid, err := h.registrationGRPC.StartOperation(ctxWithTrace, usecase.OperationCheckUserPassword, data)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
		}

		response := make(map[string]interface{})
		response["info"] = operation_uuid.String()

		return c.JSON(http.StatusAccepted, response)
	}

}

func (h RegistrationHandlers) GetOperationTree() echo.HandlerFunc {

	return func(c echo.Context) error {

		span, ctxWithTrace := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "RegistrationHandlers.GetOperationTree")
		defer span.Finish()

		operation_info := &models.OperationID{}
		if err := h.safeReadRequest(c, operation_info); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		}

		operation_id_str := operation_info.Operation_ID
		opreation_id, err := uuid.Parse(operation_id_str)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		}

		data, err := h.registrationUC.GetOperationTree(ctxWithTrace, opreation_id)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
		}

		response := make(map[string]interface{})
		response["info"] = data

		return c.JSON(http.StatusAccepted, response)
	}
}

func (h RegistrationHandlers) GetOperationListBetween() echo.HandlerFunc {

	return func(c echo.Context) error {

		span, ctxWithTrace := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "RegistrationHandlers.GetOperationTree")
		defer span.Finish()

		operation_info := &models.TimeBetween{}
		if err := h.safeReadRequest(c, operation_info); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		}

		time_begin, err := time.Parse("02-01-2006 15:04:05", operation_info.Time_begin)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		}

		time_end, err := time.Parse("02-01-2006 15:04:05", operation_info.Time_end)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		}

		data, err := h.registrationUC.GerOperationListBetween(ctxWithTrace, time_begin, time_end)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
		}

		response := make(map[string]interface{})
		result := make(map[string]interface{})
		result["operations"] = data

		response["info"] = result

		return c.JSON(http.StatusAccepted, response)
	}

}

func NewRegistrationHandlers(cfg *config.Config, logger logger.Logger, usecase registration.UseCase, grpc registration.RegistrationGRPCHandlers) registration.Handlers {
	return &RegistrationHandlers{cfg: cfg, logger: logger, registrationGRPC: grpc, registrationUC: usecase}
}
