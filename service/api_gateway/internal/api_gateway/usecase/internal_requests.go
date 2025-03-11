package usecase

var (
	RequestCheckSignIn        string = "http://{{.Host}}:{{.Port}}/api/v1/registration/check_password"
	RequestGetUserDataByLogin string = "http://{{.Host}}:{{.Port}}/api/v1/registration/user_data_by_login"
	RequestGetOperationResult string = "http://{{.Host}}:{{.Port}}/api/v1/registration/get_operation_status"
	RequestCreateUser         string = "http://{{.Host}}:{{.Port}}/api/v1/registration/create_user"
)
