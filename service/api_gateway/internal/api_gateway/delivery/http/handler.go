package http

import (
	"fmt"
	"github.com/GCFactory/dbo-system/platform/pkg/httpErrors"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/api_gateway/config"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ApiGatewayHandlers struct {
	cfg     *config.Config
	useCase api_gateway.UseCase
	logger  logger.Logger
}

func (h ApiGatewayHandlers) safeReadRequest(c echo.Context, v interface{}) error {
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

func (h ApiGatewayHandlers) SignInPage() echo.HandlerFunc {
	return func(c echo.Context) error {

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
	}
}

func (h ApiGatewayHandlers) SignUpPage() echo.HandlerFunc {
	return func(c echo.Context) error {

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
	}
}

func (h ApiGatewayHandlers) SignIn() echo.HandlerFunc {
	return func(c echo.Context) error {

		operation_info := &models.SignInInfo{}
		if err := h.safeReadRequest(c, operation_info); err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
			}
			utils.LogResponseError(c, h.logger, err)
			return c.HTML(http.StatusBadRequest, error_page)
		}

		res, err := h.useCase.SignIn(operation_info)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			utils.LogResponseError(c, h.logger, err)
			return c.HTML(http.StatusBadRequest, error_page)
		}

		data := make(map[string]interface{})
		data["res"] = res

		return c.JSON(http.StatusOK, data)
	}
}

func (h ApiGatewayHandlers) SignUp() echo.HandlerFunc {
	return func(c echo.Context) error {

		operation_info := &models.SignUpInfo{}
		if err := h.safeReadRequest(c, operation_info); err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil))
			}
			utils.LogResponseError(c, h.logger, err)
			return c.HTML(http.StatusBadRequest, error_page)
		}

		_, err := h.useCase.SignUp(operation_info)
		if err != nil {
			error_page, err := h.useCase.CreateErrorPage(err.Error())
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewRestError(http.StatusInternalServerError, err.Error(), nil))
			}
			utils.LogResponseError(c, h.logger, err)
			return c.HTML(http.StatusBadRequest, error_page)
		}

		//TODO: тут создать страницу пользователя

		return c.JSON(http.StatusOK, "")
	}
}

func NewApiGatewayHandlers(cfg *config.Config, logger logger.Logger, usecase api_gateway.UseCase) api_gateway.Handlers {
	return &ApiGatewayHandlers{cfg: cfg, logger: logger, useCase: usecase}
}
