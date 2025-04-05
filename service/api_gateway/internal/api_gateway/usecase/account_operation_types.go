package usecase

const (
	AccountOperationTypeOpen       string = "Open account"
	AccountOperationTypeGetCredits string = "Credits"
	AccountOperationTypeClose      string = "Close account"
	AccountOperationAddCache       string = "Add cache"
	AccountOperationWidthCache     string = "Width cache"
)

var RequiredOperationTypeData = map[string][]string{
	AccountOperationTypeOpen: {},
	AccountOperationTypeGetCredits: {
		"account_id",
		"user_id",
	},
	AccountOperationTypeClose: {
		"account_id",
		"user_id",
	},
	AccountOperationAddCache: {
		"account_id",
		"user_id",
	},
	AccountOperationWidthCache: {
		"account_id",
		"user_id",
	},
}

func ValidateOperationTypeData(operation_name string, operation_data map[string]interface{}) bool {

	result := true

	if required_data, ok := RequiredOperationTypeData[operation_name]; ok {

		for _, required_data_name := range required_data {
			if _, is_exist := operation_data[required_data_name]; !is_exist {
				result = false
				break
			}
		}

	} else {
		result = false
	}

	return result

}
