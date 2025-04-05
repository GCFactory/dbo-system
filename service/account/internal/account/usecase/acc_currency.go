package usecase

var PossibleAccMoneyValueStr = map[string]uint8{
	"000": AccMoneyValueUnknown,
	"840": AccMoneyValueDollar,
	"810": AccMoneyValueRub,
	"978": AccMoneyValueEuro,
	"999": AccMoneyValueError,
}
