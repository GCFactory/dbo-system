package http

import (
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/registration/config"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration"
	"github.com/labstack/echo/v4"
)

type RegistrationHandlers struct {
	cfg              *config.Config
	registrationGRPC registration.RegistrationGRPCHandlers
	registrationUC   registration.UseCase
	logger           logger.Logger
}

func (h RegistrationHandlers) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func (h RegistrationHandlers) GetOperationStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func NewRegistrationHandlers(cfg *config.Config, logger logger.Logger, usecase registration.UseCase, grpc registration.RegistrationGRPCHandlers) registration.Handlers {
	return &RegistrationHandlers{cfg: cfg, logger: logger, registrationGRPC: grpc, registrationUC: usecase}
}
