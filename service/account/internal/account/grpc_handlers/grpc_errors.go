package grpc_handlers

import (
	"errors"
	"github.com/GCFactory/dbo-system/service/account/internal/account/usecase"
)

var (
	ErrorWrongOperationName = errors.New("Wrong operation name!")
	ErrorInvalidInputData   = errors.New("Wrong input data!")

	GRPCErrors = map[error]uint32{
		usecase.ErrorNotEnoughMoneyAmount: 1,

		usecase.ErrorWrongAccReservedStatus: 10,
		usecase.ErrorWrongAccCreatedStatus:  11,
		usecase.ErrorWrongAccOpenStatus:     12,
		usecase.ErrorWrongAccCloseStatus:    13,
		usecase.ErrorWrongAccBlockStatus:    14,

		usecase.ErrorAccIsExisting: 20,
		usecase.ErrorNoFoundAcc:    21,

		ErrorWrongOperationName: 100,

		usecase.ErrorWrongCulcNumberLen: 110,

		usecase.ErrorWrongBICLen: 120,

		usecase.ErrorWrongAccCorrNumberLen: 130,
		usecase.ErrorWrongAccCorrNumber:    131,

		usecase.ErrorWrongAccKPPLen: 140,
		usecase.ErrorWrongAccKPP:    141,

		usecase.ErrorWrongOwnerLen: 150,
		usecase.ErrorWrongOwner:    151,

		usecase.ErrorWrongActivityLen: 160,
		usecase.ErrorWrongActivity:    161,

		usecase.ErrorWrongCurrencyLen:   170,
		usecase.ErrorWrongCurrencyValue: 171,

		usecase.ErrorWrongPaymentSystemLen: 180,
		usecase.ErrorWrongPaymentSystem:    181,

		usecase.ErrorWrongAccCountryLen: 190,
		usecase.ErrorWrongAccCountry:    191,

		usecase.ErrorWrongAccCountryRegionLen: 200,
		usecase.ErrorWrongAccCountryRegion:    201,

		usecase.ErrorWrongAccMainOfficeLen: 210,
		usecase.ErrorWrongAccMainOffice:    211,

		usecase.ErrorWrongAccBankNumberLen: 220,
		usecase.ErrorWrongAccBankNumber:    221,

		usecase.ErrorWrongAccCorrFirstGroupNumber: 230,

		usecase.ErrorWrongAccCorrBankNumber: 240,

		usecase.ErrorWrongAccCorrOwnerLen: 250,
		usecase.ErrorWrongAccCorrOwner:    251,

		usecase.ErrorReasonIsExisting: 1000,

		usecase.ErrorOverflowAmount: 1010,

		usecase.ErrorUpdateAmountValue: 1020,
		usecase.ErrorUpdateAccStatus:   1021,
		usecase.ErrorCreateAcc:         1022,

		ErrorInvalidInputData: 1030,
	}
)

func GetErrorCode(err error) uint32 {
	var result uint32
	result, ok := GRPCErrors[err]

	if !ok {
		result = 4294967295
	}

	return result
}
