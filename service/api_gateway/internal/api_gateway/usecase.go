package api_gateway

import (
	"context"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/models"
	"github.com/google/uuid"
	"time"
)

type UseCase interface {
	CreateToken(ctx context.Context, token_id uuid.UUID, live_time time.Duration) (*models.Token, error)
	CheckExistingToken(ctx context.Context, token_id uuid.UUID) (bool, error)
	UpdateToken(ctx context.Context, token_id uuid.UUID, new_expire_time time.Duration) error
	DeleteToken(ctx context.Context, token_id uuid.UUID) error
	//
	CreateSignInPage() (string, error)
	CreateErrorPage(error string) (string, error)
	CreateSignUpPage() (string, error)
}
