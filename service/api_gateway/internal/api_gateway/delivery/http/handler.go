package http

import (
	"context"
	"fmt"
	"github.com/GCFactory/dbo-system/platform/pkg/httpErrors"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/api_gateway/config"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type ApiGatewayHandlers struct {
	cfg     *config.Config
	useCase api_gateway.UseCase
	logger  logger.Logger
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

func (h ApiGatewayHandlers) SignIn() echo.HandlerFunc {
	return func(c echo.Context) error {

		operation_info := &models.SignInInfo{}
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

		home_page, token, err := h.useCase.SignIn(operation_info)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			utils.LogResponseError(c, h.logger, err)
			return c.HTML(http.StatusBadRequest, error_page)
		}

		err = h.CreateCookie(c, token)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			utils.LogResponseError(c, h.logger, err)
			return c.HTML(http.StatusBadRequest, error_page)
		}

		//c.Set("Content-Type", "text/html; charset=utf-8")

		return c.HTML(http.StatusOK, home_page)
	}
}

func (h ApiGatewayHandlers) SignUp() echo.HandlerFunc {
	return func(c echo.Context) error {

		operation_info := &models.SignUpInfo{}
		if err := h.safeReadQueryParamsRequest(c, operation_info); err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			utils.LogResponseError(c, h.logger, err)
			return c.HTML(http.StatusBadRequest, error_page)
		}

		home_page, token, err := h.useCase.SignUp(operation_info)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			utils.LogResponseError(c, h.logger, err)
			return c.HTML(http.StatusBadRequest, error_page)
		}

		err = h.CreateCookie(c, token)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			utils.LogResponseError(c, h.logger, err)
			return c.HTML(http.StatusBadRequest, error_page)
		}

		return c.HTML(http.StatusOK, home_page)
	}
}

func (h ApiGatewayHandlers) SignOut() echo.HandlerFunc {
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
			err = h.useCase.DeleteToken(context.Background(), token_id)
			if err != nil {
				error_page, err := h.useCase.CreateErrorPage(err.Error())
				if err != nil {
					utils.LogResponseError(c, h.logger, err)
					return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
				}
				return c.HTML(http.StatusInternalServerError, error_page)
			}
		}
		sign_in_page, err := h.useCase.CreateSignInPage()
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			return c.HTML(http.StatusInternalServerError, error_page)
		}

		cookie := &http.Cookie{
			Name:   CookieTokenName,
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}

		c.SetCookie(cookie)

		return c.HTML(http.StatusOK, sign_in_page)
	}
}

func (h ApiGatewayHandlers) CreateCookie(c echo.Context, token *models.Token) error {

	cookie, err := c.Cookie(CookieTokenName)
	if err != echo.ErrCookieNotFound &&
		err != http.ErrNoCookie {
		return err
	}

	if cookie == nil {
		cookie = &http.Cookie{
			Name:     CookieTokenName,
			Value:    token.ID.String(),
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			Expires:  time.Now().Add(token.Live_time),
		}
	} else {
		cookie.Value = token.ID.String()
		cookie.Expires = time.Now().Add(token.Live_time)
	}

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

		is_ok, err := h.useCase.CheckExistingToken(context.Background(), token_id)
		if err != nil {
			return false, uuid.Nil, err
		}
		return is_ok, token_id, nil

	} else {
		return false, uuid.Nil, nil
	}

}

func NewApiGatewayHandlers(cfg *config.Config, logger logger.Logger, usecase api_gateway.UseCase) api_gateway.Handlers {
	return &ApiGatewayHandlers{cfg: cfg, logger: logger, useCase: usecase}
}
