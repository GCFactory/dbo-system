package api_gateway

import (
	"context"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/models"
	"github.com/google/uuid"
	"time"
)

type UseCase interface {
	CreateToken(ctx context.Context, token_id uuid.UUID, live_time time.Duration, token_value uuid.UUID) (*models.Token, error)
	CheckExistingToken(ctx context.Context, token_id uuid.UUID) (bool, error)
	GetTokenValue(ctx context.Context, token_id uuid.UUID) (uuid.UUID, error)
	UpdateToken(ctx context.Context, token_id uuid.UUID, new_expire_time time.Duration) error
	DeleteToken(ctx context.Context, token_id uuid.UUID) error
	//
	CreateSignInPage() (string, error)
	CreateErrorPage(error string) (string, error)
	CreateSignUpPage() (string, error)
	CreateUserPage(user_id uuid.UUID) (string, error)
	CreateOpenAccountPage(user_id uuid.UUID) (string, error)
	CreateAccountCreditsPage(user_id uuid.UUID, account_id uuid.UUID) (string, error)
	CreateCloseAccountPage(user_id uuid.UUID, account_id uuid.UUID) (string, error)
	CreateAddAccountCachePage(user_id uuid.UUID, account_id uuid.UUID) (string, error)
	CreateWidthAccountCachePage(user_id uuid.UUID, account_id uuid.UUID) (string, error)
	//
	SignIn(login_info *models.SignInInfo) (*models.Token, error)
	SignUp(sign_up_info *models.SignUpInfo) (*models.Token, error)
	CreateAccount(user_id uuid.UUID, account_info *models.AccountInfo) error
}
