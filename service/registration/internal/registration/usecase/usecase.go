package usecase

import (
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type registrationUC struct {
	cfg              *config.Config
	registrationRepo registration.Repository
	logger           logger.Logger
}

func (regUC *registrationUC) CreateSaga(ctx context.Context, saga models.Saga) error {
	// TODO: реализовать
	return nil
}

func (regUC *registrationUC) AddSagaEvent(ctx context.Context, saga_uuid uuid.UUID, event_name string) error {
	// TODO: реализовать
	return nil
}

func (regUC *registrationUC) RemoveSagaEvent(ctx context.Context, saga_uuid uuid.UUID, event_name string) error {
	// TODO: реализовать
	return nil
}

func (regUC *registrationUC) FallBackSaga(ctx context.Context, saga_uuid uuid.UUID) error {
	// TODO: реализовать
	return nil
}

func NewRegistrationUseCase(cfg *config.Config, registration_repo registration.Repository, log logger.Logger) registration.UseCase {
	return &registrationUC{cfg: cfg, registrationRepo: registration_repo, logger: log}
}
