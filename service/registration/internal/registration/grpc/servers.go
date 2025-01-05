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
	ServerTopicUsersConsumer    string = "users_cons"
	ServerTopicUsersProducerRes string = "users_res"
	ServerTopicUsersProducerErr string = "users_err"
)

var PosiibleServerTopics = map[uint8][]string{
	ServerTypeUsers: {
		ServerTopicUsersConsumer,
		ServerTopicUsersProducerRes,
		ServerTopicUsersProducerErr,
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
	OperationCreateUser = "add_user"
)

var PossibleServersOperations = map[uint8][]string{
	ServerTypeUsers: {
		OperationCreateUser,
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
	},
}

func ValidateOperationsData(operationName string, data map[string]interface{}) (res bool) {
	res = true

	operation_fields, ok := RequiredOperationsFields[operationName]
	if ok {
		for field, _ := range data {
			if !slices.Contains(operation_fields, field) {
				res = false
				break
			}
		}
	}

	return res
}
