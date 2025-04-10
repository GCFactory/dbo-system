package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway/usecase/notifications"
	"github.com/rabbitmq/amqp091-go"
	"net/url"
	"strconv"
	"strings"

	//"errors"
	"fmt"
	"github.com/GCFactory/dbo-system/service/api_gateway/config"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/api_gateway/usecase/html"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/models"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"io"
	"net/http"
	"text/template"
	"time"
)

type apiGateWayUseCase struct {
	cfg                    *config.Config
	repo                   api_gateway.Repository
	registrationServerInfo *models.InternalServerInfo
	usersServerInfo        *models.InternalServerInfo
	accountsServerInfo     *models.InternalServerInfo
	notificationServerInfo *models.InternalServerInfo
	totpServerInfo         *models.InternalServerInfo
	graphImagesPath        string
	qrImagesPath           string
	rmqChan                *amqp091.Channel
	rmqQueue               amqp091.Queue
}

var TokenLiveTime = time.Minute
var TokenFirstAuthLiveTime = time.Minute * 5

func (uc *apiGateWayUseCase) AddTokenFirstAuth(ctx context.Context, token *models.TokenFirstAuth) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "apiGateWayUseCase.AddTokenFirstAuth")
	defer span.Finish()

	err := uc.repo.AddTokenFirstAuth(ctxWithTrace, token)
	if err != nil {
		return err
	}

	return nil

}

func (uc *apiGateWayUseCase) GetTokenFirstAuth(ctx context.Context, tokenName string) (*models.TokenFirstAuth, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "apiGateWayUseCase.GetTokenFirstAuth")
	defer span.Finish()

	token, err := uc.repo.GetTokenFirstAuth(ctxWithTrace, tokenName)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (uc *apiGateWayUseCase) DeleteTokenFirstAuth(ctx context.Context, tokenName string) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "apiGateWayUseCase.DeleteTokenFirstAuth")
	defer span.Finish()

	err := uc.repo.DeleteTokenFirstAuth(ctxWithTrace, tokenName)
	if err != nil {
		return err
	}

	return nil

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

	template_home_page_request, err := template.New("RequestUserPage").Parse(html.RequestUserPage)
	if err != nil {
		return "", nil
	}

	err = template_home_page_request.Execute(&writer, &request_data)
	if err != nil {
		return "", err
	}

	request_home_page := writer.String()
	writer.Reset()

	requests := models.HtmlSignInRequests{
		SignInRequest:     request_sign_in,
		SignUpPageRequest: request_sign_up_page,
		HomePageRequest:   request_home_page,
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

	template_home_page_request, err := template.New("RequestUserPage").Parse(html.RequestUserPage)
	if err != nil {
		return "", nil
	}

	err = template_home_page_request.Execute(&writer, &request_data)
	if err != nil {
		return "", err
	}

	request_home_page := writer.String()
	writer.Reset()

	requests := models.HtmlSignUpPage{
		SignUpRequest:   request_sign_up,
		HomePageRequest: request_home_page,
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

	template_sing_in_page_request, err := template.New("RequestSignInPage").Parse(html.RequestSignInPage)
	if err != nil {
		return "", err
	}

	err = template_sing_in_page_request.Execute(&writer, &curr_server_data)
	if err != nil {
		return "", err
	}

	sign_in_page_request := writer.String()
	writer.Reset()

	template_turn_on_totp_page_request, err := template.New("RequestTurnOnTotpPage").Parse(html.RequestTurnOnTotpPage)
	if err != nil {
		return "", err
	}

	err = template_turn_on_totp_page_request.Execute(&writer, &curr_server_data)
	if err != nil {
		return "", err
	}

	turn_on_totp_page_request := writer.String()
	writer.Reset()

	template_turn_off_totp_page_request, err := template.New("RequestTurnOffTotpPage").Parse(html.RequestTurnOffTotpPage)
	if err != nil {
		return "", err
	}

	err = template_turn_off_totp_page_request.Execute(&writer, &curr_server_data)
	if err != nil {
		return "", err
	}

	turn_off_totp_page_request := writer.String()
	writer.Reset()

	authorityDate := strings.Split(user_data.PassportAuthorityDate, "T")[0]
	birthDate := strings.Split(user_data.BirthDate, "T")[0]

	user_page_info := &models.HomePage{
		UserId:              user_id.String(),
		Login:               user_data.Login,
		SignInPageRequest:   sign_in_page_request,
		SignOutRequest:      sign_out_request,
		Surname:             user_data.Surname,
		Name:                user_data.Name,
		Patronymic:          user_data.Patronymic,
		INN:                 user_data.Inn,
		PassportCode:        user_data.PassportSeries + " " + user_data.PassportNumber,
		BirthDate:           birthDate,
		BirthLocation:       user_data.BirthLocation,
		PickUpPoint:         user_data.PassportPickUpPoint,
		Authority:           user_data.PassportAuthority,
		AuthorityDate:       authorityDate,
		RegistrationAddress: user_data.PassportRegistrationAddress,
		Email:               user_data.Email,
		ListOfAccounts:      "",
		IsUseTotp:           user_data.UsingTotp,
		RequestTurnOnTotp:   turn_on_totp_page_request,
		RequestTurnOffTotp:  turn_off_totp_page_request,
	}

	accounts := ""
	for _, account_id := range user_data.Accounts {
		account_data, err := uc.GetAccountDataRequest(account_id)
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
			AccountId:           account_id.String(),
			GetCreditsRequest:   "",
			AddCacheRequest:     "",
			ReduceCacheRequest:  "",
			CloseAccountRequest: "",
			Disabled:            false,
		}

		if account_html_data.Status != "Opened" {
			account_html_data.Disabled = true
		}

		template_get_account_credits_request, err := template.New("RequestAccountCreditsPage").Parse(html.RequestAccountCreditsPage)
		if err != nil {
			return "", nil
		}

		err = template_get_account_credits_request.Execute(&writer, &curr_server_data)
		if err != nil {
			return "", nil
		}

		account_html_data.GetCreditsRequest = writer.String()
		writer.Reset()

		templage_close_account_request, err := template.New("RequestAccountClosePage").Parse(html.RequestAccountClosePage)
		if err != nil {
			return "", nil
		}

		err = templage_close_account_request.Execute(&writer, &curr_server_data)
		if err != nil {
			return "", nil
		}

		account_html_data.CloseAccountRequest = writer.String()
		writer.Reset()

		templage_add_account_cache_request, err := template.New("RequestAddAccountCachePage").Parse(html.RequestAddAccountCachePage)
		if err != nil {
			return "", nil
		}

		err = templage_add_account_cache_request.Execute(&writer, &curr_server_data)
		if err != nil {
			return "", nil
		}

		account_html_data.AddCacheRequest = writer.String()
		writer.Reset()

		templage_width_account_cache_request, err := template.New("RequestWidthAccountCachePage").Parse(html.RequestWidthAccountCachePage)
		if err != nil {
			return "", nil
		}

		err = templage_width_account_cache_request.Execute(&writer, &curr_server_data)
		if err != nil {
			return "", nil
		}

		account_html_data.ReduceCacheRequest = writer.String()
		writer.Reset()

		err = template_account_raw.Execute(&writer, &account_html_data)
		if err != nil {
			return "", err
		}

		account_raw := writer.String()
		writer.Reset()

		accounts += account_raw + "\n"
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

func (uc *apiGateWayUseCase) CreateOpenAccountPage(user_id uuid.UUID) (string, error) {

	user_data, err := uc.GetUserDataRequest(user_id)
	if err != nil {
		return "", err
	}

	additional_data := make(map[string]interface{})
	additional_data["login"] = user_data.Login

	result, err := uc.CreateOperationPage(AccountOperationTypeOpen, additional_data)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (uc *apiGateWayUseCase) CreateAccountCreditsPage(user_id uuid.UUID, account_id uuid.UUID) (string, error) {
	user_data, err := uc.GetUserDataRequest(user_id)
	if err != nil {
		return "", err
	}

	additional_data := make(map[string]interface{})
	additional_data["login"] = user_data.Login
	additional_data["user_id"] = user_id.String()
	additional_data["account_id"] = account_id.String()

	result, err := uc.CreateOperationPage(AccountOperationTypeGetCredits, additional_data)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (uc *apiGateWayUseCase) CreateCloseAccountPage(user_id uuid.UUID, account_id uuid.UUID) (string, error) {
	user_data, err := uc.GetUserDataRequest(user_id)
	if err != nil {
		return "", err
	}

	additional_data := make(map[string]interface{})
	additional_data["login"] = user_data.Login
	additional_data["user_id"] = user_id.String()
	additional_data["account_id"] = account_id.String()

	result, err := uc.CreateOperationPage(AccountOperationTypeClose, additional_data)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (uc *apiGateWayUseCase) CreateAddAccountCachePage(user_id uuid.UUID, account_id uuid.UUID) (string, error) {
	user_data, err := uc.GetUserDataRequest(user_id)
	if err != nil {
		return "", err
	}

	additional_data := make(map[string]interface{})
	additional_data["login"] = user_data.Login
	additional_data["user_id"] = user_id.String()
	additional_data["account_id"] = account_id.String()

	result, err := uc.CreateOperationPage(AccountOperationAddCache, additional_data)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (uc *apiGateWayUseCase) CreateWidthAccountCachePage(user_id uuid.UUID, account_id uuid.UUID) (string, error) {
	user_data, err := uc.GetUserDataRequest(user_id)
	if err != nil {
		return "", err
	}

	additional_data := make(map[string]interface{})
	additional_data["login"] = user_data.Login
	additional_data["user_id"] = user_id.String()
	additional_data["account_id"] = account_id.String()

	result, err := uc.CreateOperationPage(AccountOperationWidthCache, additional_data)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (uc *apiGateWayUseCase) CreateOperationPage(operation_type string, additional_data map[string]interface{}) (string, error) {

	result := ""
	var err error = nil

	if ValidateOperationTypeData(operation_type, additional_data) {

		login, ok := additional_data["login"]
		if !ok {
			return "", ErrorNoOperationData
		}

		var buffer bytes.Buffer

		account_operation_page_info := &models.AccountOperationPage{
			OperationName:  operation_type,
			Operation:      "",
			Login:          login.(string),
			SignOutRequest: "",
			ReturnRequest:  "",
		}

		curr_server_data := &models.RequestData{
			Port: uc.cfg.HTTPServer.Port[1:],
		}

		switch operation_type {
		case AccountOperationTypeOpen:
			{
				template_operation_open_account_request, err := template.New("RequestOpenAccount").Parse(html.RequestOpenAccount)
				if err != nil {
					return "", err
				}

				err = template_operation_open_account_request.Execute(&buffer, &curr_server_data)
				if err != nil {
					return "", err
				}

				operation_data := &models.AccountOperationData{
					OperationRequest: buffer.String(),
				}
				buffer.Reset()

				template_operation_open_account, err := template.New("AccountOperationCreateAccount").Parse(html.AccountOperationCreateAccount)
				if err != nil {
					return "", err
				}

				err = template_operation_open_account.Execute(&buffer, &operation_data)
				if err != nil {
					return "", err
				}

				account_operation_page_info.Operation = buffer.String()
				buffer.Reset()

			}
		case AccountOperationTypeGetCredits,
			AccountOperationTypeClose,
			AccountOperationAddCache,
			AccountOperationWidthCache:
			{

				acc_id_str, ok := additional_data["account_id"]
				if !ok {
					return "", ErrorNoOperationData
				}

				account_id, err := uuid.Parse(acc_id_str.(string))
				if err != nil {
					return "", err
				}

				acc_data, err := uc.GetAccountDataRequest(account_id)
				if err != nil {
					return "", err
				}

				if operation_type == AccountOperationTypeGetCredits {
					account_credits := &models.AccountOperationCreditsData{
						Name:       acc_data.Name,
						Status:     acc_data.Status,
						BIC:        acc_data.BIC,
						CIO:        acc_data.CIO,
						Amount:     acc_data.Cache,
						CorrNumber: acc_data.CorrNumber,
						CulcNumber: acc_data.CulcNumber,
					}

					template_operation_account_credits, err := template.New("AccountOperationGetCredits").Parse(html.AccountOperationGetCredits)
					if err != nil {
						return "", err
					}

					err = template_operation_account_credits.Execute(&buffer, &account_credits)
					if err != nil {
						return "", err
					}

					account_operation_page_info.Operation = buffer.String()
					buffer.Reset()

				} else if operation_type == AccountOperationTypeClose {

					operation_info := &models.AccountOperationData{
						OperationRequest: "",
						AccountId:        account_id.String(),
					}

					template_close_account_request, err := template.New("RequestAccountClose").Parse(html.RequestAccountClose)
					if err != nil {
						return "", err
					}

					err = template_close_account_request.Execute(&buffer, &curr_server_data)
					if err != nil {
						return "", err
					}

					operation_info.OperationRequest = buffer.String()
					buffer.Reset()

					template_operation_close_account, err := template.New("AccountOperationCloseAccount").Parse(html.AccountOperationCloseAccount)
					if err != nil {
						return "", err
					}

					err = template_operation_close_account.Execute(&buffer, &operation_info)
					if err != nil {
						return "", err
					}

					account_operation_page_info.Operation = buffer.String()
					buffer.Reset()

				} else if operation_type == AccountOperationAddCache {

					operation_info := &models.AccountOperationData{
						OperationRequest: "",
						AccountId:        account_id.String(),
					}

					template_account_add_cache_request, err := template.New("RequestAddAccountCache").Parse(html.RequestAddAccountCache)
					if err != nil {
						return "", err
					}

					err = template_account_add_cache_request.Execute(&buffer, &curr_server_data)
					if err != nil {
						return "", err
					}

					operation_info.OperationRequest = buffer.String()
					buffer.Reset()

					template_operation_account_add_cache, err := template.New("AccountOperationAddCache").Parse(html.AccountOperationAddCache)
					if err != nil {
						return "", err
					}

					err = template_operation_account_add_cache.Execute(&buffer, &operation_info)
					if err != nil {
						return "", err
					}

					account_operation_page_info.Operation = buffer.String()
					buffer.Reset()

				} else {

					operation_info := &models.AccountOperationData{
						OperationRequest: "",
						AccountId:        account_id.String(),
					}

					template_account_width_cache_request, err := template.New("RequestWidthAccountCache").Parse(html.RequestWidthAccountCache)
					if err != nil {
						return "", err
					}

					err = template_account_width_cache_request.Execute(&buffer, &curr_server_data)
					if err != nil {
						return "", err
					}

					operation_info.OperationRequest = buffer.String()
					buffer.Reset()

					template_operation_account_width_cache, err := template.New("AccountOperationWidthCache").Parse(html.AccountOperationWidthCache)
					if err != nil {
						return "", err
					}

					err = template_operation_account_width_cache.Execute(&buffer, &operation_info)
					if err != nil {
						return "", err
					}

					account_operation_page_info.Operation = buffer.String()
					buffer.Reset()

				}

			}
		default:
			{
				err = ErrorUnknownAccountOperationType
			}

		}

		if err == nil {

			template_operation_page, err := template.New("AccountOperationPage").Parse(html.AccountOperationPage)
			if err != nil {
				return "", err
			}

			template_sign_out_request, err := template.New("RequestSignOut").Parse(html.RequestSignOut)
			if err != nil {
				return "", err
			}

			err = template_sign_out_request.Execute(&buffer, &curr_server_data)
			if err != nil {
				return "", err
			}

			account_operation_page_info.SignOutRequest = buffer.String()
			buffer.Reset()

			template_request_user_page, err := template.New("RequestUserPage").Parse(html.RequestUserPage)
			if err != nil {
				return "", err
			}

			err = template_request_user_page.Execute(&buffer, &curr_server_data)
			if err != nil {
				return "", err
			}

			account_operation_page_info.ReturnRequest = buffer.String()
			buffer.Reset()

			err = template_operation_page.Execute(&buffer, &account_operation_page_info)
			if err != nil {
				return "", err
			}

			result = buffer.String()

		}
	} else {
		err = ErrorValidationAccountOperationData
	}

	return result, err
}

func (uc *apiGateWayUseCase) CreateTurnOnTotpPage(userId uuid.UUID) (string, error) {

	page, err := uc.createTotpOperationPage(userId, TotpOperationTypeTurnOn)
	if err != nil {
		return "", err
	}

	return page, nil
}

func (uc *apiGateWayUseCase) CreateTurnOffTotpPage(userId uuid.UUID) (string, error) {

	page, err := uc.createTotpOperationPage(userId, TotpOperationTypeTurnOff)
	if err != nil {
		return "", err
	}

	return page, nil

}

func (uc *apiGateWayUseCase) CreateTotpQrPage(userId uuid.UUID) (string, error) {

	page, err := uc.createTotpOperationPage(userId, TotpOperationTypeQr)
	if err != nil {
		return "", err
	}

	return page, nil

}

func (uc *apiGateWayUseCase) CreateTotpCheckPage() (string, error) {

	var buffer bytes.Buffer

	curr_server_data := &models.RequestData{
		Port: uc.cfg.HTTPServer.Port[1:],
	}

	totpCheckPageData := &models.TotpCheckPage{}
	returnRequestTemplate, err := template.New("RequestSignInPage").Parse(html.RequestSignInPage)
	if err != nil {
		return "", err
	}

	err = returnRequestTemplate.Execute(&buffer, &curr_server_data)
	if err != nil {
		return "", err
	}

	totpCheckPageData.ReturnRequest = buffer.String()
	buffer.Reset()

	templateValidateTotpCode, err := template.New("RequestCheckTotp").Parse(html.RequestCheckTotp)
	if err != nil {
		return "", err
	}

	err = templateValidateTotpCode.Execute(&buffer, &curr_server_data)
	if err != nil {
		return "", err
	}
	totpCheckPageData.OperationRequest = buffer.String()
	buffer.Reset()

	templateTotpCheckPage, err := template.New("TotpCheckPage").Parse(html.TotpCheckPage)
	if err != nil {
		return "", err
	}

	err = templateTotpCheckPage.Execute(&buffer, &totpCheckPageData)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (uc *apiGateWayUseCase) createTotpOperationPage(userId uuid.UUID, operation string) (string, error) {

	userData, err := uc.GetUserDataRequest(userId)
	if err != nil {
		return "", err
	}

	totpOperationInfo := &models.TotpOperationPage{
		OperationName:  operation,
		Login:          userData.Login,
		Operation:      "",
		ReturnRequest:  "",
		SignOutRequest: "",
	}

	curr_server_data := &models.RequestData{
		Port: uc.cfg.HTTPServer.Port[1:],
	}

	var buffer bytes.Buffer

	template_sign_out_request, err := template.New("RequestSignOut").Parse(html.RequestSignOut)
	if err != nil {
		return "", err
	}

	err = template_sign_out_request.Execute(&buffer, &curr_server_data)
	if err != nil {
		return "", err
	}

	totpOperationInfo.SignOutRequest = buffer.String()
	buffer.Reset()

	template_request_user_page, err := template.New("RequestUserPage").Parse(html.RequestUserPage)
	if err != nil {
		return "", err
	}

	err = template_request_user_page.Execute(&buffer, &curr_server_data)
	if err != nil {
		return "", err
	}

	totpOperationInfo.ReturnRequest = buffer.String()
	buffer.Reset()

	switch operation {
	case TotpOperationTypeTurnOn:
		{
			templateGetTotpCodeRequest, err := template.New("RequestTurnOnTotp").Parse(html.RequestTurnOnTotp)
			if err != nil {
				return "", err
			}

			err = templateGetTotpCodeRequest.Execute(&buffer, &curr_server_data)
			if err != nil {
				return "", err
			}

			operationData := &models.TotpOperationData{
				OperationRequest: buffer.String(),
			}
			buffer.Reset()

			templateGetTotpCode, err := template.New("TotpOperationOpen").Parse(html.TotpOperationOpen)
			if err != nil {
				return "", err
			}

			err = templateGetTotpCode.Execute(&buffer, &operationData)
			if err != nil {
				return "", err
			}

			totpOperationInfo.Operation = buffer.String()
			buffer.Reset()

		}
	case TotpOperationTypeTurnOff:
		{
			templateTurnOffTotpRequest, err := template.New("RequestTurnOffTotp").Parse(html.RequestTurnOffTotp)
			if err != nil {
				return "", err
			}

			err = templateTurnOffTotpRequest.Execute(&buffer, &curr_server_data)
			if err != nil {
				return "", err
			}

			operationData := &models.TotpOperationData{
				OperationRequest: buffer.String(),
			}
			buffer.Reset()

			templateTurnOffTotpOperation, err := template.New("TotpOperationClose").Parse(html.TotpOperationClose)
			if err != nil {
				return "", err
			}

			err = templateTurnOffTotpOperation.Execute(&buffer, &operationData)
			if err != nil {
				return "", err
			}

			totpOperationInfo.Operation = buffer.String()
			buffer.Reset()
		}
	case TotpOperationTypeQr:
		{

			templateTotpQrOperation, err := template.New("TotpOperationQr").Parse(html.TotpOperationQr)
			if err != nil {
				return "", err
			}

			totpUrl, err := uc.getUserTotpUrl(userId)
			if err != nil {
				return "", err
			}

			qrImageName, err := uc.createTotpQr(userId, totpUrl)
			if err != nil {
				return "", err
			}

			totpQrOperationInfo := &models.TotpQrImagePath{
				ImagePath: "qr/" + qrImageName,
			}

			err = templateTotpQrOperation.Execute(&buffer, &totpQrOperationInfo)
			if err != nil {
				return "", err
			}

			totpOperationInfo.Operation = buffer.String()
			buffer.Reset()
		}
	default:
		return "", ErrorUnknownTotpOperationType
	}

	templateTotpOperationPage, err := template.New("TotpOperationPage").Parse(html.TotpOperationPage)
	if err != nil {
		return "", err
	}

	err = templateTotpOperationPage.Execute(&buffer, &totpOperationInfo)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil

}

func (uc *apiGateWayUseCase) CheckTotp(userId uuid.UUID, code string) error {

	return uc.checkTotpRequest(userId, code)

}

func (uc *apiGateWayUseCase) checkTotpRequest(userId uuid.UUID, code string) error {

	templateTotpCheckRequest, err := template.New("RequestTotpValidate").Parse(RequestTotpValidate)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer

	err = templateTotpCheckRequest.Execute(&buffer, uc.totpServerInfo)
	if err != nil {
		return err
	}

	totpCheckRequest := buffer.String()
	buffer.Reset()

	request_get_operations_body := &models.TotpCodeValidateBody{
		UserId:   userId,
		TotpCode: code,
	}

	request_body, err := json.Marshal(&request_get_operations_body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, totpCheckRequest, bytes.NewBuffer(request_body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.totpServerInfo.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		resp_body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var resp_data = &models.OperationResponse{}

		err = json.Unmarshal(resp_body, &resp_data)
		if err != nil {
			return err
		}

		return errors.New(resp_data.Info)
	}

}

func (uc *apiGateWayUseCase) createTotpQr(userId uuid.UUID, url string) (string, error) {

	fileName := userId.String() + time.Now().Format("02-01-2006_15:04:05")
	qrFileName, err := createQrCode(url, fileName, uc.qrImagesPath)
	if err != nil {
		return "", err
	}

	return qrFileName, nil
}

func (uc *apiGateWayUseCase) CreateAdminPage(begin string, end string) (string, error) {

	template_admin_page, err := template.New("AdminPage").Parse(html.AdminPage)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer

	admin_page_data := &models.AdminPageData{
		GetOperationsRequest: "",
		Operations:           "",
	}

	template_get_operations_request, err := template.New("RequestAdminPage").Parse(html.RequestAdminPage)
	if err != nil {
		return "", err
	}

	curr_server_data := &models.RequestData{
		Port: uc.cfg.HTTPServer.Port[1:],
	}

	err = template_get_operations_request.Execute(&buffer, &curr_server_data)
	if err != nil {
		return "", err
	}
	admin_page_data.GetOperationsRequest = buffer.String()
	buffer.Reset()

	operations := ""

	operations_id_list, err := uc.getListOfOperations(begin, end)
	if err != nil {
		return "", err
	}

	for _, operation_id := range operations_id_list.Operations {

		operation_data, _ := uc.GetOperationData(operation_id)

		if operation_data != nil {

			operation_tree, err := uc.getOperationTree(operation_id)
			if err != nil {
				return "", err
			}

			graph_file_name := operation_id.String() + "_" + time.Now().Format("02-01-2006_15:04:05")
			graph_image_path, err := CreateGraph(operation_tree, uc.graphImagesPath, graph_file_name)
			if err != nil {
				return "", err
			}

			template_admin_operation, err := template.New("AdminOperation").Parse(html.AdminOperation)
			if err != nil {
				return "", err
			}

			admin_operation_data := &models.AdminOperationData{
				Id:        operation_id,
				Name:      operation_tree.OperationName,
				Status:    operation_data.Info,
				Begin:     "",
				End:       "",
				ImagePath: "graph/" + graph_image_path,
			}

			additional_data := operation_data.AdditionalInfo.(map[string]interface{})

			if time_begin, ok := additional_data["time_begin"]; ok {
				tmp, err := time.Parse(time.RFC3339Nano, time_begin.(string))
				if err != nil {
					return "", err
				}
				tmp_str := tmp.Format("02-01-2006 15:04:05")
				admin_operation_data.Begin = tmp_str
			}
			if time_end, ok := additional_data["time_end"]; ok {
				tmp, err := time.Parse(time.RFC3339Nano, time_end.(string))
				if err != nil {
					return "", err
				}
				tmp_str := tmp.Format("02-01-2006 15:04:05")
				admin_operation_data.End = tmp_str
			}

			err = template_admin_operation.Execute(&buffer, &admin_operation_data)
			if err != nil {
				return "", err
			}

			operations += buffer.String() + "\n"
			buffer.Reset()
		}

	}

	admin_page_data.Operations = operations

	err = template_admin_page.Execute(&buffer, &admin_page_data)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (uc *apiGateWayUseCase) SignIn(login_info *models.SignInInfo) (*models.Token, error) {

	user_data, err := uc.GetUserDataByLoginRequest(login_info.Login)
	if err != nil {
		return nil, err
	}

	is_ok, hasTotp, err := uc.CheckUserPasswordRequest(login_info.Login, login_info.Password)
	if err != nil {
		return nil, err
	}

	hasTotp = hasTotp

	if is_ok {

		token, err := uc.CreateToken(context.Background(), uuid.New(), TokenLiveTime, user_data.Id)
		if err != nil {
			return nil, err
		}

		return token, nil
	}

	return nil, ErrorWrongPassword

}

func (uc *apiGateWayUseCase) SignUp(sign_up_info *models.SignUpInfo) (*models.Token, error) {

	user_id, err := uc.CreateUserRequest(sign_up_info)
	if err != nil {
		return nil, err
	}

	_ = uc.CreateNotificationSignUp(context.Background(), user_id)

	token, err := uc.CreateToken(context.Background(), uuid.New(), TokenLiveTime, user_id)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (uc *apiGateWayUseCase) CreateNotificationSignUp(ctx context.Context, userId uuid.UUID) error {

	template_welcome_notify, err := template.New("NotificationWelcome").Parse(notifications.NotificationWelcome)
	if err != nil {
		return err
	}

	userInfo, err := uc.GetUserDataRequest(userId)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer

	userNotifyInfo := &models.WelcomeMessage{
		Login: userInfo.Login,
		Name:  userInfo.Surname + " " + userInfo.Name + " " + userInfo.Patronymic,
	}

	err = template_welcome_notify.Execute(&buffer, userNotifyInfo)
	if err != nil {
		return err
	}

	err = uc.createNotification(ctx, userId, notifications.NotificationLvlAll, buffer.String())
	if err != nil {
		return err
	}

	return nil

}

func (uc *apiGateWayUseCase) CreateNotificationSignIn(ctx context.Context, userId uuid.UUID) error {

	template_message_sign_in, err := template.New("NotificationSignIn").Parse(notifications.NotificationSignIn)
	if err != nil {
		return err
	}

	userInfo, err := uc.GetUserDataRequest(userId)
	if err != nil {
		return err
	}

	msgInfo := &models.SignInMessage{
		Login: userInfo.Login,
	}

	var buffer bytes.Buffer

	err = template_message_sign_in.Execute(&buffer, &msgInfo)
	if err != nil {
		return err
	}

	err = uc.createNotification(ctx, userId, notifications.NotificationLvlEmail, buffer.String())
	if err != nil {
		return err
	}

	return nil

}

func (uc *apiGateWayUseCase) CreateAccount(user_id uuid.UUID, account_info *models.AccountInfo) error {

	return uc.openAccountRequest(user_id, account_info)

}

func (uc *apiGateWayUseCase) CloseAccount(user_id uuid.UUID, account_id uuid.UUID) error {
	return uc.closeAccountRequest(user_id, account_id)
}

func (uc *apiGateWayUseCase) AddAccountCache(user_id uuid.UUID, account_id uuid.UUID, money float64) error {
	return uc.addAccountCacheRequest(user_id, account_id, money)
}

func (uc *apiGateWayUseCase) WidthAccountCache(user_id uuid.UUID, account_id uuid.UUID, money float64) error {
	return uc.widthAccountCacheRequest(user_id, account_id, money)
}

func (uc *apiGateWayUseCase) widthAccountCacheRequest(user_id uuid.UUID, account_id uuid.UUID, money float64) error {

	template_request_width_account_cache, err := template.New("RequestWidthAccountCache").Parse(RequestWidthAccountCache)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer

	err = template_request_width_account_cache.Execute(&buffer, uc.registrationServerInfo)
	if err != nil {
		return err
	}

	request_width_account_cache := buffer.String()
	buffer.Reset()

	request_width_account_cache_body := &models.AddAccountCacheBody{
		UserId:    user_id,
		AccountId: account_id,
		CacheDiff: money,
	}

	request_body, err := json.Marshal(&request_width_account_cache_body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, request_width_account_cache, bytes.NewBuffer(request_body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.registrationServerInfo.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var resp_data = &models.OperationResponse{}

	err = json.Unmarshal(resp_body, &resp_data)
	if err != nil {
		return err
	}

	operation_id_str := resp_data.Info

	operation_id, err := uuid.Parse(operation_id_str)
	if err != nil {
		return err
	}

	_, err = uc.GetOperationData(operation_id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *apiGateWayUseCase) addAccountCacheRequest(user_id uuid.UUID, account_id uuid.UUID, money float64) error {

	template_request_add_account_cache, err := template.New("RequestAddAccountCache").Parse(RequestAddAccountCache)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer

	err = template_request_add_account_cache.Execute(&buffer, uc.registrationServerInfo)
	if err != nil {
		return err
	}

	request_add_account_cache := buffer.String()
	buffer.Reset()

	request_add_account_cache_body := &models.AddAccountCacheBody{
		UserId:    user_id,
		AccountId: account_id,
		CacheDiff: money,
	}

	request_body, err := json.Marshal(&request_add_account_cache_body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, request_add_account_cache, bytes.NewBuffer(request_body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.registrationServerInfo.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var resp_data = &models.OperationResponse{}

	err = json.Unmarshal(resp_body, &resp_data)
	if err != nil {
		return err
	}

	operation_id_str := resp_data.Info

	operation_id, err := uuid.Parse(operation_id_str)
	if err != nil {
		return err
	}

	_, err = uc.GetOperationData(operation_id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *apiGateWayUseCase) closeAccountRequest(user_id uuid.UUID, account_id uuid.UUID) error {

	template_request_close_account, err := template.New("RequestCloseAccount").Parse(RequestCloseAccount)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer

	err = template_request_close_account.Execute(&buffer, uc.registrationServerInfo)
	if err != nil {
		return err
	}

	request_close_account := buffer.String()
	buffer.Reset()

	request_close_account_body := &models.CloseAccountBody{
		UserId:    user_id,
		AccountId: account_id,
	}

	request_body, err := json.Marshal(&request_close_account_body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, request_close_account, bytes.NewBuffer(request_body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.registrationServerInfo.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var resp_data = &models.OperationResponse{}

	err = json.Unmarshal(resp_body, &resp_data)
	if err != nil {
		return err
	}

	operation_id_str := resp_data.Info

	operation_id, err := uuid.Parse(operation_id_str)
	if err != nil {
		return err
	}

	_, err = uc.GetOperationData(operation_id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *apiGateWayUseCase) openAccountRequest(user_id uuid.UUID, account_info *models.AccountInfo) error {

	template_request_open_account, err := template.New("RequestOpenAccount").Parse(RequestOpenAccount)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer

	err = template_request_open_account.Execute(&buffer, uc.registrationServerInfo)
	if err != nil {
		return err
	}

	request_open_account := buffer.String()
	buffer.Reset()

	request_open_account_body := &models.OpenAccountBody{
		UserId:        user_id,
		BIC:           account_info.BIC,
		CIO:           account_info.CIO,
		CulcNumber:    account_info.CulcNumber,
		CorrNumber:    account_info.CorrNumber,
		AccName:       account_info.Name,
		ReserveReason: "",
	}

	request_body, err := json.Marshal(&request_open_account_body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, request_open_account, bytes.NewBuffer(request_body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.registrationServerInfo.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var resp_data = &models.OperationResponse{}

	err = json.Unmarshal(resp_body, &resp_data)
	if err != nil {
		return err
	}

	operation_id_str := resp_data.Info

	operation_id, err := uuid.Parse(operation_id_str)
	if err != nil {
		return err
	}

	_, err = uc.GetOperationData(operation_id)
	if err != nil {
		return err
	}

	return nil

}

func (uc *apiGateWayUseCase) getListOfOperations(start string, end string) (*models.ListOfOperations, error) {

	template_request_get_operations, err := template.New("RequestGetListOfOperations").Parse(RequestGetListOfOperations)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer

	err = template_request_get_operations.Execute(&buffer, uc.registrationServerInfo)
	if err != nil {
		return nil, err
	}

	request_get_operations := buffer.String()
	buffer.Reset()

	request_get_operations_body := &models.OperationListRequestBody{
		TimeBegin: start,
		TimeEnd:   end,
	}

	request_body, err := json.Marshal(&request_get_operations_body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, request_get_operations, bytes.NewBuffer(request_body))
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

	var resp_data = &models.OperationListRequestResultBody{}

	err = json.Unmarshal(resp_body, &resp_data)
	if err != nil {
		return nil, err
	}

	info := resp_data.Info.(map[string]interface{})

	result := &models.ListOfOperations{
		Operations: make([]uuid.UUID, 0),
	}

	if operations, ok := info["operations"]; ok {
		list_of_operations := operations.([]interface{})
		for _, operation_id_str := range list_of_operations {
			operation_id, err := uuid.Parse(operation_id_str.(string))
			if err != nil {
				return nil, err
			}
			result.Operations = append(result.Operations, operation_id)
		}
	}

	return result, nil
}

func (uc *apiGateWayUseCase) getOperationTree(operation_id uuid.UUID) (*models.OperationTree, error) {

	template_request_get_operation_tree, err := template.New("RequestGetOperationTree").Parse(RequestGetOperationTree)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer

	err = template_request_get_operation_tree.Execute(&buffer, uc.registrationServerInfo)
	if err != nil {
		return nil, err
	}

	request_get_operation_tree := buffer.String()
	buffer.Reset()

	request_get_operation_tree_body := &models.OperationTreeRequestBody{
		OperationId: operation_id,
	}

	request_body, err := json.Marshal(&request_get_operation_tree_body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, request_get_operation_tree, bytes.NewBuffer(request_body))
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

	var resp_data = &models.OperationTreeRequestResultBody{}

	err = json.Unmarshal(resp_body, &resp_data)
	if err != nil {
		return nil, err
	}

	info := resp_data.Info.(map[string]interface{})

	result := &models.OperationTree{
		OperationName:  "",
		SagaList:       make(map[uuid.UUID]*models.SagaTree, 0),
		EventList:      make(map[uuid.UUID]*models.EventTree, 0),
		SagaDependList: make([]*models.SagaDependTree, 0),
	}

	if saga, ok := info["saga"]; ok && saga != nil {
		saga_list := saga.([]interface{})
		for _, saga_info_inter := range saga_list {
			saga_info := saga_info_inter.(map[string]interface{})
			saga_tree := &models.SagaTree{
				Events: make([]uuid.UUID, 0),
			}

			if name, ok := saga_info["name"]; ok {
				saga_tree.Name = name.(string)
			}
			if status, ok := saga_info["status"]; ok {
				saga_tree.Status = status.(float64)
			}
			if id_str, ok := saga_info["id"]; ok {
				id, err := uuid.Parse(id_str.(string))
				if err != nil {
					return nil, err
				}
				saga_tree.Id = id
			}
			if list_of_events, ok := saga_info["events"]; ok {
				for _, event_id_str := range list_of_events.([]interface{}) {
					event_id, err := uuid.Parse(event_id_str.(string))
					if err != nil {
						return nil, err
					}
					saga_tree.Events = append(saga_tree.Events, event_id)
				}
			}

			result.SagaList[saga_tree.Id] = saga_tree

		}
	}

	if event, ok := info["events"]; ok && event != nil {
		for _, event_inter := range event.([]interface{}) {
			event_info := event_inter.(map[string]interface{})
			event_tree := &models.EventTree{}

			if name, ok := event_info["name"]; ok {
				event_tree.Name = name.(string)
			}
			if status, ok := event_info["status"]; ok {
				event_tree.Status = status.(float64)
			}
			if id_str, ok := event_info["id"]; ok {
				id, err := uuid.Parse(id_str.(string))
				if err != nil {
					return nil, err
				}
				event_tree.Id = id
			}
			if id_rollback_str, ok := info["roll_back_id"]; ok {
				id, err := uuid.Parse(id_rollback_str.(string))
				if err != nil {
					return nil, err
				}
				event_tree.RollBackId = id
			}

			result.EventList[event_tree.Id] = event_tree
		}
	}

	if depends, ok := info["saga_depend"]; ok && depends != nil {
		for _, saga_depends := range depends.([]interface{}) {
			saga_depend_info := saga_depends.(map[string]interface{})
			saga_depend_tree := &models.SagaDependTree{}

			if parent_id_str, ok := saga_depend_info["parent_id"]; ok {
				id, err := uuid.Parse(parent_id_str.(string))
				if err != nil {
					return nil, err
				}
				saga_depend_tree.ParentId = id
			}

			if child_id_str, ok := saga_depend_info["child_id"]; ok {
				id, err := uuid.Parse(child_id_str.(string))
				if err != nil {
					return nil, err
				}
				saga_depend_tree.ChildId = id
			}

			result.SagaDependList = append(result.SagaDependList, saga_depend_tree)
		}
	}

	if name, ok := info["operation_name"]; ok {
		result.OperationName = name.(string)
	}

	return result, nil
}

func (uc *apiGateWayUseCase) GetAccountDataRequest(account_id uuid.UUID) (*models.AccountInfo, error) {

	template_request_get_account_data, err := template.New("GetAccountData").Parse(GetAccountData)
	if err != nil {
		return nil, err
	}

	var writer bytes.Buffer

	err = template_request_get_account_data.Execute(&writer, uc.accountsServerInfo)
	if err != nil {
		return nil, err
	}
	request_get_account_data := writer.String()
	writer.Reset()

	urlParams := url.Values{}
	urlParams.Add("acc_id", account_id.String())

	fullUrl := fmt.Sprintf("%s?%s", request_get_account_data, urlParams.Encode())

	req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.accountsServerInfo.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var resp_data = &models.GetAccountDataResponse{}

	err = json.Unmarshal(resp_body, &resp_data)
	if err != nil {
		return nil, err
	}

	result := &models.AccountInfo{
		Id:         account_id,
		Name:       resp_data.Acc_name,
		CorrNumber: resp_data.Acc_corr_number,
		CulcNumber: resp_data.Acc_culc_number,
		CIO:        resp_data.Acc_cio,
		BIC:        resp_data.Acc_bic,
		Cache:      resp_data.Acc_money_amount,
	}

	status_str := "Unknown"
	switch resp_data.Acc_status {
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

	return result, nil
}

func (uc *apiGateWayUseCase) GetUserDataRequest(user_id uuid.UUID) (*models.UserInfo, error) {

	template_request_get_user_data, err := template.New("GetUserData").Parse(GetUserData)
	if err != nil {
		return nil, err
	}

	var writer bytes.Buffer

	err = template_request_get_user_data.Execute(&writer, uc.usersServerInfo)
	if err != nil {
		return nil, err
	}
	request_get_user_data := writer.String()
	writer.Reset()

	params := url.Values{}
	params.Add("user_id", user_id.String())

	fullUrl := fmt.Sprintf("%s?%s", request_get_user_data, params.Encode())

	req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.usersServerInfo.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var resp_get_user_data = &models.GetUserDataResponse{}

	err = json.Unmarshal(resp_body, &resp_get_user_data)
	if err != nil {
		return nil, err
	}

	result := &models.UserInfo{
		Email:                       "",
		Login:                       resp_get_user_data.UserLogin,
		Patronymic:                  resp_get_user_data.PassportData.FCS.Patronymic,
		Name:                        resp_get_user_data.PassportData.FCS.Name,
		Surname:                     resp_get_user_data.PassportData.FCS.Surname,
		Accounts:                    resp_get_user_data.Accounts,
		BirthLocation:               resp_get_user_data.PassportData.BirthLocation,
		BirthDate:                   resp_get_user_data.PassportData.BirthDate,
		Id:                          user_id,
		PassportRegistrationAddress: resp_get_user_data.PassportData.RegistrationAdress,
		PassportAuthorityDate:       resp_get_user_data.PassportData.AuthorityDate,
		PassportAuthority:           resp_get_user_data.PassportData.Authority,
		PassportPickUpPoint:         resp_get_user_data.PassportData.PickUpPoint,
		PassportNumber:              resp_get_user_data.PassportData.Number,
		PassportSeries:              resp_get_user_data.PassportData.Series,
		Inn:                         resp_get_user_data.UserInn,
		UsingTotp:                   resp_get_user_data.UsingTotp,
	}

	template_request_get_user_notif_settings, err := template.New("GetUserNotificationSettings").Parse(GetUserNotificationSettings)
	if err != nil {
		return nil, err
	}

	writer.Reset()

	err = template_request_get_user_notif_settings.Execute(&writer, uc.notificationServerInfo)
	if err != nil {
		return nil, err
	}
	request_get_user_notif_settings := writer.String()
	writer.Reset()

	params = url.Values{}
	params.Add("user_id", user_id.String())

	fullUrl = fmt.Sprintf("%s?%s", request_get_user_notif_settings, params.Encode())

	req, err = http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client = &http.Client{
		Timeout: uc.usersServerInfo.TimeWaitResponse,
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	resp_body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var resp_data_notify_settings = &models.GetUserNotifySettingsResponse{}

	err = json.Unmarshal(resp_body, &resp_data_notify_settings)
	if err != nil {
		return nil, err
	}

	result.Email = resp_data_notify_settings.Email

	return result, nil
}

func (uc *apiGateWayUseCase) CheckUserPasswordRequest(login string, password string) (bool, bool, error) {

	template_request_check_user_password, err := template.New("RequestCheckUserPassword").Parse(RequestCheckUserPassword)
	if err != nil {
		return false, false, err
	}

	var writer bytes.Buffer

	err = template_request_check_user_password.Execute(&writer, &uc.usersServerInfo)
	if err != nil {
		return false, false, err
	}
	request_check_user_password := writer.String()
	writer.Reset()

	request_check_user_password_body := &models.CheckUserPasswordBody{
		Password: password,
		Login:    login,
	}

	request_body, err := json.Marshal(&request_check_user_password_body)
	if err != nil {
		return false, false, err
	}

	req, err := http.NewRequest(http.MethodPost, request_check_user_password, bytes.NewBuffer(request_body))
	if err != nil {
		return false, false, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.usersServerInfo.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, false, err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, false, err
	}

	if resp.StatusCode == http.StatusOK {
		var resp_data = &models.CheckPasswordResponse{}

		err = json.Unmarshal(resp_body, &resp_data)
		if err != nil {
			return false, false, err
		}

		return true, resp_data.TotpUsage, nil
	} else {
		var resp_data = &models.OperationResponse{}

		err = json.Unmarshal(resp_body, &resp_data)
		if err != nil {
			return false, false, err
		}
		return false, false, errors.New(resp_data.Info)
	}

}

func (uc *apiGateWayUseCase) getUserTotpUrl(userId uuid.UUID) (string, error) {

	templateRequestGetTotpUrl, err := template.New("RequestGetTotpUrl").Parse(RequestGetTotpUrl)
	if err != nil {
		return "", err
	}

	var writer bytes.Buffer

	err = templateRequestGetTotpUrl.Execute(&writer, &uc.totpServerInfo)
	if err != nil {
		return "", err
	}
	requestgetTotpUrl := writer.String()
	writer.Reset()

	requestGetTotpUrlBody := &models.GetTotpUrlBody{
		UserId: userId,
	}

	request_body, err := json.Marshal(&requestGetTotpUrlBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodGet, requestgetTotpUrl, bytes.NewBuffer(request_body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.totpServerInfo.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusOK {
		var resp_data = &models.TotpInfo{}

		err = json.Unmarshal(resp_body, &resp_data)
		if err != nil {
			return "", err
		}

		return resp_data.TotpUrl, nil
	} else {
		var resp_data = &models.OperationResponse{}

		err = json.Unmarshal(resp_body, &resp_data)
		if err != nil {
			return "", err
		}
		return "", errors.New(resp_data.Info)
	}

}

func (uc *apiGateWayUseCase) TurnOnTotp(userId uuid.UUID) error {

	userInfo, err := uc.GetUserDataRequest(userId)
	if err != nil {
		return err
	}

	totpInfo, err := uc.enrollTotpRequest(userId, userInfo.Login)
	if err != nil {
		return err
	}

	totpInfo.TotpUsage = true
	err = uc.updateTotpUserInfo(userId, totpInfo)
	if err != nil {
		return err
	}

	_ = uc.createNotificationTurnOnTotp(context.Background(), userId)

	return nil
}

func (uc *apiGateWayUseCase) createNotificationTurnOnTotp(ctx context.Context, userId uuid.UUID) error {

	templateMessageTurnOnTotp, err := template.New("NotificationTurnOnTotp").Parse(notifications.NotificationTurnOnTotp)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer

	userInfo, err := uc.GetUserDataRequest(userId)
	if err != nil {
		return err
	}

	messageData := &models.TurnOnTotpMessage{
		Login: userInfo.Login,
	}

	err = templateMessageTurnOnTotp.Execute(&buffer, &messageData)
	if err != nil {
		return err
	}

	err = uc.createNotification(ctx, userId, notifications.NotificationLvlEmail, buffer.String())
	if err != nil {
		return err
	}

	return nil

}

func (uc *apiGateWayUseCase) createNotificationTurnOffTotp(ctx context.Context, userId uuid.UUID) error {

	templateMessageTurnOffTotp, err := template.New("NotificationTurnOffTotp").Parse(notifications.NotificationTurnOffTotp)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer

	userInfo, err := uc.GetUserDataRequest(userId)
	if err != nil {
		return err
	}

	messageData := &models.TurnOffTotpMessage{
		Login: userInfo.Login,
	}

	err = templateMessageTurnOffTotp.Execute(&buffer, &messageData)
	if err != nil {
		return err
	}

	err = uc.createNotification(ctx, userId, notifications.NotificationLvlAll, buffer.String())
	if err != nil {
		return err
	}

	return nil

}

func (uc *apiGateWayUseCase) TurnOffTotp(userId uuid.UUID) error {

	userInfo, err := uc.GetUserDataRequest(userId)
	if err != nil {
		return err
	}

	if userInfo.UsingTotp {
		totpInfo, err := uc.GetUserTotpInfo(userId)
		if err != nil {
			return err
		}

		err = uc.disableTotp(userId, totpInfo.TotpId)
		if err != nil {
			return err
		}

		totpInfo.TotpUsage = false
		totpInfo.TotpId = uuid.Nil

		err = uc.updateTotpUserInfo(userId, totpInfo)
		if err != nil {
			return err
		}

		_ = uc.createNotificationTurnOffTotp(context.Background(), userId)

	}

	return nil
}

func (uc *apiGateWayUseCase) disableTotp(userId uuid.UUID, totpId uuid.UUID) error {

	templateRequestTurnOffTotp, err := template.New("RequestTurnOffTotp").Parse(RequestTurnOffTotp)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer

	err = templateRequestTurnOffTotp.Execute(&buffer, &uc.totpServerInfo)
	if err != nil {
		return err
	}

	requestEnrollTotp := buffer.String()
	buffer.Reset()

	requestBody := &models.TotpDisactivateRequestBody{
		UserId: userId,
		TotpId: totpId,
	}

	request_body, err := json.Marshal(&requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, requestEnrollTotp, bytes.NewBuffer(request_body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.totpServerInfo.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {

		return nil

	} else {
		var resp_data = &models.OperationResponse{}

		err = json.Unmarshal(resp_body, &resp_data)
		if err != nil {
			return err
		}
		return errors.New(resp_data.Info)
	}
}

func (uc *apiGateWayUseCase) GetUserTotpInfo(userId uuid.UUID) (*models.TotpInfo, error) {

	templateRequestUserTotpInfo, err := template.New("RequestGetUserTotpInfo").Parse(RequestGetUserTotpInfo)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer

	err = templateRequestUserTotpInfo.Execute(&buffer, &uc.usersServerInfo)
	if err != nil {
		return nil, err
	}

	requestUserTotpInfo := buffer.String()
	buffer.Reset()

	requestBody := &models.GetUserTotpInfoBody{
		UserId: userId,
	}

	request_body, err := json.Marshal(&requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, requestUserTotpInfo, bytes.NewBuffer(request_body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.usersServerInfo.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var resp_data = &models.TotpInfo{}

		err = json.Unmarshal(resp_body, &resp_data)
		if err != nil {
			return nil, err
		}
		return resp_data, nil

	} else {
		var resp_data = &models.OperationResponse{}

		err = json.Unmarshal(resp_body, &resp_data)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(resp_data.Info)
	}
}

func (uc *apiGateWayUseCase) updateTotpUserInfo(userId uuid.UUID, totpInfo *models.TotpInfo) error {

	templateRequestUpdateTotpUsersInfo, err := template.New("RequestUpdateTotpInfo").Parse(RequestUpdateTotpInfo)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer

	err = templateRequestUpdateTotpUsersInfo.Execute(&buffer, &uc.usersServerInfo)
	if err != nil {
		return err
	}

	requestUpdateTotpUsersInfo := buffer.String()
	buffer.Reset()

	requestBody := &models.UpdateTotpUsersInfoBody{
		UserId:    userId,
		TotpId:    totpInfo.TotpId,
		TotpUsage: totpInfo.TotpUsage,
	}

	request_body, err := json.Marshal(&requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, requestUpdateTotpUsersInfo, bytes.NewBuffer(request_body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.usersServerInfo.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil

	} else {
		var resp_data = &models.OperationResponse{}

		err = json.Unmarshal(resp_body, &resp_data)
		if err != nil {
			return err
		}
		return errors.New(resp_data.Info)
	}
}

func (uc *apiGateWayUseCase) enrollTotpRequest(userId uuid.UUID, userName string) (*models.TotpInfo, error) {

	template_request_enroll_totp, err := template.New("RequestCreateTotp").Parse(RequestCreateTotp)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer

	err = template_request_enroll_totp.Execute(&buffer, &uc.totpServerInfo)
	if err != nil {
		return nil, err
	}

	requestEnrollTotp := buffer.String()
	buffer.Reset()

	requestBody := &models.TotpEnrollRequestBody{
		UserId:   userId,
		UserName: userName,
	}

	request_body, err := json.Marshal(&requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, requestEnrollTotp, bytes.NewBuffer(request_body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.totpServerInfo.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusCreated {
		var resp_data = &models.TotpEnrollResponse{}

		err = json.Unmarshal(resp_body, &resp_data)
		if err != nil {
			return nil, err
		}

		result := &models.TotpInfo{
			TotpId:  resp_data.TotpId,
			TotpUrl: resp_data.TotpUrl,
		}

		return result, nil

	} else {
		var resp_data = &models.OperationResponse{}

		err = json.Unmarshal(resp_body, &resp_data)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(resp_data.Info)
	}
}

func (uc *apiGateWayUseCase) GetUserDataByLoginRequest(login string) (*models.UserInfo, error) {

	template_request_get_user_data_by_login, err := template.New("RequestGetUserDataByLogin").Parse(RequestGetUserDataByLogin)
	if err != nil {
		return nil, err
	}

	var writer bytes.Buffer

	err = template_request_get_user_data_by_login.Execute(&writer, &uc.usersServerInfo)
	if err != nil {
		return nil, err
	}
	request_get_user_data_by_login := writer.String()
	writer.Reset()

	urlParams := url.Values{}
	urlParams.Add("login", login)

	fullUrl := fmt.Sprintf("%s?%s", request_get_user_data_by_login, urlParams.Encode())

	req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: uc.usersServerInfo.TimeWaitResponse,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var resp_data = &models.GetUserDataByLoginResponse{}

	err = json.Unmarshal(resp_body, &resp_data)
	if err != nil {
		return nil, err
	}

	result := &models.UserInfo{
		Id: resp_data.UserId,
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
			AuthorityDate:       user_info.PassportAuthorityDate + " 00:00:00",
			BirthDate:           user_info.BirthDate + " 00:00:00",
		},
		UserEmail: user_info.Email,
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

	if resp_data.Info == "" {
		err_str := "Operation_error, code: " + strconv.Itoa(resp_data.Status)
		return uuid.Nil, errors.New(err_str)
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

	user_id, err := uuid.Parse(additional_info["user_id"].(string))
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

func (uc *apiGateWayUseCase) createNotification(ctx context.Context, userId uuid.UUID, notificationLvl string, message string) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "apiGateWayUseCase.CreateNotification")
	defer span.Finish()

	headers := make(amqp.Table)
	headers[HeaderUserId] = userId.String()
	headers[HeaderNotificationLvl] = notificationLvl
	err := uc.rmqChan.PublishWithContext(ctxWithTrace,
		"",               // exchange
		uc.rmqQueue.Name, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			Headers:     headers,
		})
	if err != nil {
		return err
	}

	return nil
}

func NewApiGatewayUseCase(cfg *config.Config, repo api_gateway.Repository, registration_server_info *models.InternalServerInfo,
	usersServerInfo *models.InternalServerInfo, accountsServerInfo *models.InternalServerInfo,
	notificationServerInfo *models.InternalServerInfo, totpServerInfo *models.InternalServerInfo,
	graphImagesPath string, qrImagesPath string, rmqChan *amqp091.Channel,
	rmqQueue amqp091.Queue) api_gateway.UseCase {
	return &apiGateWayUseCase{cfg: cfg, repo: repo, registrationServerInfo: registration_server_info, graphImagesPath: graphImagesPath,
		rmqQueue: rmqQueue, rmqChan: rmqChan, accountsServerInfo: accountsServerInfo, usersServerInfo: usersServerInfo,
		notificationServerInfo: notificationServerInfo, totpServerInfo: totpServerInfo, qrImagesPath: qrImagesPath}
}
