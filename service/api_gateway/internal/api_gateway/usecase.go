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
	AddTokenFirstAuth(ctx context.Context, token *models.TokenFirstAuth) error
	GetTokenFirstAuth(ctx context.Context, tokenName string) (*models.TokenFirstAuth, error)
	DeleteTokenFirstAuth(ctx context.Context, tokenName string) error
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
	CreateTurnOnTotpPage(userId uuid.UUID) (string, error)
	CreateTurnOffTotpPage(userId uuid.UUID) (string, error)
	CreateTotpQrPage(userId uuid.UUID) (string, error)
	CreateTotpCheckPage() (string, error)
	CreateAdminPage(begin string, end string) (string, error)
	//
	SignIn(login_info *models.SignInInfo) (*models.Token, error)
	SignUp(sign_up_info *models.SignUpInfo) (*models.Token, error)
	CreateAccount(user_id uuid.UUID, account_info *models.AccountInfo) error
	CloseAccount(user_id uuid.UUID, account_id uuid.UUID) error
	AddAccountCache(user_id uuid.UUID, account_id uuid.UUID, money float64) error
	WidthAccountCache(user_id uuid.UUID, account_id uuid.UUID, money float64) error
	TurnOnTotp(userId uuid.UUID) error
	TurnOffTotp(userId uuid.UUID) error
	CheckTotp(userId uuid.UUID, code string) error
	//
	GetUserTotpInfo(userId uuid.UUID) (*models.TotpInfo, error)
	CreateNotificationSignUp(ctx context.Context, userId uuid.UUID) error
	CreateNotificationSignIn(ctx context.Context, userId uuid.UUID) error
}
