package http

import (
	"context"
	"fmt"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/notification/internal/models"
	"github.com/GCFactory/dbo-system/service/notification/internal/notification"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
	"strconv"
)

type ApiGatewayHandlers struct {
	cfg     *config.Config
	useCase notification.UseCase
	logger  logger.Logger
}

func (h ApiGatewayHandlers) safeReadBodyRequest(c echo.Context, v interface{}) error {
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

func (h ApiGatewayHandlers) safeReadQueryParamsRequest(c echo.Context, v interface{}) error {
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

func (h ApiGatewayHandlers) UpdateUserSettings() echo.HandlerFunc {
	return func(c echo.Context) error {

		result := &models.DefaultHttpResponse{
			Status: http.StatusOK,
			Info:   "",
		}

		operationInfo := &models.UpdateUserSettings{}
		err := h.safeReadBodyRequest(c, operationInfo)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusBadRequest
			result.Info = err.Error()
			return c.JSON(http.StatusBadRequest, result)
		}

		userSettings := &models.UserNotificationInfo{
			UserUuid:   operationInfo.UserId,
			EmailUsage: operationInfo.EmailNotification,
			Email:      operationInfo.Email,
		}
		err = h.useCase.UpdateUserSettings(context.Background(), userSettings)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusBadRequest
			result.Info = err.Error()
			return c.JSON(http.StatusBadRequest, result)
		}

		result.Info = "Success"
		return c.JSON(http.StatusOK, result)
	}
}

func (h ApiGatewayHandlers) GetUserSettings() echo.HandlerFunc {
	return func(c echo.Context) error {

		result := &models.DefaultHttpResponse{
			Status: http.StatusOK,
			Info:   "",
		}

		operationInfo := &models.GetUserSettings{}
		err := h.safeReadQueryParamsRequest(c, operationInfo)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusBadRequest
			result.Info = err.Error()
			return c.JSON(http.StatusBadRequest, result)
		}

		userId, err := uuid.Parse(operationInfo.UserId)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusBadRequest
			result.Info = err.Error()
			return c.JSON(http.StatusBadRequest, result)
		}

		userSettings, err := h.useCase.GetUserSettings(context.Background(), userId)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusNotFound
			result.Info = err.Error()
			return c.JSON(http.StatusNotFound, result)
		}

		return c.JSON(http.StatusOK, userSettings)
	}
}

func NewNotificationHandlers(cfg *config.Config, useCase notification.UseCase, logger logger.Logger) notification.HttpHandlers {
	return &ApiGatewayHandlers{cfg: cfg, logger: logger, useCase: useCase}
}
