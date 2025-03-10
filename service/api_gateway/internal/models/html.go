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
