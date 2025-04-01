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
}

const (
	CookieTokenName string = "token"
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

func (h ApiGatewayHandlers) SignInPage() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		if !is_ok {
			page, err := h.useCase.CreateSignInPage()
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
			page, err := h.useCase.CreateUserPage(user_id)
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
}

func (h ApiGatewayHandlers) SignUpPage() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c)
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
		err := h.safeReadBodyRequest(c, operation_info)
		if err != nil {
			operation_result.Error = err.Error()
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, operation_result)
		}

		token, err := h.useCase.SignIn(operation_info)
		if err != nil {
			operation_result.Error = err.Error()
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, operation_result)
		}

		err = h.CreateCookie(c, token)
		if err != nil {
			operation_result.Error = err.Error()
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, operation_result)
		}

		operation_result.Success = true

		return c.JSON(http.StatusOK, operation_result)
	}
}

func (h ApiGatewayHandlers) SignUp() echo.HandlerFunc {
	return func(c echo.Context) error {

		operation_info := &models.SignUpInfo{}
		operation_result := &models.PostRequestStatus{
			Success: false,
		}
		if err := h.safeReadBodyRequest(c, operation_info); err != nil {
			operation_result.Error = err.Error()
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, operation_result)
		}

		token, err := h.useCase.SignUp(operation_info)
		if err != nil {
			operation_result.Error = err.Error()
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, operation_result)
		}

		err = h.CreateCookie(c, token)
		if err != nil {
			operation_result.Error = err.Error()
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, operation_result)
		}

		return c.JSON(http.StatusOK, operation_result)
	}
}

func (h ApiGatewayHandlers) SignOut() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c)
		operation_result := &models.PostRequestStatus{
			Success: false,
		}
		if err != nil {
			operation_result.Error = err.Error()
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, operation_result)
		}

		if is_ok && token_id != uuid.Nil {
			err = h.useCase.DeleteToken(context.Background(), token_id)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}
		}

		cookie := new(http.Cookie)
		cookie.Name = CookieTokenName
		cookie.Value = ""
		cookie.Expires = time.Now().Add(-usecase.TokenLiveTime)
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.Secure = false

		c.SetCookie(cookie)

		operation_result.Success = true

		return c.JSON(http.StatusOK, operation_result)
	}
}

func (h ApiGatewayHandlers) OpenAccount() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c)
		operation_result := &models.PostRequestStatus{
			Success: false,
		}
		if err != nil {
			operation_result.Error = err.Error()
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, operation_result)
		}
		operation_info := &models.AccountInfo{}
		if err := h.safeReadBodyRequest(c, operation_info); err != nil {
			operation_result.Error = err.Error()
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, operation_result)
		}

		if is_ok && token_id != uuid.Nil {
			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			err = h.UpdateCookie(c)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			err = h.useCase.CreateAccount(user_id, operation_info)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}
		}

		operation_result.Success = true

		return c.JSON(http.StatusOK, operation_result)
	}
}

func (h ApiGatewayHandlers) CloseAccount() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c)
		operation_result := &models.PostRequestStatus{
			Success: false,
		}
		if err != nil {
			operation_result.Error = err.Error()
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, operation_result)
		}
		operation_info := &models.AccountInfoRequest{}
		if err := h.safeReadBodyRequest(c, operation_info); err != nil {
			operation_result.Error = err.Error()
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, operation_result)
		}

		if is_ok && token_id != uuid.Nil {
			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			err = h.UpdateCookie(c)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			account_id, err := uuid.Parse(operation_info.AccountId)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			err = h.useCase.CloseAccount(user_id, account_id)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}
		}

		operation_result.Success = true

		return c.JSON(http.StatusOK, operation_result)
	}
}

func (h ApiGatewayHandlers) AddAccountCache() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c)
		operation_result := &models.PostRequestStatus{
			Success: false,
		}
		if err != nil {
			operation_result.Error = err.Error()
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, operation_result)
		}
		operation_info := &models.AccountChangeMoneyRequestBody{}
		if err := h.safeReadBodyRequest(c, operation_info); err != nil {
			operation_result.Error = err.Error()
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, operation_result)
		}

		if is_ok && token_id != uuid.Nil {
			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			err = h.UpdateCookie(c)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			account_id, err := uuid.Parse(operation_info.AccountId)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			money, err := strconv.ParseFloat(operation_info.Money, 64)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			err = h.useCase.AddAccountCache(user_id, account_id, money)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}
		}

		operation_result.Success = true

		return c.JSON(http.StatusOK, operation_result)
	}
}

func (h ApiGatewayHandlers) WidthAccountCache() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c)
		operation_result := &models.PostRequestStatus{
			Success: false,
		}
		if err != nil {
			operation_result.Error = err.Error()
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, operation_result)
		}
		operation_info := &models.AccountChangeMoneyRequestBody{}
		if err := h.safeReadBodyRequest(c, operation_info); err != nil {
			operation_result.Error = err.Error()
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, operation_result)
		}

		if is_ok && token_id != uuid.Nil {
			err = h.useCase.UpdateToken(context.Background(), token_id, usecase.TokenLiveTime)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			err = h.UpdateCookie(c)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			user_id, err := h.useCase.GetTokenValue(context.Background(), token_id)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			account_id, err := uuid.Parse(operation_info.AccountId)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			money, err := strconv.ParseFloat(operation_info.Money, 64)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}

			err = h.useCase.WidthAccountCache(user_id, account_id, money)
			if err != nil {
				operation_result.Error = err.Error()
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, operation_result)
			}
		}

		operation_result.Success = true

		return c.JSON(http.StatusOK, operation_result)
	}
}

func (h ApiGatewayHandlers) OpenAccountPage() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c)
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

			err = h.UpdateCookie(c)
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
			sign_in_page, err := h.useCase.CreateSignInPage()
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, sign_in_page)
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

		is_ok, token_id, err := h.CheckToken(c)
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

			err = h.UpdateCookie(c)
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
			sign_in_page, err := h.useCase.CreateSignInPage()
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, sign_in_page)
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

		is_ok, token_id, err := h.CheckToken(c)
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

			err = h.UpdateCookie(c)
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
			sign_in_page, err := h.useCase.CreateSignInPage()
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, sign_in_page)
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

		is_ok, token_id, err := h.CheckToken(c)
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

			err = h.UpdateCookie(c)
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
			sign_in_page, err := h.useCase.CreateSignInPage()
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, sign_in_page)
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

		is_ok, token_id, err := h.CheckToken(c)
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

			err = h.UpdateCookie(c)
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
			sign_in_page, err := h.useCase.CreateSignInPage()
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, sign_in_page)
		}
	}
}

func (h ApiGatewayHandlers) HomePage() echo.HandlerFunc {
	return func(c echo.Context) error {

		is_ok, token_id, err := h.CheckToken(c)
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

			err = h.UpdateCookie(c)
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
			sign_in_page, err := h.useCase.CreateSignInPage()
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}

			return c.HTML(http.StatusOK, sign_in_page)
		}
	}
}

func (h ApiGatewayHandlers) CreateCookie(c echo.Context, token *models.Token) error {
	cookie := new(http.Cookie)
	cookie.Name = CookieTokenName
	cookie.Value = token.ID.String()
	cookie.Expires = time.Now().Add(usecase.TokenLiveTime)
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = false

	c.SetCookie(cookie)

	return nil
}

func (h ApiGatewayHandlers) CheckToken(c echo.Context) (bool, uuid.UUID, error) {

	cookie, err := c.Cookie(CookieTokenName)
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

func (h ApiGatewayHandlers) UpdateCookie(c echo.Context) error {
	old_cookie, err := c.Cookie(CookieTokenName)
	if err == echo.ErrCookieNotFound ||
		err == http.ErrNoCookie {
		return err
	} else if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = CookieTokenName
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

func NewApiGatewayHandlers(cfg *config.Config, logger logger.Logger, graphFolder string, usecase api_gateway.UseCase) api_gateway.Handlers {
	return &ApiGatewayHandlers{cfg: cfg, logger: logger, useCase: usecase, folderGraph: graphFolder}
}
