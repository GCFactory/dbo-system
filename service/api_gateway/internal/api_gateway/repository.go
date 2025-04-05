package api_gateway

import (
	"context"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/models"
	"github.com/google/uuid"
	"time"
)

type Repository interface {
	AddToken(ctx context.Context, token *models.Token) error
	GetToken(ctx context.Context, token_id uuid.UUID) (*models.Token, error)
	UpdateToken(tx context.Context, token_id uuid.UUID, new_live_time time.Duration) error
	DeleteToken(ctx context.Context, token_id uuid.UUID) error
	//
	AddTokenFirstAuth(ctx context.Context, token *models.TokenFirstAuth) error
	GetTokenFirstAuth(ctx context.Context, tokenName string) (*models.TokenFirstAuth, error)
	DeleteTokenFirstAuth(ctx context.Context, tokenName string) error
}
