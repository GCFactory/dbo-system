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
	cfg                    *config.Config
	repo                   api_gateway.Repository
	registrationServerInfo *models.RegistrationServerInfo
}

func (uc *apiGateWayUseCase) CreateToken(ctx context.Context, token_id uuid.UUID, live_time time.Duration, token_value uuid.UUID) (*models.Token, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "apiGateWayUseCase.CreateToken")
	defer span.Finish()

	token, err := uc.repo.GetToken(ctxWithTrace, token_id)
	if err == nil && token != nil {
		return nil, ErrorTokenAlreadyExist
	}

	token = &models.Token{
		ID:        token_id,
		Live_time: live_time,
		Data:      token_value,
	}

	err = uc.repo.AddToken(ctxWithTrace, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (uc *apiGateWayUseCase) GetTokenValue(ctx context.Context, token_id uuid.UUID) (uuid.UUID, error) {

	token, err := uc.repo.GetToken(ctx, token_id)
	if err != nil {
		return uuid.Nil, err
	}

	return token.Data, nil

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

func (uc *apiGateWayUseCase) CreateUserPage(user_id uuid.UUID) (string, error) {

	user_data, err := uc.GetUserDataRequest(user_id)
	if err != nil {
		return "", err
	}

	user_data = user_data

	curr_server_data := &models.RequestData{
		Port: uc.cfg.HTTPServer.Port[1:],
	}

	template_sign_out_request, err := template.New("SignOutRequest").Parse(html.RequestSignOut)
	if err != nil {
		return "", err
	}

	var writer bytes.Buffer

	err = template_sign_out_request.Execute(&writer, &curr_server_data)
	if err != nil {
		return "", err
	}

	sign_out_request := writer.String()
	writer.Reset()

	user_page_info := &models.HomePage{
		UserId:              user_id.String(),
		Login:               user_data.Login,
		SignOutRequest:      sign_out_request,
		Surname:             user_data.Surname,
		Name:                user_data.Name,
		Patronymic:          user_data.Patronymic,
		INN:                 user_data.Inn,
		PassportCode:        user_data.PassportSeries + " " + user_data.PassportNumber,
		BirthDate:           user_data.BirthDate,
		BirthLocation:       user_data.BirthLocation,
		PickUpPoint:         user_data.PassportPickUpPoint,
		Authority:           user_data.PassportAuthority,
		AuthorityDate:       user_data.PassportAuthorityDate,
		RegistrationAddress: user_data.PassportRegistrationAddress,
		ListOfAccounts:      "",
	}

	accounts := ""
	for _, account_id := range user_data.Accounts {
		account_data, err := uc.GetAccountDataRequest(user_id, account_id)
		if err != nil {
			return "", err
		}

		template_account_raw, err := template.New("HomePageAccount").Parse(html.HomePageAccount)
		if err != nil {
			return "", err
		}

		account_html_data := &models.HomePageAccountDescription{
			Name:                account_data.Name,
			Status:              account_data.Status,
			Cache:               fmt.Sprint(account_data.Cache),
			GetCreditsRequest:   "",
			AddCacheRequest:     "",
			ReduceCacheRequest:  "",
			CloseAccountRequest: "",
		}

		err = template_account_raw.Execute(&writer, &account_html_data)
		if err != nil {
			return "", err
		}

		account_raw := writer.String()
		writer.Reset()

		accounts += account_raw + "\n"

		//	TODO: доделать запросы
	}

	user_page_info.ListOfAccounts = accounts

	template_open_account_page_request, err := template.New("RequestOpenAccountPage").Parse(html.RequestOpenAccountPage)
	if err != nil {
		return "", err
	}
	err = template_open_account_page_request.Execute(&writer, &curr_server_data)
	if err != nil {
		return "", err
	}
	open_account_page_request := writer.String()
	writer.Reset()

	user_page_info.CreateAccountRequest = open_account_page_request

	template_user_page, err := template.New("UserPage").Parse(html.HomePage)
	if err != nil {
		return "", nil
	}

	err = template_user_page.Execute(&writer, &user_page_info)
	if err != nil {
		return "", err
	}

	user_page := writer.String()

	return user_page, nil

}

func (uc *apiGateWayUseCase) SignIn(login_info *models.SignInInfo) (string, *models.Token, error) {

	user_data, err := uc.GetUserDataByLoginRequest(login_info.Login)
	if err != nil {
		return "", nil, err
	}

	is_ok, err := uc.CheckUserPasswordRequest(user_data.Id, login_info.Password)
	if err != nil {
		return "", nil, err
	}

	if is_ok {
		home_page, err := uc.CreateUserPage(user_data.Id)
		if err != nil {
			return "", nil, err
		}

		token, err := uc.CreateToken(context.Background(), uuid.New(), time.Minute, user_data.Id)
		if err != nil {
			return "", nil, err
		}

		return home_page, token, nil
	}

	return "", nil, ErrorWrongPassword

}

func (uc *apiGateWayUseCase) SignUp(sign_up_info *models.SignUpInfo) (string, *models.Token, error) {

	user_id, err := uc.CreateUserRequest(sign_up_info)
	if err != nil {
		return "", nil, err
	}

	home_page, err := uc.CreateUserPage(user_id)
	if err != nil {
		return "", nil, err
	}

	token, err := uc.CreateToken(context.Background(), uuid.New(), time.Minute, user_id)
	if err != nil {
		return "", nil, err
	}

	return home_page, token, nil
}

func (uc *apiGateWayUseCase) GetAccountDataRequest(user_id uuid.UUID, account_id uuid.UUID) (*models.AccountInfo, error) {

	template_request_get_account_data, err := template.New("GetAccountData").Parse(GetAccountData)
	if err != nil {
		return nil, err
	}

	var writer bytes.Buffer

	err = template_request_get_account_data.Execute(&writer, uc.registrationServerInfo)
	if err != nil {
		return nil, err
	}
	request_get_account_data := writer.String()
	writer.Reset()

	request_get_account_data_body := &models.GetAccountDataBody{
		UserId:    user_id,
		AccountId: account_id,
	}

	request_body, err := json.Marshal(&request_get_account_data_body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, request_get_account_data, bytes.NewBuffer(request_body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.registrationServerInfo.TimeWaitResponse,
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

	if operation_data.AdditionalInfo == nil {
		return nil, ErrorNoAccountData
	}

	additional_data := operation_data.AdditionalInfo.(map[string]interface{})

	additional_data = additional_data

	result := &models.AccountInfo{
		Id: account_id,
	}

	if cache_amount, ok := additional_data["acc_cache"]; ok {
		result.Cache = cache_amount.(float64)
	}
	if name, ok := additional_data["acc_name"]; ok {
		result.Name = name.(string)
	}
	if status, ok := additional_data["acc_status"]; ok {
		status_val := status.(float64)
		status_str := "Unknown"
		switch status_val {
		case 10:
			{
				status_str = "Reserved"
			}
		case 20:
			{
				status_str = "Created"
			}
		case 30:
			{
				status_str = "Opened"
			}
		case 40:
			{
				status_str = "Closed"
			}
		case 50:
			{
				status_str = "Blocked"
			}
		}
		result.Status = status_str
	}

	return result, nil
}

func (uc *apiGateWayUseCase) GetUserDataRequest(user_id uuid.UUID) (*models.UserInfo, error) {

	template_request_get_user_data, err := template.New("GetUserData").Parse(GetUserData)
	if err != nil {
		return nil, err
	}

	var writer bytes.Buffer

	err = template_request_get_user_data.Execute(&writer, uc.registrationServerInfo)
	if err != nil {
		return nil, err
	}
	request_check_user_password := writer.String()
	writer.Reset()

	request_get_user_data_body := &models.CheckUserPasswordBody{
		UserId: user_id,
	}

	request_body, err := json.Marshal(&request_get_user_data_body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, request_check_user_password, bytes.NewBuffer(request_body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.registrationServerInfo.TimeWaitResponse,
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

	if operation_data.AdditionalInfo == nil {
		return nil, ErrorNoUserData
	}

	additional_data := operation_data.AdditionalInfo.(map[string]interface{})

	result := &models.UserInfo{}

	if name, ok := additional_data["passport_first_name"]; ok && name != nil {
		result.Name = name.(string)
	}
	if surname, ok := additional_data["passport_first_surname"]; ok && surname != nil {
		result.Surname = surname.(string)
	}
	if patronimic, ok := additional_data["passport_first_patronimic"]; ok && patronimic != nil {
		result.Patronymic = patronimic.(string)
	}
	if passport_series, ok := additional_data["passport_series"]; ok && passport_series != nil {
		result.PassportSeries = passport_series.(string)
	}
	if passport_number, ok := additional_data["passport_number"]; ok && passport_number != nil {
		result.PassportNumber = passport_number.(string)
	}
	if user_id, ok := additional_data["user_id"]; ok && user_id != nil {
		user_id_str := user_id.(string)
		id, err := uuid.Parse(user_id_str)
		if err != nil {
			return nil, err
		}
		result.Id = id
	}
	if birth_date, ok := additional_data["passport_birth_date"]; ok && birth_date != nil {
		result.BirthDate = birth_date.(string)
	}
	if birth_location, ok := additional_data["passport_birth_location"]; ok && birth_location != nil {
		result.BirthLocation = birth_location.(string)
	}
	if passport_authority, ok := additional_data["passport_authority"]; ok && passport_authority != nil {
		result.PassportAuthority = passport_authority.(string)
	}
	if authority_date, ok := additional_data["passport_authority_date"]; ok && authority_date != nil {
		result.PassportAuthorityDate = authority_date.(string)
	}
	if passport_pick_up_point, ok := additional_data["passport_pick_up_point"]; ok && passport_pick_up_point != nil {
		result.PassportPickUpPoint = passport_pick_up_point.(string)
	}
	if registration_address, ok := additional_data["passport_registration_address"]; ok && registration_address != nil {
		result.PassportRegistrationAddress = registration_address.(string)
	}
	if inn, ok := additional_data["inn"]; ok && inn != nil {
		result.Inn = inn.(string)
	}
	if login, ok := additional_data["user_login"]; ok && login != nil {
		result.Login = login.(string)
	}
	if accounts, ok := additional_data["accounts"]; ok && accounts != nil {
		accounts_list := accounts.([]interface{})
		result.Accounts = make([]uuid.UUID, 0)
		for _, account_id_str := range accounts_list {
			account_id, err := uuid.Parse(account_id_str.(string))
			if err != nil {
				return nil, err
			}
			result.Accounts = append(result.Accounts, account_id)
		}
	}

	return result, nil
}

func (uc *apiGateWayUseCase) CheckUserPasswordRequest(user_id uuid.UUID, password string) (bool, error) {

	template_request_check_user_password, err := template.New("RequestCheckUserPassword").Parse(RequestCheckUserPassword)
	if err != nil {
		return false, err
	}

	var writer bytes.Buffer

	err = template_request_check_user_password.Execute(&writer, &uc.registrationServerInfo)
	if err != nil {
		return false, err
	}
	request_check_user_password := writer.String()
	writer.Reset()

	request_check_user_password_body := &models.CheckUserPasswordBody{
		Password: password,
		UserId:   user_id,
	}

	request_body, err := json.Marshal(&request_check_user_password_body)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest(http.MethodPost, request_check_user_password, bytes.NewBuffer(request_body))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.registrationServerInfo.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var resp_data = &models.OperationResponse{}

	err = json.Unmarshal(resp_body, &resp_data)
	if err != nil {
		return false, err
	}

	operation_id_str := resp_data.Info

	operation_id, err := uuid.Parse(operation_id_str)
	if err != nil {
		return false, err
	}

	operation_data, err := uc.GetOperationData(operation_id)
	if err != nil {
		return false, err
	}

	if operation_data.Info != "Success" {
		return false, nil
	}

	return true, nil
}

func (uc *apiGateWayUseCase) GetUserDataByLoginRequest(login string) (*models.UserInfo, error) {

	template_request_get_user_data_by_login, err := template.New("RequestGetUserDataByLogin").Parse(RequestGetUserDataByLogin)
	if err != nil {
		return nil, err
	}

	var writer bytes.Buffer

	err = template_request_get_user_data_by_login.Execute(&writer, &uc.registrationServerInfo)
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
		Timeout: uc.registrationServerInfo.TimeWaitResponse,
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

	template_request_create_user, err := template.New("ReuqestCreateUser").Parse(RequestCreateUser)
	if err != nil {
		return uuid.Nil, err
	}

	var writer bytes.Buffer

	err = template_request_create_user.Execute(&writer, &uc.registrationServerInfo)
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
		Timeout: uc.registrationServerInfo.TimeWaitResponse,
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

	for i := 0; i < uc.registrationServerInfo.NumRetry; i++ {
		operation_data, err := uc.GetOperationDataRequest(operation_id)
		if err != nil {
			return nil, err
		}
		if operation_data.Info == "In progress" {
			time.Sleep(uc.registrationServerInfo.WaitTimeRetry)
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

	template_request_get_operation_status, err := template.New("RequestGetOperationResult").Parse(RequestGetOperationResult)
	if err != nil {
		return nil, err
	}

	var writer bytes.Buffer

	err = template_request_get_operation_status.Execute(&writer, &uc.registrationServerInfo)
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
		Timeout: uc.registrationServerInfo.TimeWaitResponse,
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

func NewApiGatewayUseCase(cfg *config.Config, repo api_gateway.Repository, registration_server_info *models.RegistrationServerInfo) api_gateway.UseCase {
	return &apiGateWayUseCase{cfg: cfg, repo: repo, registrationServerInfo: registration_server_info}
}
