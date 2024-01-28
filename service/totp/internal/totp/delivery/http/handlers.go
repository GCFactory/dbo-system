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
		UserName string `json:"user_name"`
		UserId   string `json:"user_id"`
	}
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "auth.Login")
		defer span.Finish()

		user := &User{}
		if err := utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, t.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(errors.Wrap(err, "totpH.Enroll.utils.ReadRequest")))
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
			return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(errors.Wrap(err, "totpH.Enroll.uuid.Parse")))
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

// @Summary Verify totp URL
// @Description Verify totp URL
// @Tags TOTP
// @Accept json
// @Produce json
// @Param totp_url url string true "TOTP path url"
// @Success 200 {object} models.TOTPVerify
// @Failure 400 {object} httputil.HTTPError
// @Failure 400 {object} models.TOTPVerify
// @Router /verify [post]
func (t totpHandlers) Verify() echo.HandlerFunc {
	type Url struct {
		TotpUrl string `json:"totp_url"`
	}
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "auth.Login")
		defer span.Finish()

		input := &Url{}
		if err := utils.ReadRequest(c, input); err != nil {
			utils.LogResponseError(c, t.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(errors.Wrap(err, "totpH.Verify.utils.ReadRequest")))
		}

		if input.TotpUrl == "" {
			utils.LogResponseError(c, t.logger, errors.New("No totp_url field"))
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, "No totp_url field", nil))
		}

		url := input.TotpUrl
		totpVerify, err := t.totpUC.Verify(ctx, url)
		if err != nil {
			utils.LogResponseError(c, t.logger, err)
			if errors.Is(err, totpErrors.NoAlgorithmField) ||
				errors.Is(err, totpErrors.NoDigitsField) ||
				errors.Is(err, totpErrors.NoIssuerField) ||
				errors.Is(err, totpErrors.NoPeriodField) ||
				errors.Is(err, totpErrors.NoSecretField) {
				return c.JSON(http.StatusBadRequest, totpVerify)
			} else {
				return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(errors.Wrap(err, "totpH.Verify.utils.ReadRequest")))
			}
		}
		return c.JSON(http.StatusOK, totpVerify)
	}
}

// @Summary Validate totp code
// @Description Validate users's titp cide
// @Tags TOTP
// @Accept json
// @Produce json
// @Param user_id uuid string true "User id"
// @Param totp_code string string true "Totp user's code"
// @Success 200 {object} models.TotpValidate
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Router /verify [post]
func (t totpHandlers) Validate() echo.HandlerFunc {
	type Input struct {
		UserId   string `json:"user_id"`
		TotpCode string `json:"totp_code"`
	}
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "auth.Login")
		defer span.Finish()

		input := &Input{}
		if err := utils.ReadRequest(c, input); err != nil {
			utils.LogResponseError(c, t.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(errors.Wrap(err, "totpH.Validate.utils.ReadRequest")))
		}

		if input.UserId == "" {
			utils.LogResponseError(c, t.logger, errors.New("No user_id field"))
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, "No user_id field", nil))
		}

		if input.TotpCode == "" {
			utils.LogResponseError(c, t.logger, errors.New("No totp_code field"))
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, "No totp_code field", nil))
		}

		userId, err := uuid.Parse(input.UserId)
		if err != nil {
			utils.LogResponseError(c, t.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(errors.Wrap(err, "totpH.Validate.uuid.Parse")))
		}

		totpValidate, err := t.totpUC.Validate(ctx, userId, input.TotpCode)
		if err != nil {
			if errors.Is(err, totpErrors.NoUserTotp) {
				return c.JSON(http.StatusNotFound, totpValidate)
			} else if errors.Is(err, totpErrors.WrongTotpCode) {
				return c.JSON(http.StatusBadRequest, totpValidate)
			} else {
				return c.JSON(httpErrors.ErrorResponse(err))
			}
		}

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
