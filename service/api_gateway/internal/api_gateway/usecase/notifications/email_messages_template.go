package notifications

const (
	NotificationWelcome     string = "Welcome to dbo-system, {{.Login}}({{.Name}})!"
	NotificationSignIn      string = "Dear {{.Login}}, someone log in to your account. If it's not you - call us!"
	NotificationTurnOnTotp  string = "Dear {{.Login}}, you connect totp check!"
	NotificationTurnOffTotp string = "Dear {{.Login}}, you disconnect totp check! If it's not you - call us!"
)
