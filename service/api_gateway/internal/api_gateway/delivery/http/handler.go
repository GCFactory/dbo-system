package http

import (
	"context"
	"fmt"
	"github.com/GCFactory/dbo-system/platform/pkg/httpErrors"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/api_gateway/config"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway/usecase"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"
)

type ApiGatewayHandlers struct {
	cfg         *config.Config
	useCase     api_gateway.UseCase
	logger      logger.Logger
	folderGraph string
	folderQr    string
}

const (
	CookieTokenNameMain      string = "token"
	CookieTokenNameFirstAuth string = "token_fa"
)

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

func (h ApiGatewayHandlers) safeReadFormDataRequest(c echo.Context, v interface{}) error {
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

func (h ApiGatewayHandlers) SignInPage() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, _, err := h.CheckToken(c, CookieTokenNameMain)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		if is_ok {
			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/main_page")
		} else {
			signIngPage, err := h.useCase.CreateSignInPage()
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}
			return c.HTML(http.StatusOK, signIngPage)
		}
	}
}

func (h ApiGatewayHandlers) SignUpPage() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		if !is_ok {
			page, err := h.useCase.CreateSignUpPage()
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, page)
		} else {
			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}
			home_page, err := h.useCase.CreateUserPage(user_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}
			return c.HTML(http.StatusOK, home_page)
		}
	}
}

func (h ApiGatewayHandlers) AdminPage() echo.HandlerFunc {
	return func(c echo.Context) error {

		operation_info := &models.AdminPageRequestBody{}
		if err := h.safeReadQueryParamsRequest(c, operation_info); err != nil {
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}
		}

		if operation_info.Start == "" || operation_info.End == "" {
			operation_info.Start = time.Date(
				time.Now().Year(),
				time.Now().Month(),
				time.Now().Day(),
				0,
				0,
				0,
				0,
				time.Now().Location(),
			).Format("02-01-2006 15:04:05")
			operation_info.End = time.Date(
				time.Now().Year(),
				time.Now().Month(),
				time.Now().Day(),
				23,
				59,
				59,
				0,
				time.Now().Location(),
			).Format("02-01-2006 15:04:05")
		}
		page, err := h.useCase.CreateAdminPage(operation_info.Start, operation_info.End)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		return c.HTML(http.StatusOK, page)
	}
}

func (h ApiGatewayHandlers) SignIn() echo.HandlerFunc {
	return func(c echo.Context) error {

		operation_info := &models.SignInInfo{}
		operation_result := &models.PostRequestStatus{
			Success: false,
		}
		err := h.safeReadFormDataRequest(c, operation_info)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			errPage, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				operation_result.Error = err.Error()
				return c.JSON(http.StatusInternalServerError, operation_result)
			}
			return c.HTML(http.StatusBadRequest, errPage)
		}

		token, err := h.useCase.SignIn(operation_info)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			errPage, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				operation_result.Error = err.Error()
				return c.JSON(http.StatusInternalServerError, operation_result)
			}
			return c.HTML(http.StatusBadRequest, errPage)
		}

		totpInfo, err := h.useCase.GetUserTotpInfo(token.Data)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			errPage, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				operation_result.Error = err.Error()
				return c.JSON(http.StatusInternalServerError, operation_result)
			}
			return c.HTML(http.StatusInternalServerError, errPage)
		}
		if totpInfo.TotpUsage {

			tokenFirstAuth := &models.TokenFirstAuth{
				UserId:    token.Data,
				TokenName: operation_info.Login + "_totp",
				Live_time: usecase.TokenFirstAuthLiveTime,
			}

			err = h.useCase.AddTokenFirstAuth(context.Background(), tokenFirstAuth)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			cookie := new(http.Cookie)
			cookie.Name = CookieTokenNameFirstAuth
			cookie.Value = tokenFirstAuth.TokenName
			cookie.Expires = time.Now().Add(usecase.TokenFirstAuthLiveTime)
			cookie.Path = "/"
			cookie.HttpOnly = true
			cookie.Secure = false

			c.SetCookie(cookie)

			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/totp_check")
		} else {
			err = h.CreateCookie(c, CookieTokenNameMain, token)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			_ = h.useCase.CreateNotificationSignIn(context.Background(), token.Data)

			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/main_page")
		}
	}
}

func (h ApiGatewayHandlers) SignUp() echo.HandlerFunc {
	return func(c echo.Context) error {

		operation_info := &models.SignUpInfo{}
		operation_result := &models.PostRequestStatus{
			Success: false,
		}
		err := h.safeReadFormDataRequest(c, operation_info)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			errPage, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				operation_result.Error = err.Error()
				return c.JSON(http.StatusInternalServerError, operation_result)
			}
			return c.HTML(http.StatusBadRequest, errPage)
		}

		token, err := h.useCase.SignUp(operation_info)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			errPage, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				operation_result.Error = err.Error()
				return c.JSON(http.StatusInternalServerError, operation_result)
			}
			return c.HTML(http.StatusBadRequest, errPage)
		}

		err = h.CreateCookie(c, CookieTokenNameMain, token)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			errPage, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				operation_result.Error = err.Error()
				return c.JSON(http.StatusInternalServerError, operation_result)
			}
			return c.HTML(http.StatusInternalServerError, errPage)
		}

		return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/main_page")
	}
}

func (h ApiGatewayHandlers) SignOut() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		operation_result := &models.PostRequestStatus{
			Success: false,
		}
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			errPage, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				operation_result.Error = err.Error()
				return c.JSON(http.StatusInternalServerError, operation_result)
			}
			return c.HTML(http.StatusBadRequest, errPage)
		}

		if is_ok && token_id != uuid.Nil {
			err = h.useCase.DeleteToken(context.Background(), token_id)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusBadRequest, errPage)
			}
		}

		cookie := new(http.Cookie)
		cookie.Name = CookieTokenNameMain
		cookie.Value = ""
		cookie.Expires = time.Now().Add(-usecase.TokenLiveTime)
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.Secure = false

		c.SetCookie(cookie)

		operation_result.Success = true

		return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/sign_in")
	}
}

func (h ApiGatewayHandlers) OpenAccount() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		operation_result := &models.PostRequestStatus{
			Success: false,
		}
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			errPage, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				operation_result.Error = err.Error()
				return c.JSON(http.StatusInternalServerError, operation_result)
			}
			return c.HTML(http.StatusBadRequest, errPage)
		}

		operation_info := &models.AccountInfo{}
		if err := h.safeReadFormDataRequest(c, operation_info); err != nil {
			utils.LogResponseError(c, h.logger, err)
			errPage, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				operation_result.Error = err.Error()
				return c.JSON(http.StatusInternalServerError, operation_result)
			}
			return c.HTML(http.StatusBadRequest, errPage)
		}

		if is_ok && token_id != uuid.Nil {
			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			err = h.UpdateCookie(c, CookieTokenNameMain)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusBadRequest, errPage)
			}

			err = h.useCase.CreateAccount(user_id, operation_info)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}
		}

		return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/main_page")
	}
}

func (h ApiGatewayHandlers) CloseAccount() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		operation_result := &models.PostRequestStatus{
			Success: false,
		}
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			errPage, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				operation_result.Error = err.Error()
				return c.JSON(http.StatusInternalServerError, operation_result)
			}
			return c.HTML(http.StatusBadRequest, errPage)
		}

		operation_info := &models.AccountInfoRequest{}
		err = h.safeReadFormDataRequest(c, operation_info)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			errPage, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				operation_result.Error = err.Error()
				return c.JSON(http.StatusInternalServerError, operation_result)
			}
			return c.HTML(http.StatusBadRequest, errPage)
		}

		if is_ok && token_id != uuid.Nil {
			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			err = h.UpdateCookie(c, CookieTokenNameMain)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			account_id, err := uuid.Parse(operation_info.AccountId)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusBadRequest, errPage)
			}

			err = h.useCase.CloseAccount(user_id, account_id)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}
		}

		return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/main_page")
	}
}

func (h ApiGatewayHandlers) AddAccountCache() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		operation_result := &models.PostRequestStatus{
			Success: false,
		}
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			errPage, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				operation_result.Error = err.Error()
				return c.JSON(http.StatusInternalServerError, operation_result)
			}
			return c.HTML(http.StatusBadRequest, errPage)
		}

		operation_info := &models.AccountChangeMoneyRequestBody{}
		if err := h.safeReadFormDataRequest(c, operation_info); err != nil {
			utils.LogResponseError(c, h.logger, err)
			errPage, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				operation_result.Error = err.Error()
				return c.JSON(http.StatusInternalServerError, operation_result)
			}
			return c.HTML(http.StatusBadRequest, errPage)
		}

		if is_ok && token_id != uuid.Nil {
			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			err = h.UpdateCookie(c, CookieTokenNameMain)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			account_id, err := uuid.Parse(operation_info.AccountId)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			money, err := strconv.ParseFloat(operation_info.Money, 64)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusBadRequest, errPage)
			}

			err = h.useCase.AddAccountCache(user_id, account_id, money)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}
		}

		return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/main_page")
	}
}

func (h ApiGatewayHandlers) WidthAccountCache() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		operation_result := &models.PostRequestStatus{
			Success: false,
		}
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			errPage, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				operation_result.Error = err.Error()
				return c.JSON(http.StatusInternalServerError, operation_result)
			}
			return c.HTML(http.StatusBadRequest, errPage)
		}

		operation_info := &models.AccountChangeMoneyRequestBody{}
		if err := h.safeReadFormDataRequest(c, operation_info); err != nil {
			utils.LogResponseError(c, h.logger, err)
			errPage, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				operation_result.Error = err.Error()
				return c.JSON(http.StatusInternalServerError, operation_result)
			}
			return c.HTML(http.StatusInternalServerError, errPage)
		}

		if is_ok && token_id != uuid.Nil {
			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			err = h.UpdateCookie(c, CookieTokenNameMain)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			account_id, err := uuid.Parse(operation_info.AccountId)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			money, err := strconv.ParseFloat(operation_info.Money, 64)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusBadRequest, errPage)
			}

			err = h.useCase.WidthAccountCache(user_id, account_id, money)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}
		}

		return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/main_page")
	}
}

func (h ApiGatewayHandlers) OpenAccountPage() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		if is_ok && token_id != uuid.Nil {

			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			err = h.UpdateCookie(c, CookieTokenNameMain)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			operation_page, err := h.useCase.CreateOpenAccountPage(user_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, operation_page)
		} else {
			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/sign_in")
		}
	}
}

func (h ApiGatewayHandlers) AccountCreditsPage() echo.HandlerFunc {
	return func(c echo.Context) error {

		operation_info := &models.AccountInfoRequest{}
		err := h.safeReadQueryParamsRequest(c, operation_info)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			utils.LogResponseError(c, h.logger, err)
			return c.HTML(http.StatusBadRequest, error_page)
		}

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		if is_ok && token_id != uuid.Nil {

			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			err = h.UpdateCookie(c, CookieTokenNameMain)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			account_id, err := uuid.Parse(operation_info.AccountId)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			operation_page, err := h.useCase.CreateAccountCreditsPage(user_id, account_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, operation_page)
		} else {
			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/sign_in")
		}
	}
}

func (h ApiGatewayHandlers) CloseAccountPage() echo.HandlerFunc {
	return func(c echo.Context) error {

		operation_info := &models.AccountInfoRequest{}
		err := h.safeReadQueryParamsRequest(c, operation_info)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			utils.LogResponseError(c, h.logger, err)
			return c.HTML(http.StatusBadRequest, error_page)
		}

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		if is_ok && token_id != uuid.Nil {

			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			err = h.UpdateCookie(c, CookieTokenNameMain)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			account_id, err := uuid.Parse(operation_info.AccountId)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			operation_page, err := h.useCase.CreateCloseAccountPage(user_id, account_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, operation_page)
		} else {
			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/sign_in")
		}
	}
}

func (h ApiGatewayHandlers) AddAccountCachePage() echo.HandlerFunc {
	return func(c echo.Context) error {

		operation_info := &models.AccountInfoRequest{}
		err := h.safeReadQueryParamsRequest(c, operation_info)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			utils.LogResponseError(c, h.logger, err)
			return c.HTML(http.StatusBadRequest, error_page)
		}

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		if is_ok && token_id != uuid.Nil {

			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			err = h.UpdateCookie(c, CookieTokenNameMain)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			account_id, err := uuid.Parse(operation_info.AccountId)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			operation_page, err := h.useCase.CreateAddAccountCachePage(user_id, account_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, operation_page)
		} else {
			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/sign_in")
		}
	}
}

func (h ApiGatewayHandlers) WidthAccountCachePage() echo.HandlerFunc {
	return func(c echo.Context) error {

		operation_info := &models.AccountInfoRequest{}
		err := h.safeReadQueryParamsRequest(c, operation_info)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			utils.LogResponseError(c, h.logger, err)
			return c.HTML(http.StatusBadRequest, error_page)
		}

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		if is_ok && token_id != uuid.Nil {

			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			err = h.UpdateCookie(c, CookieTokenNameMain)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			account_id, err := uuid.Parse(operation_info.AccountId)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			operation_page, err := h.useCase.CreateWidthAccountCachePage(user_id, account_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, operation_page)
		} else {
			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/sign_in")
		}
	}
}

func (h ApiGatewayHandlers) HomePage() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		if is_ok && token_id != uuid.Nil {

			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			err = h.UpdateCookie(c, CookieTokenNameMain)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			home_page, err := h.useCase.CreateUserPage(user_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, home_page)
		} else {
			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/sign_in")
		}
	}
}

func (h ApiGatewayHandlers) TurnOnTotpPage() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		if is_ok && token_id != uuid.Nil {

			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			err = h.UpdateCookie(c, CookieTokenNameMain)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			totpTurnOnPage, err := h.useCase.CreateTurnOnTotpPage(user_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, totpTurnOnPage)
		} else {
			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/sign_in")
		}
	}
}

func (h ApiGatewayHandlers) TurnOffTotpPage() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		if is_ok && token_id != uuid.Nil {

			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			err = h.UpdateCookie(c, CookieTokenNameMain)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			totpTurnOffPage, err := h.useCase.CreateTurnOffTotpPage(user_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, totpTurnOffPage)
		} else {
			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/sign_in")
		}
	}
}

func (h ApiGatewayHandlers) TurnOnTotp() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		if is_ok && token_id != uuid.Nil {

			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			err = h.UpdateCookie(c, CookieTokenNameMain)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			err = h.useCase.TurnOnTotp(user_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/totp_qr")
		} else {
			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/sign_in")
		}
	}
}

func (h ApiGatewayHandlers) TurnOffTotp() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		if is_ok && token_id != uuid.Nil {

			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			err = h.UpdateCookie(c, CookieTokenNameMain)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			err = h.useCase.TurnOffTotp(user_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/main_page")
		} else {
			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/sign_in")
		}
	}
}

func (h ApiGatewayHandlers) TotpQrPage() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		if is_ok && token_id != uuid.Nil {

			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			err = h.UpdateCookie(c, CookieTokenNameMain)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			qrPage, err := h.useCase.CreateTotpQrPage(user_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, qrPage)
		} else {
			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/sign_in")
		}
	}
}

func (h ApiGatewayHandlers) TotpCheckPage() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c, CookieTokenNameMain)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		if is_ok && token_id != uuid.Nil {
			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/main_page")
		} else {
			checkTotpPage, err := h.useCase.CreateTotpCheckPage()
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, checkTotpPage)
		}
	}
}

func (h ApiGatewayHandlers) TotpCheck() echo.HandlerFunc {
	return func(c echo.Context) error {

		operation_info := &models.TotpCheckInput{}
		operation_result := &models.PostRequestStatus{
			Success: false,
		}
		err := h.safeReadFormDataRequest(c, operation_info)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			errPage, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				operation_result.Error = err.Error()
				return c.JSON(http.StatusInternalServerError, operation_result)
			}
			return c.HTML(http.StatusBadRequest, errPage)
		}
		cookie, err := c.Cookie(CookieTokenNameFirstAuth)
		if err != nil {
			if err == echo.ErrCookieNotFound ||
				err == http.ErrNoCookie {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			} else {
				return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/sign_in")
			}
		}

		if cookie != nil {

			tokenName := cookie.Value

			tokenFirstAuth, err := h.useCase.GetTokenFirstAuth(context.Background(), tokenName)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusBadRequest, error_page)
			}

			err = h.useCase.CheckTotp(tokenFirstAuth.UserId, operation_info.TotpCode)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusBadRequest, error_page)
			}

			cookie := new(http.Cookie)
			cookie.Name = CookieTokenNameFirstAuth
			cookie.Value = ""
			cookie.Expires = time.Now().Add(-usecase.TokenFirstAuthLiveTime)
			cookie.Path = "/"
			cookie.HttpOnly = true
			cookie.Secure = false

			c.SetCookie(cookie)

			mainToken, err := h.useCase.CreateToken(context.Background(), uuid.New(), usecase.TokenLiveTime, tokenFirstAuth.UserId)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			err = h.CreateCookie(c, CookieTokenNameMain, mainToken)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				errPage, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					operation_result.Error = err.Error()
					return c.JSON(http.StatusInternalServerError, operation_result)
				}
				return c.HTML(http.StatusInternalServerError, errPage)
			}

			_ = h.useCase.CreateNotificationSignIn(context.Background(), mainToken.Data)

			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/main_page")

		} else {
			return c.Redirect(http.StatusSeeOther, "/api/v1/api_gateway/sign_in")
		}
	}
}

func (h ApiGatewayHandlers) CreateCookie(c echo.Context, cookieName string, token *models.Token) error {
	cookie := new(http.Cookie)
	cookie.Name = cookieName
	cookie.Value = token.ID.String()
	cookie.Expires = time.Now().Add(usecase.TokenLiveTime)
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = false

	c.SetCookie(cookie)

	return nil
}

func (h ApiGatewayHandlers) CheckToken(c echo.Context, cookieName string) (bool, uuid.UUID, error) {

	cookie, err := c.Cookie(cookieName)
	if err != nil {
		if err == echo.ErrCookieNotFound ||
			err == http.ErrNoCookie {
			return false, uuid.Nil, nil
		} else {
			return false, uuid.Nil, err
		}
	}

	if cookie != nil {

		token_id, err := uuid.Parse(cookie.Value)
		if err != nil {
			return false, uuid.Nil, err
		}

		is_ok, _ := h.useCase.CheckExistingToken(context.Background(), token_id)
		return is_ok, token_id, nil

	} else {
		return false, uuid.Nil, nil
	}

}

func (h ApiGatewayHandlers) UpdateCookie(c echo.Context, cookieName string) error {
	old_cookie, err := c.Cookie(cookieName)
	if err == echo.ErrCookieNotFound ||
		err == http.ErrNoCookie {
		return err
	} else if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = cookieName
	cookie.Value = old_cookie.Value
	cookie.Expires = time.Now().Add(usecase.TokenLiveTime)
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = false

	c.SetCookie(cookie)

	return nil
}

func (h ApiGatewayHandlers) GraphImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Извлекаем путь к файлу из URL
		filePath := c.Param("*")

		// Путь к директории с картинками
		imagesDir := h.folderGraph

		// Формируем полный путь к файлу
		fullPath := filepath.Join(imagesDir, filePath)

		// Проверяем, существует ли файл
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return c.String(http.StatusNotFound, "File not found")
		}

		// Отдаём файл
		return c.File(fullPath)
	}
}

func (h ApiGatewayHandlers) QrImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Извлекаем путь к файлу из URL
		filePath := c.Param("*")

		// Путь к директории с картинками
		imagesDir := h.folderQr

		// Формируем полный путь к файлу
		fullPath := filepath.Join(imagesDir, filePath)

		// Проверяем, существует ли файл
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return c.String(http.StatusNotFound, "File not found")
		}

		// Отдаём файл
		return c.File(fullPath)
	}
}

func NewApiGatewayHandlers(cfg *config.Config, logger logger.Logger, graphFolder string, folderQr string,
	usecase api_gateway.UseCase) api_gateway.Handlers {
	return &ApiGatewayHandlers{cfg: cfg, logger: logger, useCase: usecase, folderGraph: graphFolder, folderQr: folderQr}
}
