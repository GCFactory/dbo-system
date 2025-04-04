package models

type WelcomeMessage struct {
	Login string
	Name  string
}

type SignInMessage struct {
	Login string
}

type TurnOnTotpMessage struct {
	Login string
}

type TurnOffTotpMessage struct {
	Login string
}
