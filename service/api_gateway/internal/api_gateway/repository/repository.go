package repository

import (
	"context"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/models"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	redis "github.com/redis/go-redis/v9"
	"time"
)

type apiGatewayRepo struct {
	redis *redis.Client
}

func (repo *apiGatewayRepo) AddToken(ctx context.Context, token *models.Token) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "apiGatewayRepo.AddToken")
	defer span.Finish()

	err := repo.redis.Set(ctxWithTrace, token.ID.String(), token.Data.String(), token.Live_time).Err()
	if err != nil {
		return ErrorAddToken
	}

	return nil
}

func (repo *apiGatewayRepo) GetToken(ctx context.Context, token_id uuid.UUID) (*models.Token, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "apiGatewayRepo.GetToken")
	defer span.Finish()

	data_str, err := repo.redis.Get(ctxWithTrace, token_id.String()).Result()
	if err != nil {
		return nil, ErrorGetTokenValue
	}

	data, err := uuid.Parse(data_str)
	if err != nil {
		return nil, err
	}

	cmd := repo.redis.TTL(ctxWithTrace, token_id.String())
	if cmd.Err() != nil {
		return nil, ErrorGetTokenExpire
	}

	time_live := cmd.Val()

	return &models.Token{
		ID:        token_id,
		Data:      data,
		Live_time: time_live,
	}, nil
}

func (repo *apiGatewayRepo) UpdateToken(ctx context.Context, token_id uuid.UUID, new_live_time time.Duration) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "apiGatewayRepo.UpdateToken")
	defer span.Finish()

	err := repo.redis.Expire(ctxWithTrace, token_id.String(), new_live_time).Err()
	if err != nil {
		return ErrorUpdateTokenExpire
	}

	return nil
}

func (repo *apiGatewayRepo) DeleteToken(ctx context.Context, token_id uuid.UUID) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "apiGatewayRepo.DeleteToken")
	defer span.Finish()

	err := repo.redis.Del(ctxWithTrace, token_id.String()).Err()
	if err != nil {
		return ErrorDeleteToken
	}

	return nil
}

func (repo *apiGatewayRepo) AddTokenFirstAuth(ctx context.Context, token *models.TokenFirstAuth) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "apiGatewayRepo.AddTokenFirstAuth")
	defer span.Finish()

	err := repo.redis.Set(ctxWithTrace, token.TokenName, token.UserId.String(), token.Live_time).Err()
	if err != nil {
		return ErrorAddToken
	}

	return nil
}

func (repo *apiGatewayRepo) GetTokenFirstAuth(ctx context.Context, tokenName string) (*models.TokenFirstAuth, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "apiGatewayRepo.GetTokenFirstAuth")
	defer span.Finish()

	data_str, err := repo.redis.Get(ctxWithTrace, tokenName).Result()
	if err != nil {
		return nil, ErrorGetTokenValue
	}

	data, err := uuid.Parse(data_str)
	if err != nil {
		return nil, err
	}

	cmd := repo.redis.TTL(ctxWithTrace, tokenName)
	if cmd.Err() != nil {
		return nil, ErrorGetTokenExpire
	}

	time_live := cmd.Val()

	return &models.TokenFirstAuth{
		TokenName: tokenName,
		UserId:    data,
		Live_time: time_live,
	}, nil
}

func (repo *apiGatewayRepo) DeleteTokenFirstAuth(ctx context.Context, tokenName string) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "apiGatewayRepo.DeleteTokenFirstAuth")
	defer span.Finish()

	err := repo.redis.Del(ctxWithTrace, tokenName).Err()
	if err != nil {
		return ErrorDeleteToken
	}

	return nil
}

func NewApiGatewayRepository(db *redis.Client) api_gateway.Repository {
	return &apiGatewayRepo{redis: db}
}
