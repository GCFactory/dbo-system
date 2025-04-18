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
	"mime/multipart"
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

func (h UsersHandlers) safeReadFormDataRequest(c echo.Context, v interface{}) error {
	val := reflect.ValueOf(v).Elem() // Получаем значение структуры

	// Получаем все form data
	form, err := c.MultipartForm()
	if err != nil {
		// Если это не multipart, пробуем получить как обычную форму
		if err := c.Request().ParseForm(); err != nil {
			return fmt.Errorf("failed to parse form data: %v", err)
		}
		form = &multipart.Form{
			Value: c.Request().Form,
		}
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)     // Получаем информацию о поле
		fieldName := field.Name          // Имя поля в структуре
		jsonTag := field.Tag.Get("json") // Получаем значение тега `json`
		if jsonTag == "" {
			continue
		}
		fieldType := field.Type // Тип поля

		// Получаем значение из form data по имени поля
		values, exists := form.Value[jsonTag]
		if !exists || len(values) == 0 {
			continue
		}
		formValue := values[0] // Берем первое значение (для массивов нужно обрабатывать иначе)

		// Устанавливаем значение поля в зависимости от его типа
		switch fieldType.Kind() {
		case reflect.String:
			val.Field(i).SetString(formValue) // Устанавливаем строку
		case reflect.Int, reflect.Int64:
			intValue, err := strconv.Atoi(formValue) // Преобразуем строку в int
			if err != nil {
				return fmt.Errorf("invalid value for field %s: %v", fieldName, err)
			}
			val.Field(i).SetInt(int64(intValue)) // Устанавливаем int
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(formValue) // Преобразуем строку в bool
			if err != nil {
				return fmt.Errorf("invalid value for field %s: %v", fieldName, err)
			}
			val.Field(i).SetBool(boolValue) // Устанавливаем bool
		case reflect.Slice:
			// Обработка массивов/слайсов
			if fieldType.Elem().Kind() == reflect.String {
				val.Field(i).Set(reflect.ValueOf(values))
			} else {
				return fmt.Errorf("unsupported slice type: %s", fieldType.Elem().Kind())
			}
		default:
			return fmt.Errorf("unsupported field type: %s", fieldType.Kind())
		}
	}

	// Валидация структуры
	validate := validator.New()
	if err := validate.Struct(v); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("validation error in field '%s': condition '%s' not met",
				err.Field(), err.Tag())
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
			UsingTotp: userData.User.UsingTotp,
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

func (h UsersHandlers) CheckUserPassw() echo.HandlerFunc {
	return func(c echo.Context) error {

		result := &models.DefaultHttpResponse{
			Status: http.StatusOK,
			Info:   "",
		}

		operationInfo := &models.CheckPasswordRequest{}
		err := h.safeReadBodyRequest(c, operationInfo)
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

		err = h.useCase.CheckUserPassword(context.Background(), user.User_uuid, operationInfo.Password)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusBadRequest
			result.Info = err.Error()
			return c.JSON(http.StatusBadRequest, result)
		}

		loginInfo := &models.CheckPasswordResponse{
			UserId:    user.User_uuid,
			TotpUsage: user.UsingTotp,
		}

		return c.JSON(http.StatusOK, loginInfo)
	}
}

func (h UsersHandlers) GetUserTotpInfo() echo.HandlerFunc {
	return func(c echo.Context) error {

		result := &models.DefaultHttpResponse{
			Status: http.StatusOK,
			Info:   "",
		}

		operationInfo := &models.GetUserTotpDataRequest{}
		err := h.safeReadBodyRequest(c, operationInfo)
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

		userInfo, err := h.useCase.GetUserData(context.Background(), userId)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusNotFound
			result.Info = err.Error()
			return c.JSON(http.StatusNotFound, result)
		}

		totpInfo := &models.GetUserTotpDataResponse{
			TotpId:    userInfo.User.TotpId,
			TotpUsage: userInfo.User.UsingTotp,
		}
		return c.JSON(http.StatusOK, totpInfo)
	}
}

func (h UsersHandlers) UpdateTotpInfo() echo.HandlerFunc {
	return func(c echo.Context) error {

		result := &models.DefaultHttpResponse{
			Status: http.StatusOK,
			Info:   "",
		}

		operationInfo := &models.UpdateTotpInfoRequest{}
		err := h.safeReadBodyRequest(c, operationInfo)
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

		totpId, err := uuid.Parse(operationInfo.TotpId)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusBadRequest
			result.Info = err.Error()
			return c.JSON(http.StatusBadRequest, result)
		}

		_, err = h.useCase.GetUserData(context.Background(), userId)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusNotFound
			result.Info = err.Error()
			return c.JSON(http.StatusNotFound, result)
		}

		err = h.useCase.UpdateTotpInfo(context.Background(), userId, totpId, operationInfo.TotpUsage)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			result.Status = http.StatusInternalServerError
			result.Info = err.Error()
			return c.JSON(http.StatusInternalServerError, result)
		}

		return c.JSON(http.StatusOK, nil)
	}
}

func NewUsersHandlers(cfg *config.Config, useCase users.UseCase, logger logger.Logger) users.HttpHandlers {
	return &UsersHandlers{cfg: cfg, logger: logger, useCase: useCase}
}
