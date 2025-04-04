package models

import (
	"github.com/google/uuid"
	"time"
)

type UserInfo struct {
	Id                          uuid.UUID
	Login                       string      `json:"login" validate:"required"`
	Password                    string      `json:"password" validate:"required"`
	PassportSeries              string      `json:"passport_series" validate:"required"`
	PassportNumber              string      `json:"passport_number" validate:"required"`
	Name                        string      `json:"name" validate:"required"`
	Surname                     string      `json:"surname" validate:"required"`
	Patronymic                  string      `json:"patronymic" validate:"required"`
	BirthDate                   string      `json:"birth_date" validate:"required,datetime=02-01-2006 15:04:05,min=1"`
	BirthLocation               string      `json:"birth_location" validate:"required"`
	PassportPickUpPoint         string      `json:"passport_pick_up_point" validate:"required"`
	PassportAuthority           string      `json:"passport_authority" validate:"required"`
	PassportAuthorityDate       string      `json:"passport_authority_date" validate:"required,datetime=02-01-2006 15:04:05,min=1"`
	PassportRegistrationAddress string      `json:"passport_registration_address" validate:"required"`
	Inn                         string      `json:"inn" validate:"required,number,min=20,max=1"`
	Accounts                    []uuid.UUID `json:"accounts"`
	Email                       string      `json:"email" validate:"email"`
	UsingTotp                   bool        `json:"using_totp" validate:"boolean"`
}

type AccountInfo struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	Cache      float64   `json:"cache"`
	BIC        string    `json:"bic"`
	CIO        string    `json:"cio"`
	CulcNumber string    `json:"culc_number"`
	CorrNumber string    `json:"corr_number"`
}

type InternalServerInfo struct {
	Host             string
	Port             string
	NumRetry         int
	WaitTimeRetry    time.Duration
	TimeWaitResponse time.Duration
}

type GetUserDataByLoginBody struct {
	UserLogin string `json:"user_login"`
}

type GetOperationResultBody struct {
	OperationId uuid.UUID `json:"operation_id"`
}

type CreateUserBodyUserData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CreateUserBodyPassport struct {
	Series              string `json:"series"`
	Number              string `json:"number"`
	Name                string `json:"name"`
	Surname             string `json:"surname"`
	Patronymic          string `json:"patronymic"`
	BirthDate           string `json:"birth_date"`
	BirthLocation       string `json:"birth_location"`
	PickUpPoint         string `json:"pick_up_point"`
	Authority           string `json:"authority"`
	AuthorityDate       string `json:"authority_date"`
	RegistrationAddress string `json:"registration_address"`
}

type CreateUserBody struct {
	UserInn   string                  `json:"user_inn"`
	UserData  *CreateUserBodyUserData `json:"user_data"`
	Passport  *CreateUserBodyPassport `json:"passport"`
	UserEmail string                  `json:"user_email"`
}

type CheckUserPasswordBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetUserDataBoyd struct {
	UserId uuid.UUID `json:"user_id"`
}

type GetAccountDataBody struct {
	UserId    uuid.UUID `json:"user_id"`
	AccountId uuid.UUID `json:"account_id"`
}

type CloseAccountBody struct {
	UserId    uuid.UUID `json:"user_id"`
	AccountId uuid.UUID `json:"account_id"`
}

type AddAccountCacheBody struct {
	UserId    uuid.UUID `json:"user_id"`
	AccountId uuid.UUID `json:"account_id"`
	CacheDiff float64   `json:"cache_diff"`
}

type OpenAccountBody struct {
	UserId        uuid.UUID `json:"user_id"`
	AccName       string    `json:"acc_name" validate:"required"`
	CulcNumber    string    `json:"culc_number" validate:"required"`
	CorrNumber    string    `json:"corr_number" validate:"required"`
	BIC           string    `json:"bic" validate:"required"`
	CIO           string    `json:"cio" validate:"required"`
	ReserveReason string    `json:"reserve_reason"`
}

type OperationResponse struct {
	Status         int         `json:"status"`
	Info           string      `json:"info"`
	AdditionalInfo interface{} `json:"additional_info"`
}

type OperationListRequestBody struct {
	TimeBegin string `json:"time_begin"`
	TimeEnd   string `json:"time_end"`
}

type OperationListRequestResultBody struct {
	Status string      `json:"status"`
	Info   interface{} `json:"info"`
}

type OperationTreeRequestBody struct {
	OperationId uuid.UUID `json:"operation_id"`
}

type OperationTreeRequestResultBody struct {
	Status string      `json:"status"`
	Info   interface{} `json:"info"`
}

type ListOfOperations struct {
	Operations []uuid.UUID
}

type SagaTree struct {
	Id     uuid.UUID
	Status float64
	Name   string
	Events []uuid.UUID
}

type EventTree struct {
	Id         uuid.UUID
	Name       string
	Status     float64
	RollBackId uuid.UUID
}

type SagaDependTree struct {
	ParentId uuid.UUID
	ChildId  uuid.UUID
}

type OperationTree struct {
	OperationName  string
	SagaList       map[uuid.UUID]*SagaTree
	EventList      map[uuid.UUID]*EventTree
	SagaDependList []*SagaDependTree
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

type GetUserNotifySettingsResponse struct {
	UserId     uuid.UUID `json:"user_id"`
	Email      string    `json:"email"`
	EmailUsage bool      `json:"email_usage"`
}

type GetAccountDataResponse struct {
	Acc_uuid         uuid.UUID `json:"acc_uuid" db:"acc_uuid" validate:"len=36 required unique"`
	Acc_name         string    `json:"acc_name" db:"acc_name" validate:"required"`
	Acc_status       uint8     `json:"acc_status" db:"acc_status" validate:"min=0 max=255 required oneof=0 10 22 30 40 50 255"`
	Acc_culc_number  string    `json:"acc_culc_number" db:"acc_culc_number" validate:"len=20 required"`
	Acc_corr_number  string    `json:"acc_corr_number" db:"acc_corr_number" validate:"len=20 required"`
	Acc_bic          string    `json:"acc_bic" db:"acc_bic" validate:"len=9 required"`
	Acc_cio          string    `json:"acc_cio" db:"acc_cio" validate:"len=9 required"`
	Acc_money_value  uint8     `json:"acc_money_value" db:"acc_money_value" validate:"min=0 max=1 oneof=0 1 2 3"`
	Acc_money_amount float64   `json:"acc_money_amount" db:"acc_money_amount" validate:"required"`
	Reason           string    `json:"reserve_reason" db:"reserve_reason" validate:"required"`
}

type GetUserDataByLoginResponse struct {
	UserId uuid.UUID `json:"user_id"`
}

type CheckPasswordResponse struct {
	UserId    uuid.UUID `json:"user_id"`
	TotpUsage bool      `json:"totp_usage"`
}

type TotpEnrollRequestBody struct {
	UserId   uuid.UUID `json:"user_id"`
	UserName string    `json:"user_name"`
}

type TotpEnrollResponse struct {
	TotpId     uuid.UUID `json:"totp_id"`
	TotpSecret string    `json:"totp_secret"`
	TotpUrl    string    `json:"totp_url"`
}

type TotpInfo struct {
	TotpId    uuid.UUID `json:"totp_id"`
	TotpUrl   string    `json:"totp_url"`
	TotpUsage bool      `json:"totp_usage"`
}

type UpdateTotpUsersInfoBody struct {
	UserId    uuid.UUID `json:"user_id"`
	TotpId    uuid.UUID `json:"totp_id"`
	TotpUsage bool      `json:"totp_usage"`
}

type GetUserTotpInfoBody struct {
	UserId uuid.UUID `json:"user_id"`
}

type TotpDisactivateRequestBody struct {
	UserId uuid.UUID `json:"user_id"`
	TotpId uuid.UUID `json:"totp_id"`
}

type GetTotpUrlBody struct {
	UserId uuid.UUID `json:"user_id"`
}

type TotpCodeValidateBody struct {
	UserId   uuid.UUID `json:"user_id"`
	TotpCode string    `json:"totp_code"`
}
