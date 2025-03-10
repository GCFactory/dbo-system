package http

import (
	"fmt"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type RegistrationHandlers struct {
	cfg     *config.Config
	useCase api_gateway.UseCase
	logger  logger.Logger
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

func (h RegistrationHandlers) SignInPage() echo.HandlerFunc {
	return func(c echo.Context) error {
		//span, ctxWithTrace := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "RegistrationHandlers.GetOperationStatus")
		//defer span.Finish()

		//operation_info := &models.OperationID{}
		//if err := h.safeReadRequest(c, operation_info); err != nil {
		//	utils.LogResponseError(c, h.logger, err)
		//	return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
		//}
		//
		//operation_id := uuid.MustParse(operation_info.Operation_ID)
		//
		//operation_status, err := h.registrationUC.GetOperationStatus(ctxWithTrace, operation_id)
		//if err != nil {
		//	utils.LogResponseError(c, h.logger, err)
		//	return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
		//}
		//
		//response := make(map[string]interface{})
		//response["info"] = operation_status
		//
		//operation_data, err := h.registrationUC.GetOperationResultData(ctxWithTrace, operation_id)
		//if err != nil {
		//	return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
		//}
		//response["additional_info"] = operation_data

		//return c.JSON(http.StatusOK, response)
		return nil
	}
}

func NewApiGatewayHandlers(cfg *config.Config, logger logger.Logger, usecase api_gateway.UseCase) api_gateway.Handlers {
	return &RegistrationHandlers{cfg: cfg, logger: logger, useCase: usecase}
}
