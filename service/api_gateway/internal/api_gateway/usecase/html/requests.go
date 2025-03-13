package html

var (
	RequestSignIn                string = "http://localhost:{{.Port}}/api/v1/api_gateway/sign_in/sign_in"
	RequestSignUpPage            string = "http://localhost:{{.Port}}/api/v1/api_gateway/sign_up"
	RequestSignInPage            string = "http://localhost:{{.Port}}/api/v1/api_gateway/sign_in"
	RequestSignUp                string = "http://localhost:{{.Port}}/api/v1/api_gateway/sign_up/sign_up"
	RequestSignOut               string = "http://localhost:{{.Port}}/api/v1/api_gateway/sign_out"
	RequestOpenAccountPage       string = "http://localhost:{{.Port}}/api/v1/api_gateway/open_account"
	RequestUserPage              string = "http://localhost:{{.Port}}/api/v1/api_gateway/main_page"
	RequestOpenAccount           string = "http://localhost:{{.Port}}/api/v1/api_gateway/open_account/open_account"
	RequestAccountCreditsPage    string = "http://localhost:{{.Port}}/api/v1/api_gateway/get_account_info"
	RequestAccountClosePage      string = "http://localhost:{{.Port}}/api/v1/api_gateway/close_account"
	RequestAddAccountCachePage   string = "http://localhost:{{.Port}}/api/v1/api_gateway/adding_account"
	RequestWidthAccountCachePage string = "http://localhost:{{.Port}}/api/v1/api_gateway/width_account"
	RequestAccountClose          string = "http://localhost:{{.Port}}/api/v1/api_gateway/close_account/close_account"
	RequestAddAccountCache       string = "http://localhost:{{.Port}}/api/v1/api_gateway/adding_account/adding_account"
	RequestWidthAccountCache     string = "http://localhost:{{.Port}}/api/v1/api_gateway/width_account/width_account"
)
