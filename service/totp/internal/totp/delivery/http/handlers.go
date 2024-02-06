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

type httpError struct {
	_ httpErrors.RestError
}

type totpHandlers struct {
	cfg    *config.Config
	totpUC totp.UseCase
	logger logger.Logger
}

func NewTOTPHandlers(cfg *config.Config, totpUC totp.UseCase, log logger.Logger) totp.Handlers {
	return &totpHandlers{cfg: cfg, totpUC: totpUC, logger: log}
}

// @Summary		Enroll new user device
// @Description	Create new totp data for user
// @Tags			TOTP
// @Accept			json
// @Produce		json
// @Param			user_name	body		string	true	"User account name"
// @Param			user_id		body		string	true	"User account uuid"
// @Success		201			{object}	models.TOTPEnroll
// @Failure		403			{object}	httpError
// @Router			/enroll [post]
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

// @Summary		Verify totp URL
// @Description	Verify totp URL
// @Tags			TOTP
// @Accept			json
// @Produce		json
// @Param			totp_url	body		string	true	"TOTP path url"
// @Success		200			{object}	models.TOTPVerify
// @Failure		400			{object}	httpError
// @Failure		400			{object}	models.TOTPVerify
// @Router			/verify [post]
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

// @Summary		Validate totp code
// @Description	Validate users's totp cide
// @Tags			TOTP
// @Accept			json
// @Produce		json
// @Param			user_id		body		string	true	"User id"
// @Param			totp_code	body		string	true	"Totp user's code"
// @Success		200			{object}	models.TOTPValidate
// @Failure		400			{object}	httpError
// @Failure		404			{object}	httpError
// @Router			/validate [post]
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
			if errors.Is(err, totpErrors.NoUserId) {
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

// @Summary		Enable totp code
// @Description	Enable users's totp cide or selected totp code
// @Tags			TOTP
// @Accept			json
// @Produce		json
// @Param			user_id	body		string	true	"User id"
// @Param			totp_id	body		string	true	"Totp id"
// @Success		200		{object}	models.TOTPEnable
// @Failure		400		{object}	httpError
// @Failure		404		{object}	httpError
// @Router			/enable [post]
func (t totpHandlers) Enable() echo.HandlerFunc {
	type Input struct {
		UserId string `json:"user_id"`
		TotpId string `json:"totp_id"`
	}
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "auth.Login")
		defer span.Finish()

		input := &Input{}
		if err := utils.ReadRequest(c, input); err != nil {
			utils.LogResponseError(c, t.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(errors.Wrap(err, "totpH.Enable.utils.ReadRequest")))
		}

		if input.TotpId == "" && input.UserId == "" {
			utils.LogResponseError(c, t.logger, errors.New("No user_id and totp_id fields"))
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, "No user_id and totp_id fields", nil))
		}

		var userId uuid.UUID = uuid.Nil
		var totpId uuid.UUID = uuid.Nil

		var err error

		if input.TotpId != "" {
			totpId, err = uuid.Parse(input.TotpId)
			if err != nil {
				utils.LogResponseError(c, t.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(errors.Wrap(err, "totpH.Enable.uuid(totp_id).Parse")))
			}
		}
		if input.UserId != "" {
			userId, err = uuid.Parse(input.UserId)
			if err != nil {
				utils.LogResponseError(c, t.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(errors.Wrap(err, "totpH.Enable.uuid(user_id).Parse")))
			}
		}

		totpEnable, err := t.totpUC.Enable(ctx, totpId, userId)
		if err != nil {
			if errors.Is(err, totpErrors.NoTotpId) ||
				errors.Is(err, totpErrors.NoId) {
				return c.JSON(http.StatusNotFound, totpEnable)
			} else if errors.Is(err, totpErrors.TotpIsActive) {
				return c.JSON(http.StatusBadRequest, totpEnable)
			} else {
				return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(errors.Wrap(err, err.Error())))
			}
		}
		return c.JSON(http.StatusOK, totpEnable)
	}
}

// @Summary		Disable totp code
// @Description	Disable users's totp cide or selected totp code
// @Tags			TOTP
// @Accept			json
// @Produce		json
// @Param			user_id	body		string	false	"User id"
// @Param			totp_id	body		string	false	"Totp id"
// @Success		200		{object}	models.TOTPDisable
// @Failure		400		{object}	httpError
// @Failure		404		{object}	httpError
// @Router			/disable [post]
func (t totpHandlers) Disable() echo.HandlerFunc {
	type Input struct {
		UserId string `json:"user_id"`
		TotpId string `json:"totp_id"`
	}
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "auth.Login")
		defer span.Finish()

		input := &Input{}
		if err := utils.ReadRequest(c, input); err != nil {
			utils.LogResponseError(c, t.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(errors.Wrap(err, "totpH.Disable.utils.ReadRequest")))
		}

		if input.TotpId == "" && input.UserId == "" {
			utils.LogResponseError(c, t.logger, errors.New("No user_id and totp_id fields"))
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, "No user_id and totp_id fields", nil))
		}

		var userId uuid.UUID = uuid.Nil
		var totpId uuid.UUID = uuid.Nil

		var err error

		if input.TotpId != "" {
			totpId, err = uuid.Parse(input.TotpId)
			if err != nil {
				utils.LogResponseError(c, t.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(errors.Wrap(err, "totpH.Disable.uuid(totp_id).Parse")))
			}
		}
		if input.UserId != "" {
			userId, err = uuid.Parse(input.UserId)
			if err != nil {
				utils.LogResponseError(c, t.logger, err)
				return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(errors.Wrap(err, "totpH.Disable.uuid(user_id).Parse")))
			}
		}

		totpDisable, err := t.totpUC.Disable(ctx, totpId, userId)
		if err != nil {
			if errors.Is(err, totpErrors.NoTotpId) ||
				errors.Is(err, totpErrors.NoId) {
				return c.JSON(http.StatusNotFound, totpDisable)
			} else if errors.Is(err, totpErrors.TotpIsDisabled) {
				return c.JSON(http.StatusBadRequest, totpDisable)
			} else {
				return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(errors.Wrap(err, err.Error())))
			}
		}
		return c.JSON(http.StatusOK, totpDisable)
	}
}
