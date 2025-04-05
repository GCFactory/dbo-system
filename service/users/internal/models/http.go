package models

import "github.com/google/uuid"

type DefaultHttpResponse struct {
	Status int    `json:"status" validate:"required"`
	Info   string `json:"info" validate:"required"`
}

type GetUserData struct {
	UserId string `json:"user_id" validate:"required"`
}

type GetUserDataResponse_Passport_FCS struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type GetUserDataResponse_Passport struct {
	Series             string                            `json:"series"`
	Number             string                            `json:"number"`
	BirthDate          string                            `json:"birth_date"`
	BirthLocation      string                            `json:"birth_location"`
	PickUpPoint        string                            `json:"pick_up_point"`
	Authority          string                            `json:"authority"`
	AuthorityDate      string                            `json:"authority_date"`
	RegistrationAdress string                            `json:"registration_adress"`
	FCS                *GetUserDataResponse_Passport_FCS `json:"fcs"`
}

type GetUserDataResponse struct {
	UserInn      string                        `json:"user_inn"`
	UserId       uuid.UUID                     `json:"user_id"`
	UserLogin    string                        `json:"user_login"`
	PassportData *GetUserDataResponse_Passport `json:"passport_data"`
	Accounts     []uuid.UUID                   `json:"accounts"`
	UsingTotp    bool                          `json:"using_totp"`
}

type GetUserAccount struct {
	UserId string `json:"user_id"`
}

type GetUserAccountResponse struct {
	Accounts []uuid.UUID `json:"accounts"`
}

type GetUserDataByLogin struct {
	Login string `json:"login"`
}

type GetUserDataByLoginResponse struct {
	UserId uuid.UUID `json:"user_id"`
}

type CheckPasswordRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CheckPasswordResponse struct {
	UserId    uuid.UUID `json:"user_id"`
	TotpUsage bool      `json:"totp_usage"`
}

type GetUserTotpDataRequest struct {
	UserId string `json:"user_id"`
}

type GetUserTotpDataResponse struct {
	TotpId    uuid.UUID `json:"totp_id"`
	TotpUsage bool      `json:"totp_usage"`
}

type UpdateTotpInfoRequest struct {
	UserId    string `json:"user_id"`
	TotpId    string `json:"totp_id"`
	TotpUsage bool   `json:"totp_usage"`
}
