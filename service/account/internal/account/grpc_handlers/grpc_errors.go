package grpc_handlers

import (
	"errors"
	"github.com/GCFactory/dbo-system/service/account/internal/account/usecase"
)

var (
	ErrorWrongOperationName = errors.New("Wrong operation name!")
)

var GRPCErrors = map[uint32]error{
	1: usecase.ErrorNotEnoughMoneyAmount,

	10: usecase.ErrorWrongAccReservedStatus,
	11: usecase.ErrorWrongAccCreatedStatus,
	12: usecase.ErrorWrongAccOpenStatus,
	13: usecase.ErrorWrongAccCloseStatus,
	14: usecase.ErrorWrongAccBlockStatus,

	20: usecase.ErrorAccIsExisting,
	21: usecase.ErrorNoFoundAcc,

	100: ErrorWrongOperationName,

	110: usecase.ErrorWrongCulcNumberLen,

	120: usecase.ErrorWrongBICLen,

	130: usecase.ErrorWrongAccCorrNumberLen,
	131: usecase.ErrorWrongAccCorrNumber,

	140: usecase.ErrorWrongAccKPPLen,
	141: usecase.ErrorWrongAccKPP,

	150: usecase.ErrorWrongOwnerLen,
	151: usecase.ErrorWrongOwner,

	160: usecase.ErrorWrongActivityLen,
	161: usecase.ErrorWrongActivity,

	170: usecase.ErrorWrongCurrencyLen,
	171: usecase.ErrorWrongCurrencyValue,

	180: usecase.ErrorWrongPaymentSystemLen,
	181: usecase.ErrorWrongPaymentSystem,

	190: usecase.ErrorWrongAccCountryLen,
	191: usecase.ErrorWrongAccCountry,

	200: usecase.ErrorWrongAccCountryRegionLen,
	201: usecase.ErrorWrongAccCountryRegion,

	210: usecase.ErrorWrongAccMainOfficeLen,
	211: usecase.ErrorWrongAccMainOffice,

	220: usecase.ErrorWrongAccBankNumberLen,
	221: usecase.ErrorWrongAccBankNumber,

	230: usecase.ErrorWrongAccCorrFirstGroupNumber,

	240: usecase.ErrorWrongAccCorrBankNumber,

	250: usecase.ErrorWrongAccCorrOwnerLen,
	251: usecase.ErrorWrongAccCorrOwner,

	1000: usecase.ErrorReasonIsExisting,

	1010: usecase.ErrorOverflowAmount,

	1020: usecase.ErrorUpdateAmountValue,
	1021: usecase.ErrorUpdateAccStatus,

	1022: usecase.ErrorCreateAcc,
}
