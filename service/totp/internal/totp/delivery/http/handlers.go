package http

import (
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/httpErrors"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/GCFactory/dbo-system/service/totp/internal/totp"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type totpHandlers struct {
	cfg    *config.Config
	totpUC totp.UseCase
	logger logger.Logger
}

func NewTOTPHandlers(cfg *config.Config, totpUC totp.UseCase, log logger.Logger) totp.Handlers {
	return &totpHandlers{cfg: cfg, totpUC: totpUC, logger: log}
}

// Enroll user
// @Summary Enroll new user device
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} models.TOTPEnroll
// @Router /enroll [post]
func (t totpHandlers) Enroll() echo.HandlerFunc {
	type User struct {
		UserId string `json:"user_id"`
	}
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "auth.Login")
		defer span.Finish()

		user := &User{}
		if err := utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, t.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		t.logger.Info(user.UserId)
		userId, err := uuid.Parse(user.UserId)
		if err != nil {
			utils.LogResponseError(c, t.logger, err)
			return c.JSON(http.StatusBadRequest, err)
		}

		totpEnroll, err := t.totpUC.Enroll(ctx, models.TOTPConfig{
			UserId: userId,
		})
		if err != nil {
			utils.LogResponseError(c, t.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, totpEnroll)
	}
}

// По url
// Verify otp URL
// @Summary Verify otp URL
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} models.TOTPVerify
// @Router /verify [post]
func (t totpHandlers) Verify() echo.HandlerFunc {
	type Url struct {
		URL string `json:"url"`
	}
	return func(c echo.Context) error {
		//Вопросы
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "auth.Login")
		defer span.Finish()

		input := &Url{}
		if err := utils.ReadRequest(c, input); err != nil {
			utils.LogResponseError(c, t.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		t.logger.Info(input.URL)

		url := input.URL
		totpVerify, _ := t.totpUC.Verify(ctx, url)
		return c.JSON(http.StatusOK, totpVerify)
	}
}

// По коду
func (t totpHandlers) Validate() echo.HandlerFunc {
	type Input struct {
		Id   string `json:"id"`
		Code string `json:"code"`
	}
	return func(c echo.Context) error {
		//Вопросы
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "auth.Login")
		defer span.Finish()

		input := &Input{}
		if err := utils.ReadRequest(c, input); err != nil {
			utils.LogResponseError(c, t.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		t.logger.Info(input.Code)
		userId, err := uuid.Parse(input.Id)
		if err != nil {
			utils.LogResponseError(c, t.logger, err)
			return c.JSON(http.StatusBadRequest, err)
		}

		totpValidate, _ := t.totpUC.Validate(ctx, userId, input.Code)
		return c.JSON(http.StatusOK, totpValidate)
	}
}

func (t totpHandlers) Enable() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, "")
	}
}

func (t totpHandlers) Disable() echo.HandlerFunc {
	type Input struct {
		Id string `json:"id"`
	}
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "auth.Login")
		defer span.Finish()

		input := &Input{}
		if err := utils.ReadRequest(c, input); err != nil {
			utils.LogResponseError(c, t.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		t.logger.Info(input.Id)
		userId, err := uuid.Parse(input.Id)
		if err != nil {
			utils.LogResponseError(c, t.logger, err)
			return c.JSON(http.StatusBadRequest, err)
		}

		totpDisable, _ := t.totpUC.Disable(ctx, userId)
		return c.JSON(http.StatusOK, totpDisable)
	}
}
