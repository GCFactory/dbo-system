package usecase

const (
	PaymentSystemStraight uint8 = 0
	PaymentSystemIndirect uint8 = 1
	PaymentSystemCB       uint8 = 2
	PaymentSystemUnknown  uint8 = 255
)

var PossiblePaymentSystems = [...]uint8{
	PaymentSystemStraight,
	PaymentSystemIndirect,
	PaymentSystemCB,
}
