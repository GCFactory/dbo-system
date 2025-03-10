package html

var (
	RequestSignIn     string = "http://localhost:{{.Port}}/api/v1/api_gateway/sign_in/sign_in"
	RequestSignUpPage string = "http://localhost{{.Port}}/api/v1/api_gateway/sign_up"
	RequestSignInPage string = "http://localhost:{{.Port}}/api/v1/api_gateway/sign_in"
	RequestSignUp     string = "http://localhost:{{.Port}}/api/v1/api_gateway/sign_up/sign_up"
)
