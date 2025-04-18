package grpc

import "slices"

const (
	ServerTypeUnknown      uint8 = 0
	ServerTypeUsers        uint8 = 1
	ServerTypeAccounts     uint8 = 2
	ServerTypeNotification uint8 = 3
)

var PossibleServerTypes = []uint8{
	ServerTypeUsers,
	ServerTypeAccounts,
	ServerTypeNotification,
}

func ValidateServer(serverType uint8) (res bool) {
	res = false

	if slices.Contains(PossibleServerTypes, serverType) {
		res = true
	}

	return res
}

const (
	ServerTopicUsersConsumer        string = "users_cons"
	ServerTopicUsersProducerRes     string = "users_res"
	ServerTopicUsersProducerErr     string = "users_err"
	ServerTopicAccountsConsumer     string = "account_cons"
	ServerTopicAccountsProducerRes  string = "account_res"
	ServerTopicAccountsProducerErr  string = "account_err"
	ServerTopicNotificationConsumer string = "notification_cons"
	ServerTopicNotificationRes      string = "notification_res"
	ServerTopicNotificationErr      string = "notification_err"
)

var PosiibleServerTopics = map[uint8][]string{
	ServerTypeUsers: {
		ServerTopicUsersConsumer,
		ServerTopicUsersProducerRes,
		ServerTopicUsersProducerErr,
	},
	ServerTypeAccounts: {
		ServerTopicAccountsConsumer,
		ServerTopicAccountsProducerRes,
		ServerTopicAccountsProducerErr,
	},
	ServerTypeNotification: {
		ServerTopicNotificationConsumer,
		ServerTopicNotificationRes,
		ServerTopicNotificationErr,
	},
}

func ValidateServerTopic(serverType uint8, topic string) (res bool) {
	res = false

	if ValidateServer(serverType) {
		server_topics, ok := PosiibleServerTopics[serverType]
		if ok {
			if slices.Contains(server_topics, topic) {
				res = true
			}
		}
	}

	return res
}

const (
	OperationGetUserData                    string = "get_user_data"
	OperationCreateUser                     string = "add_user"
	OperationAddAccountToUser               string = "add_user_account"
	OperationReserveAcc                     string = "reserve_acc"
	OperationCreateAcc                      string = "create_acc"
	OperationOpenAcc                        string = "open_acc"
	OperationGetAccountData                 string = "get_acc_data"
	OperationAddAccountCache                string = "adding_acc"
	OperationWidthAccountCache              string = "width_acc"
	OperationCloseAccount                   string = "close_acc"
	OperationRemoveAccount                  string = "remove_acc"
	OperationRemoveUserAccount              string = "remove_user_account"
	OperationUpdateUserPassword             string = "update_user_password"
	OperationGetUserDataByLogin             string = "get_user_data_by_login"
	OperationCheckUSerPassword              string = "check_user_password"
	OperationCreateUserNotificationSettings string = "add_user_notification_settings"
	OperationDeleteUserNotificationSettings string = "remove_user_notification_settings"
	OperationDeleteUser                     string = "remove_user"
)

var PossibleServersOperations = map[uint8][]string{
	ServerTypeUsers: {
		OperationCreateUser,
		OperationGetUserData,
		OperationAddAccountToUser,
		OperationRemoveUserAccount,
		OperationUpdateUserPassword,
		OperationGetUserDataByLogin,
		OperationCheckUSerPassword,
		OperationDeleteUser,
	},
	ServerTypeAccounts: {
		OperationReserveAcc,
		OperationCreateAcc,
		OperationOpenAcc,
		OperationGetAccountData,
		OperationAddAccountCache,
		OperationWidthAccountCache,
		OperationCloseAccount,
		OperationRemoveAccount,
	},
	ServerTypeNotification: {
		OperationCreateUserNotificationSettings,
		OperationDeleteUserNotificationSettings,
	},
}

func ValidateServersOperation(serverType uint8, operationName string) (res bool) {
	res = false

	if ValidateServer(serverType) {
		server_operations, ok := PossibleServersOperations[serverType]
		if ok {
			if slices.Contains(server_operations, operationName) {
				res = true
			}
		}

	}

	return res
}

func GetServerByOperation(operationName string) (server uint8) {
	server = 0

	for server_type, server_operations := range PossibleServersOperations {
		if slices.Contains(server_operations, operationName) {
			server = server_type
			break
		}
	}

	return server
}

// Список данных, необходимых для операции
var RequiredOperationsFields = map[string][]string{
	OperationCreateUser: {
		"user_inn",
		"passport_number",
		"passport_series",
		"name",
		"surname",
		"patronimic",
		"birth_date",
		"birth_location",
		"pick_up_point",
		"authority",
		"authority_date",
		"registration_adress",
		"login",
		"password",
	},
	OperationGetUserData: {
		"user_id",
	},
	OperationReserveAcc: {
		"acc_name",
		"culc_number",
		"corr_number",
		"bic",
		"cio",
		"reserve_reason",
	},
	OperationCreateAcc: {
		"acc_id",
	},
	OperationOpenAcc: {
		"acc_id",
	},
	OperationAddAccountToUser: {
		"user_id",
		"acc_id",
	},
	OperationGetAccountData: {
		"acc_id",
	},
	OperationAddAccountCache: {
		"acc_id",
		"cache_diff",
	},
	OperationWidthAccountCache: {
		"acc_id",
		"cache_diff",
	},
	OperationCloseAccount: {
		"acc_id",
	},
	OperationRemoveAccount: {
		"acc_id",
	},
	OperationRemoveUserAccount: {
		"user_id",
		"acc_id",
	},
	OperationUpdateUserPassword: {
		"user_id",
		"new_password",
	},
	OperationGetUserDataByLogin: {
		"user_login",
	},
	OperationCheckUSerPassword: {
		"user_id",
		"password",
	},
	OperationCreateUserNotificationSettings: {
		"user_id",
		"email_notification",
		"email",
	},
	OperationDeleteUserNotificationSettings: {
		"user_id",
	},
	OperationDeleteUser: {
		"user_id",
	},
}

func ValidateOperationsData(operationName string, data map[string]interface{}) (res bool) {
	res = true

	operation_fields, ok := RequiredOperationsFields[operationName]
	if ok {
		for _, field := range operation_fields {
			if _, ok := data[field]; ok == false {
				res = false
				break
			}
		}
	}

	return res
}
