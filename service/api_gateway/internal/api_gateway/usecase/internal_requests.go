package usecase

var (
	RequestCheckUserPassword    string = "http://{{.Host}}:{{.Port}}/api/v1/registration/check_password"
	RequestGetUserDataByLogin   string = "http://{{.Host}}:{{.Port}}/api/v1/users/get_user_data_by_login"
	RequestGetOperationResult   string = "http://{{.Host}}:{{.Port}}/api/v1/registration/get_operation_status"
	RequestCreateUser           string = "http://{{.Host}}:{{.Port}}/api/v1/registration/create_user"
	GetUserData                 string = "http://{{.Host}}:{{.Port}}/api/v1/users/get_user_data"
	GetAccountData              string = "http://{{.Host}}:{{.Port}}/api/v1/account/get_account_data"
	RequestOpenAccount          string = "http://{{.Host}}:{{.Port}}/api/v1/registration/open_account"
	RequestCloseAccount         string = "http://{{.Host}}:{{.Port}}/api/v1/registration/close_acc"
	RequestAddAccountCache      string = "http://{{.Host}}:{{.Port}}/api/v1/registration/add_account_cache"
	RequestWidthAccountCache    string = "http://{{.Host}}:{{.Port}}/api/v1/registration/width_account_cache"
	RequestGetListOfOperations  string = "http://{{.Host}}:{{.Port}}/api/v1/registration/get_operations_range"
	RequestGetOperationTree     string = "http://{{.Host}}:{{.Port}}/api/v1/registration/get_operation_tree_data"
	GetUserNotificationSettings string = "http://{{.Host}}:{{.Port}}/api/v1/notification/get_user_notification_settings"
)
