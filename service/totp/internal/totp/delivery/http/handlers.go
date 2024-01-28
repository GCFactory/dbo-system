package http

import (
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/httpErrors"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/GCFactory/dbo-system/service/totp/internal/totp"
	totpErrors "github.com/GCFactory/dbo-system/service/totp/pkg/errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
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

// @Summary Enroll new user device
// @Description Create new totp data for user
// @Tags TOTP
// @Accept json
// @Produce json
// @Param user_name string string true "User account name"
// @Param user_id uuid string true "User account uuid"
// @Success 201 {object} models.TOTPEnroll
// @Failure 403 {object} httputil.HTTPError
// @Router /enroll [post]
func (t totpHandlers) Enroll() echo.HandlerFunc {
	type User struct {
		UserName string `json:"user_name" `
		UserId   string `json:"user_id"`
	}
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "auth.Login")
		defer span.Finish()

		user := &User{}
		if err := utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, t.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(errors.Wrap(err, "totpH.utils.ReadRequest")))
		}

		if user.UserId == "" {
			utils.LogResponseError(c, t.logger, errors.New("No user_id field"))
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, "No user_id field", nil))
		}

		if user.UserName == "" {
			utils.LogResponseError(c, t.logger, errors.New("No user_name field"))
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, "No user_name field", nil))
		}

		userId, err := uuid.Parse(user.UserId)
		if err != nil {
			utils.LogResponseError(c, t.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(errors.Wrap(err, "totpH.uuid.Parse")))
		}

		totpEnroll, err := t.totpUC.Enroll(ctx, models.TOTPConfig{
			UserId:      userId,
			AccountName: user.UserName,
		})
		if err != nil {
			utils.LogResponseError(c, t.logger, err)
			if errors.Is(err, totpErrors.ActiveTotp) {
				return c.JSON(http.StatusForbidden, httpErrors.NewRestError(http.StatusForbidden, err.Error(), nil))
			} else {
				return c.JSON(http.StatusInternalServerError, err)
			}
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

		totpEnable, _ := t.totpUC.Enable(ctx, userId)
		return c.JSON(http.StatusOK, totpEnable)
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
