package usecase

const (
	AccCountryTest    uint8 = 0
	AccCountryRus     uint8 = 4
	AccCountryUnknown uint8 = 255
)

var PossibleAccCountry = [...]uint8{
	AccCountryRus,
	AccCountryTest,
}
