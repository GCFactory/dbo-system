package usecase

const (
	AccMoneyValueUnknown uint8 = 0
	AccMoneyValueDollar  uint8 = 1
	AccMoneyValueRub     uint8 = 2
	AccMoneyValueEuro    uint8 = 3
	AccMoneyValueError   uint8 = 255
)

var PossibleAccMoneyValue = [...]uint8{
	AccMoneyValueUnknown,
	AccMoneyValueDollar,
	AccMoneyValueRub,
	AccMoneyValueEuro,
	AccMoneyValueError,
}
