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

type RegistrationUserInfo struct {
	User_INN string   `json:"user_inn" validate:"required,len=20,number"`
	Passport Passport `json:"passport" validate:"required"`
}

type OperationID struct {
	Operation_ID string `json:"operation_id" validate:"required,len=36,uuid4"`
}
