package usecase

const (
	AccStatusUnknown  uint8 = 0
	AccStatusReserved uint8 = 10
	AccStatusCreated  uint8 = 20
	AccStatusOpen     uint8 = 30
	AccStatusClose    uint8 = 40
	AccStatusBlocked  uint8 = 50
	AccStatusError    uint8 = 255
)

var PossibleAccStatus = [...]uint8{
	AccStatusUnknown,
	AccStatusReserved,
	AccStatusCreated,
	AccStatusOpen,
	AccStatusClose,
	AccStatusBlocked,
	AccStatusError,
}
