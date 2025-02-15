package usecase

const (
	AdditionalCheckUserHasAccount      string = "additional_check_user_has_account"
	AdditionalCheckAccountStatusIsOpen string = "additional_check_account_status_is_open"
)

var PossibleChecks = []string{
	AdditionalCheckUserHasAccount,
	AdditionalCheckAccountStatusIsOpen,
}

func ValidateAdditionalCheckName(check_name string) (result bool) {

	result = false

	for _, existing_check_name := range PossibleChecks {

		if existing_check_name == check_name {

			result = true
			break

		}

	}

	return result

}

var ListOfRequiredCheckData = map[string][]string{
	AdditionalCheckUserHasAccount: []string{
		"accounts",
		"acc_id",
	},
	AdditionalCheckAccountStatusIsOpen: []string{
		"acc_status",
	},
}

func ValidateCheckData(check_name string, data map[string]interface{}) (result bool) {

	result = true

	if ValidateAdditionalCheckName(check_name) {

		check_fields, ok := ListOfRequiredCheckData[check_name]
		if ok {

			for _, check_field := range check_fields {

				_, is_exist := data[check_field]

				if !is_exist {
					result = false
					break
				}

			}

		} else {
			result = false
		}

	} else {
		result = false
	}

	return result

}

var ListOfAdditionalSagaEventsChecks = map[uint8]map[string][]string{
	SagaGroupAddAccountCache: map[string][]string{
		EventTypeGetAccountData: []string{
			AdditionalCheckUserHasAccount,
		},
		EventTypeAddAccountCache: []string{
			AdditionalCheckAccountStatusIsOpen,
		},
	},
}

func AdditionalValidation(saga_group uint8, event_type string, data map[string]interface{}) (err error) {

	err = nil

	if ValidateSagaGroup(saga_group) {

		if ValidateEventType(event_type) {

			saga_group_additional_validation, ok := ListOfAdditionalSagaEventsChecks[saga_group]
			if ok {

				event_additional_checks, is_exist := saga_group_additional_validation[event_type]
				if is_exist {

					for _, additional_check_name := range event_additional_checks {

						switch additional_check_name {
						case AdditionalCheckAccountStatusIsOpen:
							{
								if !CheckAccountIsOpen(data) {
									return ErrorCheckAccountIsOpen
								}
							}
						case AdditionalCheckUserHasAccount:
							{
								if !CheckUserHasAccount(data) {
									return ErrorCheckUserHasAccount
								}
							}
						default:
							{
								return ErrorUnknownAdditionalValidationFunction
							}

						}

					}

				}

			}

		} else {

			err = ErrorInvalidEventType

		}

	} else {

		err = ErrorInvalidSagGroup

	}

	return err

}

func CheckAccountIsOpen(data map[string]interface{}) (result bool) {

	result = false

	if ValidateCheckData(AdditionalCheckAccountStatusIsOpen, data) {

		account_status := int(data["acc_status"].(float64))

		if account_status == 30 {
			result = true
		}

	}

	return result

}

func CheckUserHasAccount(data map[string]interface{}) (result bool) {

	result = false

	if ValidateCheckData(AdditionalCheckUserHasAccount, data) {

		user_accounts := data["accounts"].([]interface{})
		account_id := data["acc_id"].(string)

		for _, user_account := range user_accounts {

			if user_account.(string) == account_id {
				result = true
				break
			}

		}

	}

	return result

}
