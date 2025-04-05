package http

import (
	"fmt"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/httpErrors"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/GCFactory/dbo-system/service/totp/internal/totp"
	totpErrors "github.com/GCFactory/dbo-system/service/totp/pkg/errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"net/http"
	"time"
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
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpH.Enroll")
		defer span.Finish()

		user := &User{}
		result := &models.DefaultHttpRequest{
			Info:   "",
			Status: http.StatusCreated,
		}
		if err := utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, t.logger, err)
			result.Status = http.StatusBadRequest
			result.Info = err.Error()
			return c.JSON(http.StatusBadRequest, result)
		}

		if user.UserId == "" {
			utils.LogResponseError(c, t.logger, errors.New("No user_id field"))
			result.Status = http.StatusBadRequest
			result.Info = "No user_id field"
			return c.JSON(http.StatusBadRequest, result)
		}

		if user.UserName == "" {
			utils.LogResponseError(c, t.logger, errors.New("No user_name field"))
			result.Status = http.StatusBadRequest
			result.Info = "No user_name field"
			return c.JSON(http.StatusBadRequest, result)
		}

		userId, err := uuid.Parse(user.UserId)
		if err != nil {
			utils.LogResponseError(c, t.logger, err)
			result.Status = http.StatusInternalServerError
			result.Info = err.Error()
			return c.JSON(http.StatusInternalServerError, result)
		}

		totpEnroll, err := t.totpUC.Enroll(ctx, models.TOTPConfig{
			UserId:      userId,
			AccountName: user.UserName,
		})
		if err != nil {
			utils.LogResponseError(c, t.logger, err)

			if errors.Is(err, totpErrors.ActiveTotp) {
				result.Status = http.StatusForbidden
				result.Info = err.Error()
			} else if errors.Is(err, totpErrors.NoUserName) ||
				errors.Is(err, totpErrors.NoUserId) {
				result.Status = http.StatusBadRequest
				result.Info = err.Error()
			} else {
				result.Status = http.StatusInternalServerError
				result.Info = err.Error()
			}
			return c.JSON(result.Status, result)
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
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpUC.Verify")
		defer span.Finish()

		input := &Url{}
		result := &models.DefaultHttpRequest{
			Info:   "",
			Status: http.StatusOK,
		}
		if err := utils.ReadRequest(c, input); err != nil {
			utils.LogResponseError(c, t.logger, err)
			result.Status = http.StatusInternalServerError
			result.Info = err.Error()
			return c.JSON(http.StatusInternalServerError, result)
		}

		if input.TotpUrl == "" {
			utils.LogResponseError(c, t.logger, errors.New("No totp_url field"))
			result.Status = http.StatusBadRequest
			result.Info = "No totp_url field"
			return c.JSON(http.StatusInternalServerError, result)
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
				result.Status = http.StatusBadRequest
				result.Info = totpVerify.Status
			} else {
				result.Status = http.StatusInternalServerError
				result.Info = err.Error()
			}
			return c.JSON(result.Status, result)
		}

		result.Info = totpVerify.Status
		return c.JSON(http.StatusOK, result)
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
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpUC.Validate")
		defer span.Finish()

		input := &Input{}
		result := &models.DefaultHttpRequest{
			Info:   "",
			Status: http.StatusOK,
		}
		if err := utils.ReadRequest(c, input); err != nil {
			utils.LogResponseError(c, t.logger, err)
			result.Status = http.StatusInternalServerError
			result.Info = err.Error()
			return c.JSON(http.StatusInternalServerError, result)
		}

		if input.UserId == "" {
			utils.LogResponseError(c, t.logger, errors.New("No user_id field"))
			result.Status = http.StatusBadRequest
			result.Info = "No user_id field"
			return c.JSON(http.StatusBadRequest, result)
		}

		if input.TotpCode == "" {
			utils.LogResponseError(c, t.logger, errors.New("No totp_code field"))
			result.Status = http.StatusBadRequest
			result.Info = "No totp_code field"
			return c.JSON(http.StatusBadRequest, result)
		}

		userId, err := uuid.Parse(input.UserId)
		if err != nil {
			utils.LogResponseError(c, t.logger, err)
			result.Status = http.StatusInternalServerError
			result.Info = err.Error()
			return c.JSON(http.StatusInternalServerError, result)
		}

		totpValidate, err := t.totpUC.Validate(ctx, userId, input.TotpCode, time.Now())
		if err != nil {
			if errors.Is(err, totpErrors.NoUserId) {
				result.Status = http.StatusNotFound
				result.Info = totpValidate.Status
			} else if errors.Is(err, totpErrors.WrongTotpCode) {
				result.Status = http.StatusBadRequest
				result.Info = totpValidate.Status
			} else {
				result.Status = http.StatusInternalServerError
				result.Info = err.Error()
			}
			return c.JSON(result.Status, result)
		}

		result.Info = totpValidate.Status
		return c.JSON(http.StatusOK, result)
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
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpUC.Enable")
		defer span.Finish()

		input := &Input{}
		result := &models.DefaultHttpRequest{
			Info:   "",
			Status: http.StatusOK,
		}
		if err := utils.ReadRequest(c, input); err != nil {
			utils.LogResponseError(c, t.logger, err)
			result.Status = http.StatusInternalServerError
			result.Info = err.Error()
			return c.JSON(http.StatusInternalServerError, result)
		}

		if input.TotpId == "" && input.UserId == "" {
			utils.LogResponseError(c, t.logger, errors.New("No user_id and totp_id fields"))
			result.Status = http.StatusBadRequest
			result.Info = "No user_id and totp_id fields"
			return c.JSON(http.StatusBadRequest, result)
		}

		var userId uuid.UUID = uuid.Nil
		var totpId uuid.UUID = uuid.Nil

		var err error

		if input.TotpId != "" {
			totpId, err = uuid.Parse(input.TotpId)
			if err != nil {
				utils.LogResponseError(c, t.logger, err)
				result.Status = http.StatusInternalServerError
				result.Info = err.Error()
				return c.JSON(http.StatusInternalServerError, result)
			}
		}
		if input.UserId != "" {
			userId, err = uuid.Parse(input.UserId)
			if err != nil {
				utils.LogResponseError(c, t.logger, err)
				result.Status = http.StatusInternalServerError
				result.Info = err.Error()
				return c.JSON(http.StatusInternalServerError, result)
			}
		}

		totpEnable, err := t.totpUC.Enable(ctx, totpId, userId)
		if err != nil {
			if errors.Is(err, totpErrors.NoTotpId) ||
				errors.Is(err, totpErrors.NoId) {
				result.Status = http.StatusNotFound
				result.Info = totpEnable.Status
			} else if errors.Is(err, totpErrors.TotpIsActive) {
				result.Status = http.StatusBadRequest
				result.Info = totpEnable.Status
			} else {
				result.Status = http.StatusInternalServerError
				result.Info = err.Error()
			}
			return c.JSON(result.Status, result)
		}
		result.Info = totpEnable.Status
		return c.JSON(http.StatusOK, result)
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
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpUC.Disable")
		defer span.Finish()

		input := &Input{}
		result := &models.DefaultHttpRequest{
			Info:   "",
			Status: http.StatusOK,
		}
		if err := utils.ReadRequest(c, input); err != nil {
			utils.LogResponseError(c, t.logger, err)
			result.Status = http.StatusInternalServerError
			result.Info = err.Error()
			return c.JSON(http.StatusInternalServerError, result)
		}

		if input.TotpId == "" && input.UserId == "" {
			utils.LogResponseError(c, t.logger, errors.New("No user_id and totp_id fields"))
			result.Status = http.StatusBadRequest
			result.Info = "No user_id and totp_id fields"
			return c.JSON(http.StatusBadRequest, result)
		}

		var userId uuid.UUID = uuid.Nil
		var totpId uuid.UUID = uuid.Nil

		var err error

		if input.TotpId != "" {
			totpId, err = uuid.Parse(input.TotpId)
			if err != nil {
				utils.LogResponseError(c, t.logger, err)
				result.Status = http.StatusInternalServerError
				result.Info = err.Error()
				return c.JSON(http.StatusInternalServerError, result)
			}
		}
		if input.UserId != "" {
			userId, err = uuid.Parse(input.UserId)
			if err != nil {
				utils.LogResponseError(c, t.logger, err)
				result.Status = http.StatusInternalServerError
				result.Info = err.Error()
				return c.JSON(http.StatusInternalServerError, result)
			}
		}

		totpDisable, err := t.totpUC.Disable(ctx, totpId, userId)
		if err != nil {
			if errors.Is(err, totpErrors.NoTotpId) ||
				errors.Is(err, totpErrors.NoId) {
				result.Status = http.StatusNotFound
				result.Info = totpDisable.Status
			} else if errors.Is(err, totpErrors.TotpIsDisabled) {
				result.Status = http.StatusBadRequest
				result.Info = totpDisable.Status
			} else {
				result.Status = http.StatusInternalServerError
				result.Info = err.Error()
			}
			return c.JSON(result.Status, result)
		}

		result.Info = totpDisable.Status
		return c.JSON(http.StatusOK, result)
	}
}

func (t totpHandlers) Url() echo.HandlerFunc {

	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "totpUC.Url")
		defer span.Finish()

		operation_result := &models.DefaultHttpRequest{}
		operation_info := &models.GetTotpUrlBody{}
		err := t.safeReadBodyRequest(c, operation_info)
		if err != nil {
			operation_result.Status = http.StatusBadRequest
			operation_result.Info = err.Error()
			return c.JSON(operation_result.Status, operation_result)
		}

		totpInfo, err := t.totpUC.Url(ctx, operation_info.UserId)
		if err != nil {
			operation_result.Status = http.StatusInternalServerError
			operation_result.Info = err.Error()
			return c.JSON(operation_result.Status, operation_result)
		}

		result := &models.GetTotpUrlResponse{
			TotpUrl: totpInfo.TotpUrl,
		}

		return c.JSON(http.StatusOK, result)
	}
}

func (h totpHandlers) safeReadBodyRequest(c echo.Context, v interface{}) error {
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
