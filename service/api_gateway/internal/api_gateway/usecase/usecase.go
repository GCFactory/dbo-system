package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	//"errors"
	"fmt"
	"github.com/GCFactory/dbo-system/service/api_gateway/config"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway/usecase/html"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/models"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"io"
	"net/http"
	"text/template"
	"time"
)

type apiGateWayUseCase struct {
	cfg  *config.Config
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

func (uc *apiGateWayUseCase) CreateUserPage(user_info *models.UserInfo) (string, error) {

	return "", nil

}

func (uc *apiGateWayUseCase) SignIn(login_info *models.SignInInfo) (bool, error) {

	_, err := uc.GetUserDataByLoginRequest(login_info.Login)
	if err != nil {
		return false, err
	}

	// TODO: тут сделать проверку пароля

	return false, nil

}

func (uc *apiGateWayUseCase) SignUp(sign_up_info *models.SignUpInfo) (uuid.UUID, error) {

	user_id, err := uc.CreateUserRequest(sign_up_info)
	if err != nil {
		return uuid.Nil, err
	}

	return user_id, nil
}

func (uc *apiGateWayUseCase) GetUserDataByLoginRequest(login string) (*models.UserInfo, error) {

	registration_server_info := &models.RegistrationServerInfo{
		Host:             uc.cfg.Registration.Host,
		Port:             uc.cfg.Registration.Port,
		NumRetry:         uc.cfg.Registration.Retry,
		WaitTimeRetry:    time.Duration(time.Second.Nanoseconds() * int64(uc.cfg.Registration.TimeWaitRetry)),
		TimeWaitResponse: time.Duration(time.Second.Nanoseconds() * int64(uc.cfg.Registration.TimeWaitResponse)),
	}

	template_request_get_user_data_by_login, err := template.New("RequestGetUserDataByLogin").Parse(RequestGetUserDataByLogin)
	if err != nil {
		return nil, err
	}

	var writer bytes.Buffer

	err = template_request_get_user_data_by_login.Execute(&writer, &registration_server_info)
	if err != nil {
		return nil, err
	}
	request_get_user_data_by_login := writer.String()
	writer.Reset()

	request_get_user_data_by_login_body := &models.GetUserDataByLoginBody{
		UserLogin: login,
	}

	request_body, err := json.Marshal(&request_get_user_data_by_login_body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, request_get_user_data_by_login, bytes.NewBuffer(request_body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: registration_server_info.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var resp_data = &models.OperationResponse{}

	err = json.Unmarshal(resp_body, &resp_data)
	if err != nil {
		return nil, err
	}

	operation_id_str := resp_data.Info

	operation_id, err := uuid.Parse(operation_id_str)
	if err != nil {
		return nil, err
	}

	operation_data, err := uc.GetOperationData(operation_id)
	if err != nil {
		return nil, err
	}

	additional_info := operation_data.AdditionalInfo.(map[string]interface{})

	user_id, err := uuid.Parse(additional_info["user_id"].(string))
	if err != nil {
		return nil, err
	}

	result := &models.UserInfo{
		Id: user_id,
	}

	return result, nil

}

func (uc *apiGateWayUseCase) CreateUserRequest(user_info *models.SignUpInfo) (uuid.UUID, error) {

	registration_server_info := &models.RegistrationServerInfo{
		Host:             uc.cfg.Registration.Host,
		Port:             uc.cfg.Registration.Port,
		NumRetry:         uc.cfg.Registration.Retry,
		WaitTimeRetry:    time.Duration(time.Second.Nanoseconds() * int64(uc.cfg.Registration.TimeWaitRetry)),
		TimeWaitResponse: time.Duration(time.Second.Nanoseconds() * int64(uc.cfg.Registration.TimeWaitResponse)),
	}

	template_request_create_user, err := template.New("ReuqestCreateUser").Parse(RequestCreateUser)
	if err != nil {
		return uuid.Nil, err
	}

	var writer bytes.Buffer

	err = template_request_create_user.Execute(&writer, &registration_server_info)
	if err != nil {
		return uuid.Nil, err
	}
	request_create_user := writer.String()
	writer.Reset()

	request_get_user_data_by_login_body := &models.CreateUserBody{
		UserInn: user_info.Inn,
		UserData: &models.CreateUserBodyUserData{
			Login:    user_info.Login,
			Password: user_info.Password,
		},
		Passport: &models.CreateUserBodyPassport{
			Name:                user_info.Name,
			Surname:             user_info.Surname,
			Patronymic:          user_info.Patronymic,
			Series:              user_info.PassportSeries,
			Number:              user_info.PassportNumber,
			Authority:           user_info.PassportAuthority,
			BirthLocation:       user_info.BirthLocation,
			PickUpPoint:         user_info.PassportPickUpPoint,
			RegistrationAddress: user_info.PassportRegistrationAddress,
			AuthorityDate:       user_info.PassportAuthorityDate,
			BirthDate:           user_info.BirthDate,
		},
	}

	request_body, err := json.Marshal(&request_get_user_data_by_login_body)
	if err != nil {
		return uuid.Nil, err
	}

	req, err := http.NewRequest(http.MethodPost, request_create_user, bytes.NewBuffer(request_body))
	if err != nil {
		return uuid.Nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: registration_server_info.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return uuid.Nil, err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return uuid.Nil, err
	}

	var resp_data = &models.OperationResponse{}

	err = json.Unmarshal(resp_body, &resp_data)
	if err != nil {
		return uuid.Nil, err
	}

	operation_id_str := resp_data.Info

	operation_id, err := uuid.Parse(operation_id_str)
	if err != nil {
		return uuid.Nil, err
	}

	operation_data, err := uc.GetOperationData(operation_id)
	if err != nil {
		return uuid.Nil, err
	}

	additional_info := operation_data.AdditionalInfo.(map[string]interface{})

	user_id, err := uuid.Parse(additional_info["user_uuid"].(string))
	if err != nil {
		return uuid.Nil, err
	}

	return user_id, nil
}

func (uc *apiGateWayUseCase) GetOperationData(operation_id uuid.UUID) (*models.OperationResponse, error) {

	registration_server_info := &models.RegistrationServerInfo{
		Host:             uc.cfg.Registration.Host,
		Port:             uc.cfg.Registration.Port,
		NumRetry:         uc.cfg.Registration.Retry,
		WaitTimeRetry:    time.Duration(time.Second.Nanoseconds() * int64(uc.cfg.Registration.TimeWaitRetry)),
		TimeWaitResponse: time.Duration(time.Second.Nanoseconds() * int64(uc.cfg.Registration.TimeWaitResponse)),
	}

	for i := 0; i < registration_server_info.NumRetry; i++ {
		operation_data, err := uc.GetOperationDataRequest(operation_id)
		if err != nil {
			return nil, err
		}
		if operation_data.Info == "In progress" {
			time.Sleep(registration_server_info.WaitTimeRetry)
			continue
		} else if operation_data.Info == "Success" {
			return operation_data, nil
		} else if operation_data.Info == "Failed" {

			additional_info := operation_data.AdditionalInfo.(map[string]interface{})

			error_string := ""

			errors_list := additional_info["errors"].(map[string]interface{})
			for _, saga := range errors_list {
				for _, event := range saga.(map[string]interface{}) {
					error_msg, ok := event.(map[string]interface{})["info"]
					if ok {
						error_string += error_msg.(string) + "\n"
					}
				}
			}

			err := errors.New(error_string)
			return operation_data, err
		}
	}

	return nil, ErrorOperationProcessedYet
}

func (uc *apiGateWayUseCase) GetOperationDataRequest(operation_id uuid.UUID) (*models.OperationResponse, error) {

	registration_server_info := &models.RegistrationServerInfo{
		Host:             uc.cfg.Registration.Host,
		Port:             uc.cfg.Registration.Port,
		NumRetry:         uc.cfg.Registration.Retry,
		WaitTimeRetry:    time.Duration(time.Second.Nanoseconds() * int64(uc.cfg.Registration.TimeWaitRetry)),
		TimeWaitResponse: time.Duration(time.Second.Nanoseconds() * int64(uc.cfg.Registration.TimeWaitResponse)),
	}

	template_request_get_operation_status, err := template.New("RequestGetOperationResult").Parse(RequestGetOperationResult)
	if err != nil {
		return nil, err
	}

	var writer bytes.Buffer

	err = template_request_get_operation_status.Execute(&writer, registration_server_info)
	if err != nil {
		return nil, err
	}

	request_get_operation_status := writer.String()

	request_body := &models.GetOperationResultBody{
		OperationId: operation_id,
	}

	request_body_bytes, err := json.Marshal(&request_body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, request_get_operation_status, bytes.NewBuffer(request_body_bytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: registration_server_info.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(resp_body))

	var resp_data = &models.OperationResponse{}

	err = json.Unmarshal(resp_body, &resp_data)
	if err != nil {
		return nil, err
	}

	return resp_data, nil
}

func NewApiGatewayUseCase(cfg *config.Config, repo api_gateway.Repository) api_gateway.UseCase {
	return &apiGateWayUseCase{cfg: cfg, repo: repo}
}
