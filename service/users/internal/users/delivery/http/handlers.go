package http

import (
	"context"
	"fmt"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/users/internal/models"
	"github.com/GCFactory/dbo-system/service/users/internal/users"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
	"strconv"
)

type UsersHandlers struct {
	cfg     *config.Config
	useCase users.UseCase
	logger  logger.Logger
}

func (h UsersHandlers) safeReadBodyRequest(c echo.Context, v interface{}) error {
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

func (h UsersHandlers) safeReadQueryParamsRequest(c echo.Context, v interface{}) error {
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

func (h UsersHandlers) GetUserData() echo.HandlerFunc {
	return func(c echo.Context) error {

		result := &models.DefaultHttpResponse{
			Status: http.StatusOK,
			Info:   "",
		}

		operationInfo := &models.GetUserData{}
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

		userData, err := h.useCase.GetUserData(context.Background(), userId)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusBadRequest
			result.Info = err.Error()
			return c.JSON(http.StatusBadRequest, result)
		}

		userDataPrepare := &models.GetUserDataResponse{
			UserId:    userId,
			UserInn:   userData.User.User_inn,
			UserLogin: userData.User.User_login,
			Accounts:  userData.User.User_accounts.Data,
			PassportData: &models.GetUserDataResponse_Passport{
				Number:             userData.Passport.Passport_number,
				Series:             userData.Passport.Passport_series,
				RegistrationAdress: userData.Passport.Registration_adress,
				PickUpPoint:        userData.Passport.Pick_up_point,
				BirthDate:          userData.Passport.Birth_date,
				BirthLocation:      userData.Passport.Birth_location,
				AuthorityDate:      userData.Passport.Authority_date,
				Authority:          userData.Passport.Authority,
				FCS: &models.GetUserDataResponse_Passport_FCS{
					Name:       userData.Passport.Name,
					Surname:    userData.Passport.Surname,
					Patronymic: userData.Passport.Patronimic,
				},
			},
		}

		return c.JSON(http.StatusOK, userDataPrepare)
	}
}

func (h UsersHandlers) GetUserAccounts() echo.HandlerFunc {
	return func(c echo.Context) error {

		result := &models.DefaultHttpResponse{
			Status: http.StatusOK,
			Info:   "",
		}

		operationInfo := &models.GetUserAccount{}
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

		userAccounts, err := h.useCase.GetUserAccounts(context.Background(), userId)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusBadRequest
			result.Info = err.Error()
			return c.JSON(http.StatusBadRequest, result)
		}

		accounts := &models.GetUserAccountResponse{
			Accounts: userAccounts.Data,
		}
		return c.JSON(http.StatusOK, accounts)
	}
}

func (h UsersHandlers) GetUserDataByLogin() echo.HandlerFunc {
	return func(c echo.Context) error {

		result := &models.DefaultHttpResponse{
			Status: http.StatusOK,
			Info:   "",
		}

		operationInfo := &models.GetUserDataByLogin{}
		err := h.safeReadQueryParamsRequest(c, operationInfo)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusBadRequest
			result.Info = err.Error()
			return c.JSON(http.StatusBadRequest, result)
		}

		user, err := h.useCase.GetUserDataByLogin(context.Background(), operationInfo.Login)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusBadRequest
			result.Info = err.Error()
			return c.JSON(http.StatusBadRequest, result)
		}

		userInfo := &models.GetUserDataByLoginResponse{
			UserId: user.User_uuid,
		}
		return c.JSON(http.StatusOK, userInfo)
	}

}

func NewUsersHandlers(cfg *config.Config, useCase users.UseCase, logger logger.Logger) users.HttpHandlers {
	return &UsersHandlers{cfg: cfg, logger: logger, useCase: useCase}
}
