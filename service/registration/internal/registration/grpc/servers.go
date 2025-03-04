package grpc

import "slices"

const (
	ServerTypeUnknown  uint8 = 0
	ServerTypeUsers    uint8 = 1
	ServerTypeAccounts uint8 = 2
)

var PossibleServerTypes = []uint8{
	ServerTypeUsers,
	ServerTypeAccounts,
}

func ValidateServer(serverType uint8) (res bool) {
	res = false

	if slices.Contains(PossibleServerTypes, serverType) {
		res = true
	}

	return res
}

const (
	ServerTopicUsersConsumer       string = "users_cons"
	ServerTopicUsersProducerRes    string = "users_res"
	ServerTopicUsersProducerErr    string = "users_err"
	ServerTopicAccountsConsumer    string = "account_cons"
	ServerTopicAccountsProducerRes string = "account_res"
	ServerTopicAccountsProducerErr string = "account_err"
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
	OperationGetUserData        string = "get_user_data"
	OperationCreateUser         string = "add_user"
	OperationAddAccountToUser   string = "add_user_account"
	OperationReserveAcc         string = "reserve_acc"
	OperationCreateAcc          string = "create_acc"
	OperationOpenAcc            string = "open_acc"
	OperationGetAccountData     string = "get_acc_data"
	OperationAddAccountCache    string = "adding_acc"
	OperationWidthAccountCache  string = "width_acc"
	OperationCloseAccount       string = "close_acc"
	OperationRemoveAccount      string = "remove_acc"
	OperationRemoveUserAccount  string = "remove_user_account"
	OperationUpdateUserPassword string = "update_user_password"
)

var PossibleServersOperations = map[uint8][]string{
	ServerTypeUsers: {
		OperationCreateUser,
		OperationGetUserData,
		OperationAddAccountToUser,
		OperationRemoveUserAccount,
		OperationUpdateUserPassword,
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
