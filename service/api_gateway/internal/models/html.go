package models

type RequestData struct {
	Port string
}
type HtmlSignInRequests struct {
	SignInRequest     string
	SignUpPageRequest string
}

type HtmlErrorPage struct {
	Message           string
	SignInPageRequest string
}

type HtmlSignUpPage struct {
	SignUpRequest string
}

type SignInInfo struct {
	Login    string `json:"login" validation:"required"`
	Password string `json:"password" validation:"required"`
}

type SignUpInfo struct {
	Login                       string `json:"login" validation:"required"`
	Password                    string `json:"password" validation:"required"`
	PassportSeries              string `json:"passport_series" validation:"required"`
	PassportNumber              string `json:"passport_number" validation:"required"`
	Name                        string `json:"name" validation:"required"`
	Surname                     string `json:"surname" validation:"required"`
	Patronymic                  string `json:"patronymic" validation:"required"`
	BirthDate                   string `json:"birth_date" validation:"required,datetime=02-01-2006 15:04:05,min=1"`
	BirthLocation               string `json:"birth_location" validation:"required"`
	PassportPickUpPoint         string `json:"passport_pick_up_point" validation:"required"`
	PassportAuthority           string `json:"passport_authority" validation:"required"`
	PassportAuthorityDate       string `json:"passport_authority_date" validation:"required,datetime=02-01-2006 15:04:05,min=1"`
	PassportRegistrationAddress string `json:"passport_registration_address" validation:"required"`
	Inn                         string `json:"inn" validation:"required,number,min=20,max=1"`
}
