package http

import (
	"context"
	"fmt"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/httpErrors"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/account/internal/account"
	"github.com/GCFactory/dbo-system/service/account/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
	"strconv"
)

type httpError struct {
	_ httpErrors.RestError
}

type accHandlers struct {
	cfg    *config.Config
	accUC  account.UseCase
	logger logger.Logger
}

func NewACCHandlers(cfg *config.Config, accUC account.UseCase, log logger.Logger) account.Handlers {
	return &accHandlers{cfg: cfg, accUC: accUC, logger: log}
}

func (h accHandlers) GetAccountData() echo.HandlerFunc {
	return func(c echo.Context) error {

		result := &models.DefaultHttpResponse{
			Status: http.StatusOK,
			Info:   "",
		}

		operationInfo := &models.GetAccountInfo{}
		err := h.safeReadQueryParamsRequest(c, operationInfo)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusBadRequest
			result.Info = err.Error()
			return c.JSON(http.StatusBadRequest, result)
		}

		accId, err := uuid.Parse(operationInfo.AccId)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusBadRequest
			result.Info = err.Error()
			return c.JSON(http.StatusBadRequest, result)
		}

		accInfo, err := h.accUC.GetAccInfo(context.Background(), accId)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusNotFound
			result.Info = err.Error()
			return c.JSON(http.StatusNotFound, result)
		}

		return c.JSON(http.StatusOK, accInfo)
	}
}

func (h accHandlers) safeReadBodyRequest(c echo.Context, v interface{}) error {
	var err error = nil
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("Validation error: %v", r)
				h.logger.Error(err.Error())
			}
		}()
		err = utils.ReadRequest(c, v)
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

func (h accHandlers) safeReadQueryParamsRequest(c echo.Context, v interface{}) error {
	val := reflect.ValueOf(v).Elem() // Получаем значение структуры

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)     // Получаем информацию о поле
		fieldName := field.Name          // Имя поля в структуре
		jsonTag := field.Tag.Get("json") // Получаем значение тега `json`
		if jsonTag == "" {
			continue
		}
		fieldType := field.Type // Тип поля

		// Получаем значение query-параметра по имени поля
		queryParam := c.QueryParam(jsonTag)

		// Если параметр не передан, пропускаем поле
		if queryParam == "" {
			continue
		}

		// Устанавливаем значение поля в зависимости от его типа
		switch fieldType.Kind() {
		case reflect.String:
			val.Field(i).SetString(queryParam) // Устанавливаем строку
		case reflect.Int, reflect.Int64:
			intValue, err := strconv.Atoi(queryParam) // Преобразуем строку в int
			if err != nil {
				return fmt.Errorf("invalid value for field %s: %v", fieldName, err)
			}
			val.Field(i).SetInt(int64(intValue)) // Устанавливаем int
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(queryParam) // Преобразуем строку в bool
			if err != nil {
				return fmt.Errorf("invalid value for field %s: %v", fieldName, err)
			}
			val.Field(i).SetBool(boolValue) // Устанавливаем bool
		default:
			return fmt.Errorf("unsupported field type: %s", fieldType.Kind())
		}
	}

	validate := validator.New()
	err := validate.Struct(v)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validation_err := fmt.Errorf("Ошибка в поле '%s': условие '%s' не выполнено\n", err.Field(), err.Tag())
			return validation_err
		}
	}

	return nil
}
