package models

type Passport struct {
	Series               string `json:"series" validate:"required,len=4,number"`
	Number               string `json:"number" validate:"required,len=6,number"`
	Name                 string `json:"name" validate:"required,min=1"`
	Surname              string `json:"surname" validate:"required,min=1"`
	Patronymic           string `json:"patronymic" validate:"required,min=1"`
	Birth_date           string `json:"birth_date" validate:"required,datetime=02-01-2006 15:04:05,min=1"`
	Birth_location       string `json:"birth_location" validate:"required,min=1"`
	Pick_up_point        string `json:"pick_up_point" validate:"required,min=1"`
	Authority            string `json:"authority" validate:"required,len=7"`
	Authority_date       string `json:"authority_date" validate:"required,datetime=02-01-2006 15:04:05,min=1"`
	Registration_address string `json:"registration_address" validate:"required,min=1"`
}

type UserData struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegistrationUserInfo struct {
	User_INN string   `json:"user_inn" validate:"required,len=20,number"`
	Passport Passport `json:"passport" validate:"required"`
	User     UserData `json:"user_data" validate:"required"`
}

type OperationID struct {
	Operation_ID string `json:"operation_id" validate:"required,len=36,uuid4"`
}

type OpenAccountInfo struct {
	User_ID        string `json:"user_id" validate:"required,uuid4"`
	Culc_number    string `json:"culc_number" validate:"required,number,len=20"`
	Corr_number    string `json:"corr_number" validate:"required,number,len=20"`
	BIC            string `json:"bic" validate:"required,number,len=9"`
	CIO            string `json:"cio" validate:"required,number,min=9,max=11"`
	Reserve_reason string `json:"reserve_reason"`
}

type AddAccountCache struct {
	User_ID    string  `json:"user_id" validate:"required,uuid4"`
	Account_ID string  `json:"account_id" validate:"required,uuid4"`
	Cache_diff float32 `json:"cache_diff" validate:"required,number,min=0"`
}

type WidthAccountCache struct {
	User_ID    string  `json:"user_id" validate:"required,uuid4"`
	Account_ID string  `json:"account_id" validate:"required,uuid4"`
	Cache_diff float32 `json:"cache_diff" validate:"required,number,min=0"`
}

type CloseAccount struct {
	User_ID    string `json:"user_id" validate:"required,uuid4"`
	Account_ID string `json:"account_id" validate:"required,uuid4"`
}

type GetUserData struct {
	User_ID string `json:"user_id" validate:"required,uuid4"`
}

type GetAccountData struct {
	User_ID    string `json:"user_id" validate:"required,uuid4"`
	Account_ID string `json:"account_id" validate:"required,uuid4"`
}

type UpdateUserPassword struct {
	User_ID      string `json:"user_id" validate:"required,uuid4"`
	New_password string `json:"new_password" validate:"required"`
}

type GetUserDataByLogin struct {
	User_login string `json:"user_login" validate:"required"`
}

type CheckUserPassword struct {
	User_ID  string `json:"user_id" validate:"required,uuid4"`
	Password string `json:"password" validate:"required"`
}
