package models

type RequestData struct {
	Port string
}
type HtmlSignInRequests struct {
	SignInRequest     string
	HomePageRequest   string
	SignUpPageRequest string
}

type HtmlErrorPage struct {
	Message           string
	SignInPageRequest string
}

type HtmlSignUpPage struct {
	SignUpRequest   string
	HomePageRequest string
}

type SignInInfo struct {
	Login    string `json:"login"  validate:"required" msg:"Error validation login"`
	Password string `json:"password" validate:"required" msg:"Error validation password"`
}

type AccountInfoRequest struct {
	AccountId string `json:"account_id" validate:"required" `
}

type SignUpInfo struct {
	Login                       string `json:"login" validate:"required" msg:"Error validation login"`
	Password                    string `json:"password" validate:"required" msg:"Error validation password"`
	PassportSeries              string `json:"passport_series" validate:"required" msg:"Error validation passport_series"`
	PassportNumber              string `json:"passport_number" validate:"required" msg:"Error validation passport_number"`
	Name                        string `json:"name" validate:"required" msg:"Error validation name"`
	Surname                     string `json:"surname" validate:"required" msg:"Error validation surname"`
	Patronymic                  string `json:"patronymic" validate:"required" msg:"Error validation patronymic"`
	BirthDate                   string `json:"birth_date" validate:"required,datetime=02-01-2006,min=1" msg:"Error validation birth_date"`
	BirthLocation               string `json:"birth_location" validate:"required" msg:"Error validation birth_location"`
	PassportPickUpPoint         string `json:"passport_pick_up_point" validate:"required" msg:"Error validation passport_pick_up_point"`
	PassportAuthority           string `json:"passport_authority" validate:"required" msg:"Error validation passport_authority"`
	PassportAuthorityDate       string `json:"passport_authority_date" validate:"required,datetime=02-01-2006,min=1" msg:"Error validation passport_authority_date"`
	PassportRegistrationAddress string `json:"passport_registration_address" validate:"required" msg:"Error validation passport_registration_address"`
	Inn                         string `json:"inn" validate:"required,number,min=20,max=20" msg:"Error validation inn"`
}

type PostRequestStatus struct {
	Success bool   `json:"success"`
	Info    string `json:"info"`
	Error   string `json:"error"`
}

type HomePage struct {
	Login                string
	SignOutRequest       string
	SignInPageRequest    string
	Surname              string
	Name                 string
	Patronymic           string
	INN                  string
	PassportCode         string
	BirthDate            string
	BirthLocation        string
	PickUpPoint          string
	Authority            string
	AuthorityDate        string
	RegistrationAddress  string
	CreateAccountRequest string
	UserId               string
	ListOfAccounts       string
}

type HomePageAccountDescription struct {
	Name                string
	Status              string
	Cache               string
	AccountId           string
	GetCreditsRequest   string
	AddCacheRequest     string
	ReduceCacheRequest  string
	CloseAccountRequest string
}

type AccountOperationPage struct {
	Login             string
	SignOutRequest    string
	SignInPageRequest string
	OperationName     string
	Operation         string
	ReturnRequest     string
}

type AccountOperationData struct {
	OperationRequest string
	AccountId        string
}

type AccountOperationCreditsData struct {
	Name       string
	Status     string
	Amount     float64
	CulcNumber string
	CorrNumber string
	BIC        string
	CIO        string
}
