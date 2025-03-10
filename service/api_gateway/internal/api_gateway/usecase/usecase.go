package usecase

import (
	"bytes"
	"context"
	platformConfig "github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway/usecase/html"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/models"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"text/template"
	"time"
)

type apiGateWayUseCase struct {
	cfg  *platformConfig.Config
	repo api_gateway.Repository
}

func (uc *apiGateWayUseCase) CreateToken(ctx context.Context, token_id uuid.UUID, live_time time.Duration) (*models.Token, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "apiGateWayUseCase.CreateToken")
	defer span.Finish()

	token, err := uc.repo.GetToken(ctxWithTrace, token_id)
	if err == nil && token != nil {
		return nil, ErrorTokenAlreadyExist
	}

	token = &models.Token{
		ID:          token_id,
		Live_time:   live_time,
		Date_expire: time.Now().Add(live_time),
	}

	err = uc.repo.AddToken(ctxWithTrace, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (uc *apiGateWayUseCase) CheckExistingToken(ctx context.Context, token_id uuid.UUID) (bool, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "apiGateWayUseCase.CheckExistingToken")
	defer span.Finish()

	_, err := uc.repo.GetToken(ctxWithTrace, token_id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (uc *apiGateWayUseCase) UpdateToken(ctx context.Context, token_id uuid.UUID, new_expire_time time.Duration) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "apiGateWayUseCase.UpdateToken")
	defer span.Finish()

	err := uc.repo.UpdateToken(ctxWithTrace, token_id, new_expire_time)
	if err != nil {
		return err
	}

	return nil
}

func (uc *apiGateWayUseCase) DeleteToken(ctx context.Context, token_id uuid.UUID) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "apiGateWayUseCase.DeleteToken")
	defer span.Finish()

	err := uc.repo.DeleteToken(ctxWithTrace, token_id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *apiGateWayUseCase) CreateSignInPage() (string, error) {

	template_request_sign_in, err := template.New("SignInRequest").Parse(html.RequestSignIn)
	if err != nil {
		return "", nil
	}

	var writer bytes.Buffer

	request_data := &models.RequestData{
		Port: uc.cfg.HTTPServer.Port[1:],
	}

	err = template_request_sign_in.Execute(&writer, request_data)
	if err != nil {
		return "", err
	}
	request_sign_in := writer.String()
	writer.Reset()

	template_request_sign_up_page, err := template.New("SignUpPageRequest").Parse(html.RequestSignUpPage)
	if err != nil {
		return "", err
	}

	err = template_request_sign_up_page.Execute(&writer, request_data)
	if err != nil {
		return "", err
	}
	request_sign_up_page := writer.String()
	writer.Reset()

	template_page, err := template.New("SignInPage").Parse(html.SignInPage)
	if err != nil {
		return "", err
	}

	requests := models.HtmlSignInRequests{
		SignInRequest:     request_sign_in,
		SignUpPageRequest: request_sign_up_page,
	}

	err = template_page.Execute(&writer, &requests)
	if err != nil {
		return "", err
	}

	result := writer.String()

	return result, nil

}

func (uc *apiGateWayUseCase) CreateErrorPage(error string) (string, error) {

	template_request_sign_in_page, err := template.New("SignInPageRequest").Parse(html.RequestSignInPage)
	if err != nil {
		return "", err
	}

	var writer bytes.Buffer

	request_data := &models.RequestData{
		Port: uc.cfg.HTTPServer.Port[1:],
	}

	err = template_request_sign_in_page.Execute(&writer, &request_data)
	if err != nil {
		return "", err
	}

	request_sign_in_page := writer.String()
	writer.Reset()

	template_page, err := template.New("ErrorPage").Parse(html.ErrorPage)
	if err != nil {
		return "", err
	}

	requests := models.HtmlErrorPage{
		Message:           error,
		SignInPageRequest: request_sign_in_page,
	}

	err = template_page.Execute(&writer, &requests)
	if err != nil {
		return "", err
	}

	result := writer.String()

	return result, nil

}

func (uc *apiGateWayUseCase) CreateSignUpPage() (string, error) {

	template_request_sign_up, err := template.New("SignUpRequest").Parse(html.RequestSignUp)
	if err != nil {
		return "", err
	}

	var writer bytes.Buffer

	request_data := &models.RequestData{
		Port: uc.cfg.HTTPServer.Port[1:],
	}

	err = template_request_sign_up.Execute(&writer, &request_data)
	if err != nil {
		return "", err
	}

	request_sign_up := writer.String()
	writer.Reset()

	template_page, err := template.New("SignUpPage").Parse(html.SignUpPage)
	if err != nil {
		return "", err
	}

	requests := models.HtmlSignUpPage{
		SignUpRequest: request_sign_up,
	}

	err = template_page.Execute(&writer, &requests)
	if err != nil {
		return "", err
	}

	result := writer.String()

	return result, nil

}

func NewApiGatewayUseCase(cfg *platformConfig.Config, repo api_gateway.Repository) api_gateway.UseCase {
	return &apiGateWayUseCase{cfg: cfg, repo: repo}
}
