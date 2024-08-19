package http

import (
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/httpErrors"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/account/internal/account"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type httpError struct {
	_ httpErrors.RestError
}

type accHandlers struct {
	cfg    *config.Config
	accUC  account.UseCase
	logger logger.Logger
}

func NewACCHandlers(cfg *config.Config, accUC account.UseCase, log logger.Logger) account.Handlers {
	return &accHandlers{cfg: cfg, accUC: accUC, logger: log}
}

func (t accHandlers) TestFunc() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, _ := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "accH.testFunc")
		defer span.Finish()

		//	TODO: реализовать тут

		return c.JSON(http.StatusCreated, nil)
	}
}
